# generate

```go
import "github.com/ice-bergtech/kr8/pkg/generate"
```

Package generate implements the logic for generating output files based on input data.

Combines a directory of cluster configurations with a directory of components \(along with some Jsonnet libs\) to generate output files.

The package prepares a Jsonnet VM and loads the necessary libraries and extvars. A new VM is created for each component.

## Index

- [func CalculateClusterComponentList\(clusterComponents map\[string\]gjson.Result, filters util.PathFilterOptions, existingClusterComponents \[\]string\) \[\]string](<#CalculateClusterComponentList>)
- [func CheckComponentCache\(cache \*kr8\_cache.DeploymentCache, compSpec kr8\_types.Kr8ComponentSpec, config string, componentName string, baseDir string, logger zerolog.Logger\) \(bool, \*kr8\_cache.ComponentCache, error\)](<#CheckComponentCache>)
- [func CheckIfUpdateNeeded\(outFile string, outStr string\) \(bool, error\)](<#CheckIfUpdateNeeded>)
- [func CleanOutputDir\(outputFileMap map\[string\]bool, componentOutputDir string\) error](<#CleanOutputDir>)
- [func CleanupOldComponentDirs\(existingComponents \[\]string, clusterComponents map\[string\]gjson.Result, kr8Spec \*kr8\_types.Kr8ClusterSpec, logger zerolog.Logger\)](<#CleanupOldComponentDirs>)
- [func CompileClusterConfiguration\(clusterName, clusterDir string, kr8Opts types.Kr8Opts, vmConfig types.VMConfig, generateDirOverride string, logger zerolog.Logger\) \(\*kr8\_types.Kr8ClusterSpec, map\[string\]gjson.Result, error\)](<#CompileClusterConfiguration>)
- [func CreateClusterGenerateDirs\(kr8Spec kr8\_types.Kr8ClusterSpec\) \(\[\]string, error\)](<#CreateClusterGenerateDirs>)
- [func GatherClusterConfig\(clusterName, clusterDir string, kr8Opts types.Kr8Opts, vmConfig types.VMConfig, generateDirOverride string, filters util.PathFilterOptions, clusterParamsFile string, logger zerolog.Logger\) \(\*kr8\_types.Kr8ClusterSpec, \[\]string, string, error\)](<#GatherClusterConfig>)
- [func GenProcessCluster\(clusterName string, clusterdir string, baseDir string, generateDirOverride string, kr8Opts types.Kr8Opts, clusterParamsFile string, filters util.PathFilterOptions, vmConfig types.VMConfig, pool \*ants.Pool, enableCache bool, logger zerolog.Logger\) error](<#GenProcessCluster>)
- [func GenProcessComponent\(vmConfig types.VMConfig, componentName string, kr8Spec kr8\_types.Kr8ClusterSpec, kr8Opts types.Kr8Opts, config string, allConfig \*SafeString, filters util.PathFilterOptions, paramsFile string, cache \*kr8\_cache.DeploymentCache, logger zerolog.Logger\) \(bool, \*kr8\_cache.ComponentCache, error\)](<#GenProcessComponent>)
- [func GenerateIncludesFiles\(includesFiles \[\]kr8\_types.Kr8ComponentSpecIncludeObject, kr8Spec kr8\_types.Kr8ClusterSpec, kr8Opts types.Kr8Opts, config string, componentName string, compPath string, componentOutputDir string, jvm \*jsonnet.VM, logger zerolog.Logger\) \(map\[string\]bool, error\)](<#GenerateIncludesFiles>)
- [func GetAllClusterParams\(clusterDir string, vmConfig types.VMConfig, jvm \*jsonnet.VM, logger zerolog.Logger\) error](<#GetAllClusterParams>)
- [func GetClusterComponentParamsThreadsafe\(allConfig \*SafeString, config string, vmConfig types.VMConfig, kr8Spec kr8\_types.Kr8ClusterSpec, filters util.PathFilterOptions, paramsFile string, jvm \*jsonnet.VM, logger zerolog.Logger\) error](<#GetClusterComponentParamsThreadsafe>)
- [func GetClusterParams\(clusterDir string, vmConfig types.VMConfig, logger zerolog.Logger\) \(map\[string\]string, error\)](<#GetClusterParams>)
- [func GetComponentFiles\(compSpec kr8\_types.Kr8ComponentSpec\) \[\]string](<#GetComponentFiles>)
- [func GetComponentPath\(config string, componentName string\) string](<#GetComponentPath>)
- [func LoadClusterCache\(kr8Spec \*kr8\_types.Kr8ClusterSpec, logger zerolog.Logger\) \(\*kr8\_cache.DeploymentCache, string\)](<#LoadClusterCache>)
- [func ProcessComponentFinalizer\(compSpec kr8\_types.Kr8ComponentSpec, componentOutputDir string, outputFileMap map\[string\]bool\) error](<#ProcessComponentFinalizer>)
- [func ProcessFile\(inputFile string, outputFile string, kr8Spec kr8\_types.Kr8ClusterSpec, componentName string, config string, incInfo kr8\_types.Kr8ComponentSpecIncludeObject, jvm \*jsonnet.VM, logger zerolog.Logger\) \(string, error\)](<#ProcessFile>)
- [func RenderComponents\(config string, vmConfig types.VMConfig, kr8Spec kr8\_types.Kr8ClusterSpec, cache \*kr8\_cache.DeploymentCache, compList \[\]string, clusterParamsFile string, pool \*ants.Pool, kr8Opts types.Kr8Opts, filters util.PathFilterOptions, logger zerolog.Logger\) \(map\[string\]kr8\_cache.ComponentCache, error\)](<#RenderComponents>)
- [func SetupBaseComponentJvm\(vmconfig types.VMConfig, config string, kr8Spec kr8\_types.Kr8ClusterSpec\) \(\*jsonnet.VM, error\)](<#SetupBaseComponentJvm>)
- [func SetupComponentVM\(vmConfig types.VMConfig, config string, kr8Spec kr8\_types.Kr8ClusterSpec, componentName string, compSpec kr8\_types.Kr8ComponentSpec, allConfig \*SafeString, filters util.PathFilterOptions, paramsFile string, kr8Opts types.Kr8Opts, logger zerolog.Logger\) \(\*jsonnet.VM, string, error\)](<#SetupComponentVM>)
- [func SetupJvmForComponent\(jvm \*jsonnet.VM, config string, kr8Spec kr8\_types.Kr8ClusterSpec, componentName string\)](<#SetupJvmForComponent>)
- [func StoreClusterComponentCache\(kr8Spec \*kr8\_types.Kr8ClusterSpec, baseDir string, config string, cacheResults map\[string\]kr8\_cache.ComponentCache, cacheFilePath string\) error](<#StoreClusterComponentCache>)
- [func ValidateOrCreateCache\(cache \*kr8\_cache.DeploymentCache, config string, logger zerolog.Logger\) \*kr8\_cache.DeploymentCache](<#ValidateOrCreateCache>)
- [type SafeString](<#SafeString>)


<a name="CalculateClusterComponentList"></a>
## func [CalculateClusterComponentList](<https://github.com:icebergtech/kr8/blob/main/pkg/generate/generate.go#L70-L74>)

```go
func CalculateClusterComponentList(clusterComponents map[string]gjson.Result, filters util.PathFilterOptions, existingClusterComponents []string) []string
```

Calculates which components should be generated based on filters. Only processes specified component if it's defined in the cluster. Processes components in string sorted order. Sorts out orphaned, generated components directories.

<a name="CheckComponentCache"></a>
## func [CheckComponentCache](<https://github.com:icebergtech/kr8/blob/main/pkg/generate/generate.go#L180-L187>)

```go
func CheckComponentCache(cache *kr8_cache.DeploymentCache, compSpec kr8_types.Kr8ComponentSpec, config string, componentName string, baseDir string, logger zerolog.Logger) (bool, *kr8_cache.ComponentCache, error)
```



<a name="CheckIfUpdateNeeded"></a>
## func [CheckIfUpdateNeeded](<https://github.com:icebergtech/kr8/blob/main/pkg/generate/file_system.go#L71>)

```go
func CheckIfUpdateNeeded(outFile string, outStr string) (bool, error)
```

Check if a file needs updating based on its current contents and the new contents.

<a name="CleanOutputDir"></a>
## func [CleanOutputDir](<https://github.com:icebergtech/kr8/blob/main/pkg/generate/file_system.go#L13>)

```go
func CleanOutputDir(outputFileMap map[string]bool, componentOutputDir string) error
```



<a name="CleanupOldComponentDirs"></a>
## func [CleanupOldComponentDirs](<https://github.com:icebergtech/kr8/blob/main/pkg/generate/generate.go#L555-L560>)

```go
func CleanupOldComponentDirs(existingComponents []string, clusterComponents map[string]gjson.Result, kr8Spec *kr8_types.Kr8ClusterSpec, logger zerolog.Logger)
```



<a name="CompileClusterConfiguration"></a>
## func [CompileClusterConfiguration](<https://github.com:icebergtech/kr8/blob/main/pkg/generate/generate.go#L578-L584>)

```go
func CompileClusterConfiguration(clusterName, clusterDir string, kr8Opts types.Kr8Opts, vmConfig types.VMConfig, generateDirOverride string, logger zerolog.Logger) (*kr8_types.Kr8ClusterSpec, map[string]gjson.Result, error)
```

Build the list of cluster parameter files to combine by walking folder tree leaf to root.

<a name="CreateClusterGenerateDirs"></a>
## func [CreateClusterGenerateDirs](<https://github.com:icebergtech/kr8/blob/main/pkg/generate/file_system.go#L45>)

```go
func CreateClusterGenerateDirs(kr8Spec kr8_types.Kr8ClusterSpec) ([]string, error)
```

Create the root cluster output directory. Returns a list of cluster component output directories that already existed.

<a name="GatherClusterConfig"></a>
## func [GatherClusterConfig](<https://github.com:icebergtech/kr8/blob/main/pkg/generate/generate.go#L472-L480>)

```go
func GatherClusterConfig(clusterName, clusterDir string, kr8Opts types.Kr8Opts, vmConfig types.VMConfig, generateDirOverride string, filters util.PathFilterOptions, clusterParamsFile string, logger zerolog.Logger) (*kr8_types.Kr8ClusterSpec, []string, string, error)
```



<a name="GenProcessCluster"></a>
## func [GenProcessCluster](<https://github.com:icebergtech/kr8/blob/main/pkg/generate/generate.go#L410-L422>)

```go
func GenProcessCluster(clusterName string, clusterdir string, baseDir string, generateDirOverride string, kr8Opts types.Kr8Opts, clusterParamsFile string, filters util.PathFilterOptions, vmConfig types.VMConfig, pool *ants.Pool, enableCache bool, logger zerolog.Logger) error
```

The root function for generating a cluster. Prepares and builds the cluster config. Build and processes the list of components.

<a name="GenProcessComponent"></a>
## func [GenProcessComponent](<https://github.com:icebergtech/kr8/blob/main/pkg/generate/generate.go#L98-L109>)

```go
func GenProcessComponent(vmConfig types.VMConfig, componentName string, kr8Spec kr8_types.Kr8ClusterSpec, kr8Opts types.Kr8Opts, config string, allConfig *SafeString, filters util.PathFilterOptions, paramsFile string, cache *kr8_cache.DeploymentCache, logger zerolog.Logger) (bool, *kr8_cache.ComponentCache, error)
```

Root function for processing a kr8 component. Processes a component through a jsonnet VM to generate output files.

<a name="GenerateIncludesFiles"></a>
## func [GenerateIncludesFiles](<https://github.com:icebergtech/kr8/blob/main/pkg/generate/generate.go#L362-L372>)

```go
func GenerateIncludesFiles(includesFiles []kr8_types.Kr8ComponentSpecIncludeObject, kr8Spec kr8_types.Kr8ClusterSpec, kr8Opts types.Kr8Opts, config string, componentName string, compPath string, componentOutputDir string, jvm *jsonnet.VM, logger zerolog.Logger) (map[string]bool, error)
```

Generates the list of includes files for a component. Processes each includes file using the component's config. Returns an error if there's an issue with ANY includes file.

<a name="GetAllClusterParams"></a>
## func [GetAllClusterParams](<https://github.com:icebergtech/kr8/blob/main/pkg/generate/generate.go#L306>)

```go
func GetAllClusterParams(clusterDir string, vmConfig types.VMConfig, jvm *jsonnet.VM, logger zerolog.Logger) error
```

Combine all the cluster params into a single object indexed by cluster name.

<a name="GetClusterComponentParamsThreadsafe"></a>
## func [GetClusterComponentParamsThreadsafe](<https://github.com:icebergtech/kr8/blob/main/pkg/generate/generate.go#L323-L332>)

```go
func GetClusterComponentParamsThreadsafe(allConfig *SafeString, config string, vmConfig types.VMConfig, kr8Spec kr8_types.Kr8ClusterSpec, filters util.PathFilterOptions, paramsFile string, jvm *jsonnet.VM, logger zerolog.Logger) error
```

Include full render of all component params for cluster. Only do this if we have not already cached it and don't already have it stored.

<a name="GetClusterParams"></a>
## func [GetClusterParams](<https://github.com:icebergtech/kr8/blob/main/pkg/generate/generate.go#L45>)

```go
func GetClusterParams(clusterDir string, vmConfig types.VMConfig, logger zerolog.Logger) (map[string]string, error)
```



<a name="GetComponentFiles"></a>
## func [GetComponentFiles](<https://github.com:icebergtech/kr8/blob/main/pkg/generate/generate.go#L216>)

```go
func GetComponentFiles(compSpec kr8_types.Kr8ComponentSpec) []string
```



<a name="GetComponentPath"></a>
## func [GetComponentPath](<https://github.com:icebergtech/kr8/blob/main/pkg/generate/generate.go#L301>)

```go
func GetComponentPath(config string, componentName string) string
```



<a name="LoadClusterCache"></a>
## func [LoadClusterCache](<https://github.com:icebergtech/kr8/blob/main/pkg/generate/generate.go#L519-L522>)

```go
func LoadClusterCache(kr8Spec *kr8_types.Kr8ClusterSpec, logger zerolog.Logger) (*kr8_cache.DeploymentCache, string)
```



<a name="ProcessComponentFinalizer"></a>
## func [ProcessComponentFinalizer](<https://github.com:icebergtech/kr8/blob/main/pkg/generate/generate.go#L164-L168>)

```go
func ProcessComponentFinalizer(compSpec kr8_types.Kr8ComponentSpec, componentOutputDir string, outputFileMap map[string]bool) error
```



<a name="ProcessFile"></a>
## func [ProcessFile](<https://github.com:icebergtech/kr8/blob/main/pkg/generate/file_processing.go#L74-L83>)

```go
func ProcessFile(inputFile string, outputFile string, kr8Spec kr8_types.Kr8ClusterSpec, componentName string, config string, incInfo kr8_types.Kr8ComponentSpecIncludeObject, jvm *jsonnet.VM, logger zerolog.Logger) (string, error)
```

Process an includes file. Based on the extension, the file is processed differently.

- .jsonnet: Imported and processed using jsonnet VM.
- .yml, .yaml: Imported and processed through native function ParseYaml.
- .tpl, .tmpl: Processed using component config and Sprig templating.

<a name="RenderComponents"></a>
## func [RenderComponents](<https://github.com:icebergtech/kr8/blob/main/pkg/generate/generate.go#L623-L634>)

```go
func RenderComponents(config string, vmConfig types.VMConfig, kr8Spec kr8_types.Kr8ClusterSpec, cache *kr8_cache.DeploymentCache, compList []string, clusterParamsFile string, pool *ants.Pool, kr8Opts types.Kr8Opts, filters util.PathFilterOptions, logger zerolog.Logger) (map[string]kr8_cache.ComponentCache, error)
```

Renders a list of components with a given Kr8ClusterSpec configuration. Each component is processed by a process thread from a thread pool.

<a name="SetupBaseComponentJvm"></a>
## func [SetupBaseComponentJvm](<https://github.com:icebergtech/kr8/blob/main/pkg/generate/vm_helpers.go#L37-L41>)

```go
func SetupBaseComponentJvm(vmconfig types.VMConfig, config string, kr8Spec kr8_types.Kr8ClusterSpec) (*jsonnet.VM, error)
```

This function sets up the JVM for a given component. It registers native functions, sets up post\-processing, and prunes parameters as required. It's faster to create this VM for each component, rather than re\-use. Default postprocessor just copies input to output.

<a name="SetupComponentVM"></a>
## func [SetupComponentVM](<https://github.com:icebergtech/kr8/blob/main/pkg/generate/generate.go#L244-L255>)

```go
func SetupComponentVM(vmConfig types.VMConfig, config string, kr8Spec kr8_types.Kr8ClusterSpec, componentName string, compSpec kr8_types.Kr8ComponentSpec, allConfig *SafeString, filters util.PathFilterOptions, paramsFile string, kr8Opts types.Kr8Opts, logger zerolog.Logger) (*jsonnet.VM, string, error)
```

Setup and configures a jsonnet VM for processing kr8 resources. Creates a new VM and does the following:

- loads cluster and component config
- loads jsonnet library files
- loads external file references

<a name="SetupJvmForComponent"></a>
## func [SetupJvmForComponent](<https://github.com:icebergtech/kr8/blob/main/pkg/generate/vm_helpers.go#L19-L24>)

```go
func SetupJvmForComponent(jvm *jsonnet.VM, config string, kr8Spec kr8_types.Kr8ClusterSpec, componentName string)
```

This function sets up component\-specific external code in the JVM. It makes the component config available to the jvm under the \`kr8\` extVar.

<a name="StoreClusterComponentCache"></a>
## func [StoreClusterComponentCache](<https://github.com:icebergtech/kr8/blob/main/pkg/generate/generate.go#L541-L547>)

```go
func StoreClusterComponentCache(kr8Spec *kr8_types.Kr8ClusterSpec, baseDir string, config string, cacheResults map[string]kr8_cache.ComponentCache, cacheFilePath string) error
```

Generates and stores cluster cache provided the config and component caches.

<a name="ValidateOrCreateCache"></a>
## func [ValidateOrCreateCache](<https://github.com:icebergtech/kr8/blob/main/pkg/generate/generate.go#L684-L688>)

```go
func ValidateOrCreateCache(cache *kr8_cache.DeploymentCache, config string, logger zerolog.Logger) *kr8_cache.DeploymentCache
```



<a name="SafeString"></a>
## type [SafeString](<https://github.com:icebergtech/kr8/blob/main/pkg/generate/generate.go#L38-L43>)

A thread\-safe string that can be used to store and retrieve configuration data.

```go
type SafeString struct {
    // contains filtered or unexported fields
}
```