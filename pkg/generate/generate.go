// Package generate implements the logic for generating output files based on input data.
//
// Combines a directory of cluster configurations
// with a directory of components
// (along with some Jsonnet libs)
// to generate output files.
//
// The package prepares a Jsonnet VM and loads the necessary libraries and extvars.
// A new VM is created for each component.
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
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/tidwall/gjson"
	"golang.org/x/exp/maps"

	jnetvm "github.com/ice-bergtech/kr8/pkg/jnetvm"
	"github.com/ice-bergtech/kr8/pkg/kr8_types"
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

func GetClusterParams(clusterDir string, vmConfig types.VMConfig, logger zerolog.Logger) (map[string]string, error) {
	// get list of all clusters, render cluster level params for all of them
	allClusterParams := make(map[string]string)
	allClusters, err := util.GetClusterFilenames(clusterDir)
	if err := util.LogErrorIfCheck("Error getting list of clusters", err, logger); err != nil {
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

// Root function for processing a kr8 component.
// Processes a component through a jsonnet VM to generate output files.
func GenProcessComponent(
	vmConfig types.VMConfig,
	componentName string,
	kr8Spec kr8_types.Kr8ClusterSpec,
	kr8Opts types.Kr8Opts,
	config string,
	allConfig *safeString,
	filters util.PathFilterOptions,
	paramsFile string,
	logger zerolog.Logger,
) error {
	logger.Info().Msg("Processing component")
	// get kr8_spec from component's params
	compSpec, err := kr8_types.CreateComponentSpec(gjson.Get(config, componentName+".kr8_spec"), logger)
	if err := util.LogErrorIfCheck("Error creating component spec", err, logger); err != nil {
		return err
	}

	// it's faster to create this VM for each component, rather than re-use
	jvm, compPath, err := SetupAndConfigureVM(
		vmConfig,
		config,
		kr8Spec,
		componentName,
		compSpec,
		allConfig,
		filters,
		paramsFile,
		kr8Opts,
		logger,
	)
	if err := util.LogErrorIfCheck("Error setting up JVM for component", err, logger); err != nil {
		return err
	}

	componentOutputDir := filepath.Join(kr8Spec.GenerateDir, kr8Spec.Name, componentName)
	// create component dir if needed
	if _, err := os.Stat(componentOutputDir); os.IsNotExist(err) {
		err := os.MkdirAll(componentOutputDir, 0750)
		if err := util.LogErrorIfCheck("Error creating component directory", err, logger); err != nil {
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
		logger,
	)
	if err := util.LogErrorIfCheck("Error generating includes files", err, logger); err != nil {
		return err
	}

	// purge any yaml files in the output dir that were not generated
	if !compSpec.DisableOutputDirClean {
		return CleanOutputDir(outputFileMap, componentOutputDir)
	}

	return nil
}

// Setup and configures a jsonnet VM for processing kr8 resources.
// Creates a new VM and does the following:
//   - loads cluster and component config
//   - loads jsonnet library files
//   - loads external file references
func SetupAndConfigureVM(
	vmConfig types.VMConfig,
	config string,
	kr8Spec kr8_types.Kr8ClusterSpec,
	componentName string,
	compSpec kr8_types.Kr8ComponentSpec,
	allConfig *safeString,
	filters util.PathFilterOptions,
	paramsFile string,
	kr8Opts types.Kr8Opts,
	logger zerolog.Logger,
) (*jsonnet.VM, string, error) {
	jvm, err := SetupJvmForComponent(vmConfig, config, kr8Spec, componentName)
	if err := util.LogErrorIfCheck("error setting up JVM for component", err, logger); err != nil {
		return nil, "", err
	}
	// include full render of all component params
	if compSpec.Kr8_allParams {
		// only do this if we have not already cached it and don't already have it stored
		if err := getAllComponentParamsThreadsafe(
			allConfig,
			config,
			vmConfig,
			kr8Spec,
			filters,
			paramsFile,
			jvm,
			logger,
		); err != nil {
			return nil, "", util.LogErrorIfCheck("error getting all component params", err, logger)
		}
	}
	if compSpec.Kr8_allClusters {
		// add kr8_allclusters extcode with every cluster's cluster level params
		if err := getAllClusterParams(kr8Opts.ClusterDir, vmConfig, jvm, logger); err != nil {
			return nil, "", util.LogErrorIfCheck("error getting all cluster params", err, logger)
		}
	}

	compPath := filepath.Clean(gjson.Get(config, "_components."+componentName+".path").String())
	// jPathResults always includes base lib. Add jpaths from spec if set
	loadJPathsIntoVM(compSpec, compPath, kr8Opts.BaseDir, jvm, logger)
	// file imports
	if err := loadExtFilesIntoVars(compSpec, compPath, kr8Spec, kr8Opts, componentName, jvm, logger); err != nil {
		return nil, "", util.LogErrorIfCheck("error loading ext files into vars", err, logger)
	}

	return jvm, compPath, err
}

// Combine all the cluster params into a single object indexed by cluster name.
func getAllClusterParams(clusterDir string, vmConfig types.VMConfig, jvm *jsonnet.VM, logger zerolog.Logger) error {
	allClusterParamsObject := "{ "
	params, err := GetClusterParams(clusterDir, vmConfig, logger)
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
	vmConfig types.VMConfig,
	kr8Spec kr8_types.Kr8ClusterSpec,
	filters util.PathFilterOptions,
	paramsFile string,
	jvm *jsonnet.VM,
	logger zerolog.Logger,
) error {
	allConfig.mu.Lock()
	if allConfig.config == "" {
		if filters.Components == "" {
			allConfig.config = config
		} else {
			var err error
			allConfig.config, err = jnetvm.JsonnetRenderClusterParams(
				vmConfig,
				kr8Spec.Name,
				[]string{},
				paramsFile,
				false,
			)
			if err != nil {
				allConfig.mu.Unlock()

				return util.LogErrorIfCheck("Error rendering cluster params", err, logger)
			}
		}
	}
	jvm.ExtCode("kr8_allparams", allConfig.config)
	allConfig.mu.Unlock()

	return nil
}

// Generates the list of includes files for a component.
// Processes each includes file using the component's config.
// Returns an error if there's an issue with ANY includes file.
func GenerateIncludesFiles(
	includesFiles []kr8_types.Kr8ComponentSpecIncludeObject,
	kr8Spec kr8_types.Kr8ClusterSpec,
	kr8Opts types.Kr8Opts,
	config string,
	componentName string,
	compPath string,
	componentOutputDir string,
	jvm *jsonnet.VM,
	logger zerolog.Logger,
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
			logger.With().Str("includes files", include.File).Logger(),
		)
		if err != nil {
			return nil, util.LogErrorIfCheck("error processing includes file", err, logger)
		}
	}

	return outputFileMap, nil
}

// The root function for generating a cluster.
// Prepares and builds the cluster config.
// Build and processes the list of components.
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
	logger zerolog.Logger,
) error {
	logger.Debug().Str("cluster", clusterName).Msg("Processing cluster")
	clusterPaths, err := util.GetClusterPaths(clusterdir, clusterName)
	if err != nil {
		return err
	}
	// get list of components for cluster
	params := util.GetClusterParamsFilenames(clusterdir, clusterPaths)
	renderedKr8Spec, err := jnetvm.JsonnetRenderFiles(vmConfig, params, "._kr8_spec", false, "", "kr8_spec")
	if err := util.LogErrorIfCheck("error rendering cluster `_kr8_spec`", err, logger); err != nil {
		return err
	}

	// get kr8 settings for cluster
	kr8Spec, err := kr8_types.CreateClusterSpec(clusterName, gjson.Parse(renderedKr8Spec),
		kr8Opts, generateDirOverride,
		logger,
	)
	if err := util.LogErrorIfCheck("error creating kr8Spec", err, logger); err != nil {
		return err
	}

	generatedCompList, err := setupClusterGenerateDirs(kr8Spec)
	if err := util.LogErrorIfCheck("error creating generate dirs", err, logger); err != nil {
		return err
	}

	renderedCompSpec, err := jnetvm.JsonnetRenderFiles(vmConfig, params, "._components", true, "", "clustercomponents")
	if err := util.LogErrorIfCheck("error rendering cluster components list", err, logger); err != nil {
		return err
	}

	clusterComponents := gjson.Parse(renderedCompSpec).Map()

	// determine list of components to process
	compList := buildComponentList(generatedCompList, clusterComponents, kr8Spec.ClusterOutputDir, kr8Spec.Name, filters)

	config, err := jnetvm.JsonnetRenderClusterParams(
		vmConfig,
		kr8Spec.Name,
		compList,
		clusterParamsFile,
		false,
	)
	if err := util.LogErrorIfCheck("error rendering cluster params", err, logger); err != nil {
		return err
	}

	// render full params for cluster for all selected components
	return renderComponents(config, vmConfig, kr8Spec, compList, clusterParamsFile, pool, kr8Opts, filters, logger)
}

// Renders a list of components with a given Kr8ClusterSpec configuration.
// Each component is processed by a process thread from a thread pool.
func renderComponents(
	config string,
	vmConfig types.VMConfig,
	kr8Spec kr8_types.Kr8ClusterSpec,
	compList []string,
	clusterParamsFile string,
	pool *ants.Pool,
	kr8Opts types.Kr8Opts,
	filters util.PathFilterOptions,
	logger zerolog.Logger,
) error {
	var allConfig safeString

	var waitGroup sync.WaitGroup
	for _, componentName := range compList {
		waitGroup.Add(1)
		cName := componentName
		_ = pool.Submit(func() {
			defer waitGroup.Done()
			sublogger := logger.With().Str("component", componentName).Logger()
			if err := GenProcessComponent(
				vmConfig,
				cName,
				kr8Spec,
				kr8Opts,
				config,
				&allConfig,
				filters,
				clusterParamsFile,
				sublogger,
			); err != nil {
				sublogger.Error().
					Err(err).
					Msg("Failed to process component")
			}
		})
	}
	waitGroup.Wait()

	return nil
}
