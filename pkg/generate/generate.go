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

	jnetvm "github.com/ice-bergtech/kr8p/pkg/jnetvm"
	types "github.com/ice-bergtech/kr8p/pkg/types"
	util "github.com/ice-bergtech/kr8p/pkg/util"
)

// A thread-safe string that can be used to store and retrieve configuration data.
type safeString struct {
	// mu is a mutex that ensures thread-safe access to the struct field
	mu sync.Mutex
	// config is a string that stores the configuration data
	config string
}

func GetClusterParams(clusterDir string, vmConfig types.VMConfig) (map[string]string, error) {
	// get list of all clusters, render cluster level params for all of them
	allClusterParams := make(map[string]string)
	allClusters, err := util.GetClusterFilenames(clusterDir)
	if err := util.GenErrorIfCheck("Error getting list of clusters", err); err != nil {
		return nil, err
	}
	log.Debug().Msg("Found " + strconv.Itoa(len(allClusters)) + " clusters")

	for _, c := range allClusters {
		allClusterParams[c.Name], err = jnetvm.JsonnetRenderClusterParamsOnly(vmConfig, c.Name, "", false)
		if err != nil {
			return nil, err
		}
	}

	return allClusterParams, nil
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
) error {
	log.Info().Str("cluster", kr8Spec.Name).
		Str("component", componentName).
		Msg("Process component")

	// get kr8_spec from component's params
	compSpec, err := types.CreateComponentSpec(gjson.Get(config, componentName+".kr8_spec"))
	if err := util.GenErrorIfCheck("Error creating component spec", err); err != nil {
		return err
	}

	// it's faster to create this VM for each component, rather than re-use
	jvm, compPath, err := SetupAndConfigureVM(
		vmconfig,
		config,
		kr8Spec,
		componentName,
		compSpec,
		allConfig,
		filters,
		paramsFile,
		kr8Opts,
	)
	if err := util.GenErrorIfCheck("Error setting up JVM for component", err); err != nil {
		return err
	}

	componentOutputDir := filepath.Join(kr8Spec.ClusterDir, componentName)
	// create component dir if needed
	if _, err := os.Stat(componentOutputDir); os.IsNotExist(err) {
		err := os.MkdirAll(componentOutputDir, 0750)
		if err := util.GenErrorIfCheck("Error creating component directory", err); err != nil {
			return err
		}
	}

	// generate each included file
	outputFileMap, err := GenerateIncludesFiles(
		compSpec.Includes,
		kr8Spec,
		kr8Opts,
		config,
		componentName,
		compPath,
		componentOutputDir,
		jvm,
	)
	if err := util.GenErrorIfCheck("Error generating includes files", err); err != nil {
		return err
	}

	// purge any yaml files in the output dir that were not generated
	if !compSpec.DisableOutputDirClean {
		return CleanOutputDir(outputFileMap, componentOutputDir)
	}

	return nil
}

func SetupAndConfigureVM(
	vmconfig types.VMConfig,
	config string,
	kr8Spec types.Kr8ClusterSpec,
	componentName string,
	compSpec types.Kr8ComponentSpec,
	allConfig *safeString,
	filters util.PathFilterOptions,
	paramsFile string,
	kr8Opts types.Kr8Opts,
) (*jsonnet.VM, string, error) {
	jvm, err := SetupJvmForComponent(vmconfig, config, kr8Spec, componentName)
	if err := util.GenErrorIfCheck("error setting up JVM for component", err); err != nil {
		return nil, "", err
	}
	// include full render of all component params
	if compSpec.Kr8_allparams {
		// only do this if we have not already cached it and don't already have it stored
		if err := getAllComponentParamsThreadsafe(
			allConfig,
			config,
			vmconfig,
			kr8Spec,
			filters,
			paramsFile,
			jvm,
		); err != nil {
			return nil, "", util.GenErrorIfCheck("error getting all component params", err)
		}
	}
	if compSpec.Kr8_allclusters {
		// add kr8_allclusters extcode with every cluster's cluster level params
		if err := getAllClusterParams(kr8Spec.ClusterDir, vmconfig, jvm); err != nil {
			return nil, "", util.GenErrorIfCheck("error getting all cluster params", err)
		}
	}

	compPath := filepath.Clean(gjson.Get(config, "_components."+componentName+".path").String())
	// jPathResults always includes base lib. Add jpaths from spec if set
	loadJPathsIntoVM(compSpec, compPath, kr8Spec.ClusterDir, jvm)
	// file imports
	if err := loadExtFilesIntoVars(compSpec, compPath, kr8Spec, kr8Opts, componentName, jvm); err != nil {
		return nil, "", util.GenErrorIfCheck("error loading ext files into vars", err)
	}

	return jvm, compPath, err
}

// combine all the cluster params into a single object indexed by cluster name.
func getAllClusterParams(clusterDir string, vmconfig types.VMConfig, jvm *jsonnet.VM) error {
	allClusterParamsObject := "{ "
	params, err := GetClusterParams(clusterDir, vmconfig)
	if err != nil {
		return err
	}
	for cl, clp := range params {
		allClusterParamsObject = allClusterParamsObject + "'" + cl + "': " + clp + ","
	}
	allClusterParamsObject += "}"
	jvm.ExtCode("kr8_allclusters", allClusterParamsObject)

	return nil
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
) error {
	allConfig.mu.Lock()
	if allConfig.config == "" {
		if filters.Components == "" {
			allConfig.config = config
		} else {
			var err error
			allConfig.config, err = jnetvm.JsonnetRenderClusterParams(
				vmconfig,
				kr8Spec.Name,
				[]string{},
				paramsFile,
				false,
			)
			if err != nil {
				allConfig.mu.Unlock()

				return util.GenErrorIfCheck("Error rendering cluster params", err)
			}
		}
	}
	jvm.ExtCode("kr8_allparams", allConfig.config)
	allConfig.mu.Unlock()

	return nil
}

func GenerateIncludesFiles(
	includesFiles []types.Kr8ComponentSpecIncludeObject,
	kr8Spec types.Kr8ClusterSpec,
	kr8Opts types.Kr8Opts,
	config string,
	componentName string,
	compPath string,
	componentOutputDir string,
	jvm *jsonnet.VM,
) (map[string]bool, error) {
	outputFileMap := make(map[string]bool)
	for _, include := range includesFiles {
		if include.DestName == "" {
			if kr8Spec.GenerateShortNames {
				sBase := filepath.Base(include.File)
				include.DestName = sBase[0 : len(sBase)-len(filepath.Ext(include.File))]
			} else {
				// replaces slashes with _ in multi-dir paths and replace extension with yaml
				include.DestName = strings.ReplaceAll(
					include.File[0:len(include.File)-len(filepath.Ext(include.File))],
					"/", "_",
				)
			}
		}
		err := processIncludesFile(
			jvm,
			config,
			kr8Spec,
			kr8Opts,
			componentName,
			compPath,
			componentOutputDir,
			include,
			outputFileMap,
		)
		if err != nil {
			return nil, util.GenErrorIfCheck("error processing includes file", err)
		}
	}

	return outputFileMap, nil
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
) error {
	log.Debug().Str("cluster", clusterName).Msg("Processing cluster")
	clusterPaths, err := util.GetClusterPaths(clusterdir, clusterName)
	if err != nil {
		return err
	}
	// get list of components for cluster
	params := util.GetClusterParamsFilenames(clusterdir, clusterPaths)
	renderedKr8Spec, err := jnetvm.JsonnetRenderFiles(vmConfig, params, "._kr8_spec", false, "", "kr8_spec")
	if err := util.GenErrorIfCheck("error rendering cluster `_kr8_spec`", err); err != nil {
		return err
	}

	// get kr8p settings for cluster
	kr8Spec, err := types.CreateClusterSpec(clusterName, gjson.Parse(renderedKr8Spec),
		kr8Opts, generateDirOverride,
	)
	if err := util.GenErrorIfCheck("error creating kr8Spec", err); err != nil {
		return err
	}

	generatedCompList, err := setupClusterGenerateDirs(kr8Spec)
	if err := util.GenErrorIfCheck("error creating generate dirs", err); err != nil {
		return err
	}

	renderedCompSpec, err := jnetvm.JsonnetRenderFiles(vmConfig, params, "._components", true, "", "clustercomponents")
	if err := util.GenErrorIfCheck("error rendering cluster components list", err); err != nil {
		return err
	}

	clusterComponents := gjson.Parse(renderedCompSpec).Map()

	// determine list of components to process
	compList := buildComponentList(generatedCompList, clusterComponents, kr8Spec.ClusterDir, kr8Spec.Name, filters)

	config, err := jnetvm.JsonnetRenderClusterParams(
		vmConfig,
		kr8Spec.Name,
		compList,
		clusterParamsFile,
		false,
	)
	if err := util.GenErrorIfCheck("error rendering cluster params", err); err != nil {
		return err
	}

	// render full params for cluster for all selected components
	return renderComponents(config, vmConfig, kr8Spec, compList, clusterParamsFile, pool, kr8Opts, filters)
}

func renderComponents(
	config string,
	vmConfig types.VMConfig,
	kr8Spec types.Kr8ClusterSpec,
	compList []string,
	clusterParamsFile string,
	pool *ants.Pool,
	kr8Opts types.Kr8Opts,
	filters util.PathFilterOptions,
) error {
	var allconfig safeString

	var waitGroup sync.WaitGroup
	for _, componentName := range compList {
		waitGroup.Add(1)
		cName := componentName
		_ = pool.Submit(func() {
			defer waitGroup.Done()
			if err := GenProcessComponent(
				vmConfig,
				cName,
				kr8Spec,
				kr8Opts,
				config,
				&allconfig,
				filters,
				clusterParamsFile,
			); err != nil {
				log.Error().
					Str("cluster", componentName).
					Str("component", componentName).
					Err(err).
					Msg("Failed to process component")
			}
		})
	}
	waitGroup.Wait()

	return nil
}
