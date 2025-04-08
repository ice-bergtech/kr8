package generate

import (
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"

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
