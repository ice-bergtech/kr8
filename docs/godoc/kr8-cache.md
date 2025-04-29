# kr8\_cache

```go
import "github.com/ice-bergtech/kr8/pkg/kr8_cache"
```

Package kr8\_cache defines the structure for kr8\+ cluster\-component resource caching. Cache is based on cluster\-level config, component config, and component file reference hashes.

## Index

- [type ClusterCache](<#ClusterCache>)
  - [func CreateClusterCache\(config string\) \*ClusterCache](<#CreateClusterCache>)
  - [func \(cache \*ClusterCache\) CheckClusterCache\(config string, logger zerolog.Logger\) bool](<#ClusterCache.CheckClusterCache>)
- [type ComponentCache](<#ComponentCache>)
  - [func CreateComponentCache\(config string, componentPath string, listFiles \[\]string\) \(\*ComponentCache, error\)](<#CreateComponentCache>)
  - [func \(cache \*ComponentCache\) CheckComponentCache\(config string, componentName string, componentPath string, files \[\]string, logger zerolog.Logger\) bool](<#ComponentCache.CheckComponentCache>)
- [type DeploymentCache](<#DeploymentCache>)
  - [func LoadClusterCache\(cacheFile string\) \(\*DeploymentCache, error\)](<#LoadClusterCache>)
  - [func \(cache \*DeploymentCache\) CheckClusterCache\(config string, logger zerolog.Logger\) bool](<#DeploymentCache.CheckClusterCache>)
  - [func \(cache \*DeploymentCache\) CheckClusterComponentCache\(config string, componentName string, componentPath string, files \[\]string, logger zerolog.Logger\) bool](<#DeploymentCache.CheckClusterComponentCache>)
  - [func \(cache \*DeploymentCache\) WriteCache\(outFile string\) error](<#DeploymentCache.WriteCache>)


<a name="ClusterCache"></a>
## type [ClusterCache](<https://github.com:icebergtech/kr8/blob/main/pkg/kr8_cache/cache.go#L103-L108>)

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
### func [CreateClusterCache](<https://github.com:icebergtech/kr8/blob/main/pkg/kr8_cache/cache.go#L111>)

```go
func CreateClusterCache(config string) *ClusterCache
```

Stores the cluster kr8\_spec and cluster config as cluster\-level cache.

<a name="ClusterCache.CheckClusterCache"></a>
### func \(\*ClusterCache\) [CheckClusterCache](<https://github.com:icebergtech/kr8/blob/main/pkg/kr8_cache/cache.go#L120>)

```go
func (cache *ClusterCache) CheckClusterCache(config string, logger zerolog.Logger) bool
```

Compares current cluster config represented as a json string to the cache. Returns true if cache is valid.

<a name="ComponentCache"></a>
## type [ComponentCache](<https://github.com:icebergtech/kr8/blob/main/pkg/kr8_cache/cache.go#L137-L142>)



```go
type ComponentCache struct {
    // Raw component config string
    ComponentConfig string `json:"component_config"`
    // Map of filenames to file hashes
    ComponentFiles map[string]string `json:"component_files"`
}
```

<a name="CreateComponentCache"></a>
### func [CreateComponentCache](<https://github.com:icebergtech/kr8/blob/main/pkg/kr8_cache/cache.go#L144>)

```go
func CreateComponentCache(config string, componentPath string, listFiles []string) (*ComponentCache, error)
```



<a name="ComponentCache.CheckComponentCache"></a>
### func \(\*ComponentCache\) [CheckComponentCache](<https://github.com:icebergtech/kr8/blob/main/pkg/kr8_cache/cache.go#L160-L166>)

```go
func (cache *ComponentCache) CheckComponentCache(config string, componentName string, componentPath string, files []string, logger zerolog.Logger) bool
```



<a name="DeploymentCache"></a>
## type [DeploymentCache](<https://github.com:icebergtech/kr8/blob/main/pkg/kr8_cache/cache.go#L43-L48>)

Object that contains the cache for a single cluster.

```go
type DeploymentCache struct {
    ClusterConfig *ClusterCache `json:"cluster_config"`
    // Map of cache entries for cluster components.
    // Depends on ClusterConfig cache being valid to be considered valid.
    ComponentConfigs map[string]ComponentCache `json:"component_config"`
}
```

<a name="LoadClusterCache"></a>
### func [LoadClusterCache](<https://github.com:icebergtech/kr8/blob/main/pkg/kr8_cache/cache.go#L17>)

```go
func LoadClusterCache(cacheFile string) (*DeploymentCache, error)
```

Load cluster cache from a specified cache file.

<a name="DeploymentCache.CheckClusterCache"></a>
### func \(\*DeploymentCache\) [CheckClusterCache](<https://github.com:icebergtech/kr8/blob/main/pkg/kr8_cache/cache.go#L71>)

```go
func (cache *DeploymentCache) CheckClusterCache(config string, logger zerolog.Logger) bool
```



<a name="DeploymentCache.CheckClusterComponentCache"></a>
### func \(\*DeploymentCache\) [CheckClusterComponentCache](<https://github.com:icebergtech/kr8/blob/main/pkg/kr8_cache/cache.go#L80-L86>)

```go
func (cache *DeploymentCache) CheckClusterComponentCache(config string, componentName string, componentPath string, files []string, logger zerolog.Logger) bool
```



<a name="DeploymentCache.WriteCache"></a>
### func \(\*DeploymentCache\) [WriteCache](<https://github.com:icebergtech/kr8/blob/main/pkg/kr8_cache/cache.go#L50>)

```go
func (cache *DeploymentCache) WriteCache(outFile string) error
```

