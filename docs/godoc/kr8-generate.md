# generate

```go
import "github.com/ice-bergtech/kr8/pkg/generate"
```

## Index

- [func CheckIfUpdateNeeded\(outFile string, outStr string\) bool](<#CheckIfUpdateNeeded>)
- [func CleanOutputDir\(outputFileMap map\[string\]bool, componentOutputDir string\)](<#CleanOutputDir>)
- [func GenProcessCluster\(clusterName string, clusterdir string, baseDir string, generateDirOverride string, kr8Opts types.Kr8Opts, clusterParamsFile string, filters util.PathFilterOptions, vmConfig types.VMConfig, pool \*ants.Pool\)](<#GenProcessCluster>)
- [func GenProcessComponent\(vmconfig types.VMConfig, componentName string, kr8Spec types.Kr8ClusterSpec, kr8Opts types.Kr8Opts, config string, allConfig \*safeString, filters util.PathFilterOptions, paramsFile string\)](<#GenProcessComponent>)
- [func GenerateIncludesFiles\(includesFiles \[\]interface\{\}, kr8Spec types.Kr8ClusterSpec, kr8Opts types.Kr8Opts, config string, componentName string, compPath string, componentOutputDir string, jvm \*jsonnet.VM\) map\[string\]bool](<#GenerateIncludesFiles>)
- [func GetClusterParams\(clusterDir string, vmConfig types.VMConfig\) map\[string\]string](<#GetClusterParams>)
- [func ProcessFile\(inputFile string, outputFile string, kr8Spec types.Kr8ClusterSpec, componentName string, config string, incInfo types.Kr8ComponentSpecIncludeObject, jvm \*jsonnet.VM\) string](<#ProcessFile>)
- [func SetupJvmForComponent\(vmconfig types.VMConfig, config string, kr8Spec types.Kr8ClusterSpec, componentName string\) \(\*jsonnet.VM, error\)](<#SetupJvmForComponent>)


<a name="CheckIfUpdateNeeded"></a>
## func [CheckIfUpdateNeeded](<https://github.com/ice-bergtech/kr8/blob/main/pkg/generate/file_system.go#L54>)

```go
func CheckIfUpdateNeeded(outFile string, outStr string) bool
```

Check if a file needs updating based on its current contents and the new contents.

<a name="CleanOutputDir"></a>
## func [CleanOutputDir](<https://github.com/ice-bergtech/kr8/blob/main/pkg/generate/file_system.go#L12>)

```go
func CleanOutputDir(outputFileMap map[string]bool, componentOutputDir string)
```



<a name="GenProcessCluster"></a>
## func [GenProcessCluster](<https://github.com/ice-bergtech/kr8/blob/main/pkg/generate/generate.go#L262-L272>)

```go
func GenProcessCluster(clusterName string, clusterdir string, baseDir string, generateDirOverride string, kr8Opts types.Kr8Opts, clusterParamsFile string, filters util.PathFilterOptions, vmConfig types.VMConfig, pool *ants.Pool)
```



<a name="GenProcessComponent"></a>
## func [GenProcessComponent](<https://github.com/ice-bergtech/kr8/blob/main/pkg/generate/generate.go#L96-L105>)

```go
func GenProcessComponent(vmconfig types.VMConfig, componentName string, kr8Spec types.Kr8ClusterSpec, kr8Opts types.Kr8Opts, config string, allConfig *safeString, filters util.PathFilterOptions, paramsFile string)
```



<a name="GenerateIncludesFiles"></a>
## func [GenerateIncludesFiles](<https://github.com/ice-bergtech/kr8/blob/main/pkg/generate/generate.go#L201-L210>)

```go
func GenerateIncludesFiles(includesFiles []interface{}, kr8Spec types.Kr8ClusterSpec, kr8Opts types.Kr8Opts, config string, componentName string, compPath string, componentOutputDir string, jvm *jsonnet.VM) map[string]bool
```



<a name="GetClusterParams"></a>
## func [GetClusterParams](<https://github.com/ice-bergtech/kr8/blob/main/pkg/generate/generate.go#L31>)

```go
func GetClusterParams(clusterDir string, vmConfig types.VMConfig) map[string]string
```



<a name="ProcessFile"></a>
## func [ProcessFile](<https://github.com/ice-bergtech/kr8/blob/main/pkg/generate/file_processing.go#L67-L75>)

```go
func ProcessFile(inputFile string, outputFile string, kr8Spec types.Kr8ClusterSpec, componentName string, config string, incInfo types.Kr8ComponentSpecIncludeObject, jvm *jsonnet.VM) string
```

Process an includes file. Based on the extension, it will process it differently.

.jsonnet: Imported and processed using jsonnet VM.

.yml, .yaml: Imported and processed through native function ParseYaml.

.tpl, .tmpl: Processed using component config and Sprig templating.

<a name="SetupJvmForComponent"></a>
## func [SetupJvmForComponent](<https://github.com/ice-bergtech/kr8/blob/main/pkg/generate/vm_helpers.go#L19-L24>)

```go
func SetupJvmForComponent(vmconfig types.VMConfig, config string, kr8Spec types.Kr8ClusterSpec, componentName string) (*jsonnet.VM, error)
```

This function sets up the JVM for a given component. It registers native functions, sets up post\-processing, and prunes parameters as required. It's faster to create this VM for each component, rather than re\-use. Default postprocessor just copies input to output.