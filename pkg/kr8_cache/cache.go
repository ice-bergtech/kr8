// Package kr8_cache defines the structure for kr8+ cluster-component resource caching.
// Cache is based on cluster-level config, component config, and component file reference hashes.
package kr8_cache

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/ice-bergtech/kr8/pkg/util"
	"github.com/rs/zerolog"
	"github.com/tidwall/gjson"
)

// Load cluster cache from a specified cache file.
func LoadClusterCache(cacheFile string) (*DeploymentCache, error) {
	fCache, err := os.Open(filepath.Clean(cacheFile))
	if err != nil {
		return nil, err
	}
	defer fCache.Close()

	text := []byte{}
	_, err = fCache.Read(text)
	if err != nil {
		return nil, err
	}

	result := DeploymentCache{}
	err = json.Unmarshal(text, &result)
	if err != nil {
		return nil, err
	}

	return &DeploymentCache{
		ClusterConfig:    &ClusterCache{},
		ComponentConfigs: make(map[string]ComponentCache),
	}, nil
}

// Object that contains the cache for a single cluster.
type DeploymentCache struct {
	ClusterConfig *ClusterCache `json:"cluster_config"`
	// Map of cache entries for cluster components.
	// Depends on ClusterConfig cache being valid to be considered valid.
	ComponentConfigs map[string]ComponentCache `json:"component_config"`
}

func (cache *DeploymentCache) WriteCache(outFile string) error {
	// confirm cluster-level configuration matches the cache
	text, err := json.Marshal(cache)
	if err != nil {
		return err
	}
	f, err := os.Create(filepath.Clean(outFile))
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(text)

	return err
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
) bool {
	// first confirm cluster-level configuration matches the cache
	result := cache.CheckClusterCache(config, logger)
	if !result {
		return result
	}

	componentCache, ok := cache.ComponentConfigs[componentName]
	if !ok {
		return false
	}

	return result && componentCache.CheckComponentCache(config, componentName, componentPath, files, logger)
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
		Kr8_Spec: gjson.Get(config, "._kr8_spec").Raw,
		Cluster:  gjson.Get(config, "._cluster").Raw,
	}
}

// Compares current cluster config represented as a json string to the cache.
// Returns true if cache is valid.
func (cache *ClusterCache) CheckClusterCache(config string, logger zerolog.Logger) bool {
	// compare cluster (non-component) configuration to cached cluster
	if cache.Kr8_Spec != gjson.Get(config, "._kr8_spec").Raw {
		logger.Info().Msg("_kr8_spec differs from cache")

		return false
	}

	if cache.Cluster != gjson.Get(config, "._cluster").Raw {
		logger.Info().Msg("_cluster differs from cache")

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

func CreateComponentCache(config string, componentPath string, listFiles []string) *ComponentCache {
	cacheResult := ComponentCache{}
	for _, file := range listFiles {
		currentHash, err := util.HashFile(filepath.Join(componentPath, file))
		if err != nil {
			return nil
		}
		cacheResult.ComponentFiles[file] = currentHash
	}

	return &cacheResult
}

func (cache *ComponentCache) CheckComponentCache(
	config string,
	componentName string,
	componentPath string,
	files []string,
	logger zerolog.Logger,
) bool {
	// compare cluster-level component config
	if cache.ComponentConfig != gjson.Get(config, componentName).Raw {
		logger.Info().Msg("component config differs from cache")
		// invalidate component cache
	}
	if len(cache.ComponentFiles) != len(files) {
		return false
	}
	for _, file := range files {
		hash, ok := cache.ComponentFiles[file]
		// didn't find file in cache
		if !ok {
			return false
		}
		currentHash, err := util.HashFile(filepath.Join(componentPath, file))
		if err != nil {
			logger.Warn().Err(err).Msg("issue hashing file, cache invalid")

			return false
		}
		if hash != currentHash {
			return false
		}
	}

	return true
}
