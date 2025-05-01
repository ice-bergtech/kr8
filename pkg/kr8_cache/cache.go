// Package kr8_cache defines the structure for kr8+ cluster-component resource caching.
// Cache is based on cluster-level config, component config, and component file reference hashes.
package kr8_cache

import (
	"encoding/base64"
	"encoding/json"
	"path/filepath"
	"reflect"

	"github.com/ice-bergtech/kr8/pkg/types"
	"github.com/ice-bergtech/kr8/pkg/util"
	"github.com/rs/zerolog"
	"github.com/tidwall/gjson"
)

// Load cluster cache from a specified cache file.
// Assumes cache is gzipped, but falls back to plaintext if there's an error.
func LoadClusterCache(cacheFile string) (*DeploymentCache, error) {
	text, err := util.ReadGzip(cacheFile)
	if err != nil {
		text, err = util.ReadFile(cacheFile)
		if err != nil {
			return nil, err
		}
	}

	//nolint:exhaustruct
	result := DeploymentCache{}
	err = json.Unmarshal(text, &result)
	if err != nil {
		return nil, err
	}

	if result.ClusterConfig == nil || result.ComponentConfigs == nil {
		return nil, types.Kr8Error{Message: "cache missing fields", Value: result}
	}

	return &result, nil
}

// Object that contains the cache for a single cluster.
type DeploymentCache struct {
	// A struct containing cluster-level cache values
	ClusterConfig *ClusterCache `json:"cluster_config"`
	// Map of cache entries for cluster components.
	// Depends on ClusterConfig cache being valid to be considered valid.
	ComponentConfigs map[string]ComponentCache `json:"component_config"`
	LibraryCache     *LibraryCache             `json:"library_cache"`
}

func InitDeploymentCache(config string, baseDir string, cacheResults map[string]ComponentCache) *DeploymentCache {
	cache := DeploymentCache{
		ClusterConfig:    CreateClusterCache(config),
		ComponentConfigs: cacheResults,
		LibraryCache:     CreateLibraryCache(baseDir),
	}

	return &cache
}

func (cache *DeploymentCache) WriteCache(outFile string, compress bool) error {
	// confirm cluster-level configuration matches the cache
	var text []byte
	text, err := json.Marshal(cache)
	if err != nil {
		return err
	}

	if compress {
		return util.WriteGzip(text, outFile)
	}

	return util.WriteFile(text, outFile)
}

func (cache *DeploymentCache) CheckClusterCache(config string, logger zerolog.Logger) bool {
	// confirm cluster-level configuration matches the cache
	if cache.ClusterConfig != nil {
		return cache.ClusterConfig.CheckClusterCache(config, logger)
	}

	return false
}

func (cache *DeploymentCache) CheckClusterComponentCache(
	config string,
	componentName string,
	componentPath string,
	files []string,
	logger zerolog.Logger,
) (bool, *ComponentCache, error) {
	currentState, err := CreateComponentCache(config, files)
	if err != nil {
		return false, currentState, err
	}

	// first confirm cluster-level configuration matches the cache
	if !cache.CheckClusterCache(config, logger) {
		return false, currentState, nil
	}

	componentCache, ok := cache.ComponentConfigs[componentName]
	if !ok {
		return false, currentState, nil
	}

	result := currentState.ComponentConfig == componentCache.ComponentConfig &&
		reflect.DeepEqual(currentState.ComponentFiles, componentCache.ComponentFiles)

	return result, currentState, nil
}

type LibraryCache struct {
	Directory string            `json:"directory"`
	Entries   map[string]string `json:"entries"`
}

// This is cluster-level cache that applies to all components.
// If it is deemed invalid, the component cache is also invalid.
type ClusterCache struct {
	// Raw cluster _kr8_spec object
	Kr8_Spec string `json:"kr8_spec"`
	// Raw cluster _cluster object
	Cluster string `json:"cluster"`
}

// Stores the cluster kr8_spec and cluster config as cluster-level cache.
func CreateClusterCache(config string) *ClusterCache {
	return &ClusterCache{
		Kr8_Spec: base64.RawStdEncoding.EncodeToString([]byte(gjson.Get(config, "_kr8_spec").Raw)),
		Cluster:  base64.RawStdEncoding.EncodeToString([]byte(gjson.Get(config, "_cluster").Raw)),
	}
}

func CreateLibraryCache(baseDir string) *LibraryCache {
	result := LibraryCache{
		Directory: baseDir,
		Entries:   map[string]string{},
	}

	files, err := util.BuildDirFileList(filepath.Join(baseDir, "lib"))
	if err != nil {
		return &result
	}
	result.Entries = make(map[string]string, len(files))
	for _, file := range files {
		hash, err := util.HashFile(file)
		if err != nil {
			result.Entries[file] = "error: " + err.Error()
		} else {
			result.Entries[file] = hash
		}
	}

	return &result
}

// Compares current cluster config represented as a json string to the cache.
// Returns true if cache is valid.
func (cache *ClusterCache) CheckClusterCache(config string, logger zerolog.Logger) bool {
	currentState := CreateClusterCache(config)
	// compare cluster (non-component) configuration to cached cluster
	if cache.Kr8_Spec != currentState.Kr8_Spec {
		logger.Debug().Msg("_kr8_spec differs from cache")

		return false
	}

	if cache.Cluster != currentState.Cluster {
		logger.Debug().Msg("_cluster differs from cache")

		return false
	}

	return true
}

type ComponentCache struct {
	// Raw component config string
	ComponentConfig string `json:"component_config"`
	// Map of filenames to file hashes
	ComponentFiles map[string]string `json:"component_files"`
}

func CreateComponentCache(config string, listFiles []string) (*ComponentCache, error) {
	cacheResult := ComponentCache{
		ComponentConfig: base64.RawStdEncoding.EncodeToString([]byte(config)),
		ComponentFiles:  map[string]string{},
	}
	for _, file := range listFiles {
		currentHash, err := util.HashFile(file)
		if err != nil {
			return nil, err
		}
		cacheResult.ComponentFiles[file] = currentHash
	}

	return &cacheResult, nil
}

func (cache *ComponentCache) CheckComponentCache(
	config string,
	componentName string,
	componentPath string,
	files []string,
	logger zerolog.Logger,
) (bool, *ComponentCache) {
	currentState, err := CreateComponentCache(config, files)
	if err != nil {
		return false, nil
	}
	// compare cluster-level component config
	if cache.ComponentConfig != currentState.ComponentConfig {
		logger.Info().Msg("component config differs from cache")

		return false, currentState
	}
	if len(cache.ComponentFiles) != len(currentState.ComponentFiles) {
		return false, currentState
	}
	for _, file := range currentState.ComponentFiles {
		hash, ok := cache.ComponentFiles[file]
		// didn't find file in cache
		if !ok {
			return false, currentState
		}
		currentHash, err := util.HashFile(filepath.Join(componentPath, file))
		if err != nil {
			logger.Warn().Err(err).Msg("issue hashing file, cache invalid")

			return false, currentState
		}
		if hash != currentHash {
			return false, currentState
		}
	}

	return true, currentState
}
