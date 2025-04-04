# types

```go
import "github.com/ice-bergtech/kr8/pkg/types"
```

## Index

- [type CmdJsonnetOptions](<#CmdJsonnetOptions>)
- [type ExtFileVar](<#ExtFileVar>)
- [type KomposeConvertOptions](<#KomposeConvertOptions>)
  - [func Create\(inputFiles \[\]string, outDir string, cmp Kr8ComponentJsonnet\) \*KomposeConvertOptions](<#Create>)
  - [func \(k KomposeConvertOptions\) Convert\(\) \(interface\{\}, error\)](<#KomposeConvertOptions.Convert>)
  - [func \(k KomposeConvertOptions\) GenKomposePkgOpts\(\) \*kobject.ConvertOptions](<#KomposeConvertOptions.GenKomposePkgOpts>)
  - [func \(k KomposeConvertOptions\) Validate\(\) error](<#KomposeConvertOptions.Validate>)
- [type Kr8Cluster](<#Kr8Cluster>)
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


<a name="CmdJsonnetOptions"></a>
## type [CmdJsonnetOptions](<https://github.com/ice-bergtech/kr8/blob/main/pkg/types/types.go#L183-L190>)



```go
type CmdJsonnetOptions struct {
    Prune         bool
    Cluster       string
    ClusterParams string
    Component     string
    Format        string
    Color         bool
}
```

<a name="ExtFileVar"></a>
## type [ExtFileVar](<https://github.com/ice-bergtech/kr8/blob/main/pkg/types/types.go#L160>)

Map of external files to load into jsonnet vm as external variables Keys are the variable names, values are the paths to the files to load as strings into the jsonnet vm To reference the variable in jsonnet code, use std.extvar\("variable\_name"\)

```go
type ExtFileVar map[string]string
```

<a name="KomposeConvertOptions"></a>
## type [KomposeConvertOptions](<https://github.com/ice-bergtech/kr8/blob/main/pkg/types/kompose.go#L13-L87>)

A struct describing a compose file that will be processed by kompose to produce kubernetes manifests.

Based on https://github.com/kubernetes/kompose/blob/main/cmd/convert.go

```go
type KomposeConvertOptions struct {
    // Kubernetes: Set the output controller ("deployment"|"daemonSet"|"replicationController")
    Controller string

    // The kubecfg (?) profile to use, can use multiple profiles
    Profiles []string

    // List of compose file filenames.
    // Filenames should be in the format `[docker-]compose.ym[a]l`
    InputFiles []string
    // Specify a file name or directory to save objects to (if path does not exist, a file will be created)
    OutFile string
    // Generate a Helm chart for converted objects
    CreateChart bool
    // Add kompose annotations to generated resource
    WithKomposeAnnotation bool
    // Generate resource files into YAML format
    GenerateYaml bool
    // Spaces length to indent generated yaml files
    GenerateYAMLIndent int
    // Generate resource files into JSON format
    GenerateJSON bool
    // Print converted objects to stdout
    GenerateToStdout bool

    // Set the type of build ("local"|"build-config"(OpenShift only)|"none")
    Build string

    // Specify the namespace of the generated resources`)
    Namespace string
    // Specify the number of replicas in the generated resource spec
    Replicas int
    // Convert docker-compose secrets into files instead of symlinked directories
    SecretsAsFiles bool
    // Use Empty Volumes. Do not generate PVCs
    EmptyVols bool
    // Volumes to be generated ("persistentVolumeClaim"|"emptyDir"|"hostPath" | "configMap")
    Volumes string
    // Specify the size of pvc storage requests in the generated resource spec
    PVCRequestSize string
    // Specify whether to generate network policies or not
    NetworkPolicies bool

    // Create multiple containers grouped by 'kompose.service.group' label
    MultipleContainerMode bool
    // Group multiple service to create single workload by `label`(`kompose.service.group`) or `volume`(shared volumes)
    ServiceGroupMode string
    // Using with --service-group-mode=volume to specific a final service name for the group
    ServiceGroupName string

    // ??
    Provider string
    // ??
    StoreManifest bool
    // ??
    Server string

    // OpenShift: ??
    OSCreateDeploymentConfig bool
    // Openshift: Specify source repository for buildconfig (default remote origin)
    OSBuildRepo string
    // Openshift: Use an insecure Docker repository for OpenShift ImageStream
    OSInsecureRepository bool
    // Openshift: Specify repository branch to use for buildconfig (default master)
    OSBuildBranch string

    // Whether to push built docker image to remote registry.
    ImagePush bool
    // Command used to build to image.  Used with PushCommand
    ImageBuildCommand string
    // Command used to push image
    ImagePushCommand string
    // Specify registry for pushing image, which will override registry from image name
    ImagePushRegistry string
}
```

<a name="Create"></a>
### func [Create](<https://github.com/ice-bergtech/kr8/blob/main/pkg/types/kompose.go#L90>)

```go
func Create(inputFiles []string, outDir string, cmp Kr8ComponentJsonnet) *KomposeConvertOptions
```

Initialie Kompose options with sensible defaults

<a name="KomposeConvertOptions.Convert"></a>
### func \(KomposeConvertOptions\) [Convert](<https://github.com/ice-bergtech/kr8/blob/main/pkg/types/kompose.go#L203>)

```go
func (k KomposeConvertOptions) Convert() (interface{}, error)
```

Converts a Docker Compose file described by k into a set of kubernetes manifests.

<a name="KomposeConvertOptions.GenKomposePkgOpts"></a>
### func \(KomposeConvertOptions\) [GenKomposePkgOpts](<https://github.com/ice-bergtech/kr8/blob/main/pkg/types/kompose.go#L122>)

```go
func (k KomposeConvertOptions) GenKomposePkgOpts() *kobject.ConvertOptions
```

Generates a ConvertOptions struct that kompose expects from our commented KomposeConvertOptions

References:

https://pkg.go.dev/github.com/kubernetes/kompose@v1.35.0/pkg/kobject#ConvertOptions

https://github.com/kubernetes/kompose/blob/v1.35.0/pkg/app/app.go#L166

<a name="KomposeConvertOptions.Validate"></a>
### func \(KomposeConvertOptions\) [Validate](<https://github.com/ice-bergtech/kr8/blob/main/pkg/types/kompose.go#L191>)

```go
func (k KomposeConvertOptions) Validate() error
```

Validates a set of options for converting a Kubernetes manifest to a Docker Compose file.

<a name="Kr8Cluster"></a>
## type [Kr8Cluster](<https://github.com/ice-bergtech/kr8/blob/main/pkg/types/types.go#L13-L16>)

An object that stores variables that can be referenced by components

```go
type Kr8Cluster struct {
    Name string `json:"name"`
    Path string `json:"-"`
}
```

<a name="Kr8ClusterComponentRef"></a>
## type [Kr8ClusterComponentRef](<https://github.com/ice-bergtech/kr8/blob/main/pkg/types/types.go#L30-L33>)

A reference to a component folder that contains a params.jsonnet file. This is used in the cluster jsonnet file to reference components.

```go
type Kr8ClusterComponentRef struct {
    // The path to a component folder that contains a params.jsonnet file
    Path string `json:"path"`
}
```

<a name="Kr8ClusterJsonnet"></a>
## type [Kr8ClusterJsonnet](<https://github.com/ice-bergtech/kr8/blob/main/pkg/types/types.go#L20-L27>)

The specification for a clusters.jsonnet file This file contains configuration for clusters, including

```go
type Kr8ClusterJsonnet struct {
    // kr8 configuration for how to process the cluster
    ClusterSpec Kr8ClusterSpec `json:"_kr8_spec"`
    // Cluster Level configuration that components can reference
    Cluster Kr8Cluster `json:"_cluster"`
    // Distictly named components.
    Components map[string]Kr8ClusterComponentRef `json:"_components"`
}
```

<a name="Kr8ClusterSpec"></a>
## type [Kr8ClusterSpec](<https://github.com/ice-bergtech/kr8/blob/main/pkg/types/types.go#L37-L50>)

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
### func [CreateClusterSpec](<https://github.com/ice-bergtech/kr8/blob/main/pkg/types/types.go#L54>)

```go
func CreateClusterSpec(clusterName string, spec gjson.Result, baseDir string, genDirOverride string) (Kr8ClusterSpec, error)
```

This function creates a cluster spec from the given cluster name, spec, base directory, and generate directory override. If genDirOverride is empty, the value of generate\_dir from the spec is used.

<a name="Kr8ComponentJsonnet"></a>
## type [Kr8ComponentJsonnet](<https://github.com/ice-bergtech/kr8/blob/main/pkg/types/types.go#L82-L91>)

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
## type [Kr8ComponentSpec](<https://github.com/ice-bergtech/kr8/blob/main/pkg/types/types.go#L95-L108>)

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
### func [CreateComponentSpec](<https://github.com/ice-bergtech/kr8/blob/main/pkg/types/types.go#L111>)

```go
func CreateComponentSpec(spec gjson.Result) (Kr8ComponentSpec, error)
```

Extracts a component spec from a jsonnet object.

<a name="Kr8ComponentSpecIncludeFile"></a>
## type [Kr8ComponentSpecIncludeFile](<https://github.com/ice-bergtech/kr8/blob/main/pkg/types/types.go#L163-L166>)

A struct describing an included file that will be processed to produce a file

```go
type Kr8ComponentSpecIncludeFile interface {
    Kr8ComponentSpecIncludeObject
    // contains filtered or unexported methods
}
```

<a name="Kr8ComponentSpecIncludeObject"></a>
## type [Kr8ComponentSpecIncludeObject](<https://github.com/ice-bergtech/kr8/blob/main/pkg/types/types.go#L171-L181>)

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
## type [VMConfig](<https://github.com/ice-bergtech/kr8/blob/main/pkg/types/types.go#L193-L200>)

VMConfig describes configuration to initialize Jsonnet VM with

```go
type VMConfig struct {
    // Jpaths is a list of paths to search for Jsonnet libraries (libsonnet files)
    Jpaths []string `json:"jpath" yaml:"jpath"`
    // ExtVars is a list of external variables to pass to Jsonnet VM
    ExtVars []string `json:"ext_str_file" yaml:"ext_str_files"`
    // base directory for the project
    BaseDir string `json:"base_dir" yaml:"base_dir"`
}
```