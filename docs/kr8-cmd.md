# cmd

```go
import "github.com/ice-bergtech/kr8/cmd"
```

## Index

- [Variables](<#variables>)
- [func Execute\(version string\)](<#Execute>)
- [func GenerateClusterJsonnet\(cSpec Kr8ClusterSpec, dstDir string\) error](<#GenerateClusterJsonnet>)
- [func GenerateComponentJsonnet\(componentOptions cmdInitOptions, dstDir string\) error](<#GenerateComponentJsonnet>)
- [func GenerateLib\(fetch bool, dstDir string\)](<#GenerateLib>)
- [func GenerateReadme\(dstDir string, cmdOptions cmdInitOptions, clusterSpec Kr8ClusterSpec\)](<#GenerateReadme>)
- [func GetDefaultFormatOptions\(\) formatter.Options](<#GetDefaultFormatOptions>)
- [func JsonnetVM\(vmconfig VMConfig\) \(\*jsonnet.VM, error\)](<#JsonnetVM>)
- [func Pretty\(input string, colorOutput bool\) string](<#Pretty>)
- [func RegisterNativeFuncs\(vm \*jsonnet.VM\)](<#RegisterNativeFuncs>)
- [type CmdGetOptions](<#CmdGetOptions>)
- [type CmdJsonnetOptions](<#CmdJsonnetOptions>)
- [type ExtFileVar](<#ExtFileVar>)
- [type KomposeConvertOptions](<#KomposeConvertOptions>)
- [type Kr8ClusterComponentRef](<#Kr8ClusterComponentRef>)
- [type Kr8ClusterJsonnet](<#Kr8ClusterJsonnet>)
- [type Kr8ClusterSpec](<#Kr8ClusterSpec>)
  - [func CreateClusterSpec\(clusterName string, spec gjson.Result, baseDir string, genDirOverride string\) \(Kr8ClusterSpec, error\)](<#CreateClusterSpec>)
- [type Kr8ComponentJsonnet](<#Kr8ComponentJsonnet>)
- [type Kr8ComponentSpec](<#Kr8ComponentSpec>)
  - [func CreateComponentSpec\(spec gjson.Result\) \(Kr8ComponentSpec, error\)](<#CreateComponentSpec>)
- [type Kr8ComponentSpecIncludeFile](<#Kr8ComponentSpecIncludeFile>)
- [type Kr8ComponentSpecIncludeObject](<#Kr8ComponentSpecIncludeObject>)
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
## func [Execute](<https://github.com/ice-bergtech/kr8/blob/main/cmd/root.go#L31>)

```go
func Execute(version string)
```

Execute adds all child commands to the root command sets flags appropriately. This is called by main.main\(\). It only needs to happen once to the rootCmd.

<a name="GenerateClusterJsonnet"></a>
## func [GenerateClusterJsonnet](<https://github.com/ice-bergtech/kr8/blob/main/cmd/init.go#L198>)

```go
func GenerateClusterJsonnet(cSpec Kr8ClusterSpec, dstDir string) error
```

Generate a cluster.jsonnet file based on the provided Kr8ClusterSpec and store it in the specified directory.

<a name="GenerateComponentJsonnet"></a>
## func [GenerateComponentJsonnet](<https://github.com/ice-bergtech/kr8/blob/main/cmd/init.go#L213>)

```go
func GenerateComponentJsonnet(componentOptions cmdInitOptions, dstDir string) error
```

Generate default component kr8\_spec values and store in params.jsonnet Based on the type: jsonnet: create a component.jsonnet file and reference it from the params.jsonnet file yml: leave a note in the params.jsonnet file about where and how the yml files can be referenced chart: generate a simple taskfile that handles vendoring the chart data

<a name="GenerateLib"></a>
## func [GenerateLib](<https://github.com/ice-bergtech/kr8/blob/main/cmd/init.go#L278>)

```go
func GenerateLib(fetch bool, dstDir string)
```



<a name="GenerateReadme"></a>
## func [GenerateReadme](<https://github.com/ice-bergtech/kr8/blob/main/cmd/init.go#L283>)

```go
func GenerateReadme(dstDir string, cmdOptions cmdInitOptions, clusterSpec Kr8ClusterSpec)
```



<a name="GetDefaultFormatOptions"></a>
## func [GetDefaultFormatOptions](<https://github.com/ice-bergtech/kr8/blob/main/cmd/format.go#L21>)

```go
func GetDefaultFormatOptions() formatter.Options
```

Configures the default options for the jsonnet formatter

<a name="JsonnetVM"></a>
## func [JsonnetVM](<https://github.com/ice-bergtech/kr8/blob/main/cmd/jsonnet.go#L52>)

```go
func JsonnetVM(vmconfig VMConfig) (*jsonnet.VM, error)
```



<a name="Pretty"></a>
## func [Pretty](<https://github.com/ice-bergtech/kr8/blob/main/cmd/json.go#L10>)

```go
func Pretty(input string, colorOutput bool) string
```

Pretty formats the input jsonnet string with indentation and optional color output.

<a name="RegisterNativeFuncs"></a>
## func [RegisterNativeFuncs](<https://github.com/ice-bergtech/kr8/blob/main/cmd/jsonnet.go#L147>)

```go
func RegisterNativeFuncs(vm *jsonnet.VM)
```

Registers additional native functions in the jsonnet VM These functions are used to extend the functionality of jsonnet Adds on to functions part of the jsonnet stdlib: https://jsonnet.org/ref/stdlib.html

<a name="CmdGetOptions"></a>
## type [CmdGetOptions](<https://github.com/ice-bergtech/kr8/blob/main/cmd/get.go#L37-L50>)

Holds the options for the get command.

```go
type CmdGetOptions struct {
    // ClusterParams provides a way to provide cluster params as a single file. This can be combined with --cluster to override the cluster.
    ClusterParams string
    // If true, just prints result instead of placing in table
    NoTable bool
    // Field to display from the resource
    FieldName string
    // Cluster to get resources from
    Cluster string
    // Component to get resources from
    Component string
    // Param to display from the resource
    ParamField string
}
```

<a name="CmdJsonnetOptions"></a>
## type [CmdJsonnetOptions](<https://github.com/ice-bergtech/kr8/blob/main/cmd/jsonnet.go#L316-L322>)



```go
type CmdJsonnetOptions struct {
    Prune         bool
    Cluster       string
    ClusterParams string
    Component     string
    Format        string
}
```

<a name="ExtFileVar"></a>
## type [ExtFileVar](<https://github.com/ice-bergtech/kr8/blob/main/cmd/types.go#L160>)

Map of external files to load into jsonnet vm as external variables Keys are the variable names, values are the paths to the files to load as strings into the jsonnet vm To reference the variable in jsonnet code, use std.extvar\("variable\_name"\)

```go
type ExtFileVar map[string]string
```

<a name="KomposeConvertOptions"></a>
## type [KomposeConvertOptions](<https://github.com/ice-bergtech/kr8/blob/main/cmd/types.go#L185-L262>)

A struct describing a compose file that will be processed by kompose to produce kubernetes manifests Based on https://github.com/kubernetes/kompose/blob/main/cmd/convert.go

```go
type KomposeConvertOptions struct {
    // Kubernetes: Create a Helm chart for converted objects
    CreateChart bool
    // Kubernetes: Set the output controller ("deployment"|"daemonSet"|"replicationController")
    Controller       string
    IsDaemonSetFlag  bool
    IsDeploymentFlag bool

    // Print converted objects to stdout
    ToStdout bool

    // Openshift: Specify source repository for buildconfig (default remote origin)
    BuildRepo string
    // Openshift: Specify repository branch to use for buildconfig (default master)
    BuildBranch string
    // Set the type of build ("local"|"build-config"(OpenShift only)|"none")
    Build string

    Profiles  []string
    PushImage bool
    // Specify registry for pushing image, which will override registry from image name
    PushImageRegistry string
    // Generate resource files into YAML format
    GenerateYaml bool
    // Generate resource files into JSON format
    GenerateJSON  bool
    StoreManifest bool
    // Use Empty Volumes. Do not generate PVCs
    EmptyVols bool
    // Volumes to be generated ("persistentVolumeClaim"|"emptyDir"|"hostPath" | "configMap")
    Volumes string
    // Specify the size of pvc storage requests in the generated resource spec
    PVCRequestSize string
    // Use an insecure Docker repository for OpenShift ImageStream
    InsecureRepository bool
    // Specify the number of replicas in the generated resource spec
    Replicas   int
    InputFiles []string
    // Specify a file name or directory to save objects to (if path does not exist, a file will be created)
    OutFile  string
    Provider string
    // Specify the namespace of the generated resources`)
    Namespace string

    IsReplicationControllerFlag bool
    IsReplicaSetFlag            bool
    IsDeploymentConfigFlag      bool
    IsNamespaceFlag             bool

    BuildCommand string
    PushCommand  string

    Server string

    // Spaces length to indent generated yaml files
    YAMLIndent int

    // Add kompose annotations to generated resource
    WithKomposeAnnotation bool

    // Create multiple containers grouped by 'kompose.service.group' label
    MultipleContainerMode bool
    // Group multiple service to create single workload by `label`(`kompose.service.group`) or `volume`(shared volumes)
    ServiceGroupMode string
    // Using with --service-group-mode=volume to specific a final service name for the group
    ServiceGroupName string

    // Convert docker-compose secrets into files instead of symlinked directories
    SecretsAsFiles bool
    // Specify whether to generate network policies or not
    GenerateNetworkPolicies bool
}
```

<a name="Kr8ClusterComponentRef"></a>
## type [Kr8ClusterComponentRef](<https://github.com/ice-bergtech/kr8/blob/main/cmd/types.go#L30-L33>)

A reference to a component folder that contains a params.jsonnet file. This is used in the cluster jsonnet file to reference components.

```go
type Kr8ClusterComponentRef struct {
    // The path to a component folder that contains a params.jsonnet file
    Path string `json:"path"`
}
```

<a name="Kr8ClusterJsonnet"></a>
## type [Kr8ClusterJsonnet](<https://github.com/ice-bergtech/kr8/blob/main/cmd/types.go#L20-L27>)

The specification for a clusters.jsonnet file This file contains configuration for clusters, including

```go
type Kr8ClusterJsonnet struct {
    // kr8 configuration for how to process the cluster
    ClusterSpec Kr8ClusterSpec `json:"_kr8_spec"`
    // Cluster Level configuration that components can reference
    Cluster kr8Cluster `json:"_cluster"`
    // Distictly named components.
    Components map[string]Kr8ClusterComponentRef `json:"_components"`
}
```

<a name="Kr8ClusterSpec"></a>
## type [Kr8ClusterSpec](<https://github.com/ice-bergtech/kr8/blob/main/cmd/types.go#L37-L50>)

The specification for how to process a cluster. This is used in the cluster jsonnet file to configure how kr8 should process the cluster.

```go
type Kr8ClusterSpec struct {
    // The name of the cluster
    Name string `json:"-"`
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
}
```

<a name="CreateClusterSpec"></a>
### func [CreateClusterSpec](<https://github.com/ice-bergtech/kr8/blob/main/cmd/types.go#L54>)

```go
func CreateClusterSpec(clusterName string, spec gjson.Result, baseDir string, genDirOverride string) (Kr8ClusterSpec, error)
```

This function creates a cluster spec from the given cluster name, spec, base directory, and generate directory override. If genDirOverride is empty, the value of generate\_dir from the spec is used.

<a name="Kr8ComponentJsonnet"></a>
## type [Kr8ComponentJsonnet](<https://github.com/ice-bergtech/kr8/blob/main/cmd/types.go#L82-L91>)

The specification for component's params.jsonnet file It contains all the configuration and variables used to generate component resources This configuration is often modified from the cluster config to add cluster\-specific configuration

```go
type Kr8ComponentJsonnet struct {
    // Component-specific configuration for how kr8 should process the component (required)
    Kr8Spec Kr8ComponentSpec `json:"kr8_spec"`
    // The default namespace to deploy the component to
    Namespace string `json:"namespace"`
    // A unique name for the component
    ReleaseName string `json:"release_name"`
    // Component version string (optional)
    Version string `json:"version"`
}
```

<a name="Kr8ComponentSpec"></a>
## type [Kr8ComponentSpec](<https://github.com/ice-bergtech/kr8/blob/main/cmd/types.go#L95-L108>)

The kr8\_spec object in a cluster config file This configures how kr8 processes the component

```go
type Kr8ComponentSpec struct {
    // If true, includes the parameters of the current cluster when generating this component
    Kr8_allparams bool `json:"enable_kr8_allparams"`
    // If true, includes the parameters of all other clusters when generating this component
    Kr8_allclusters bool `json:"enable_kr8_allclusters"`
    // If false, all non-generated files present in the output directory will be removed
    DisableOutputDirClean bool `json:"disable_output_clean"`
    // A list of filenames to include as jsonnet vm external vars
    ExtFiles ExtFileVar `json:"extfiles"`
    // Additional jsonnet libs to the jsonnet vm, component-path scoped
    JPaths []string `json:"jpaths"`
    // A list of filenames to include and output as files
    Includes []interface{} `json:"includes"`
}
```

<a name="CreateComponentSpec"></a>
### func [CreateComponentSpec](<https://github.com/ice-bergtech/kr8/blob/main/cmd/types.go#L111>)

```go
func CreateComponentSpec(spec gjson.Result) (Kr8ComponentSpec, error)
```

Extracts a component spec from a jsonnet object.

<a name="Kr8ComponentSpecIncludeFile"></a>
## type [Kr8ComponentSpecIncludeFile](<https://github.com/ice-bergtech/kr8/blob/main/cmd/types.go#L163-L166>)

A struct describing an included file that will be processed to produce a file

```go
type Kr8ComponentSpecIncludeFile interface {
    Kr8ComponentSpecIncludeObject
    // contains filtered or unexported methods
}
```

<a name="Kr8ComponentSpecIncludeObject"></a>
## type [Kr8ComponentSpecIncludeObject](<https://github.com/ice-bergtech/kr8/blob/main/cmd/types.go#L171-L181>)

An includes object which configures how kr8 includes an object It allows configuring the included file's destination directory and file name The input file will be processed differently depending on the filetype

```go
type Kr8ComponentSpecIncludeObject struct {
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

<a name="VMConfig"></a>
## type [VMConfig](<https://github.com/ice-bergtech/kr8/blob/main/cmd/jsonnet.go#L45-L50>)

VMConfig describes configuration to initialize Jsonnet VM with

```go
type VMConfig struct {
    // Jpaths is a list of paths to search for Jsonnet libraries (libsonnet files)
    Jpaths []string `json:"jpath" yaml:"jpath"`
    // ExtVars is a list of external variables to pass to Jsonnet VM
    ExtVars []string `json:"ext_str_file" yaml:"ext_str_files"`
}
```