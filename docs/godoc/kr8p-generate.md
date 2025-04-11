# generate

```go
import "github.com/ice-bergtech/kr8p/pkg/generate"
```

## Index

- [func CheckIfUpdateNeeded\(outFile string, outStr string\) \(bool, error\)](<#CheckIfUpdateNeeded>)
- [func CleanOutputDir\(outputFileMap map\[string\]bool, componentOutputDir string\) error](<#CleanOutputDir>)
- [func GenProcessCluster\(clusterName string, clusterdir string, baseDir string, generateDirOverride string, kr8Opts types.Kr8pOpts, clusterParamsFile string, filters util.PathFilterOptions, vmConfig types.VMConfig, pool \*ants.Pool\) error](<#GenProcessCluster>)
- [func GenProcessComponent\(vmconfig types.VMConfig, componentName string, kr8Spec types.Kr8pClusterSpec, kr8Opts types.Kr8pOpts, config string, allConfig \*safeString, filters util.PathFilterOptions, paramsFile string\) error](<#GenProcessComponent>)
- [func GenerateIncludesFiles\(includesFiles \[\]types.Kr8pComponentSpecIncludeObject, kr8Spec types.Kr8pClusterSpec, kr8Opts types.Kr8pOpts, config string, componentName string, compPath string, componentOutputDir string, jvm \*jsonnet.VM\) \(map\[string\]bool, error\)](<#GenerateIncludesFiles>)
- [func GetClusterParams\(clusterDir string, vmConfig types.VMConfig\) \(map\[string\]string, error\)](<#GetClusterParams>)
- [func ProcessFile\(inputFile string, outputFile string, kr8Spec types.Kr8pClusterSpec, componentName string, config string, incInfo types.Kr8pComponentSpecIncludeObject, jvm \*jsonnet.VM\) \(string, error\)](<#ProcessFile>)
- [func SetupAndConfigureVM\(vmconfig types.VMConfig, config string, kr8Spec types.Kr8pClusterSpec, componentName string, compSpec types.Kr8pComponentSpec, allConfig \*safeString, filters util.PathFilterOptions, paramsFile string, kr8Opts types.Kr8pOpts\) \(\*jsonnet.VM, string, error\)](<#SetupAndConfigureVM>)
- [func SetupJvmForComponent\(vmconfig types.VMConfig, config string, kr8Spec types.Kr8pClusterSpec, componentName string\) \(\*jsonnet.VM, error\)](<#SetupJvmForComponent>)


<a name="CheckIfUpdateNeeded"></a>
## func CheckIfUpdateNeeded

```go
func CheckIfUpdateNeeded(outFile string, outStr string) (bool, error)
```

Check if a file needs updating based on its current contents and the new contents.

<a name="CleanOutputDir"></a>
## func CleanOutputDir

```go
func CleanOutputDir(outputFileMap map[string]bool, componentOutputDir string) error
```



<a name="GenProcessCluster"></a>
## func GenProcessCluster

```go
func GenProcessCluster(clusterName string, clusterdir string, baseDir string, generateDirOverride string, kr8Opts types.Kr8pOpts, clusterParamsFile string, filters util.PathFilterOptions, vmConfig types.VMConfig, pool *ants.Pool) error
```



<a name="GenProcessComponent"></a>
## func GenProcessComponent

```go
func GenProcessComponent(vmconfig types.VMConfig, componentName string, kr8Spec types.Kr8pClusterSpec, kr8Opts types.Kr8pOpts, config string, allConfig *safeString, filters util.PathFilterOptions, paramsFile string) error
```



<a name="GenerateIncludesFiles"></a>
## func GenerateIncludesFiles

```go
func GenerateIncludesFiles(includesFiles []types.Kr8pComponentSpecIncludeObject, kr8Spec types.Kr8pClusterSpec, kr8Opts types.Kr8pOpts, config string, componentName string, compPath string, componentOutputDir string, jvm *jsonnet.VM) (map[string]bool, error)
```



<a name="GetClusterParams"></a>
## func GetClusterParams

```go
func GetClusterParams(clusterDir string, vmConfig types.VMConfig) (map[string]string, error)
```



<a name="ProcessFile"></a>
## func ProcessFile

```go
func ProcessFile(inputFile string, outputFile string, kr8Spec types.Kr8pClusterSpec, componentName string, config string, incInfo types.Kr8pComponentSpecIncludeObject, jvm *jsonnet.VM) (string, error)
```

Process an includes file. Based on the extension, it will process it differently.

.jsonnet: Imported and processed using jsonnet VM.

.yml, .yaml: Imported and processed through native function ParseYaml.

.tpl, .tmpl: Processed using component config and Sprig templating.

<a name="SetupAndConfigureVM"></a>
## func SetupAndConfigureVM

```go
func SetupAndConfigureVM(vmconfig types.VMConfig, config string, kr8Spec types.Kr8pClusterSpec, componentName string, compSpec types.Kr8pComponentSpec, allConfig *safeString, filters util.PathFilterOptions, paramsFile string, kr8Opts types.Kr8pOpts) (*jsonnet.VM, string, error)
```



<a name="SetupJvmForComponent"></a>
## func SetupJvmForComponent

```go
func SetupJvmForComponent(vmconfig types.VMConfig, config string, kr8Spec types.Kr8pClusterSpec, componentName string) (*jsonnet.VM, error)
```

This function sets up the JVM for a given component. It registers native functions, sets up post\-processing, and prunes parameters as required. It's faster to create this VM for each component, rather than re\-use. Default postprocessor just copies input to output.