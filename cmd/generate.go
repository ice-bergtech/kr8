package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
	"text/template"

	goyaml "github.com/ghodss/yaml"
	jsonnet "github.com/google/go-jsonnet"
	"github.com/panjf2000/ants/v2"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"
)

type safeString struct {
	mu     sync.Mutex
	config string
}

var (
	allClusterParams map[string]string
)

type cmdGenerateOptions struct {
	ClusterParamsFile string
	Components        string
	Clusters          string
	GenerateDir       string
	Filters           PathFilterOptions
}

var cmdGenerateFlags cmdGenerateOptions

func init() {
	RootCmd.AddCommand(generateCmd)
	generateCmd.Flags().StringVarP(&cmdGenerateFlags.ClusterParamsFile, "clusterparams", "p", "", "provide cluster params as single file - can be combined with --cluster to override cluster")
	generateCmd.Flags().StringVarP(&cmdGenerateFlags.Clusters, "clusters", "C", "", "clusters to generate - comma separated list of cluster names and/or regular expressions ")
	generateCmd.Flags().StringVarP(&cmdGenerateFlags.Components, "components", "c", "", "components to generate - comma separated list of component names and/or regular expressions")
	generateCmd.Flags().StringVarP(&cmdGenerateFlags.GenerateDir, "generate-dir", "o", "generated", "output directory")
	generateCmd.Flags().StringVarP(&cmdGenerateFlags.Filters.Includes, "clincludes", "i", "", "filter included cluster by including clusters with matching cluster parameters - comma separate list of key/value conditions separated by = or ~ (for regex match)")
	generateCmd.Flags().StringVarP(&cmdGenerateFlags.Filters.Excludes, "clexcludes", "x", "", "filter included cluster by excluding clusters with matching cluster parameters - comma separate list of key/value conditions separated by = or ~ (for regex match)")
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate components",
	Long:  `Generate components in clusters`,

	Args: cobra.MinimumNArgs(0),
	Run:  generateCommand,
}

func generateCommand(cmd *cobra.Command, args []string) {

	// get list of all clusters, render cluster level params for all of them
	allClusterParams = make(map[string]string)
	allClusters, err := getClusters(rootConfig.ClusterDir)
	fatalErrorCheck(err, "Error getting list of clusters")
	for _, c := range allClusters.Cluster {
		allClusterParams[c.Name] = renderClusterParamsOnly(rootConfig.VMConfig, c.Name, "", false)
	}

	// This will store the list of clusters to generate components for.
	clusterList := calculateClusterIncludesExcludes(cmdGenerateFlags.Filters)

	// Setup the threading pools, one for clusters and one for clusters
	var wg sync.WaitGroup
	ants_cp, _ := ants.NewPool(rootConfig.Parallel)
	ants_cl, _ := ants.NewPool(rootConfig.Parallel)

	// Generate config for each cluster in parallel
	for _, clusterName := range clusterList {
		wg.Add(1)
		cl := clusterName
		_ = ants_cl.Submit(func() {
			defer wg.Done()
			genProcessCluster(rootConfig.VMConfig, cl, ants_cp)
		})
	}
	wg.Wait()
}

// Using the allClusterParams variable and command flags to create a list of clusters to generate
// Clusters can be filtered with "=" for equality or "~" for regex match
func calculateClusterIncludesExcludes(filters PathFilterOptions) []string {
	var clusterList []string
	for c := range allClusterParams {
		if filters.Includes != "" || filters.Excludes != "" {
			gjResult := gjson.Parse(allClusterParams[c])
			// includes
			if filters.Includes != "" {
				// filter on cluster parameters, passed in gjson path notation with either
				// "=" for equality or "~" for regex match
				var include bool
				for _, b := range strings.Split(filters.Includes, ",") {
					include = false
					// equality match
					kv := strings.SplitN(b, "=", 2)
					if len(kv) == 2 {
						if gjResult.Get(kv[0]).String() == kv[1] {
							include = true
						}
					} else {
						// regex match
						kv := strings.SplitN(b, "~", 2)
						if len(kv) == 2 {
							matched, _ := regexp.MatchString(kv[1], gjResult.Get(kv[0]).String())
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
			if filters.Excludes != "" {
				// filter on cluster parameters, passed in gjson path notation with either
				// "=" for equality or "~" for regex match
				var exclude bool
				for _, b := range strings.Split(filters.Excludes, ",") {
					exclude = false
					// equality match
					kv := strings.SplitN(b, "=", 2)
					if len(kv) == 2 {
						if gjResult.Get(kv[0]).String() == kv[1] {
							exclude = true
						}
					} else {
						// regex match
						kv := strings.SplitN(b, "~", 2)
						if len(kv) == 2 {
							matched, _ := regexp.MatchString(kv[1], gjResult.Get(kv[0]).String())
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

		if cmdGenerateFlags.Clusters == "" {
			// all clusters
			clusterList = append(clusterList, c)
		} else {
			// match --clusters list
			for _, b := range strings.Split(cmdGenerateFlags.Clusters, ",") {
				// match cluster names as anchored regex
				matched, _ := regexp.MatchString("^"+b+"$", c)
				if matched {
					clusterList = append(clusterList, c)
					break
				}
			}

		}
	}
	return clusterList
}

// Only processes specified component if it's defined in the cluster
// Processes components in string sorted order
// Sorts out orphaned, generated components directories
func buildComponentList(generatedCompList []string, clusterComponents map[string]gjson.Result, clusterDir string, clusterName string) []string {
	var compList []string
	var currentCompList []string

	if cmdGenerateFlags.Components != "" {
		for _, b := range strings.Split(cmdGenerateFlags.Components, ",") {
			for _, c := range generatedCompList {
				matched, _ := regexp.MatchString("^"+b+"$", c)
				if matched {
					currentCompList = append(currentCompList, c)
				}
			}
			for c := range clusterComponents {
				matched, _ := regexp.MatchString("^"+b+"$", c)
				if matched {
					compList = append(compList, c)
				}
			}
		}
	} else {
		for c := range clusterComponents {
			compList = append(compList, c)
		}
		currentCompList = generatedCompList
	}
	sort.Strings(compList)

	tmpMap := make(map[string]struct{}, len(clusterComponents))
	for e := range clusterComponents {
		tmpMap[e] = struct{}{}
	}

	for _, e := range currentCompList {
		if _, found := tmpMap[e]; !found {
			delComp := filepath.Join(clusterDir, e)
			os.RemoveAll(delComp)
			log.Info().Str("cluster", clusterName).
				Str("component", e).
				Msg("Deleting generated for component")
		}
	}
	return compList
}

func processJsonnet(vm *jsonnet.VM, input string, snippetFilename string) (string, error) {
	vm.ExtCode("input", input)
	j, err := vm.EvaluateAnonymousSnippet(snippetFilename, "std.extVar('process')(std.extVar('input'))")
	if err != nil {
		return "Error evaluating jsonnet snippet", err
	}

	// create output file contents in a string first, as a yaml stream
	var o []interface{}
	var outStr string
	fatalErrorCheck(json.Unmarshal([]byte(j), &o), "Error unmarshalling jsonnet output to go slice")
	for i, jObj := range o {
		buf, err := goyaml.Marshal(jObj)
		fatalErrorCheck(err, "Error marshalling jsonnet object to yaml")
		if i > 0 {
			outStr = outStr + "\n---\n"
		}
		outStr = outStr + string(buf)
	}
	return outStr, nil
}

func processTemplate(filename string, data map[string]gjson.Result) (string, error) {
	var tInput []byte
	var tmpl *template.Template
	var buffer bytes.Buffer
	var err error
	tInput, err = os.ReadFile(filename)
	if err != nil {
		return "Error loading template", err
	}
	tmpl, err = template.New("file").Parse(string(tInput))
	if err != nil {
		return "Error parsing template", err
	}
	if err = tmpl.Execute(&buffer, data); err != nil {
		return "Error executing templating", err
	}
	return buffer.String(), nil
}

func genProcessCluster(vmConfig VMConfig, clusterName string, p *ants.Pool) {
	log.Debug().Str("cluster", clusterName).Msg("Process cluster")

	// get list of components for cluster
	params := getClusterParams(rootConfig.ClusterDir, getCluster(rootConfig.ClusterDir, clusterName))
	clusterComponents := gjson.Parse(renderJsonnet(vmConfig, params, "._components", true, "", "clustercomponents")).Map()

	// get kr8 settings for cluster
	kr8Spec, err := CreateClusterSpec(clusterName, gjson.Parse(renderJsonnet(vmConfig, params, "._kr8_spec", false, "", "kr8_spec")), rootConfig.BaseDir, cmdGenerateFlags.GenerateDir)
	fatalErrorCheck(err, "Error creating kr8Spec")

	// create cluster dir
	if _, err := os.Stat(kr8Spec.ClusterDir); os.IsNotExist(err) {
		err = os.MkdirAll(kr8Spec.ClusterDir, os.ModePerm)
		fatalErrorCheck(err, "Error creating cluster generateDir")
	}

	// get list of current generated components directories
	d, err := os.Open(kr8Spec.ClusterDir)
	fatalErrorCheck(err, "Error opening clusterDir")
	defer d.Close()
	read_all_dirs := -1
	generatedCompList, err := d.Readdirnames(read_all_dirs)
	fatalErrorCheck(err, "Error reading directories")

	// determine list of components to process
	compList := buildComponentList(generatedCompList, clusterComponents, kr8Spec.ClusterDir, kr8Spec.Name)

	if len(compList) == 0 { // this needs to be moved so purging above works first
		return
	}

	// render full params for cluster for all selected components
	config := renderClusterParams(vmConfig, kr8Spec.Name, compList, cmdGenerateFlags.ClusterParamsFile, false)

	var allconfig safeString

	var wg sync.WaitGroup
	//p, _ := ants.NewPool(4)
	for _, componentName := range compList {
		wg.Add(1)
		cName := componentName
		_ = p.Submit(func() {
			defer wg.Done()
			genProcessComponent(vmConfig, cName, kr8Spec, config, &allconfig)
		})
	}
	wg.Wait()

}

func genProcessComponent(vmconfig VMConfig, componentName string, kr8Spec ClusterSpec, config string, allConfig *safeString) {
	log.Info().Str("cluster", kr8Spec.Name).
		Str("component", componentName).
		Msg("Process component")

	// get kr8_spec from component's params
	//spec := gjson.Get(config, componentName+".kr8_spec").Map()
	compPath := gjson.Get(config, "_components."+componentName+".path").String()
	compSpec, _ := CreateComponentSpec(gjson.Get(config, componentName+".kr8_spec"))

	// it's faster to create this VM for each component, rather than re-use
	vm, _ := JsonnetVM(vmconfig)
	vm.ExtCode("kr8_cluster", "std.prune("+config+"._cluster)")
	//vm.ExtCode("kr8_components", "std.prune("+config+"._components)")
	if kr8Spec.PostProcessor != "" {
		vm.ExtCode("process", kr8Spec.PostProcessor)
	} else {
		// default postprocessor just copies input
		vm.ExtCode("process", "function(input) input")
	}

	// prune params if required
	if kr8Spec.PruneParams {
		vm.ExtCode("kr8", "std.prune("+config+"."+componentName+")")
	} else {
		vm.ExtCode("kr8", config+"."+componentName)
	}

	// add kr8_allparams extcode with all component params in the cluster
	if compSpec.Kr8_allparams {
		// include full render of all component params
		allConfig.mu.Lock()
		if allConfig.config == "" {
			// only do this if we have not already cached it and don't already have it stored
			if cmdGenerateFlags.Components == "" {
				// all component params are in config
				allConfig.config = config
			} else {
				allConfig.config = renderClusterParams(vmconfig, kr8Spec.Name, []string{}, cmdGenerateFlags.ClusterParamsFile, false)
			}
		}
		vm.ExtCode("kr8_allparams", allConfig.config)
		allConfig.mu.Unlock()
	}

	// add kr8_allclusters extcode with every cluster's cluster level params
	if compSpec.Kr8_allclusters {
		// combine all the cluster params into a single object indexed by cluster name
		var allClusterParamsObject string
		allClusterParamsObject = "{ "
		for cl, clp := range allClusterParams {
			allClusterParamsObject = allClusterParamsObject + "'" + cl + "': " + clp + ","

		}
		allClusterParamsObject = allClusterParamsObject + "}"
		vm.ExtCode("kr8_allclusters", allClusterParamsObject)
	}

	// jPath always includes base lib. Add jpaths from spec if set
	jPath := []string{rootConfig.BaseDir + "/lib"}
	for _, j := range compSpec.JPaths {
		jPath = append(jPath, rootConfig.BaseDir+"/"+compPath+"/"+j)
	}
	vm.Importer(&jsonnet.FileImporter{
		JPaths: jPath,
	})

	// file imports
	for k, v := range compSpec.ExtFiles {
		vPath := rootConfig.BaseDir + "/" + compPath + "/" + v // use full path for file
		extFile, err := os.ReadFile(vPath)
		fatalErrorCheck(err, "Error importing extfiles item")
		log.Debug().Str("cluster", kr8Spec.Name).
			Str("component", componentName).
			Msg("Extfile: " + k + "=" + v)
		vm.ExtVar(k, string(extFile))
	}

	componentOutputDir := kr8Spec.ClusterDir + "/" + componentName
	// create component dir if needed
	if _, err := os.Stat(componentOutputDir); os.IsNotExist(err) {
		err := os.MkdirAll(componentOutputDir, os.ModePerm)
		fatalErrorCheck(err, "Error creating component directory")
	}

	incInfo := IncludeFileEntryStruct{
		DestExt:  "yaml",
		DestDir:  "",
		DestName: "",
	}

	outputFileMap := make(map[string]bool)
	// generate each included file
	for _, include := range compSpec.Includes {
		switch include.(type) {
		case string:
			incInfo = IncludeFileEntryStruct{
				File:     include.(string),
				DestExt:  "yaml", // default to yaml ,
				DestDir:  "",
				DestName: "",
			}
		case IncludeFileEntryStruct:
			// include is a map with multiple fields
			incInfo = include.(IncludeFileEntryStruct)
		default:
			log.Fatal().Msg("Invalid include type")
		}
		if incInfo.DestName == "" {
			if kr8Spec.GenerateShortNames {
				sBase := filepath.Base(incInfo.File)
				incInfo.DestName = sBase[0 : len(sBase)-len(filepath.Ext(include.(string)))]
			} else {
				// replaces slashes with _ in multi-dir paths and replace extension with yaml
				incInfo.DestName = strings.ReplaceAll(incInfo.File[0:len(incInfo.File)-len(filepath.Ext(include.(string)))], "/", "_")
			}
		}
		processIncludesFile(vm, config, kr8Spec, componentName, compPath, componentOutputDir, incInfo, outputFileMap)
	}

	// purge any yaml files in the output dir that were not generated
	if !compSpec.DisableOutputDirClean {
		// clean component dir
		d, err := os.Open(componentOutputDir)
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
				delFile := filepath.Join(componentOutputDir, name)
				err = os.RemoveAll(delFile)
				fatalErrorCheck(err, "")
				log.Debug().Str("cluster", kr8Spec.Name).
					Str("component", componentName).
					Msg("Deleted: " + delFile)
			}
		}
		d.Close()
	}
}

func processIncludesFile(vm *jsonnet.VM, config string, kr8Spec ClusterSpec, componentName string, componentPath string, componentOutputDir string, incInfo IncludeFileEntryStruct, outputFileMap map[string]bool) {
	file_extension := filepath.Ext(incInfo.File)

	// ensure this directory exists
	outputDir := componentOutputDir
	if incInfo.DestDir != "" {
		outputDir = kr8Spec.ClusterDir + "/" + incInfo.DestDir
	}
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		err = os.MkdirAll(outputDir, os.ModePerm)
		fatalErrorCheck(err, "Error creating alternate directory")
	}
	outputFile := outputDir + "/" + incInfo.DestName + "." + incInfo.DestExt

	// remember output filename for purging files
	outputFileMap[incInfo.DestName+"."+incInfo.DestExt] = true

	log.Debug().Str("cluster", kr8Spec.Name).
		Str("component", componentName).
		Msg("Process file: " + incInfo.File + " -> " + outputFile)

	var input string
	var outStr string
	var err error
	switch file_extension {
	case ".jsonnet":
		// file is processed as an ExtCode input, so that we can postprocess it
		// in the snippet
		input = "( import '" + rootConfig.BaseDir + "/" + componentPath + "/" + incInfo.File + "')"
		outStr, err = processJsonnet(vm, input, incInfo.File)
	case ".yml":
	case ".yaml":
		input = "std.native('parseYaml')(importstr '" + rootConfig.BaseDir + "/" + componentPath + "/" + incInfo.File + "')"
		outStr, err = processJsonnet(vm, input, incInfo.File)
	case ".tmpl":
	case ".tpl":
		// Pass component config as data for the template
		outStr, err = processTemplate(rootConfig.BaseDir+"/"+componentPath+"/"+incInfo.File, gjson.Get(config, componentName).Map())
	default:
		outStr, err = "", errors.New("unsupported file extension")
	}
	if err != nil {
		log.Fatal().Str("cluster", kr8Spec.Name).
			Str("component", componentName).
			Str("file", incInfo.File).
			Err(err).
			Msg(outStr)
	}

	// only write file if it does not exist, or the generated contents does not match what is on disk
	var updateNeeded bool
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		log.Debug().Str("cluster", kr8Spec.Name).
			Str("component", componentName).
			Msg("Creating " + outputFile)
		updateNeeded = true
	} else {
		currentContents, err := os.ReadFile(outputFile)
		fatalErrorCheck(err, "Error reading file")
		if string(currentContents) != outStr {
			updateNeeded = true
			log.Debug().Str("cluster", kr8Spec.Name).
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
