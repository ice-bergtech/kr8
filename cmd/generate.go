package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	goyaml "github.com/ghodss/yaml"
	jsonnet "github.com/google/go-jsonnet"
	"github.com/panjf2000/ants/v2"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"
	"golang.org/x/exp/maps"

	jvm "github.com/ice-bergtech/kr8/pkg/jvm"
	types "github.com/ice-bergtech/kr8/pkg/types"
	util "github.com/ice-bergtech/kr8/pkg/util"
)

// safeString is a thread-safe string that can be used to store and retrieve configuration data
type safeString struct {
	// mu is a mutex that ensures thread-safe access to the struct field
	mu sync.Mutex
	// config is a string that stores the configuration data
	config string
}

var (
	allClusterParams map[string]string
)

// stores the options for the 'generate' command.
type CmdGenerateOptions struct {
	// ClusterParamsFile is a string that stores the path to the cluster params file
	ClusterParamsFile string
	// GenerateDir is a string that stores the output directory for generated files
	GenerateDir string
	// Filters is a PathFilterOptions struct that stores the filters to apply to clusters and components when generating files
	Filters util.PathFilterOptions
}

var cmdGenerateFlags CmdGenerateOptions

func init() {
	RootCmd.AddCommand(GenerateCmd)
	GenerateCmd.Flags().StringVarP(&cmdGenerateFlags.ClusterParamsFile, "clusterparams", "p", "", "provide cluster params as single file - can be combined with --cluster to override cluster")
	GenerateCmd.Flags().StringVarP(&cmdGenerateFlags.Filters.Clusters, "clusters", "C", "", "clusters to generate - comma separated list of cluster names and/or regular expressions ")
	GenerateCmd.Flags().StringVarP(&cmdGenerateFlags.Filters.Components, "components", "c", "", "components to generate - comma separated list of component names and/or regular expressions")
	GenerateCmd.Flags().StringVarP(&cmdGenerateFlags.GenerateDir, "generate-dir", "o", "generated", "output directory")
	GenerateCmd.Flags().StringVarP(&cmdGenerateFlags.Filters.Includes, "clincludes", "i", "", "filter included cluster by including clusters with matching cluster parameters - comma separate list of key/value conditions separated by = or ~ (for regex match)")
	GenerateCmd.Flags().StringVarP(&cmdGenerateFlags.Filters.Excludes, "clexcludes", "x", "", "filter included cluster by excluding clusters with matching cluster parameters - comma separate list of key/value conditions separated by = or ~ (for regex match)")
}

var GenerateCmd = &cobra.Command{
	Use:   "generate [flags]",
	Short: "Generate components",
	Long:  `Generate components in clusters`,

	Args: cobra.MinimumNArgs(0),
	Run:  GenerateCommand,
}

// This function will generate the components for each cluster in parallel
// It uses a wait group to ensure that all clusters have been processed before exiting.
func GenerateCommand(cmd *cobra.Command, args []string) {
	// get list of all clusters, render cluster level params for all of them
	allClusterParams = make(map[string]string)
	allClusters, err := util.GetClusterFilenames(RootConfig.ClusterDir)
	util.FatalErrorCheck(err, "Error getting list of clusters")
	log.Debug().Msg("Found " + strconv.Itoa(len(allClusters)) + " clusters")

	for _, c := range allClusters {
		allClusterParams[c.Name] = jvm.JsonnetRenderClusterParamsOnly(RootConfig.VMConfig, c.Name, "", false)
	}

	var clusterList []string
	// Filter out and cluster or components we don't want to generate
	if cmdGenerateFlags.Filters.Includes != "" || cmdGenerateFlags.Filters.Excludes != "" {
		clusterList = util.CalculateClusterIncludesExcludes(allClusterParams, cmdGenerateFlags.Filters)
		log.Debug().Msg("Have " + strconv.Itoa(len(clusterList)) + " after filtering")
	} else {
		clusterList = maps.Keys(allClusterParams)
	}

	// Setup the threading pools, one for clusters and one for clusters
	var wg sync.WaitGroup
	ants_cp, _ := ants.NewPool(RootConfig.Parallel)
	ants_cl, _ := ants.NewPool(RootConfig.Parallel)

	// Generate config for each cluster in parallel
	for _, clusterName := range clusterList {
		wg.Add(1)
		cl := clusterName
		_ = ants_cl.Submit(func() {
			defer wg.Done()
			genProcessCluster(RootConfig.VMConfig, cl, ants_cp)
		})
	}
	wg.Wait()
}

// Only processes specified component if it's defined in the cluster
// Processes components in string sorted order
// Sorts out orphaned, generated components directories
func buildComponentList(generatedCompList []string, clusterComponents map[string]gjson.Result, clusterDir string, clusterName string) []string {
	var compList []string
	var currentCompRefList []string

	if cmdGenerateFlags.Filters.Components == "" {
		compList = maps.Keys(clusterComponents)
		currentCompRefList = generatedCompList
	} else {
		for _, b := range strings.Split(cmdGenerateFlags.Filters.Components, ",") {
			listFilterComp := util.Filter(generatedCompList, func(s string) bool { r, _ := regexp.MatchString("^"+b+"$", s); return r })
			currentCompRefList = append(currentCompRefList, listFilterComp...)

			listFilterCluster := util.Filter(maps.Keys(clusterComponents), func(s string) bool { r, _ := regexp.MatchString("^"+b+"$", s); return r })
			compList = append(compList, listFilterCluster...)
		}
	}
	sort.Strings(compList)

	// cleanup components that are no longer referenced
	for _, e := range currentCompRefList {
		if _, found := clusterComponents[e]; !found {
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
	util.FatalErrorCheck(json.Unmarshal([]byte(j), &o), "Error unmarshalling jsonnet output to go slice")
	for i, jObj := range o {
		buf, err := goyaml.Marshal(jObj)
		util.FatalErrorCheck(err, "Error marshalling jsonnet object to yaml")
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
	tmpl, err = template.New("file").Funcs(sprig.FuncMap()).Parse(string(tInput))
	if err != nil {
		return "Error parsing template", err
	}
	if err = tmpl.Execute(&buffer, data); err != nil {
		return "Error executing templating", err
	}
	return buffer.String(), nil
}

func genProcessCluster(vmConfig types.VMConfig, clusterName string, p *ants.Pool) {
	log.Debug().Str("cluster", clusterName).Msg("Process cluster")

	// get list of components for cluster
	params := util.GetClusterParamsFilenames(RootConfig.ClusterDir, util.GetClusterPaths(RootConfig.ClusterDir, clusterName))
	clusterComponents := gjson.Parse(jvm.JsonnetRenderFiles(vmConfig, params, "._components", true, "", "clustercomponents")).Map()

	// get kr8 settings for cluster
	kr8Spec, err := types.CreateClusterSpec(clusterName, gjson.Parse(jvm.JsonnetRenderFiles(vmConfig, params, "._kr8_spec", false, "", "kr8_spec")), RootConfig.BaseDir, cmdGenerateFlags.GenerateDir)
	util.FatalErrorCheck(err, "Error creating kr8Spec")

	// create cluster dir
	if _, err := os.Stat(kr8Spec.ClusterDir); os.IsNotExist(err) {
		err = os.MkdirAll(kr8Spec.ClusterDir, os.ModePerm)
		util.FatalErrorCheck(err, "Error creating cluster generateDir")
	}

	// get list of current generated components directories
	d, err := os.Open(kr8Spec.ClusterDir)
	util.FatalErrorCheck(err, "Error opening clusterDir")
	defer d.Close()
	read_all_dirs := -1
	generatedCompList, err := d.Readdirnames(read_all_dirs)
	util.FatalErrorCheck(err, "Error reading directories")

	// determine list of components to process
	compList := buildComponentList(generatedCompList, clusterComponents, kr8Spec.ClusterDir, kr8Spec.Name)

	if len(compList) == 0 { // this needs to be moved so purging above works first
		return
	}

	// render full params for cluster for all selected components
	config := jvm.JsonnetRenderClusterParams(vmConfig, kr8Spec.Name, compList, cmdGenerateFlags.ClusterParamsFile, false)

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

func genProcessComponent(vmconfig types.VMConfig, componentName string, kr8Spec types.Kr8ClusterSpec, config string, allConfig *safeString) {
	log.Info().Str("cluster", kr8Spec.Name).
		Str("component", componentName).
		Msg("Process component")

	// get kr8_spec from component's params
	//spec := gjson.Get(config, componentName+".kr8_spec").Map()
	compPath := gjson.Get(config, "_components."+componentName+".path").String()
	compSpec, _ := types.CreateComponentSpec(gjson.Get(config, componentName+".kr8_spec"))

	// it's faster to create this VM for each component, rather than re-use
	vm, _ := jvm.JsonnetVM(vmconfig)
	jvm.RegisterNativeFuncs(vm)
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
			if cmdGenerateFlags.Filters.Components == "" {
				// all component params are in config
				allConfig.config = config
			} else {
				allConfig.config = jvm.JsonnetRenderClusterParams(vmconfig, kr8Spec.Name, []string{}, cmdGenerateFlags.ClusterParamsFile, false)
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
	jPath := []string{RootConfig.BaseDir + "/lib"}
	for _, j := range compSpec.JPaths {
		jPath = append(jPath, RootConfig.BaseDir+"/"+compPath+"/"+j)
	}
	vm.Importer(&jsonnet.FileImporter{
		JPaths: jPath,
	})

	// file imports
	for k, v := range compSpec.ExtFiles {
		vPath := RootConfig.BaseDir + "/" + compPath + "/" + v // use full path for file
		extFile, err := os.ReadFile(vPath)
		util.FatalErrorCheck(err, "Error importing extfiles item")
		log.Debug().Str("cluster", kr8Spec.Name).
			Str("component", componentName).
			Msg("Extfile: " + k + "=" + v)
		vm.ExtVar(k, string(extFile))
	}

	componentOutputDir := kr8Spec.ClusterDir + "/" + componentName
	// create component dir if needed
	if _, err := os.Stat(componentOutputDir); os.IsNotExist(err) {
		err := os.MkdirAll(componentOutputDir, os.ModePerm)
		util.FatalErrorCheck(err, "Error creating component directory")
	}

	incInfo := types.Kr8ComponentSpecIncludeObject{
		DestExt:  "yaml",
		DestDir:  "",
		DestName: "",
	}

	outputFileMap := make(map[string]bool)
	// generate each included file
	for _, include := range compSpec.Includes {
		switch include.(type) {
		case string:
			incInfo = types.Kr8ComponentSpecIncludeObject{
				File:     include.(string),
				DestExt:  "yaml", // default to yaml ,
				DestDir:  "",
				DestName: "",
			}
		case types.Kr8ComponentSpecIncludeObject:
			// include is a map with multiple fields
			incInfo = include.(types.Kr8ComponentSpecIncludeObject)
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
		util.FatalErrorCheck(err, "")
		defer d.Close()
		names, err := d.Readdirnames(-1)
		util.FatalErrorCheck(err, "")
		for _, name := range names {
			if _, ok := outputFileMap[name]; ok {
				// file is managed
				continue
			}
			if filepath.Ext(name) == ".yaml" {
				delFile := filepath.Join(componentOutputDir, name)
				err = os.RemoveAll(delFile)
				util.FatalErrorCheck(err, "")
				log.Debug().Str("cluster", kr8Spec.Name).
					Str("component", componentName).
					Msg("Deleted: " + delFile)
			}
		}
		d.Close()
	}
}

func processIncludesFile(vm *jsonnet.VM, config string, kr8Spec types.Kr8ClusterSpec, componentName string, componentPath string, componentOutputDir string, incInfo types.Kr8ComponentSpecIncludeObject, outputFileMap map[string]bool) {
	file_extension := filepath.Ext(incInfo.File)

	// ensure this directory exists
	outputDir := componentOutputDir
	if incInfo.DestDir != "" {
		outputDir = kr8Spec.ClusterDir + "/" + incInfo.DestDir
	}
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		err = os.MkdirAll(outputDir, os.ModePerm)
		util.FatalErrorCheck(err, "Error creating alternate directory")
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
		input = "( import '" + RootConfig.BaseDir + "/" + componentPath + "/" + incInfo.File + "')"
		outStr, err = processJsonnet(vm, input, incInfo.File)
	case ".yml":
	case ".yaml":
		input = "std.native('parseYaml')(importstr '" + RootConfig.BaseDir + "/" + componentPath + "/" + incInfo.File + "')"
		outStr, err = processJsonnet(vm, input, incInfo.File)
	case ".tmpl":
	case ".tpl":
		// Pass component config as data for the template
		outStr, err = processTemplate(RootConfig.BaseDir+"/"+componentPath+"/"+incInfo.File, gjson.Get(config, componentName).Map())
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
		util.FatalErrorCheck(err, "Error reading file")
		if string(currentContents) != outStr {
			updateNeeded = true
			log.Debug().Str("cluster", kr8Spec.Name).
				Str("component", componentName).
				Msg("Updating: " + outputFile)
		}
	}
	if updateNeeded {
		f, err := os.Create(outputFile)
		util.FatalErrorCheck(err, "Error creating file")
		//defer f.Close()
		_, err = f.WriteString(outStr)
		f.Close()
		util.FatalErrorCheck(err, "Error writing to file")
	}
}
