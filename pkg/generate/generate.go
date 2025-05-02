// Package generate implements the logic for generating output files based on input data.
//
// Combines a directory of cluster configurations
// with a directory of components
// (along with some Jsonnet libs)
// to generate output files.
//
// The package prepares a Jsonnet VM and loads the necessary libraries and extVars.
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
	"github.com/tidwall/gjson"

	//nolint:exptostd
	"golang.org/x/exp/maps"

	jnetvm "github.com/ice-bergtech/kr8/pkg/jnetvm"
	"github.com/ice-bergtech/kr8/pkg/kr8_cache"
	"github.com/ice-bergtech/kr8/pkg/kr8_types"
	types "github.com/ice-bergtech/kr8/pkg/types"
	util "github.com/ice-bergtech/kr8/pkg/util"
)

// A thread-safe string that can be used to store and retrieve configuration data.
type SafeString struct {
	// mu is a mutex that ensures thread-safe access to the struct field
	mu sync.Mutex
	// config is a string that stores the configuration data
	config string
}

// Given a base directory, generates cluster-level configuration for each cluster found.
// Gets list of clusters from `util.GetClusterFilenames(clusterDir)`.
func GetClusterParams(clusterDir string, vmConfig types.VMConfig, logger zerolog.Logger) (map[string]string, error) {
	// get list of all clusters, render cluster level params for all of them
	allClusterParams := make(map[string]string)
	allClusters, err := util.GetClusterFilenames(clusterDir)
	if err := util.LogErrorIfCheck("Error getting list of clusters", err, logger); err != nil {
		return nil, err
	}

	logger.Debug().Msg("GetClusterParams Found " + strconv.Itoa(len(allClusters)) + " clusters")

	for _, c := range allClusters {
		allClusterParams[c.Name], err = jnetvm.JsonnetRenderClusterParamsOnly(vmConfig, c.Name, "", false)
		if err != nil {
			return nil, err
		}
	}

	return allClusterParams, nil
}

// Calculates which components should be generated based on filters.
// Only processes specified component if it's defined in the cluster.
// Processes components in string sorted order.
// Sorts out orphaned, generated components directories.
//
//nolint:exptostd
func CalculateClusterComponentList(
	clusterComponents map[string]gjson.Result,
	filters util.PathFilterOptions,
	existingClusterComponents []string,
) []string {
	var compList []string

	if filters.Components == "" {
		compList = maps.Keys(clusterComponents)
	} else {
		for _, filterStr := range strings.Split(filters.Components, ",") {
			compList = append(
				compList,
				util.Filter(maps.Keys(clusterComponents), func(s string) bool {
					r, _ := regexp.MatchString("^"+filterStr+"$", s)

					return r
				})...,
			)
		}
	}
	sort.Strings(compList)

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
	allConfig *SafeString,
	filters util.PathFilterOptions,
	paramsFile string,
	cache *kr8_cache.DeploymentCache,
	logger zerolog.Logger,
) (bool, *kr8_cache.ComponentCache, error) {
	logger.Info().Msg("Processing component")

	// get kr8_spec from component's params
	compSpec, err := kr8_types.CreateComponentSpec(gjson.Get(config, componentName+".kr8_spec"), logger)
	if err := util.LogErrorIfCheck("Error creating component spec", err, logger); err != nil {
		return false, nil, err
	}
	cacheValid, currentCacheState, err := CheckComponentCache(
		cache, compSpec, config,
		componentName, kr8Opts.BaseDir, logger,
	)
	if kr8Spec.EnableCache && !compSpec.DisableCache {
		if err != nil {
			logger.Error().Err(err).Msg("issue checking/creating component cache")
		}
		if cacheValid {
			logger.Info().Msg("+ Component matches cache, skipping")

			return true, currentCacheState, nil
		}
		logger.Info().Msg("- Component differs from cache, continuing")
	}

	// it's faster to create this VM for each component, rather than re-use.
	// TODO: It would be faster to create this outside this function at the cluster level.
	// Then deep copy/clone VM for each component with pre-allocated memory.
	// Slowed down by numerous memory allocs.
	jvm, compPath, err := SetupComponentVM(
		vmConfig, config, kr8Spec, componentName, compSpec,
		allConfig, filters, paramsFile, kr8Opts, logger,
	)
	if err := util.LogErrorIfCheck("Error setting up JVM for component", err, logger); err != nil {
		return false, nil, err
	}

	componentOutputDir := filepath.Join(kr8Spec.GenerateDir, kr8Spec.Name, componentName)
	// create component dir if needed
	if _, err := os.Stat(componentOutputDir); os.IsNotExist(err) {
		err := os.MkdirAll(componentOutputDir, 0750)
		if err := util.LogErrorIfCheck("Error creating component directory", err, logger); err != nil {
			return false, nil, err
		}
	}

	// generate each included file
	outputFileMap, err := GenerateIncludesFiles(
		compSpec.Includes, kr8Spec, kr8Opts, config,
		componentName, compPath, componentOutputDir, jvm, logger,
	)
	if err := util.LogErrorIfCheck("Error generating includes files", err, logger); err != nil {
		return false, nil, err
	}

	return true, currentCacheState, ProcessComponentFinalizer(compSpec, componentOutputDir, outputFileMap)
}

// Final actions performed once a component is generated.
// Cleans extra files from output dir if not disabled in component spec.
func ProcessComponentFinalizer(
	compSpec kr8_types.Kr8ComponentSpec,
	componentOutputDir string,
	outputFileMap map[string]bool,
) error {
	// purge any yaml files in the output dir that were not generated
	if !compSpec.DisableOutputDirClean {
		err := CleanOutputDir(outputFileMap, componentOutputDir)
		if err != nil {
			return err
		}
	}

	return nil
}

// Compares a component's current state to a cache entry.
// Returns an up-to-date cache entry for the component.
// If the cache pointer is nil or cache invalid, a fresh cache entry will be generated to return.
func CheckComponentCache(
	cache *kr8_cache.DeploymentCache,
	compSpec kr8_types.Kr8ComponentSpec,
	config string,
	componentName string,
	baseDir string,
	logger zerolog.Logger,
) (bool, *kr8_cache.ComponentCache, error) {
	compPath := GetComponentPath(config, componentName)
	// build list of files referenced by component
	listFiles, err := util.BuildDirFileList(filepath.Join(baseDir, compPath))
	if err != nil {
		logger.Warn().Err(err).Msg("CheckComponentCache issue walking component directory")

		return false, nil, err
	}
	// check if the component matches the cache
	if cache != nil {
		return cache.CheckClusterComponentCache(
			config,
			componentName,
			compPath,
			listFiles,
			logger,
		)
	}

	newCache, err := kr8_cache.CreateComponentCache(config, listFiles)
	if err != nil {
		return false, nil, err
	}

	return false, newCache, nil
}

func GetComponentFiles(compSpec kr8_types.Kr8ComponentSpec) []string {
	numIncludes := len(compSpec.Includes)
	numExtFiles := len(compSpec.ExtFiles)
	numJPaths := len(compSpec.JPaths)
	listFiles := make([]string, numIncludes+numExtFiles+numJPaths)

	for i, obj := range compSpec.Includes {
		listFiles[i] = obj.File
	}

	idx := 0
	for _, path := range compSpec.ExtFiles {
		listFiles[idx+numIncludes] = path
		idx++
	}

	for i, path := range compSpec.JPaths {
		listFiles[i+numIncludes+numExtFiles] = path
	}

	return listFiles
}

// Setup and configures a jsonnet VM for processing kr8 resources.
// Creates a new VM and does the following:
//   - loads cluster and component config
//   - loads jsonnet library files
//   - loads external file references
func SetupComponentVM(
	vmConfig types.VMConfig,
	config string,
	kr8Spec kr8_types.Kr8ClusterSpec,
	componentName string,
	compSpec kr8_types.Kr8ComponentSpec,
	allConfig *SafeString,
	filters util.PathFilterOptions,
	paramsFile string,
	kr8Opts types.Kr8Opts,
	logger zerolog.Logger,
) (*jsonnet.VM, string, error) {
	// Initialize a default jsonnet VM for components to build on top of
	jvm, err := SetupBaseComponentJvm(vmConfig, config, kr8Spec)
	if err != nil {
		kErr := types.Kr8Error{Message: "error initializing component jsonnet VM", Value: err}

		return nil, "", kErr
	}
	// Add component-specific config to the JVM
	SetupJvmForComponent(jvm, config, kr8Spec, componentName)
	// Check if a full render of all cluster component params should be included
	if compSpec.Kr8_allParams {
		// only do this if we have not already cached it and don't already have it stored
		if err := GetClusterComponentParamsThreadSafe(
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

	// check if a full render of ALL cluster params should be included
	if compSpec.Kr8_allClusters {
		// add kr8_allclusters extCode with every cluster's cluster level params
		if err := GetAllClusterParams(kr8Opts.ClusterDir, vmConfig, jvm, logger); err != nil {
			return nil, "", util.LogErrorIfCheck("error getting all cluster params", err, logger)
		}
	}

	// Load files referenced by the component
	compPath := GetComponentPath(config, componentName)
	// jPathResults always includes base lib. Add jPaths from spec if set
	loadLibPathsIntoVM(compSpec, compPath, kr8Opts.BaseDir, jvm, logger)
	// file imports
	if err := loadExtFilesIntoVM(compSpec, compPath, kr8Opts, jvm, logger); err != nil {
		return nil, "", util.LogErrorIfCheck("error loading ext files into vars", err, logger)
	}

	return jvm, compPath, nil
}

// Fetch a component path from raw cluster config.
func GetComponentPath(config string, componentName string) string {
	return filepath.Clean(gjson.Get(config, "_components."+componentName+".path").String())
}

// Combine all the cluster params into a single object indexed by cluster name.
func GetAllClusterParams(clusterDir string, vmConfig types.VMConfig, jvm *jsonnet.VM, logger zerolog.Logger) error {
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
func GetClusterComponentParamsThreadSafe(
	allConfig *SafeString,
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
			jvm, config,
			kr8Spec, kr8Opts,
			componentName, compPath,
			componentOutputDir, include,
			outputFileMap, logger.With().Str("includes_file", include.File).Logger(),
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
	clusterDir string,
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

	// Start by compiling the cluster-level configuration
	kr8Spec, compList, config, err := GatherClusterConfig(
		clusterName,
		clusterDir,
		kr8Opts,
		vmConfig,
		generateDirOverride,
		filters,
		clusterParamsFile,
		logger,
	)
	if err != nil {
		return err
	}

	var cacheCur *kr8_cache.DeploymentCache
	cacheFile := ""
	if kr8Spec.EnableCache {
		cacheCur, cacheFile = LoadClusterCache(kr8Spec, logger)
	} else {
		cacheCur = nil
	}

	// render full params for cluster for all selected components
	componentCacheResult, err := RenderComponents(
		config,
		vmConfig,
		*kr8Spec,
		cacheCur,
		compList,
		clusterParamsFile,
		pool,
		kr8Opts,
		filters,
		logger,
	)
	if err != nil {
		return err
	}

	// If caching is enabled, generate the cache file for the cluster.
	if kr8Spec.EnableCache {
		return kr8_cache.InitDeploymentCache(
			config,
			baseDir,
			componentCacheResult,
		).WriteCache(cacheFile, kr8Spec.CompressCache)
	}

	return nil
}

// Compiles configuration for each cluster.
// Creates and cleans output directories for generated cluster components.
// Uses the filter to determine which components to process.
// Renders the cluster-level configuration for each component.
func GatherClusterConfig(
	clusterName, clusterDir string,
	kr8Opts types.Kr8Opts,
	vmConfig types.VMConfig,
	generateDirOverride string,
	filters util.PathFilterOptions,
	clusterParamsFile string,
	logger zerolog.Logger,
) (*kr8_types.Kr8ClusterSpec, []string, string, error) {
	kr8Spec, clusterComponents, err := CompileClusterConfiguration(
		clusterName,
		clusterDir,
		kr8Opts,
		vmConfig,
		generateDirOverride,
		logger,
	)
	if err != nil {
		return nil, nil, "", err
	}

	// Setup output dirs and remove component output dirs that are no longer referenced
	existingComponents, err := CreateClusterGenerateDirs(*kr8Spec)
	if err := util.LogErrorIfCheck("error creating generate dirs", err, logger); err != nil {
		return nil, nil, "", err
	}

	CleanupOldComponentDirs(existingComponents, clusterComponents, kr8Spec, logger)

	// Determine list of components to process
	compList := CalculateClusterComponentList(clusterComponents, filters, existingComponents)

	// Use Jsonnet to render cluster-level configurations for components
	config, err := jnetvm.JsonnetRenderClusterParams(
		vmConfig,
		kr8Spec.Name,
		compList,
		clusterParamsFile,
		false,
	)
	if err := util.LogErrorIfCheck("error rendering cluster params", err, logger); err != nil {
		return nil, nil, "", err
	}

	return kr8Spec, compList, config, nil
}

// Loads cluster cache based on a cluster spec.
// If cache is disabled, a nil deployment cache pointer is returned.
func LoadClusterCache(
	kr8Spec *kr8_types.Kr8ClusterSpec,
	logger zerolog.Logger,
) (*kr8_cache.DeploymentCache, string) {
	var cache *kr8_cache.DeploymentCache
	var err error
	cacheFile := filepath.Join(kr8Spec.ClusterOutputDir, ".kr8_cache")
	if kr8Spec.EnableCache {
		cache, err = kr8_cache.LoadClusterCache(cacheFile)
		if err != nil {
			logger.Warn().Err(err).Msg("error loading cluster cache")
		} else {
			logger.Info().Str("file", cacheFile).Msg("Loaded cluster cache")
		}
	} else {
		cache = nil
	}

	return cache, cacheFile
}

// Go through each item in existingComponents and remove the file if it isn't in clusterComponents.
// Skips removing `.kr8_cache` files.
func CleanupOldComponentDirs(
	existingComponents []string,
	clusterComponents map[string]gjson.Result,
	kr8Spec *kr8_types.Kr8ClusterSpec,
	logger zerolog.Logger,
) {
	for _, component := range existingComponents {
		if _, found := clusterComponents[component]; !found {
			// Skip deleting cache files
			if component == ".kr8_cache" {
				continue
			}
			delComp := filepath.Join(kr8Spec.ClusterOutputDir, component)
			if err := os.RemoveAll(delComp); err != nil {
				logger.Error().Msg("Issue deleting generated for component " + component)
			}
			logger.Info().Str("component", component).
				Msg("Deleting generated dir for non-referenced component")
		}
	}
}

// Build the list of parameter files to combine for the final cluster config by walking folder tree leaf to root.
func CompileClusterConfiguration(
	clusterName, clusterDir string,
	kr8Opts types.Kr8Opts,
	vmConfig types.VMConfig,
	generateDirOverride string,
	logger zerolog.Logger,
) (*kr8_types.Kr8ClusterSpec, map[string]gjson.Result, error) {
	// First determine the path to the cluster.jsonnet file.
	clusterPath, err := util.GetClusterPath(clusterDir, clusterName)
	if err != nil {
		return nil, nil, err
	}
	// Gather list of configurations that apply to the cluster
	params := util.GetClusterParamsFilenames(clusterDir, clusterPath)

	// Compile the cluster kr8 configuration
	renderedKr8Spec, err := jnetvm.JsonnetRenderFiles(vmConfig, params, "._kr8_spec", false, "", "kr8_spec")
	if err := util.LogErrorIfCheck("error rendering cluster `_kr8_spec`", err, logger); err != nil {
		return nil, nil, err
	}
	// Package the cluster kr8 spec into struct
	kr8Spec, err := kr8_types.CreateClusterSpec(
		clusterName,
		gjson.Parse(renderedKr8Spec),
		kr8Opts,
		generateDirOverride,
		logger,
	)
	if err := util.LogErrorIfCheck("error creating kr8Spec", err, logger); err != nil {
		return nil, nil, err
	}

	// Compile the cluster component references
	renderedCompSpec, err := jnetvm.JsonnetRenderFiles(
		vmConfig, params,
		"._components", true,
		"", clusterName+": ._components",
	)
	if err := util.LogErrorIfCheck("error rendering cluster components list", err, logger); err != nil {
		return nil, nil, err
	}
	// Package into a map
	clusterComponents := gjson.Parse(renderedCompSpec).Map()

	return &kr8Spec, clusterComponents, nil
}

// Renders a list of components with a given Kr8ClusterSpec configuration.
// Each component is added to a sync.WaitGroup to be processed by the ants.Pool.
// Returns the cache results for all successfully generated components.
func RenderComponents(
	config string,
	vmConfig types.VMConfig,
	kr8Spec kr8_types.Kr8ClusterSpec,
	cache *kr8_cache.DeploymentCache,
	compList []string,
	clusterParamsFile string,
	pool *ants.Pool,
	kr8Opts types.Kr8Opts,
	filters util.PathFilterOptions,
	logger zerolog.Logger,
) (map[string]kr8_cache.ComponentCache, error) {
	// Get a cache object for components to reference.
	cacheObj := ValidateOrCreateCache(cache, config, logger)
	// Create a channel for components to place their final cache entries into.
	cacheResultChannel := make(chan map[string]kr8_cache.ComponentCache, len(compList))

	var allConfig SafeString
	var waitGroup sync.WaitGroup

	for _, componentName := range compList {
		waitGroup.Add(1)
		cName := componentName
		_ = pool.Submit(func() {
			defer waitGroup.Done()
			// Create a new logger for the component to use
			subLogger := logger.With().Str("component", componentName).Logger()
			success, cacheResult, err := GenProcessComponent(
				vmConfig, cName,
				kr8Spec, kr8Opts,
				config, &allConfig,
				filters, clusterParamsFile,
				cacheObj, subLogger,
			)
			if err != nil {
				subLogger.Error().
					Err(err).
					Msg("Failed to process component")
			}
			// Record cache results if component generate was successful.
			if success && cacheResult != nil {
				cacheResultChannel <- map[string]kr8_cache.ComponentCache{
					componentName: *cacheResult,
				}
			}
		})
	}
	waitGroup.Wait()
	close(cacheResultChannel)

	// Pre-allocate our cache object and fill it from the channel.
	result := make(map[string]kr8_cache.ComponentCache, len(cacheResultChannel))
	for s := range cacheResultChannel {
		// exp/maps used elsewhere in pkg.
		//nolint:exptostd
		maps.Copy(result, s)
	}

	return result, nil
}

// For provided config, validates the cache object matches.
// If the cache is valid, it is returned.
// If cache is not valid, an empty deployment cache returned.
func ValidateOrCreateCache(
	cache *kr8_cache.DeploymentCache,
	config string,
	logger zerolog.Logger,
) *kr8_cache.DeploymentCache {
	cacheObj := cache
	if cacheObj == nil || !cacheObj.CheckClusterCache(
		config,
		logger,
	) {
		cacheObj = &kr8_cache.DeploymentCache{
			ClusterConfig:    nil,
			ComponentConfigs: map[string]kr8_cache.ComponentCache{},
			LibraryCache: &kr8_cache.LibraryCache{
				Directory: "",
				Entries:   map[string]string{},
			},
		}
	}

	return cacheObj
}
