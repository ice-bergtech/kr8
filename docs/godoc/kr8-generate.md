# generate

```go
import "github.com/ice-bergtech/kr8/pkg/generate"
```

Package generate implements the logic for generating output files based on input data.

Combines a directory of cluster configurations with a directory of components \(along with some Jsonnet libs\) to generate output files.

The package prepares a Jsonnet VM and loads the necessary libraries and extvars. A new VM is created for each component.

## Index

- [func CheckIfUpdateNeeded\(outFile string, outStr string\) \(bool, error\)](<#CheckIfUpdateNeeded>)
- [func CleanOutputDir\(outputFileMap map\[string\]bool, componentOutputDir string\) error](<#CleanOutputDir>)
- [func GenProcessCluster\(clusterName string, clusterdir string, baseDir string, generateDirOverride string, kr8Opts types.Kr8Opts, clusterParamsFile string, filters util.PathFilterOptions, vmConfig types.VMConfig, pool \*ants.Pool, logger zerolog.Logger\) error](<#GenProcessCluster>)
- [func GenProcessComponent\(vmConfig types.VMConfig, componentName string, kr8Spec kr8\_types.Kr8ClusterSpec, kr8Opts types.Kr8Opts, config string, allConfig \*safeString, filters util.PathFilterOptions, paramsFile string, logger zerolog.Logger\) error](<#GenProcessComponent>)
- [func GenerateIncludesFiles\(includesFiles \[\]kr8\_types.Kr8ComponentSpecIncludeObject, kr8Spec kr8\_types.Kr8ClusterSpec, kr8Opts types.Kr8Opts, config string, componentName string, compPath string, componentOutputDir string, jvm \*jsonnet.VM, logger zerolog.Logger\) \(map\[string\]bool, error\)](<#GenerateIncludesFiles>)
- [func GetClusterParams\(clusterDir string, vmConfig types.VMConfig, logger zerolog.Logger\) \(map\[string\]string, error\)](<#GetClusterParams>)
- [func ProcessFile\(inputFile string, outputFile string, kr8Spec kr8\_types.Kr8ClusterSpec, componentName string, config string, incInfo kr8\_types.Kr8ComponentSpecIncludeObject, jvm \*jsonnet.VM, logger zerolog.Logger\) \(string, error\)](<#ProcessFile>)
- [func SetupAndConfigureVM\(vmConfig types.VMConfig, config string, kr8Spec kr8\_types.Kr8ClusterSpec, componentName string, compSpec kr8\_types.Kr8ComponentSpec, allConfig \*safeString, filters util.PathFilterOptions, paramsFile string, kr8Opts types.Kr8Opts, logger zerolog.Logger\) \(\*jsonnet.VM, string, error\)](<#SetupAndConfigureVM>)
- [func SetupJvmForComponent\(vmconfig types.VMConfig, config string, kr8Spec kr8\_types.Kr8ClusterSpec, componentName string\) \(\*jsonnet.VM, error\)](<#SetupJvmForComponent>)


<a name="CheckIfUpdateNeeded"></a>
## func [CheckIfUpdateNeeded](<https://github.com:icebergtech/kr8/blob/main/pkg/generate/file_system.go#L69>)

```go
func CheckIfUpdateNeeded(outFile string, outStr string) (bool, error)
```

Check if a file needs updating based on its current contents and the new contents.

<a name="CleanOutputDir"></a>
## func [CleanOutputDir](<https://github.com:icebergtech/kr8/blob/main/pkg/generate/file_system.go#L13>)

```go
func CleanOutputDir(outputFileMap map[string]bool, componentOutputDir string) error
```



<a name="GenProcessCluster"></a>
## func [GenProcessCluster](<https://github.com:icebergtech/kr8/blob/main/pkg/generate/generate.go#L342-L353>)

```go
func GenProcessCluster(clusterName string, clusterdir string, baseDir string, generateDirOverride string, kr8Opts types.Kr8Opts, clusterParamsFile string, filters util.PathFilterOptions, vmConfig types.VMConfig, pool *ants.Pool, logger zerolog.Logger) error
```

The root function for generating a cluster. Prepares and builds the cluster config. Build and processes the list of components.

<a name="GenProcessComponent"></a>
## func [GenProcessComponent](<https://github.com:icebergtech/kr8/blob/main/pkg/generate/generate.go#L114-L124>)

```go
func GenProcessComponent(vmConfig types.VMConfig, componentName string, kr8Spec kr8_types.Kr8ClusterSpec, kr8Opts types.Kr8Opts, config string, allConfig *safeString, filters util.PathFilterOptions, paramsFile string, logger zerolog.Logger) error
```

Root function for processing a kr8 component. Processes a component through a jsonnet VM to generate output files.

<a name="GenerateIncludesFiles"></a>
## func [GenerateIncludesFiles](<https://github.com:icebergtech/kr8/blob/main/pkg/generate/generate.go#L294-L304>)

```go
func GenerateIncludesFiles(includesFiles []kr8_types.Kr8ComponentSpecIncludeObject, kr8Spec kr8_types.Kr8ClusterSpec, kr8Opts types.Kr8Opts, config string, componentName string, compPath string, componentOutputDir string, jvm *jsonnet.VM, logger zerolog.Logger) (map[string]bool, error)
```

Generates the list of includes files for a component. Processes each includes file using the component's config. Returns an error if there's an issue with ANY includes file.

<a name="GetClusterParams"></a>
## func [GetClusterParams](<https://github.com:icebergtech/kr8/blob/main/pkg/generate/generate.go#L42>)

```go
func GetClusterParams(clusterDir string, vmConfig types.VMConfig, logger zerolog.Logger) (map[string]string, error)
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

<a name="SetupAndConfigureVM"></a>
## func [SetupAndConfigureVM](<https://github.com:icebergtech/kr8/blob/main/pkg/generate/generate.go#L187-L198>)

```go
func SetupAndConfigureVM(vmConfig types.VMConfig, config string, kr8Spec kr8_types.Kr8ClusterSpec, componentName string, compSpec kr8_types.Kr8ComponentSpec, allConfig *safeString, filters util.PathFilterOptions, paramsFile string, kr8Opts types.Kr8Opts, logger zerolog.Logger) (*jsonnet.VM, string, error)
```

Setup and configures a jsonnet VM for processing kr8 resources. Creates a new VM and does the following:

- loads cluster and component config
- loads jsonnet library files
- loads external file references

<a name="SetupJvmForComponent"></a>
## func [SetupJvmForComponent](<https://github.com:icebergtech/kr8/blob/main/pkg/generate/vm_helpers.go#L21-L26>)

```go
func SetupJvmForComponent(vmconfig types.VMConfig, config string, kr8Spec kr8_types.Kr8ClusterSpec, componentName string) (*jsonnet.VM, error)
```

This function sets up the JVM for a given component. It registers native functions, sets up post\-processing, and prunes parameters as required. It's faster to create this VM for each component, rather than re\-use. Default postprocessor just copies input to output.