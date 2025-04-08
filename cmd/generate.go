package cmd

import (
	"bytes"
	"encoding/json"
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

	jnetvm "github.com/ice-bergtech/kr8/pkg/jnetvm"
	types "github.com/ice-bergtech/kr8/pkg/types"
	util "github.com/ice-bergtech/kr8/pkg/util"
)

// A thread-safe string that can be used to store and retrieve configuration data.
type safeString struct {
	// mu is a mutex that ensures thread-safe access to the struct field
	mu sync.Mutex
	// config is a string that stores the configuration data
	config string
}

// Stores the options for the 'generate' command.
type CmdGenerateOptions struct {
	// Stores the path to the cluster params file
	ClusterParamsFile string
	// Stores the output directory for generated files
	GenerateDir string
	// Stores the filters to apply to clusters and components when generating files
	Filters util.PathFilterOptions
}

var cmdGenerateFlags CmdGenerateOptions

func init() {
	RootCmd.AddCommand(GenerateCmd)
	GenerateCmd.Flags().StringVarP(&cmdGenerateFlags.ClusterParamsFile,
		"clusterparams", "p", "",
		"provide cluster params as single file - can be combined with --cluster to override cluster")
	GenerateCmd.Flags().StringVarP(&cmdGenerateFlags.Filters.Clusters,
		"clusters", "C", "",
		"clusters to generate - comma separated list of cluster names and/or regular expressions ")
	GenerateCmd.Flags().StringVarP(&cmdGenerateFlags.Filters.Components, "components", "c", "",
		"components to generate - comma separated list of component names and/or regular expressions")
	GenerateCmd.Flags().StringVarP(&cmdGenerateFlags.GenerateDir,
		"generate-dir", "o", "generated",
		"output directory")
	GenerateCmd.Flags().StringVarP(&cmdGenerateFlags.Filters.Includes,
		"clincludes", "i", "",
		"filter included cluster by including clusters with matching cluster parameters - "+
			"comma separate list of key/value conditions separated by = or ~ (for regex match)")
	GenerateCmd.Flags().StringVarP(&cmdGenerateFlags.Filters.Excludes,
		"clexcludes", "x", "",
		"filter included cluster by excluding clusters with matching cluster parameters - "+
			"comma separate list of key/value conditions separated by = or ~ (for regex match)")
}

var GenerateCmd = &cobra.Command{
	Use:   "generate [flags]",
	Short: "Generate components",
	Long:  `Generate components in clusters`,

	Args: cobra.MinimumNArgs(0),
	Run:  GenerateCommand,
}

func GetClusterParams() map[string]string {
	// get list of all clusters, render cluster level params for all of them
	allClusterParams := make(map[string]string)
	allClusters, err := util.GetClusterFilenames(RootConfig.ClusterDir)
	util.FatalErrorCheck("Error getting list of clusters", err)
	log.Debug().Msg("Found " + strconv.Itoa(len(allClusters)) + " clusters")

	for _, c := range allClusters {
		allClusterParams[c.Name] = jnetvm.JsonnetRenderClusterParamsOnly(RootConfig.VMConfig, c.Name, "", false)
	}

	return allClusterParams
}

// This function will generate the components for each cluster in parallel.
// It uses a wait group to ensure that all clusters have been processed before exiting.
func GenerateCommand(cmd *cobra.Command, args []string) {
	// get list of all clusters, render cluster level params for all of them
	allClusterParams := GetClusterParams()

	var clusterList []string
	// Filter out and cluster or components we don't want to generate
	if cmdGenerateFlags.Filters.Includes != "" || cmdGenerateFlags.Filters.Excludes != "" {
		clusterList = util.CalculateClusterIncludesExcludes(allClusterParams, cmdGenerateFlags.Filters)
		log.Debug().Msg("Have " + strconv.Itoa(len(clusterList)) + " after filtering")
	} else {
		clusterList = maps.Keys(allClusterParams)
	}

	// Setup the threading pools, one for clusters and one for clusters
	var waitGroup sync.WaitGroup
	ants_cp, _ := ants.NewPool(RootConfig.Parallel)
	ants_cl, _ := ants.NewPool(RootConfig.Parallel)

	// Generate config for each cluster in parallel
	for _, clusterName := range clusterList {
		waitGroup.Add(1)
		cl := clusterName
		_ = ants_cl.Submit(func() {
			defer waitGroup.Done()
			genProcessCluster(RootConfig.VMConfig, cl, ants_cp)
		})
	}
	waitGroup.Wait()
}

// Only processes specified component if it's defined in the cluster.
// Processes components in string sorted order.
// Sorts out orphaned, generated components directories.
func buildComponentList(
	generatedCompList []string,
	clusterComponents map[string]gjson.Result,
	clusterDir string,
	clusterName string,
) []string {
	var compList []string
	var currentCompRefList []string

	if cmdGenerateFlags.Filters.Components == "" {
		compList = maps.Keys(clusterComponents)
		currentCompRefList = generatedCompList
	} else {
		for _, filterStr := range strings.Split(cmdGenerateFlags.Filters.Components, ",") {
			listFilterComp := util.Filter(generatedCompList, func(s string) bool {
				r, _ := regexp.MatchString("^"+filterStr+"$", s)

				return r
			})
			currentCompRefList = append(currentCompRefList, listFilterComp...)

			listFilterCluster := util.Filter(maps.Keys(clusterComponents), func(s string) bool {
				r, _ := regexp.MatchString("^"+filterStr+"$", s)

				return r
			})
			compList = append(compList, listFilterCluster...)
		}
	}
	sort.Strings(compList)

	// cleanup components that are no longer referenced
	for _, component := range currentCompRefList {
		if _, found := clusterComponents[component]; !found {
			delComp := filepath.Join(clusterDir, component)
			if err := os.RemoveAll(delComp); err != nil {
				log.Error().Msg("Issue deleting generated for component " + component)
			}
			log.Info().Str("cluster", clusterName).
				Str("component", component).
				Msg("Deleting generated for component")
		}
	}

	return compList
}

func processJsonnet(jvm *jsonnet.VM, input string, snippetFilename string) (string, error) {
	jvm.ExtCode("input", input)
	jsonStr, err := jvm.EvaluateAnonymousSnippet(snippetFilename, "std.extVar('process')(std.extVar('input'))")
	if err != nil {
		return "Error evaluating jsonnet snippet", err
	}

	// create output file contents in a string first, as a yaml stream
	var o []interface{}
	var outStr string
	util.FatalErrorCheck("Error unmarshalling jsonnet output to go slice", json.Unmarshal([]byte(jsonStr), &o))
	for _, jObj := range o {
		buf, err := goyaml.Marshal(jObj)
		util.FatalErrorCheck("Error marshalling jsonnet object to yaml", err)
		outStr += string(buf)
		// Place yml new document at end of each object
		outStr += "\n---\n"
	}

	return outStr, nil
}

func processTemplate(filename string, data map[string]gjson.Result) (string, error) {
	var tInput []byte
	var tmpl *template.Template
	var buffer bytes.Buffer
	var err error

	tInput, err = os.ReadFile(filepath.Clean(filename))
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

func genProcessCluster(vmConfig types.VMConfig, clusterName string, pool *ants.Pool) {
	log.Debug().Str("cluster", clusterName).Msg("Process cluster")

	// get list of components for cluster
	params := util.GetClusterParamsFilenames(
		RootConfig.ClusterDir,
		util.GetClusterPaths(RootConfig.ClusterDir, clusterName),
	)
	clusterComponents := gjson.Parse(
		jnetvm.JsonnetRenderFiles(vmConfig, params, "._components", true, "", "clustercomponents"),
	).Map()

	// get kr8 settings for cluster
	kr8Spec, err := types.CreateClusterSpec(clusterName, gjson.Parse(
		jnetvm.JsonnetRenderFiles(vmConfig, params, "._kr8_spec", false, "", "kr8_spec")),
		RootConfig.BaseDir,
		cmdGenerateFlags.GenerateDir,
	)
	util.FatalErrorCheck("Error creating kr8Spec", err)

	// create cluster dir
	if _, err := os.Stat(kr8Spec.ClusterDir); os.IsNotExist(err) {
		err = os.MkdirAll(kr8Spec.ClusterDir, 0750)
		util.FatalErrorCheck("Error creating cluster generateDir", err)
	}

	// get list of current generated components directories
	d, err := os.Open(kr8Spec.ClusterDir)
	util.FatalErrorCheck("Error opening clusterDir", err)
	defer d.Close()
	read_all_dirs := -1
	generatedCompList, err := d.Readdirnames(read_all_dirs)
	util.FatalErrorCheck("Error reading directories", err)

	// determine list of components to process
	compList := buildComponentList(generatedCompList, clusterComponents, kr8Spec.ClusterDir, kr8Spec.Name)

	// this needs to be moved so purging above works first
	if len(compList) == 0 {
		return
	}

	// render full params for cluster for all selected components
	config := jnetvm.JsonnetRenderClusterParams(
		vmConfig,
		kr8Spec.Name,
		compList,
		cmdGenerateFlags.ClusterParamsFile,
		false,
	)

	var allconfig safeString

	var waitGroup sync.WaitGroup
	for _, componentName := range compList {
		waitGroup.Add(1)
		cName := componentName
		_ = pool.Submit(func() {
			defer waitGroup.Done()
			genProcessComponent(vmConfig, cName, kr8Spec, config, &allconfig)
		})
	}
	waitGroup.Wait()
}

func genProcessComponent(
	vmconfig types.VMConfig,
	componentName string,
	kr8Spec types.Kr8ClusterSpec,
	config string,
	allConfig *safeString,
) {
	log.Info().Str("cluster", kr8Spec.Name).
		Str("component", componentName).
		Msg("Process component")

	// get kr8_spec from component's params
	compSpec, _ := types.CreateComponentSpec(gjson.Get(config, componentName+".kr8_spec"))

	// it's faster to create this VM for each component, rather than re-use
	jvm, _ := jnetvm.JsonnetVM(vmconfig)
	jnetvm.RegisterNativeFuncs(jvm)
	jvm.ExtCode("kr8_cluster", "std.prune("+config+"._cluster)")

	// vm.ExtCode("kr8_components", "std.prune("+config+"._components)")
	if kr8Spec.PostProcessor != "" {
		jvm.ExtCode("process", kr8Spec.PostProcessor)
	} else {
		// default postprocessor just copies input
		jvm.ExtCode("process", "function(input) input")
	}

	// prune params if required
	if kr8Spec.PruneParams {
		jvm.ExtCode("kr8", "std.prune("+config+"."+componentName+")")
	} else {
		jvm.ExtCode("kr8", config+"."+componentName)
	}

	// add kr8_allparams extcode with all component params in the cluster
	if compSpec.Kr8_allparams {
		// include full render of all component params
		// only do this if we have not already cached it and don't already have it stored
		getAllComponentParamsThreadsafe(allConfig, config, vmconfig, kr8Spec, jvm)
	}

	// add kr8_allclusters extcode with every cluster's cluster level params
	if compSpec.Kr8_allclusters {
		getAllClusterParams(jvm)
	}

	compPath := filepath.Clean(gjson.Get(config, "_components."+componentName+".path").String())

	// jPathResults always includes base lib. Add jpaths from spec if set
	loadJPathsIntoVM(compSpec, compPath, jvm)

	// file imports
	loadExtFilesIntoVars(compSpec, compPath, kr8Spec, componentName, jvm)

	componentOutputDir := filepath.Join(kr8Spec.ClusterDir, componentName)
	// create component dir if needed
	if _, err := os.Stat(componentOutputDir); os.IsNotExist(err) {
		err := os.MkdirAll(componentOutputDir, 0750)
		util.FatalErrorCheck("Error creating component directory", err)
	}

	// generate each included file
	outputFileMap := GenerateIncludesFiles(
		compSpec.Includes,
		kr8Spec,
		config,
		componentName,
		compPath,
		componentOutputDir,
		jvm,
	)

	// purge any yaml files in the output dir that were not generated
	if !compSpec.DisableOutputDirClean {
		CleanOutputDir(outputFileMap, componentOutputDir)
	}
}

// combine all the cluster params into a single object indexed by cluster name.
func getAllClusterParams(jvm *jsonnet.VM) {
	allClusterParamsObject := "{ "
	for cl, clp := range GetClusterParams() {
		allClusterParamsObject = allClusterParamsObject + "'" + cl + "': " + clp + ","
	}
	allClusterParamsObject += "}"
	jvm.ExtCode("kr8_allclusters", allClusterParamsObject)
}

// Include full render of all component params for cluster.
// Only do this if we have not already cached it and don't already have it stored.
func getAllComponentParamsThreadsafe(
	allConfig *safeString,
	config string,
	vmconfig types.VMConfig,
	kr8Spec types.Kr8ClusterSpec,
	jvm *jsonnet.VM,
) {
	allConfig.mu.Lock()
	if allConfig.config == "" {
		if cmdGenerateFlags.Filters.Components == "" {
			allConfig.config = config
		} else {
			allConfig.config = jnetvm.JsonnetRenderClusterParams(
				vmconfig,
				kr8Spec.Name,
				[]string{},
				cmdGenerateFlags.ClusterParamsFile,
				false,
			)
		}
	}
	jvm.ExtCode("kr8_allparams", allConfig.config)
	allConfig.mu.Unlock()
}

// jPathResults always includes base lib.
// Adds jpaths from spec if set.
func loadJPathsIntoVM(compSpec types.Kr8ComponentSpec, compPath string, jvm *jsonnet.VM) {
	jPathResults := []string{filepath.Join(RootConfig.BaseDir, "lib")}
	for _, jPath := range compSpec.JPaths {
		jPathResults = append(jPathResults, filepath.Join(RootConfig.BaseDir, compPath, jPath))
	}
	jvm.Importer(&jsonnet.FileImporter{
		JPaths: jPathResults,
	})
}

func loadExtFilesIntoVars(
	compSpec types.Kr8ComponentSpec,
	compPath string,
	kr8Spec types.Kr8ClusterSpec,
	componentName string,
	jvm *jsonnet.VM,
) {
	for key, val := range compSpec.ExtFiles {
		log.Debug().Str("cluster", kr8Spec.Name).
			Str("component", componentName).
			Msg("Extfile: " + key + "=" + val)
		filePath := filepath.Join(RootConfig.BaseDir, compPath, val)
		if RootConfig.BaseDir != "./" && !strings.HasPrefix(filePath, RootConfig.BaseDir) {
			util.FatalErrorCheck("Invalid file path: "+filePath, os.ErrNotExist)
		}
		extFile, err := os.ReadFile(filepath.Clean(filePath))
		util.FatalErrorCheck("Error importing extfiles item", err)
		jvm.ExtVar(key, string(extFile))
	}
}

func GenerateIncludesFiles(
	includesFiles []interface{},
	kr8Spec types.Kr8ClusterSpec,
	config string,
	componentName string,
	compPath string,
	componentOutputDir string,
	jvm *jsonnet.VM,
) map[string]bool {
	outputFileMap := make(map[string]bool)
	for _, include := range includesFiles {
		incInfo := types.Kr8ComponentSpecIncludeObject{
			DestExt:  "yaml",
			DestDir:  "",
			DestName: "",
			File:     "",
		}
		switch val := include.(type) {
		case string:
			fileName := val
			incInfo = types.Kr8ComponentSpecIncludeObject{
				File:     fileName,
				DestExt:  "yaml", // default to yaml ,
				DestDir:  "",
				DestName: "",
			}
		case types.Kr8ComponentSpecIncludeObject:
			// include is a map with multiple fields
			incInfo = val
		default:
			log.Fatal().Msg("Invalid include type")
		}
		if incInfo.DestName == "" {
			if kr8Spec.GenerateShortNames {
				sBase := filepath.Base(incInfo.File)
				incInfo.DestName = sBase[0 : len(sBase)-len(filepath.Ext(include.(string)))]
			} else {
				// replaces slashes with _ in multi-dir paths and replace extension with yaml
				incInfo.DestName = strings.ReplaceAll(
					incInfo.File[0:len(incInfo.File)-len(filepath.Ext(include.(string)))],
					"/", "_",
				)
			}
		}
		processIncludesFile(jvm, config, kr8Spec, componentName, compPath, componentOutputDir, incInfo, outputFileMap)
	}

	return outputFileMap
}

func CleanOutputDir(outputFileMap map[string]bool, componentOutputDir string) {
	// clean component dir
	d, err := os.Open(filepath.Clean(componentOutputDir))
	util.FatalErrorCheck("", err)
	// Lifetime of function
	defer d.Close()
	names, err := d.Readdirnames(-1)
	util.FatalErrorCheck("", err)
	for _, name := range names {
		if _, ok := outputFileMap[name]; ok {
			// file is managed
			continue
		}
		if filepath.Ext(name) == ".yaml" {
			delFile := filepath.Join(componentOutputDir, name)
			err = os.RemoveAll(delFile)
			util.FatalErrorCheck("", err)
			log.Debug().Msg("Deleted: " + delFile)
		}
	}
}

func processIncludesFile(
	jvm *jsonnet.VM,
	config string,
	kr8Spec types.Kr8ClusterSpec,
	componentName string,
	componentPath string,
	componentOutputDir string,
	incInfo types.Kr8ComponentSpecIncludeObject,
	outputFileMap map[string]bool,
) {
	// ensure this directory exists
	outputDir := componentOutputDir
	if incInfo.DestDir != "" {
		outputDir = filepath.Join(kr8Spec.ClusterDir, incInfo.DestDir)
	}
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		err = os.MkdirAll(outputDir, 0750)
		util.FatalErrorCheck("Error creating alternate directory", err)
	}
	outputFile := filepath.Clean(filepath.Join(outputDir, incInfo.DestName+"."+incInfo.DestExt))
	inputFile := filepath.Clean(filepath.Join(RootConfig.BaseDir, componentPath, incInfo.File))

	// remember output filename for purging files
	outputFileMap[incInfo.DestName+"."+incInfo.DestExt] = true

	outStr := ProcessFile(inputFile, outputFile, kr8Spec, componentName, config, incInfo, jvm)

	log.Debug().Str("cluster", kr8Spec.Name).Str("component", componentName).Msg("Checking if file needs updating...")

	// only write file if it does not exist, or the generated contents does not match what is on disk
	if CheckIfUpdateNeeded(outputFile, outStr) {
		f, err := os.Create(outputFile)
		util.FatalErrorCheck("Error creating file", err)
		_, err = f.WriteString(outStr)
		util.FatalErrorCheck("Error writing to file", err)
		util.FatalErrorCheck("Error closing file", f.Close())
	}
}

// Process an includes file.
// Based on the extension, it will process it differently.
//
// .jsonnet: Imported and processed using jsonnet VM.
//
// .yml, .yaml: Imported and processed through native function ParseYaml.
//
// .tpl, .tmpl: Processed using component config and Sprig templating.
func ProcessFile(
	inputFile string,
	outputFile string,
	kr8Spec types.Kr8ClusterSpec,
	componentName string,
	config string,
	incInfo types.Kr8ComponentSpecIncludeObject,
	jvm *jsonnet.VM,
) string {
	log.Debug().Str("cluster", kr8Spec.Name).
		Str("component", componentName).
		Msg("Process file: " + inputFile + " -> " + outputFile)

	file_extension := filepath.Ext(incInfo.File)

	var input string
	var outStr string
	var err error
	switch file_extension {
	case ".jsonnet":
		// file is processed as an ExtCode input, so that we can postprocess it
		// in the snippet
		input = "( import '" + inputFile + "')"
		outStr, err = processJsonnet(jvm, input, incInfo.File)
	case ".yml":
	case ".yaml":
		input = "std.native('parseYaml')(importstr '" + inputFile + "')"
		outStr, err = processJsonnet(jvm, input, incInfo.File)
	case ".tmpl":
	case ".tpl":
		// Pass component config as data for the template
		outStr, err = processTemplate(inputFile, gjson.Get(config, componentName).Map())
	default:
		outStr, err = "", os.ErrInvalid
	}
	if err != nil {
		log.Fatal().Str("cluster", kr8Spec.Name).
			Str("component", componentName).
			Str("file", incInfo.File).
			Err(err).
			Msg(outStr)
	}

	return outStr
}

// Check if a file needs updating based on its current contents and the new contents.
func CheckIfUpdateNeeded(outFile string, outStr string) bool {
	var updateNeeded bool
	outFile = filepath.Clean(outFile)
	if _, err := os.Stat(outFile); os.IsNotExist(err) {
		log.Debug().Msg("Creating " + outFile)
		updateNeeded = true
	} else {
		currentContents, err := os.ReadFile(outFile)
		util.FatalErrorCheck("Error reading file", err)
		if string(currentContents) != outStr {
			updateNeeded = true
			log.Debug().Msg("Updating: " + outFile)
		}
	}

	return updateNeeded
}
