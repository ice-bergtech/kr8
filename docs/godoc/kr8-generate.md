# generate

```go
import "github.com/ice-bergtech/kr8/pkg/generate"
```

Package generate implements the logic for generating output files based on input data.

Combines:

- a directory of cluster configurations
- a directory of components
- a optional directory of Jsonnet library files

to generate output files. Output files can be structured yaml, or files generated from go templates.

The package prepares a Jsonnet VM and loads the necessary libraries and extVars. A new VM is created for each component.

## Index

- [func CalculateClusterComponentList\(clusterComponents map\[string\]gjson.Result, filters util.PathFilterOptions\) \[\]string](<#CalculateClusterComponentList>)
- [func CheckComponentCache\(cache \*kr8\_cache.DeploymentCache, compSpec kr8\_types.Kr8ComponentSpec, config string, componentName string, baseDir string, logger zerolog.Logger\) \(bool, \*kr8\_cache.ComponentCache, error\)](<#CheckComponentCache>)
- [func CheckIfUpdateNeeded\(outFile string, outStr string\) \(bool, error\)](<#CheckIfUpdateNeeded>)
- [func CleanOutputDir\(outputFileMap map\[string\]bool, componentOutputDir string\) error](<#CleanOutputDir>)
- [func CleanupOldComponentDirs\(existingComponents \[\]string, clusterComponents map\[string\]gjson.Result, kr8Spec \*kr8\_types.Kr8ClusterSpec, logger zerolog.Logger\)](<#CleanupOldComponentDirs>)
- [func CompileClusterConfiguration\(clusterName, clusterDir string, kr8Opts types.Kr8Opts, vmConfig types.VMConfig, generateDirOverride string, lint bool, logger zerolog.Logger\) \(\*kr8\_types.Kr8ClusterSpec, map\[string\]gjson.Result, error\)](<#CompileClusterConfiguration>)
- [func CreateClusterGenerateDirs\(kr8Spec kr8\_types.Kr8ClusterSpec\) \(\[\]string, error\)](<#CreateClusterGenerateDirs>)
- [func GatherClusterConfig\(clusterName, clusterDir string, kr8Opts types.Kr8Opts, vmConfig types.VMConfig, generateDirOverride string, filters util.PathFilterOptions, clusterParamsFile string, noop bool, lint bool, logger zerolog.Logger\) \(\*kr8\_types.Kr8ClusterSpec, \[\]string, string, error\)](<#GatherClusterConfig>)
- [func GenProcessCluster\(clusterConfig \*GenerateProcessRootConfig, pool \*ants.Pool, logger zerolog.Logger\) error](<#GenProcessCluster>)
- [func GenProcessComponent\(vmConfig types.VMConfig, componentName string, kr8Spec kr8\_types.Kr8ClusterSpec, kr8Opts types.Kr8Opts, config string, allConfig \*SafeString, filters util.PathFilterOptions, paramsFile string, cache \*kr8\_cache.DeploymentCache, lint bool, logger zerolog.Logger\) \(bool, \*kr8\_cache.ComponentCache, error\)](<#GenProcessComponent>)
- [func GenerateIncludesFiles\(includesFiles \[\]kr8\_types.Kr8ComponentSpecIncludeObject, kr8Spec kr8\_types.Kr8ClusterSpec, kr8Opts types.Kr8Opts, config string, componentName string, compPath string, componentOutputDir string, jvm \*jsonnet.VM, logger zerolog.Logger\) \(map\[string\]bool, error\)](<#GenerateIncludesFiles>)
- [func GetAllClusterParams\(clusterDir string, vmConfig types.VMConfig, jvm \*jsonnet.VM, lint bool, logger zerolog.Logger\) error](<#GetAllClusterParams>)
- [func GetClusterComponentParamsThreadSafe\(allConfig \*SafeString, config string, vmConfig types.VMConfig, kr8Spec kr8\_types.Kr8ClusterSpec, filters util.PathFilterOptions, paramsFile string, jvm \*jsonnet.VM, logger zerolog.Logger\) error](<#GetClusterComponentParamsThreadSafe>)
- [func GetClusterParams\(clusterDir string, vmConfig types.VMConfig, lint bool, logger zerolog.Logger\) \(map\[string\]string, error\)](<#GetClusterParams>)
- [func GetComponentFiles\(compSpec kr8\_types.Kr8ComponentSpec\) \[\]string](<#GetComponentFiles>)
- [func GetComponentPath\(config string, componentName string\) string](<#GetComponentPath>)
- [func LoadClusterCache\(kr8Spec \*kr8\_types.Kr8ClusterSpec, logger zerolog.Logger\) \(\*kr8\_cache.DeploymentCache, string\)](<#LoadClusterCache>)
- [func ProcessComponentFinalizer\(compSpec kr8\_types.Kr8ComponentSpec, componentOutputDir string, outputFileMap map\[string\]bool\) error](<#ProcessComponentFinalizer>)
- [func ProcessFile\(inputFile string, outputFile string, kr8Spec kr8\_types.Kr8ClusterSpec, componentName string, config string, incInfo kr8\_types.Kr8ComponentSpecIncludeObject, jvm \*jsonnet.VM, logger zerolog.Logger\) \(string, error\)](<#ProcessFile>)
- [func ProcessJsonnetToYaml\(jvm \*jsonnet.VM, input string, snippetFilename string\) \(string, error\)](<#ProcessJsonnetToYaml>)
- [func ProcessTemplate\(filename string, data gjson.Result\) \(string, error\)](<#ProcessTemplate>)
- [func RenderComponents\(config string, vmConfig types.VMConfig, kr8Spec kr8\_types.Kr8ClusterSpec, cache \*kr8\_cache.DeploymentCache, compList \[\]string, clusterParamsFile string, pool \*ants.Pool, kr8Opts types.Kr8Opts, filters util.PathFilterOptions, lint bool, logger zerolog.Logger\) \(map\[string\]kr8\_cache.ComponentCache, error\)](<#RenderComponents>)
- [func SetupBaseComponentJvm\(vmconfig types.VMConfig, config string, kr8Spec kr8\_types.Kr8ClusterSpec\) \(\*jsonnet.VM, error\)](<#SetupBaseComponentJvm>)
- [func SetupComponentVM\(vmConfig types.VMConfig, config string, kr8Spec kr8\_types.Kr8ClusterSpec, componentName string, compSpec kr8\_types.Kr8ComponentSpec, allConfig \*SafeString, filters util.PathFilterOptions, paramsFile string, kr8Opts types.Kr8Opts, lint bool, logger zerolog.Logger\) \(\*jsonnet.VM, string, error\)](<#SetupComponentVM>)
- [func SetupJvmForComponent\(jvm \*jsonnet.VM, config string, kr8Spec kr8\_types.Kr8ClusterSpec, componentName string\)](<#SetupJvmForComponent>)
- [func ValidateOrCreateCache\(cache \*kr8\_cache.DeploymentCache, config string, logger zerolog.Logger\) \*kr8\_cache.DeploymentCache](<#ValidateOrCreateCache>)
- [type GenerateProcessRootConfig](<#GenerateProcessRootConfig>)
- [type SafeString](<#SafeString>)


<a name="CalculateClusterComponentList"></a>
## func [CalculateClusterComponentList](<https://github.com:icebergtech/kr8/blob/main/pkg/generate/generate.go#L80-L83>)

```go
func CalculateClusterComponentList(clusterComponents map[string]gjson.Result, filters util.PathFilterOptions) []string
```

Calculates which components should be generated based on filters. Only processes specified component if it's defined in the cluster. Processes components in string sorted order. Sorts out orphaned, generated components directories.

<a name="CheckComponentCache"></a>
## func [CheckComponentCache](<https://github.com:icebergtech/kr8/blob/main/pkg/generate/generate.go#L197-L204>)

```go
func CheckComponentCache(cache *kr8_cache.DeploymentCache, compSpec kr8_types.Kr8ComponentSpec, config string, componentName string, baseDir string, logger zerolog.Logger) (bool, *kr8_cache.ComponentCache, error)
```

Compares a component's current state to a cache entry. Returns an up\-to\-date cache entry for the component. If the cache pointer is nil or cache invalid, a fresh cache entry will be generated to return.

<a name="CheckIfUpdateNeeded"></a>
## func [CheckIfUpdateNeeded](<https://github.com:icebergtech/kr8/blob/main/pkg/generate/file_system.go#L74>)

```go
func CheckIfUpdateNeeded(outFile string, outStr string) (bool, error)
```

Check if a file needs updating based on its current contents and potential new contents.

<a name="CleanOutputDir"></a>
## func [CleanOutputDir](<https://github.com:icebergtech/kr8/blob/main/pkg/generate/file_system.go#L15>)

```go
func CleanOutputDir(outputFileMap map[string]bool, componentOutputDir string) error
```

Removes all files not present in outputFileMap from componentOutputDir. checks if each file in the directory is present in the map, ignoring the bool value.

<a name="CleanupOldComponentDirs"></a>
## func [CleanupOldComponentDirs](<https://github.com:icebergtech/kr8/blob/main/pkg/generate/generate.go#L590-L595>)

```go
func CleanupOldComponentDirs(existingComponents []string, clusterComponents map[string]gjson.Result, kr8Spec *kr8_types.Kr8ClusterSpec, logger zerolog.Logger)
```

Go through each item in existingComponents and remove the file if it isn't in clusterComponents. Skips removing \`.kr8\_cache\` files.

<a name="CompileClusterConfiguration"></a>
## func [CompileClusterConfiguration](<https://github.com:icebergtech/kr8/blob/main/pkg/generate/generate.go#L613-L620>)

```go
func CompileClusterConfiguration(clusterName, clusterDir string, kr8Opts types.Kr8Opts, vmConfig types.VMConfig, generateDirOverride string, lint bool, logger zerolog.Logger) (*kr8_types.Kr8ClusterSpec, map[string]gjson.Result, error)
```

Build the list of parameter files to combine for the final cluster config by walking folder tree leaf to root.

<a name="CreateClusterGenerateDirs"></a>
## func [CreateClusterGenerateDirs](<https://github.com:icebergtech/kr8/blob/main/pkg/generate/file_system.go#L48>)

```go
func CreateClusterGenerateDirs(kr8Spec kr8_types.Kr8ClusterSpec) ([]string, error)
```

Create the root cluster output directory. Returns a list of cluster component output directories that already existed.

<a name="GatherClusterConfig"></a>
## func [GatherClusterConfig](<https://github.com:icebergtech/kr8/blob/main/pkg/generate/generate.go#L510-L520>)

```go
func GatherClusterConfig(clusterName, clusterDir string, kr8Opts types.Kr8Opts, vmConfig types.VMConfig, generateDirOverride string, filters util.PathFilterOptions, clusterParamsFile string, noop bool, lint bool, logger zerolog.Logger) (*kr8_types.Kr8ClusterSpec, []string, string, error)
```

Compiles configuration for each cluster. Creates and cleans output directories for generated cluster components. Uses the filter to determine which components to process. Renders the cluster\-level configuration for each component.

<a name="GenProcessCluster"></a>
## func [GenProcessCluster](<https://github.com:icebergtech/kr8/blob/main/pkg/generate/generate.go#L444-L448>)

```go
func GenProcessCluster(clusterConfig *GenerateProcessRootConfig, pool *ants.Pool, logger zerolog.Logger) error
```

The root function for generating a cluster. Prepares and builds the cluster config. Build and processes the list of components.

<a name="GenProcessComponent"></a>
## func [GenProcessComponent](<https://github.com:icebergtech/kr8/blob/main/pkg/generate/generate.go#L107-L119>)

```go
func GenProcessComponent(vmConfig types.VMConfig, componentName string, kr8Spec kr8_types.Kr8ClusterSpec, kr8Opts types.Kr8Opts, config string, allConfig *SafeString, filters util.PathFilterOptions, paramsFile string, cache *kr8_cache.DeploymentCache, lint bool, logger zerolog.Logger) (bool, *kr8_cache.ComponentCache, error)
```

Root function for processing a kr8\+ component. Processes a component through a jsonnet VM to generate output files.

<a name="GenerateIncludesFiles"></a>
## func [GenerateIncludesFiles](<https://github.com:icebergtech/kr8/blob/main/pkg/generate/generate.go#L388-L398>)

```go
func GenerateIncludesFiles(includesFiles []kr8_types.Kr8ComponentSpecIncludeObject, kr8Spec kr8_types.Kr8ClusterSpec, kr8Opts types.Kr8Opts, config string, componentName string, compPath string, componentOutputDir string, jvm *jsonnet.VM, logger zerolog.Logger) (map[string]bool, error)
```

Generates the list of includes files for a component. Processes each includes file using the component's config. Returns an error if there's an issue with ANY includes file.

<a name="GetAllClusterParams"></a>
## func [GetAllClusterParams](<https://github.com:icebergtech/kr8/blob/main/pkg/generate/generate.go#L325-L331>)

```go
func GetAllClusterParams(clusterDir string, vmConfig types.VMConfig, jvm *jsonnet.VM, lint bool, logger zerolog.Logger) error
```

Combine all the cluster params into a single object indexed by cluster name.

<a name="GetClusterComponentParamsThreadSafe"></a>
## func [GetClusterComponentParamsThreadSafe](<https://github.com:icebergtech/kr8/blob/main/pkg/generate/generate.go#L348-L357>)

```go
func GetClusterComponentParamsThreadSafe(allConfig *SafeString, config string, vmConfig types.VMConfig, kr8Spec kr8_types.Kr8ClusterSpec, filters util.PathFilterOptions, paramsFile string, jvm *jsonnet.VM, logger zerolog.Logger) error
```

Include full render of all component params for cluster. Only do this if we have not already cached it and don't already have it stored.

<a name="GetClusterParams"></a>
## func [GetClusterParams](<https://github.com:icebergtech/kr8/blob/main/pkg/generate/generate.go#L49-L54>)

```go
func GetClusterParams(clusterDir string, vmConfig types.VMConfig, lint bool, logger zerolog.Logger) (map[string]string, error)
```

Given a base directory, generates cluster\-level configuration for each cluster found. Gets list of clusters from \`util.GetClusterFilenames\(clusterDir\)\`.

<a name="GetComponentFiles"></a>
## func [GetComponentFiles](<https://github.com:icebergtech/kr8/blob/main/pkg/generate/generate.go#L232>)

```go
func GetComponentFiles(compSpec kr8_types.Kr8ComponentSpec) []string
```



<a name="GetComponentPath"></a>
## func [GetComponentPath](<https://github.com:icebergtech/kr8/blob/main/pkg/generate/generate.go#L320>)

```go
func GetComponentPath(config string, componentName string) string
```

Fetch a component path from raw cluster config.

<a name="LoadClusterCache"></a>
## func [LoadClusterCache](<https://github.com:icebergtech/kr8/blob/main/pkg/generate/generate.go#L567-L570>)

```go
func LoadClusterCache(kr8Spec *kr8_types.Kr8ClusterSpec, logger zerolog.Logger) (*kr8_cache.DeploymentCache, string)
```

Loads cluster cache based on a cluster spec. If cache is disabled, a nil deployment cache pointer is returned.

<a name="ProcessComponentFinalizer"></a>
## func [ProcessComponentFinalizer](<https://github.com:icebergtech/kr8/blob/main/pkg/generate/generate.go#L178-L182>)

```go
func ProcessComponentFinalizer(compSpec kr8_types.Kr8ComponentSpec, componentOutputDir string, outputFileMap map[string]bool) error
```

Final actions performed once a component is generated. Cleans extra files from output dir if not disabled in component spec.

<a name="ProcessFile"></a>
## func [ProcessFile](<https://github.com:icebergtech/kr8/blob/main/pkg/generate/file_processing.go#L77-L86>)

```go
func ProcessFile(inputFile string, outputFile string, kr8Spec kr8_types.Kr8ClusterSpec, componentName string, config string, incInfo kr8_types.Kr8ComponentSpecIncludeObject, jvm *jsonnet.VM, logger zerolog.Logger) (string, error)
```

Process an includes file. Based on the extension, the file is processed differently.

- .jsonnet: Imported and processed using jsonnet VM.
- .yml, .yaml: Imported and processed through native function ParseYaml.
- .tpl, .tmpl: Processed using component config and Sprig templating.

<a name="ProcessJsonnetToYaml"></a>
## func [ProcessJsonnetToYaml](<https://github.com:icebergtech/kr8/blob/main/pkg/generate/file_processing.go#L128>)

```go
func ProcessJsonnetToYaml(jvm *jsonnet.VM, input string, snippetFilename string) (string, error)
```

Processes an input string through the jsonnet VM and handles extracting the output into a yaml string. snippetFilename is used for error messages.

<a name="ProcessTemplate"></a>
## func [ProcessTemplate](<https://github.com:icebergtech/kr8/blob/main/pkg/generate/file_processing.go#L161>)

```go
func ProcessTemplate(filename string, data gjson.Result) (string, error)
```

Processes a template file with the given data. Loads file, parses template, then executes template.

<a name="RenderComponents"></a>
## func [RenderComponents](<https://github.com:icebergtech/kr8/blob/main/pkg/generate/generate.go#L665-L677>)

```go
func RenderComponents(config string, vmConfig types.VMConfig, kr8Spec kr8_types.Kr8ClusterSpec, cache *kr8_cache.DeploymentCache, compList []string, clusterParamsFile string, pool *ants.Pool, kr8Opts types.Kr8Opts, filters util.PathFilterOptions, lint bool, logger zerolog.Logger) (map[string]kr8_cache.ComponentCache, error)
```

Renders a list of components with a given Kr8ClusterSpec configuration. Each component is added to a sync.WaitGroup to be processed by the ants.Pool. Returns the cache results for all successfully generated components.

<a name="SetupBaseComponentJvm"></a>
## func [SetupBaseComponentJvm](<https://github.com:icebergtech/kr8/blob/main/pkg/generate/vm_helpers.go#L37-L41>)

```go
func SetupBaseComponentJvm(vmconfig types.VMConfig, config string, kr8Spec kr8_types.Kr8ClusterSpec) (*jsonnet.VM, error)
```

This function sets up the JVM for a given component. It sets up post\-processing, and prunes parameters as required. It's faster to create this VM for each component, rather than re\-use. Default postprocessor just copies input to output.

<a name="SetupComponentVM"></a>
## func [SetupComponentVM](<https://github.com:icebergtech/kr8/blob/main/pkg/generate/generate.go#L260-L272>)

```go
func SetupComponentVM(vmConfig types.VMConfig, config string, kr8Spec kr8_types.Kr8ClusterSpec, componentName string, compSpec kr8_types.Kr8ComponentSpec, allConfig *SafeString, filters util.PathFilterOptions, paramsFile string, kr8Opts types.Kr8Opts, lint bool, logger zerolog.Logger) (*jsonnet.VM, string, error)
```

Setup and configures a jsonnet VM for processing kr8\+ resources. Creates a new VM and does the following:

- loads cluster and component config
- loads jsonnet library files
- loads external file references

<a name="SetupJvmForComponent"></a>
## func [SetupJvmForComponent](<https://github.com:icebergtech/kr8/blob/main/pkg/generate/vm_helpers.go#L19-L24>)

```go
func SetupJvmForComponent(jvm *jsonnet.VM, config string, kr8Spec kr8_types.Kr8ClusterSpec, componentName string)
```

This function sets up component\-specific external code in the JVM. It makes the component config available to the jvm under the \`kr8\` extVar.

<a name="ValidateOrCreateCache"></a>
## func [ValidateOrCreateCache](<https://github.com:icebergtech/kr8/blob/main/pkg/generate/generate.go#L730-L734>)

```go
func ValidateOrCreateCache(cache *kr8_cache.DeploymentCache, config string, logger zerolog.Logger) *kr8_cache.DeploymentCache
```

For provided config, validates the cache object matches. If the cache is valid, it is returned. If cache is not valid, an empty deployment cache returned.

<a name="GenerateProcessRootConfig"></a>
## type [GenerateProcessRootConfig](<https://github.com:icebergtech/kr8/blob/main/pkg/generate/generate.go#L428-L439>)



```go
type GenerateProcessRootConfig struct {
    ClusterName       string
    ClusterDir        string
    BaseDir           string
    GenerateDir       string
    Kr8Opts           types.Kr8Opts
    ClusterParamsFile string
    Filters           util.PathFilterOptions
    VmConfig          types.VMConfig
    Noop              bool
    Lint              bool
}
```

<a name="SafeString"></a>
## type [SafeString](<https://github.com:icebergtech/kr8/blob/main/pkg/generate/generate.go#L40-L45>)

A thread\-safe string that can be used to store and retrieve configuration data.

```go
type SafeString struct {
    // contains filtered or unexported fields
}
```