# cmd

```go
import "github.com/ice-bergtech/kr8/cmd"
```

## Index

- [Variables](<#variables>)
- [func Execute\(version string\)](<#Execute>)
- [func JsonnetVM\(vmconfig VMConfig\) \(\*jsonnet.VM, error\)](<#JsonnetVM>)
- [func Pretty\(input string, colorOutput bool\) string](<#Pretty>)
- [func RegisterNativeFuncs\(vm \*jsonnet.VM\)](<#RegisterNativeFuncs>)
- [type Cluster](<#Cluster>)
- [type ClusterComponent](<#ClusterComponent>)
- [type ClusterJsonnet](<#ClusterJsonnet>)
- [type ClusterSpec](<#ClusterSpec>)
  - [func CreateClusterSpec\(clusterName string, spec gjson.Result, baseDir string, genDirOverride string\) \(ClusterSpec, error\)](<#CreateClusterSpec>)
- [type Clusters](<#Clusters>)
- [type CmdGetOptions](<#CmdGetOptions>)
- [type CmdJsonnetOptions](<#CmdJsonnetOptions>)
- [type ComponentJsonnet](<#ComponentJsonnet>)
- [type ComponentSpec](<#ComponentSpec>)
  - [func CreateComponentSpec\(spec gjson.Result\) \(ComponentSpec, error\)](<#CreateComponentSpec>)
- [type ExtFileVar](<#ExtFileVar>)
- [type IncludeFileEntryStruct](<#IncludeFileEntryStruct>)
- [type IncludeFileSpec](<#IncludeFileSpec>)
- [type PathFilterOptions](<#PathFilterOptions>)
- [type VMConfig](<#VMConfig>)


## Variables

<a name="RootCmd"></a>RootCmd represents the base command when called without any subcommands

```go
var RootCmd = &cobra.Command{
    Use:   "kr8",
    Short: "Kubernetes config parameter framework",
    Long: `A tool to generate Kubernetes configuration from a hierarchy
	of jsonnet files`,
}
```

<a name="Version"></a>exported Version variable

```go
var Version string
```

<a name="Execute"></a>
## func [Execute](<https://github.com/ice-bergtech/kr8/blob/main/cmd/root.go#L35>)

```go
func Execute(version string)
```

Execute adds all child commands to the root command sets flags appropriately. This is called by main.main\(\). It only needs to happen once to the rootCmd.

<a name="JsonnetVM"></a>
## func [JsonnetVM](<https://github.com/ice-bergtech/kr8/blob/main/cmd/jsonnet.go#L49>)

```go
func JsonnetVM(vmconfig VMConfig) (*jsonnet.VM, error)
```



<a name="Pretty"></a>
## func [Pretty](<https://github.com/ice-bergtech/kr8/blob/main/cmd/json.go#L8>)

```go
func Pretty(input string, colorOutput bool) string
```



<a name="RegisterNativeFuncs"></a>
## func [RegisterNativeFuncs](<https://github.com/ice-bergtech/kr8/blob/main/cmd/jsonnet.go#L140>)

```go
func RegisterNativeFuncs(vm *jsonnet.VM)
```



<a name="Cluster"></a>
## type [Cluster](<https://github.com/ice-bergtech/kr8/blob/main/cmd/types.go#L13-L16>)

init a struct for a single item

```go
type Cluster struct {
    Name string `json:"name"`
    Path string `json:"-"`
}
```

<a name="ClusterComponent"></a>
## type [ClusterComponent](<https://github.com/ice-bergtech/kr8/blob/main/cmd/types.go#L32-L35>)



```go
type ClusterComponent struct {
    // The path to a component folder that contains a params.jsonnet file
    Path string `json:"path"`
}
```

<a name="ClusterJsonnet"></a>
## type [ClusterJsonnet](<https://github.com/ice-bergtech/kr8/blob/main/cmd/types.go#L23-L30>)



```go
type ClusterJsonnet struct {
    // kr8 configuration for how to process the cluster
    ClusterSpec ClusterSpec `json:"_kr8_spec"`
    // Cluster Level configuration that components can reference
    Cluster Cluster `json:"_cluster"`
    // Distictly named components.
    Components map[string]ClusterComponent `json:"_components"`
}
```

<a name="ClusterSpec"></a>
## type [ClusterSpec](<https://github.com/ice-bergtech/kr8/blob/main/cmd/types.go#L37-L50>)



```go
type ClusterSpec struct {
    // A jsonnet function that each output entry is processed through. Default `function(input) input`
    PostProcessor string `json:"postprocessor"`
    // The name of the root generate directory. Default `generated`
    GenerateDir string `json:"generate_dir"`
    // if this is true, we don't use the full file path to generate output file names
    GenerateShortNames bool `json:"generate_short_names"`
    // if this is true, we prune component parameters
    PruneParams bool `json:"prune_params"`
    // The root directory for the cluster. Default `clusters`
    ClusterDir string `json:"-"`
    // The name of the current cluster
    Name string `json:"-"`
}
```

<a name="CreateClusterSpec"></a>
### func [CreateClusterSpec](<https://github.com/ice-bergtech/kr8/blob/main/cmd/types.go#L52>)

```go
func CreateClusterSpec(clusterName string, spec gjson.Result, baseDir string, genDirOverride string) (ClusterSpec, error)
```



<a name="Clusters"></a>
## type [Clusters](<https://github.com/ice-bergtech/kr8/blob/main/cmd/types.go#L19-L21>)

init a grouping struct

```go
type Clusters struct {
    Cluster []Cluster
}
```

<a name="CmdGetOptions"></a>
## type [CmdGetOptions](<https://github.com/ice-bergtech/kr8/blob/main/cmd/get.go#L35-L42>)



```go
type CmdGetOptions struct {
    ClusterParams string
    NoTable       bool
    FieldName     string
    Cluster       string
    Component     string
    ParamField    string
}
```

<a name="CmdJsonnetOptions"></a>
## type [CmdJsonnetOptions](<https://github.com/ice-bergtech/kr8/blob/main/cmd/jsonnet.go#L277-L283>)



```go
type CmdJsonnetOptions struct {
    Prune         bool
    Cluster       string
    ClusterParams string
    Component     string
    Format        string
}
```

<a name="ComponentJsonnet"></a>
## type [ComponentJsonnet](<https://github.com/ice-bergtech/kr8/blob/main/cmd/types.go#L77-L86>)



```go
type ComponentJsonnet struct {
    // The default namespace to deploy the component to (optional)
    Namespace string `json:"namespace"`
    // A unique name for the component (optional)
    ReleaseName string `json:"release_name"`
    // Component version number (optional)
    Version string `json:"version"`
    // Component-specific configuration for kr8 (required)
    Kr8Spec ComponentSpec `json:"kr8_spec"`
}
```

<a name="ComponentSpec"></a>
## type [ComponentSpec](<https://github.com/ice-bergtech/kr8/blob/main/cmd/types.go#L89-L102>)

kr8\_spec object in cluster config

```go
type ComponentSpec struct {
    // If true, includes the parameters of the current cluster when generating this component
    Kr8_allparams bool `json:"enable_kr8_allparams"`
    // If true, includes the parameters of all other clusters when generating this component
    Kr8_allclusters bool `json:"enable_kr8_allclusters"`
    // If false, all non-generated files present in the output directory will be removed
    DisableOutputDirClean bool `json:"disable_output_clean"`
    // A list of filenames to include as jsonnet vm external vars
    ExtFiles ExtFileVar `json:"extfiles"`
    // Additional jsonnet libs to the jsonnet vm, path component scoped
    JPaths []string `json:"jpaths"`
    // A list of filenames to include and output as files
    Includes []interface{} `json:"includes"`
}
```

<a name="CreateComponentSpec"></a>
### func [CreateComponentSpec](<https://github.com/ice-bergtech/kr8/blob/main/cmd/types.go#L104>)

```go
func CreateComponentSpec(spec gjson.Result) (ComponentSpec, error)
```



<a name="ExtFileVar"></a>
## type [ExtFileVar](<https://github.com/ice-bergtech/kr8/blob/main/cmd/types.go#L153>)

file to load as a string into the jsonnet vm name to reference the variable in jsonnet code through std.extvar\(\) value of the variable, loaded from a file or provided directly

```go
type ExtFileVar map[string]string
```

<a name="IncludeFileEntryStruct"></a>
## type [IncludeFileEntryStruct](<https://github.com/ice-bergtech/kr8/blob/main/cmd/types.go#L161-L171>)



```go
type IncludeFileEntryStruct struct {
    // an input file to process
    // accepted filetypes: .jsonnet .yml .yaml .tmpl .tpl
    File string `json:"file"`
    // handle alternate output directory for file
    DestDir string `json:"dest_dir,omitempty"`
    // override destination file name
    DestName string `json:"dest_name,omitempty"`
    // override destination file extension
    DestExt string `json:"dest_ext,omitempty"`
}
```

<a name="IncludeFileSpec"></a>
## type [IncludeFileSpec](<https://github.com/ice-bergtech/kr8/blob/main/cmd/types.go#L156-L159>)

A struct describing an included file

```go
type IncludeFileSpec interface {
    IncludeFileEntryStruct
    // contains filtered or unexported methods
}
```

<a name="PathFilterOptions"></a>
## type [PathFilterOptions](<https://github.com/ice-bergtech/kr8/blob/main/cmd/format.go#L44-L47>)



```go
type PathFilterOptions struct {
    Includes string
    Excludes string
}
```

<a name="VMConfig"></a>
## type [VMConfig](<https://github.com/ice-bergtech/kr8/blob/main/cmd/jsonnet.go#L43-L47>)



```go
type VMConfig struct {
    // VMConfig is a configuration for the Jsonnet VM
    Jpaths  []string `json:"jpath" yaml:"jpath"`
    ExtVars []string `json:"ext_str_file" yaml:"ext_str_files"`
}
```