package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"text/template"

	goyaml "github.com/ghodss/yaml"
	jsonnet "github.com/google/go-jsonnet"
	"github.com/panjf2000/ants/v2"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
)

type safeString struct {
	mu     sync.Mutex
	config string
}

var (
	components       string
	clusters         string
	generateDir      string
	clIncludes       string
	clExcludes       string
	allClusterParams map[string]string
)

func fatalErrorCheck(err error, message string) {
	if err != nil {
		log.Fatal().Err(err).Msg(message)
	}
}

func processJsonnet(vm *jsonnet.VM, input string, include string) (string, error) {
	vm.ExtCode("input", input)
	j, err := vm.EvaluateAnonymousSnippet(include, "std.extVar('process')(std.extVar('input'))")
	if err != nil {
		return "Error evaluating jsonnet snippet", err
	}

	// create output file contents in a string first, as a yaml stream
	var o []interface{}
	var outStr string
	fatalErrorCheck(json.Unmarshal([]byte(j), &o), "Error unmarshalling jsonnet output to go slice")
	for _, jObj := range o {
		buf, err := goyaml.Marshal(jObj)
		fatalErrorCheck(err, "Error marshalling jsonnet object to yaml")
		outStr = outStr + string(buf) + "\n---\n"
	}
	return outStr, nil
}

func genProcessCluster(cmd *cobra.Command, clusterName string, p *ants.Pool) {
	log.Debug().Str("cluster", clusterName).Msg("Process cluster")

	// get list of components for cluster
	params := getClusterParams(clusterDir, getCluster(clusterDir, clusterName))
	clusterComponents := gjson.Parse(renderJsonnet(cmd, params, "._components", true, "", "clustercomponents")).Map()

	// get kr8 settings for cluster
	kr8Spec := gjson.Parse(renderJsonnet(cmd, params, "._kr8_spec", false, "", "kr8_spec"))
	postProcessorFunction := kr8Spec.Get("postprocessor").String()

	var clGenerateDir string
	if generateDir == "" {
		clGenerateDir = kr8Spec.Get("generate_dir").String()
		if clGenerateDir == "" {
			log.Fatal().Msg("_kr8_spec.generate_dir must be set in parameters or passed as generate-dir flag")
		}
	} else {
		clGenerateDir = generateDir
	}
	if !strings.HasPrefix(clGenerateDir, "/") {
		// if generateDir does not start with /, then it goes in baseDir
		clGenerateDir = baseDir + "/" + clGenerateDir
	}
	clusterDir := clGenerateDir + "/" + clusterName

	// if this is true, we don't use the full file path to generate output file names
	generateShortNames := kr8Spec.Get("generate_short_names").Bool()

	// if this is true, we prune component parameters
	pruneParams := kr8Spec.Get("prune_params").Bool()

	// create generateDir
	if _, err := os.Stat(clGenerateDir); os.IsNotExist(err) {
		err = os.MkdirAll(clGenerateDir, os.ModePerm)
		if err != nil {
			log.Fatal().Err(err).Msg("")
		}
	}
	// create cluster dir
	if _, err := os.Stat(kr8Spec.ClusterDir); os.IsNotExist(err) {
		err = os.MkdirAll(kr8Spec.ClusterDir, os.ModePerm)
		fatalErrorCheck(err, "Error creating cluster generateDir")
	}

	// list of current generated components directories
	d, err := os.Open(kr8Spec.ClusterDir)
	fatalErrorCheck(err, "Error opening clusterDir")
	defer d.Close()
	read_all_dirs := -1
	generatedCompList, err := d.Readdirnames(read_all_dirs)
	fatalErrorCheck(err, "Error reading directories")

	// determine list of components to process
	var compList []string
	var currentCompList []string

	if components != "" {
		// only process specified component if it's defined in the cluster
		for _, b := range strings.Split(components, ",") {
			for _, c := range generatedCompList {
				matched, _ := regexp.MatchString("^"+b+"$", c)
				if matched {
					currentCompList = append(currentCompList, c)
				}
			}
			for c, _ := range clusterComponents {
				matched, _ := regexp.MatchString("^"+b+"$", c)
				if matched {
					compList = append(compList, c)
				}
			}
		}
	} else {
		for c, _ := range clusterComponents {
			compList = append(compList, c)
		}
		currentCompList = generatedCompList
	}
	sort.Strings(compList) // process components in sorted order

	// Sort out orphaned generated components directories
	tmpMap := make(map[string]struct{}, len(clusterComponents))
	for e, _ := range clusterComponents {
		tmpMap[e] = struct{}{}
	}

	for _, e := range currentCompList {
		if _, found := tmpMap[e]; !found {
			delcomp := filepath.Join(clusterDir, e)
			os.RemoveAll(delcomp)
			log.Info().Str("cluster", clusterName).
				Str("component", e).
				Msg("Deleting generated for component")
		}
	}

	if len(compList) == 0 { // this needs to be moved so purging above works first
		return
	}

	// render full params for cluster for all selected components
	config := renderClusterParams(cmd, clusterName, compList, clusterParams, false)

	var allconfig safeString

	var wg sync.WaitGroup
	//p, _ := ants.NewPool(4)
	for _, componentName := range compList {
		wg.Add(1)
		cName := componentName
		_ = p.Submit(func() {
			defer wg.Done()
			genProcessComponent(cmd, clusterName, cName, clusterDir, clGenerateDir, config, &allconfig, postProcessorFunction, pruneParams, generateShortNames)
		})
	}
	wg.Wait()

}

func genProcessComponent(cmd *cobra.Command, clusterName string, componentName string, clusterDir string, clGenerateDir string, config string, allconfig *safeString, postProcessorFunction string, pruneParams bool, generateShortNames bool) {

	log.Info().Str("cluster", clusterName).
		Str("component", componentName).
		Msg("Process component")

	// get kr8_spec from component's params
	spec := gjson.Get(config, componentName+".kr8_spec").Map()
	compPath := gjson.Get(config, "_components."+componentName+".path").String()

	// spec is missing?
	if len(spec) == 0 {
		log.Fatal().Str("cluster", clusterName).
			Str("component", componentName).
			Msg("Component has no kr8_spec")
		return
	}

	// it's faster to create this VM for each component, rather than re-use
	vm, _ := JsonnetVM(cmd)
	vm.ExtCode("kr8_cluster", "std.prune("+config+"._cluster)")
	//vm.ExtCode("kr8_components", "std.prune("+config+"._components)")
	if postProcessorFunction != "" {
		vm.ExtCode("process", postProcessorFunction)
	} else {
		// default postprocessor just copies input
		vm.ExtCode("process", "function(input) input")
	}

	// prune params if required
	if pruneParams {
		vm.ExtCode("kr8", "std.prune("+config+"."+componentName+")")
	} else {
		vm.ExtCode("kr8", config+"."+componentName)
	}

	// add kr8_allparams extcode with all component params in the cluster
	if spec["enable_kr8_allparams"].Bool() {
		// include full render of all component params
		allconfig.mu.Lock()
		if allconfig.config == "" {
			// only do this if we have not already cached it and don't already have it stored
			if components == "" {
				// all component params are in config
				allconfig.config = config
			} else {
				allconfig.config = renderClusterParams(cmd, clusterName, []string{}, clusterParams, false)
			}
		}
		vm.ExtCode("kr8_allparams", allconfig.config)
		allconfig.mu.Unlock()
	}

	// add kr8_allclusters extcode with every cluster's cluster level params
	if spec["enable_kr8_allclusters"].Bool() {
		// combine all the cluster params into a single object indexed by cluster name
		var allClusterParamsObject string
		allClusterParamsObject = "{ "
		for cl, clp := range allClusterParams {
			allClusterParamsObject = allClusterParamsObject + "'" + cl + "': " + clp + ","

		}
		allClusterParamsObject = allClusterParamsObject + "}"
		vm.ExtCode("kr8_allclusters", allClusterParamsObject)
	}

	// jpath always includes base lib. Add jpaths from spec if set
	jpath := []string{baseDir + "/lib"}
	for _, j := range spec["jpaths"].Array() {
		jpath = append(jpath, baseDir+"/"+compPath+"/"+j.String())
	}
	vm.Importer(&jsonnet.FileImporter{
		JPaths: jpath,
	})

	// file imports
	for k, v := range spec["extfiles"].Map() {
		vPath := baseDir + "/" + compPath + "/" + v.String() // use full path for file
		extFile, err := os.ReadFile(vPath)
		fatalErrorCheck(err, "Error importing extfiles item")
		log.Debug().Str("cluster", clusterName).
			Str("component", componentName).
			Msg("Extfile: " + k + "=" + v.String())
		vm.ExtVar(k, string(extFile))
	}

	componentDir := clusterDir + "/" + componentName
	// create component dir if needed
	if _, err := os.Stat(componentDir); os.IsNotExist(err) {
		err := os.MkdirAll(componentDir, os.ModePerm)
		fatalErrorCheck(err, "Error creating component directory")
	}

	outputFileMap := make(map[string]bool)
	// generate each included file
	for _, include := range spec["includes"].Array() {
		var filename string
		var outputDir string
		var sfile string

		itype := include.Type.String()
		outputDir = componentDir
		if itype == "String" {
			// include is just a string for the filename
			filename = include.String()
		} else if itype == "JSON" {
			// include is a map with multiple fields
			inc_spec := include.Map()
			filename = inc_spec["file"].String()
			if inc_spec["dest_dir"].Exists() {
				// handle alternate output directory for file
				altdir := inc_spec["dest_dir"].String()
				// dir is always relative to generate dir
				outputDir = clGenerateDir + "/" + altdir
				// ensure this directory exists
				if _, err := os.Stat(outputDir); os.IsNotExist(err) {
					err = os.MkdirAll(outputDir, os.ModePerm)
					fatalErrorCheck(err, "Error creating alternate directory")
				}
			}
			if inc_spec["dest_name"].Exists() {
				// override destination file name
				sfile = inc_spec["dest_name"].String()
			}
		}
		file_extension := filepath.Ext(filename)
		if sfile == "" {
			if generateShortNames {
				sbase := filepath.Base(filename)
				sfile = sbase[0 : len(sbase)-len(file_extension)]
			} else {
				// replaces slashes with _ in multi-dir paths and replace extension with yaml
				sfile = strings.ReplaceAll(filename[0:len(filename)-len(file_extension)], "/", "_")
			}
		}
		outputFile := outputDir + "/" + sfile + ".yaml"
		// remember output filename for purging files
		outputFileMap[sfile+".yaml"] = true

		log.Debug().Str("cluster", clusterName).
			Str("component", componentName).
			Msg("Process file: " + filename + " -> " + outputFile)

		var input string
		switch file_extension {
		case ".jsonnet":
			// file is processed as an ExtCode input, so that we can postprocess it
			// in the snippet
			input = "( import '" + baseDir + "/" + compPath + "/" + filename + "')"
		case ".yaml":
			input = "std.native('parseYaml')(importstr '" + baseDir + "/" + compPath + "/" + filename + "')"
		default:
			log.Fatal().Str("cluster", clusterName).
				Str("component", componentName).
				Str("file", filename).
				Msg("Unsupported file extension")
		}

		vm.ExtCode("input", input)
		j, err := vm.EvaluateAnonymousSnippet(include.String(), "std.extVar('process')(std.extVar('input'))")
		if err != nil {
			log.Fatal().Str("cluster", clusterName).
				Str("component", componentName).
				Str("file", filename).Err(err).Msg("Error evaluating jsonnet snippet")
		}

		// create output file contents in a string first, as a yaml stream
		var o []interface{}
		var outStr string
		if err := json.Unmarshal([]byte(j), &o); err != nil {
			log.Fatal().Err(err).Msg("")
		}
		for _, jobj := range o {
			outStr = outStr + "---\n"
			buf, err := goyaml.Marshal(jobj)
			if err != nil {
				log.Fatal().Err(err).Msg("")
			}
			outStr = outStr + string(buf) + "\n"
		}

		// only write file if it does not exist, or the generated contents does not match what is on disk
		var updateNeeded bool
		if _, err := os.Stat(outputFile); os.IsNotExist(err) {
			log.Debug().Str("cluster", clusterName).
				Str("component", componentName).
				Msg("Creating " + outputFile)
			updateNeeded = true
		} else {
			currentContents, err := os.ReadFile(outputFile)
			fatalErrorCheck(err, "Error reading file")
			if string(currentContents) != outStr {
				updateNeeded = true
				log.Debug().Str("cluster", clusterName).
					Str("component", componentName).
					Msg("Updating: " + outputFile)
			}
		}
		if updateNeeded {
			f, err := os.Create(outputFile)
			fatalErrorCheck(err, "Error creating file")
			//defer f.Close()
			_, err = f.WriteString(outStr)
			f.Close()
			fatalErrorCheck(err, "Error writing to file")
		}
	}
	// purge any yaml files in the output dir that were not generated
	if !spec["disable_output_clean"].Bool() {
		// clean component dir
		d, err := os.Open(componentDir)
		fatalErrorCheck(err, "")
		defer d.Close()
		names, err := d.Readdirnames(-1)
		fatalErrorCheck(err, "")
		for _, name := range names {
			if _, ok := outputFileMap[name]; ok {
				// file is managed
				continue
			}
			if filepath.Ext(name) == ".yaml" {
				delFile := filepath.Join(componentDir, name)
				err = os.RemoveAll(delFile)
				fatalErrorCheck(err, "")
				log.Debug().Str("cluster", clusterName).
					Str("component", componentName).
					Msg("Deleted: " + delFile)
			}
		}
		d.Close()
	}
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate components",
	Long:  `Generate components in clusters`,

	Args: cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		var clusterList []string

		// get list of all clusters, render cluster level params for all of them
		allClusterParams = make(map[string]string)
		allClusters, err := getClusters(clusterDir)
		fatalErrorCheck(err, "Error getting list of clusters")
		for _, c := range allClusters.Cluster {
			allClusterParams[c.Name] = renderClusterParamsOnly(cmd, c.Name, "", false)
		}

		for c, _ := range allClusterParams {
			if clIncludes != "" || clExcludes != "" {
				gjresult := gjson.Parse(allClusterParams[c])
				// includes
				if clIncludes != "" {
					// filter on cluster parameters, passed in gjson path notation with either
					// "=" for equality or "~" for regex match
					var include bool
					for _, b := range strings.Split(clIncludes, ",") {
						include = false
						// equality match
						kv := strings.SplitN(b, "=", 2)
						if len(kv) == 2 {
							if gjresult.Get(kv[0]).String() == kv[1] {
								include = true
							}
						} else {
							// regex match
							kv := strings.SplitN(b, "~", 2)
							if len(kv) == 2 {
								matched, _ := regexp.MatchString(kv[1], gjresult.Get(kv[0]).String())
								if matched {
									include = true
								}
							}
						}
						if !include {
							break
						}
					}
					if !include {
						continue
					}
				}
				// excludes
				if clExcludes != "" {
					// filter on cluster parameters, passed in gjson path notation with either
					// "=" for equality or "~" for regex match
					var exclude bool
					for _, b := range strings.Split(clExcludes, ",") {
						exclude = false
						// equality match
						kv := strings.SplitN(b, "=", 2)
						if len(kv) == 2 {
							if gjresult.Get(kv[0]).String() == kv[1] {
								exclude = true
							}
						} else {
							// regex match
							kv := strings.SplitN(b, "~", 2)
							if len(kv) == 2 {
								matched, _ := regexp.MatchString(kv[1], gjresult.Get(kv[0]).String())
								if matched {
									exclude = true
								}
							}
						}
						if exclude {
							break
						}
					}
					if exclude {
						continue
					}
				}
			}

			if clusters == "" {
				// all clusters
				clusterList = append(clusterList, c)
			} else {
				// match --clusters list
				for _, b := range strings.Split(clusters, ",") {
					// match cluster names as anchored regex
					matched, _ := regexp.MatchString("^"+b+"$", c)
					if matched {
						clusterList = append(clusterList, c)
						break
					}
				}

			}
		}

		var wg sync.WaitGroup
		parallel, err := cmd.Flags().GetInt("parallel")
		fatalErrorCheck(err, "Failed to get parallel flag")
		log.Debug().Msg("Parallel set to " + strconv.Itoa(parallel))
		ants_cp, _ := ants.NewPool(parallel)
		ants_cl, _ := ants.NewPool(parallel)

		// Generate config for each cluster in parallel
		for _, clusterName := range clusterList {
			wg.Add(1)
			cl := clusterName
			_ = ants_cl.Submit(func() {
				defer wg.Done()
				genProcessCluster(cmd, cl, ants_cp)
			})
		}
		wg.Wait()
	},
}

func init() {
	RootCmd.AddCommand(generateCmd)
	generateCmd.Flags().StringVarP(&clusterParams, "clusterparams", "", "", "provide cluster params as single file - can be combined with --cluster to override cluster")
	generateCmd.Flags().StringVarP(&clusters, "clusters", "", "", "clusters to generate - comma separated list of cluster names and/or regular expressions ")
	generateCmd.Flags().StringVarP(&components, "components", "", "", "components to generate - comma separated list of component names and/or regular expressions")
	generateCmd.Flags().StringVarP(&generateDir, "generate-dir", "", "", "output directory")
	generateCmd.Flags().StringVarP(&clIncludes, "clincludes", "", "", "filter included cluster by including clusters with matching cluster parameters - comma separate list of key/value conditions separated by = or ~ (for regex match)")
	generateCmd.Flags().StringVarP(&clExcludes, "clexcludes", "", "", "filter included cluster by excluding clusters with matching cluster parameters - comma separate list of key/value conditions separated by = or ~ (for regex match)")
	generateCmd.Flags().IntP("parallel", "", runtime.GOMAXPROCS(0), "parallelism - defaults to GOMAXPROCS")
	viper.BindPFlag("clincludes", generateCmd.PersistentFlags().Lookup("clincludes"))
	viper.BindPFlag("clexcludes", generateCmd.PersistentFlags().Lookup("clexcludes"))
}
