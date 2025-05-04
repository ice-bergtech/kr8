# kr8+ Generate Caching


## Configuration

Caching is configured at the cluster `_kr8_spec` level:

```yaml
_kr8_spec+: {
    cache_enable: true, # default false
    cache_compress: true, # default true
}
```

If `cache_enable` is true, the cache file `.kr8_cache` will be generated and stored within the cluster's `generated` folder.
The cache is a json file containing base64-encoded information about the cluster.
When `cache_compress` is enabled, it is stored gzip'd, **greatly** reducing size.

The data structures for the caching functionality are described in [godoc/kr8-cache.md](../godoc/kr8-cache.md).

## Operation

Cache is created by storing kr8+ configuration in a json file on a per-cluster bases.

The cache contains information about the cluster, components, and jsonnet library folders.
Components and clusters will no be generated if their current configuration matches what is stored in the cache file.

### 1 Library Folders

If list of jsonnet library folders or files change, all clusters will be generated.

### 2 Cluster Config

Next cluster-level configuration is checked.
This includes:

* Cluster `_kr8_spec` object
* Cluster `_cluster` object

If either of these differ from the cache, all components for that cluster will be generated.

### 3 Cluster-Component Config and Component Files

Component-level caching is performed by:

* comparing the cluster-level component config for the component
* hashing all files in the component directory

If the configuration of the component hasn't changed and the files within the component's direcrory haven't changed, then it is skipped.

It does not account for files that a component references that are outside the component's root directory.

