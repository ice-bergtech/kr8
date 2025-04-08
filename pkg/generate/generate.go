package generate

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

func GetClusterParams(clusterDir string, vmConfig types.VMConfig) map[string]string {
	// get list of all clusters, render cluster level params for all of them
	allClusterParams := make(map[string]string)
	allClusters, err := util.GetClusterFilenames(clusterDir)
	util.FatalErrorCheck("Error getting list of clusters", err)
	log.Debug().Msg("Found " + strconv.Itoa(len(allClusters)) + " clusters")

	for _, c := range allClusters {
		allClusterParams[c.Name] = jnetvm.JsonnetRenderClusterParamsOnly(vmConfig, c.Name, "", false)
	}

	return allClusterParams
}

// Only processes specified component if it's defined in the cluster.
// Processes components in string sorted order.
// Sorts out orphaned, generated components directories.
func buildComponentList(
	generatedCompList []string,
	clusterComponents map[string]gjson.Result,
	clusterDir string,
	clusterName string,
	filters util.PathFilterOptions,
) []string {
	var compList []string
	var currentCompRefList []string

	if filters.Components == "" {
		compList = maps.Keys(clusterComponents)
		currentCompRefList = generatedCompList
	} else {
		for _, filterStr := range strings.Split(filters.Components, ",") {
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

func GenProcessComponent(
	vmconfig types.VMConfig,
	componentName string,
	kr8Spec types.Kr8ClusterSpec,
	kr8Opts types.Kr8Opts,
	config string,
	allConfig *safeString,
	filters util.PathFilterOptions,
	paramsFile string,
) {
	log.Info().Str("cluster", kr8Spec.Name).
		Str("component", componentName).
		Msg("Process component")

	// get kr8_spec from component's params
	compSpec, _ := types.CreateComponentSpec(gjson.Get(config, componentName+".kr8_spec"))

	// it's faster to create this VM for each component, rather than re-use
	jvm, err := SetupJvmForComponent(vmconfig, config, kr8Spec, componentName)
	util.FatalErrorCheck("Error setting up JVM for component", err)

	// add kr8_allparams extcode with all component params in the cluster
	if compSpec.Kr8_allparams {
		// include full render of all component params
		// only do this if we have not already cached it and don't already have it stored
		getAllComponentParamsThreadsafe(allConfig, config, vmconfig, kr8Spec, filters, paramsFile, jvm)
	}

	// add kr8_allclusters extcode with every cluster's cluster level params
	if compSpec.Kr8_allclusters {
		getAllClusterParams(kr8Spec.ClusterDir, vmconfig, jvm)
	}

	compPath := filepath.Clean(gjson.Get(config, "_components."+componentName+".path").String())

	// jPathResults always includes base lib. Add jpaths from spec if set
	loadJPathsIntoVM(compSpec, compPath, kr8Spec.ClusterDir, jvm)

	// file imports
	loadExtFilesIntoVars(compSpec, compPath, kr8Spec, kr8Opts, componentName, jvm)

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
		kr8Opts,
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

// This function sets up the JVM for a given component.
// It registers native functions, sets up post-processing, and prunes parameters as required.
// It's faster to create this VM for each component, rather than re-use.
// Default postprocessor just copies input to output.
func SetupJvmForComponent(
	vmconfig types.VMConfig,
	config string,
	kr8Spec types.Kr8ClusterSpec,
	componentName string,
) (*jsonnet.VM, error) {
	jvm, err := jnetvm.JsonnetVM(vmconfig)
	if err != nil {
		return nil, err
	}
	jnetvm.RegisterNativeFuncs(jvm)
	jvm.ExtCode("kr8_cluster", "std.prune("+config+"._cluster)")

	if kr8Spec.PostProcessor != "" {
		jvm.ExtCode("process", kr8Spec.PostProcessor)
	} else {
		// Default PostProcessor passes input to output
		jvm.ExtCode("process", "function(input) input")
	}

	// check if we should prune params
	if kr8Spec.PruneParams {
		jvm.ExtCode("kr8", "std.prune("+config+"."+componentName+")")
	} else {
		jvm.ExtCode("kr8", config+"."+componentName)
	}

	return jvm, nil
}

// combine all the cluster params into a single object indexed by cluster name.
func getAllClusterParams(clusterDir string, vmconfig types.VMConfig, jvm *jsonnet.VM) {
	allClusterParamsObject := "{ "
	for cl, clp := range GetClusterParams(clusterDir, vmconfig) {
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
	filters util.PathFilterOptions,
	paramsFile string,
	jvm *jsonnet.VM,
) {
	allConfig.mu.Lock()
	if allConfig.config == "" {
		if filters.Components == "" {
			allConfig.config = config
		} else {
			allConfig.config = jnetvm.JsonnetRenderClusterParams(
				vmconfig,
				kr8Spec.Name,
				[]string{},
				paramsFile,
				false,
			)
		}
	}
	jvm.ExtCode("kr8_allparams", allConfig.config)
	allConfig.mu.Unlock()
}

// jPathResults always includes base lib.
// Adds jpaths from spec if set.
func loadJPathsIntoVM(compSpec types.Kr8ComponentSpec, compPath string, baseDir string, jvm *jsonnet.VM) {
	jPathResults := []string{filepath.Join(baseDir, "lib")}
	for _, jPath := range compSpec.JPaths {
		jPathResults = append(jPathResults, filepath.Join(baseDir, compPath, jPath))
	}
	jvm.Importer(&jsonnet.FileImporter{
		JPaths: jPathResults,
	})
}

func loadExtFilesIntoVars(
	compSpec types.Kr8ComponentSpec,
	compPath string,
	kr8Spec types.Kr8ClusterSpec,
	kr8Opts types.Kr8Opts,
	componentName string,
	jvm *jsonnet.VM,
) {
	for key, val := range compSpec.ExtFiles {
		log.Debug().Str("cluster", kr8Spec.Name).
			Str("component", componentName).
			Msg("Extfile: " + key + "=" + val)
		filePath := filepath.Join(kr8Opts.BaseDir, compPath, val)
		if kr8Opts.BaseDir != "./" && !strings.HasPrefix(filePath, kr8Opts.BaseDir) {
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
	kr8Opts types.Kr8Opts,
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
		processIncludesFile(
			jvm,
			config,
			kr8Spec,
			kr8Opts,
			componentName,
			compPath,
			componentOutputDir,
			incInfo,
			outputFileMap,
		)
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
	kr8Opts types.Kr8Opts,
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
	inputFile := filepath.Clean(filepath.Join(kr8Opts.BaseDir, componentPath, incInfo.File))

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

func GenProcessCluster(
	clusterName string,
	clusterdir string,
	baseDir string,
	generateDirOverride string,
	kr8Opts types.Kr8Opts,
	clusterParamsFile string,
	filters util.PathFilterOptions,
	vmConfig types.VMConfig,
	pool *ants.Pool,
) {
	log.Debug().Str("cluster", clusterName).Msg("Process cluster")

	// get list of components for cluster
	params := util.GetClusterParamsFilenames(
		clusterdir,
		util.GetClusterPaths(clusterdir, clusterName),
	)

	// get kr8 settings for cluster
	kr8Spec, err := types.CreateClusterSpec(clusterName, gjson.Parse(
		jnetvm.JsonnetRenderFiles(vmConfig, params, "._kr8_spec", false, "", "kr8_spec")),
		kr8Opts,
		generateDirOverride,
	)
	util.FatalErrorCheck("Error creating kr8Spec", err)

	generatedCompList := setupClusterGenerateDirs(kr8Spec)

	clusterComponents := gjson.Parse(
		jnetvm.JsonnetRenderFiles(vmConfig, params, "._components", true, "", "clustercomponents"),
	).Map()

	// determine list of components to process
	compList := buildComponentList(generatedCompList, clusterComponents, kr8Spec.ClusterDir, kr8Spec.Name, filters)

	// render full params for cluster for all selected components
	config := jnetvm.JsonnetRenderClusterParams(
		vmConfig,
		kr8Spec.Name,
		compList,
		clusterParamsFile,
		false,
	)

	var allconfig safeString

	var waitGroup sync.WaitGroup
	for _, componentName := range compList {
		waitGroup.Add(1)
		cName := componentName
		_ = pool.Submit(func() {
			defer waitGroup.Done()
			GenProcessComponent(vmConfig, cName, kr8Spec, kr8Opts, config, &allconfig, filters, clusterParamsFile)
		})
	}
	waitGroup.Wait()
}

func setupClusterGenerateDirs(kr8Spec types.Kr8ClusterSpec) []string {
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

	return generatedCompList
}
