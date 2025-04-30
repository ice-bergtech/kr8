# kr8\_cache

```go
import "github.com/ice-bergtech/kr8/pkg/kr8_cache"
```

Package kr8\_cache defines the structure for kr8\+ cluster\-component resource caching. Cache is based on cluster\-level config, component config, and component file reference hashes.

## Index

- [type ClusterCache](<#ClusterCache>)
  - [func CreateClusterCache\(config string\) \*ClusterCache](<#CreateClusterCache>)
  - [func \(cache \*ClusterCache\) CheckClusterCache\(config string, libDir string, logger zerolog.Logger\) bool](<#ClusterCache.CheckClusterCache>)
- [type ComponentCache](<#ComponentCache>)
  - [func CreateComponentCache\(config string, componentPath string, listFiles \[\]string\) \(\*ComponentCache, error\)](<#CreateComponentCache>)
  - [func \(cache \*ComponentCache\) CheckComponentCache\(config string, componentName string, componentPath string, baseDir string, files \[\]string, logger zerolog.Logger\) \(bool, \*ComponentCache\)](<#ComponentCache.CheckComponentCache>)
- [type DeploymentCache](<#DeploymentCache>)
  - [func InitDeploymentCache\(config string, baseDir string, cacheResults map\[string\]ComponentCache\) \*DeploymentCache](<#InitDeploymentCache>)
  - [func LoadClusterCache\(cacheFile string\) \(\*DeploymentCache, error\)](<#LoadClusterCache>)
  - [func \(cache \*DeploymentCache\) CheckClusterCache\(config string, baseDir string, logger zerolog.Logger\) bool](<#DeploymentCache.CheckClusterCache>)
  - [func \(cache \*DeploymentCache\) CheckClusterComponentCache\(config string, componentName string, componentPath string, baseDir string, files \[\]string, logger zerolog.Logger\) \(bool, \*ComponentCache, error\)](<#DeploymentCache.CheckClusterComponentCache>)
  - [func \(cache \*DeploymentCache\) WriteCache\(outFile string, compress bool\) error](<#DeploymentCache.WriteCache>)
- [type LibraryCache](<#LibraryCache>)
  - [func CreateLibraryCache\(baseDir string\) \*LibraryCache](<#CreateLibraryCache>)


<a name="ClusterCache"></a>
## type [ClusterCache](<https://github.com:icebergtech/kr8/blob/main/pkg/kr8_cache/cache.go#L124-L129>)

This is cluster\-level cache that applies to all components. If it is deemed invalid, the component cache is also invalid.

```go
type ClusterCache struct {
    // Raw cluster _kr8_spec object
    Kr8_Spec string `json:"kr8_spec"`
    // Raw cluster _cluster object
    Cluster string `json:"cluster"`
}
```

<a name="CreateClusterCache"></a>
### func [CreateClusterCache](<https://github.com:icebergtech/kr8/blob/main/pkg/kr8_cache/cache.go#L132>)

```go
func CreateClusterCache(config string) *ClusterCache
```

Stores the cluster kr8\_spec and cluster config as cluster\-level cache.

<a name="ClusterCache.CheckClusterCache"></a>
### func \(\*ClusterCache\) [CheckClusterCache](<https://github.com:icebergtech/kr8/blob/main/pkg/kr8_cache/cache.go#L164>)

```go
func (cache *ClusterCache) CheckClusterCache(config string, libDir string, logger zerolog.Logger) bool
```

Compares current cluster config represented as a json string to the cache. Returns true if cache is valid.

<a name="ComponentCache"></a>
## type [ComponentCache](<https://github.com:icebergtech/kr8/blob/main/pkg/kr8_cache/cache.go#L182-L187>)



```go
type ComponentCache struct {
    // Raw component config string
    ComponentConfig string `json:"component_config"`
    // Map of filenames to file hashes
    ComponentFiles map[string]string `json:"component_files"`
}
```

<a name="CreateComponentCache"></a>
### func [CreateComponentCache](<https://github.com:icebergtech/kr8/blob/main/pkg/kr8_cache/cache.go#L189>)

```go
func CreateComponentCache(config string, componentPath string, listFiles []string) (*ComponentCache, error)
```



<a name="ComponentCache.CheckComponentCache"></a>
### func \(\*ComponentCache\) [CheckComponentCache](<https://github.com:icebergtech/kr8/blob/main/pkg/kr8_cache/cache.go#L205-L212>)

```go
func (cache *ComponentCache) CheckComponentCache(config string, componentName string, componentPath string, baseDir string, files []string, logger zerolog.Logger) (bool, *ComponentCache)
```



<a name="DeploymentCache"></a>
## type [DeploymentCache](<https://github.com:icebergtech/kr8/blob/main/pkg/kr8_cache/cache.go#L38-L45>)

Object that contains the cache for a single cluster.

```go
type DeploymentCache struct {
    // A struct containing cluster-level cache values
    ClusterConfig *ClusterCache `json:"cluster_config"`
    // Map of cache entries for cluster components.
    // Depends on ClusterConfig cache being valid to be considered valid.
    ComponentConfigs map[string]ComponentCache `json:"component_config"`
    LibraryCache     *LibraryCache             `json:"library_cache"`
}
```

<a name="InitDeploymentCache"></a>
### func [InitDeploymentCache](<https://github.com:icebergtech/kr8/blob/main/pkg/kr8_cache/cache.go#L47>)

```go
func InitDeploymentCache(config string, baseDir string, cacheResults map[string]ComponentCache) *DeploymentCache
```



<a name="LoadClusterCache"></a>
### func [LoadClusterCache](<https://github.com:icebergtech/kr8/blob/main/pkg/kr8_cache/cache.go#L17>)

```go
func LoadClusterCache(cacheFile string) (*DeploymentCache, error)
```

Load cluster cache from a specified cache file.

<a name="DeploymentCache.CheckClusterCache"></a>
### func \(\*DeploymentCache\) [CheckClusterCache](<https://github.com:icebergtech/kr8/blob/main/pkg/kr8_cache/cache.go#L72>)

```go
func (cache *DeploymentCache) CheckClusterCache(config string, baseDir string, logger zerolog.Logger) bool
```



<a name="DeploymentCache.CheckClusterComponentCache"></a>
### func \(\*DeploymentCache\) [CheckClusterComponentCache](<https://github.com:icebergtech/kr8/blob/main/pkg/kr8_cache/cache.go#L81-L88>)

```go
func (cache *DeploymentCache) CheckClusterComponentCache(config string, componentName string, componentPath string, baseDir string, files []string, logger zerolog.Logger) (bool, *ComponentCache, error)
```



<a name="DeploymentCache.WriteCache"></a>
### func \(\*DeploymentCache\) [WriteCache](<https://github.com:icebergtech/kr8/blob/main/pkg/kr8_cache/cache.go#L57>)

```go
func (cache *DeploymentCache) WriteCache(outFile string, compress bool) error
```



<a name="LibraryCache"></a>
## type [LibraryCache](<https://github.com:icebergtech/kr8/blob/main/pkg/kr8_cache/cache.go#L117-L120>)



```go
type LibraryCache struct {
    Directory string            `json:"directory"`
    Entries   map[string]string `json:"entries"`
}
```

<a name="CreateLibraryCache"></a>
### func [CreateLibraryCache](<https://github.com:icebergtech/kr8/blob/main/pkg/kr8_cache/cache.go#L139>)

```go
func CreateLibraryCache(baseDir string) *LibraryCache
```

