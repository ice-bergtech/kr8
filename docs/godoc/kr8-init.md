# kr8\_init

```go
import "github.com/ice-bergtech/kr8/pkg/kr8_init"
```

Package kr8\_init contains logic for initializing a kr8\+ starter repo. It is able to generate starter configs for components, clusters, and full repos.

## Index

- [func GenerateChartJsonnet\(compJson kr8\_types.Kr8ComponentJsonnet, componentOptions Kr8InitOptions, folderDir string\) error](<#GenerateChartJsonnet>)
- [func GenerateChartTaskfile\(comp kr8\_types.Kr8ComponentJsonnet, componentOptions Kr8InitOptions, folderDir string\) error](<#GenerateChartTaskfile>)
- [func GenerateClusterJsonnet\(cSpec kr8\_types.Kr8ClusterSpec, dstDir string\) error](<#GenerateClusterJsonnet>)
- [func GenerateComponentJsonnet\(componentOptions Kr8InitOptions, dstDir string\) error](<#GenerateComponentJsonnet>)
- [func GenerateLib\(fetch bool, dstDir string\) error](<#GenerateLib>)
- [func GenerateReadme\(dstDir string, cmdOptions Kr8InitOptions, clusterSpec kr8\_types.Kr8ClusterSpec\) error](<#GenerateReadme>)
- [func InitComponentChart\(dstDir string, componentOptions Kr8InitOptions, compJson kr8\_types.Kr8ComponentJsonnet\) error](<#InitComponentChart>)
- [func InitComponentJsonnet\(compJson kr8\_types.Kr8ComponentJsonnet, dstDir string, componentOptions Kr8InitOptions\) error](<#InitComponentJsonnet>)
- [func InitComponentTemplate\(compJson kr8\_types.Kr8ComponentJsonnet, dstDir string, componentOptions Kr8InitOptions\) error](<#InitComponentTemplate>)
- [func InitComponentYaml\(compJson kr8\_types.Kr8ComponentJsonnet, dstDir string, componentOptions Kr8InitOptions\) error](<#InitComponentYaml>)
- [type Kr8InitOptions](<#Kr8InitOptions>)


<a name="GenerateChartJsonnet"></a>
## func [GenerateChartJsonnet](<https://github.com:icebergtech/kr8/blob/main/pkg/kr8_init/component.go#L126-L130>)

```go
func GenerateChartJsonnet(compJson kr8_types.Kr8ComponentJsonnet, componentOptions Kr8InitOptions, folderDir string) error
```

Generates a jsonnet files that references a local helm chart.

<a name="GenerateChartTaskfile"></a>
## func [GenerateChartTaskfile](<https://github.com:icebergtech/kr8/blob/main/pkg/kr8_init/component.go#L156-L160>)

```go
func GenerateChartTaskfile(comp kr8_types.Kr8ComponentJsonnet, componentOptions Kr8InitOptions, folderDir string) error
```

Generates a go\-task taskfile that's setup to download a helm chart into a local \`vendor\` directory.

<a name="GenerateClusterJsonnet"></a>
## func [GenerateClusterJsonnet](<https://github.com:icebergtech/kr8/blob/main/pkg/kr8_init/cluster.go#L9>)

```go
func GenerateClusterJsonnet(cSpec kr8_types.Kr8ClusterSpec, dstDir string) error
```

Generate a cluster.jsonnet file based on the provided Kr8ClusterSpec and store it in the specified directory.

<a name="GenerateComponentJsonnet"></a>
## func [GenerateComponentJsonnet](<https://github.com:icebergtech/kr8/blob/main/pkg/kr8_init/component.go#L20>)

```go
func GenerateComponentJsonnet(componentOptions Kr8InitOptions, dstDir string) error
```

Generate default component kr8\_spec values and store in params.jsonnet. Based on the type:

jsonnet: create a component.jsonnet file and reference it from the params.jsonnet file

yml: leave a note in the params.jsonnet file about where and how the yml files can be referenced

chart: generate a simple taskfile that handles vendoring the chart data

<a name="GenerateLib"></a>
## func [GenerateLib](<https://github.com:icebergtech/kr8/blob/main/pkg/kr8_init/libs.go#L14>)

```go
func GenerateLib(fetch bool, dstDir string) error
```

Downloads a starter kr8 jsonnet lib from github. If fetch is true, downloads the repo in the /lib directory. If false, prints the git commands to run. Repo: https://github.com/ice-bergtech/kr8-libsonnet . return util.FetchRepoUrl\("https://github.com/ice-bergtech/kr8-libsonnet", dstDir\+"/kr8\-lib", \!fetch\).

<a name="GenerateReadme"></a>
## func [GenerateReadme](<https://github.com:icebergtech/kr8/blob/main/pkg/kr8_init/repo.go#L11>)

```go
func GenerateReadme(dstDir string, cmdOptions Kr8InitOptions, clusterSpec kr8_types.Kr8ClusterSpec) error
```

Generates a starter readme for the repo, and writes it to the destination directory.

<a name="InitComponentChart"></a>
## func [InitComponentChart](<https://github.com:icebergtech/kr8/blob/main/pkg/kr8_init/component.go#L52>)

```go
func InitComponentChart(dstDir string, componentOptions Kr8InitOptions, compJson kr8_types.Kr8ComponentJsonnet) error
```

Initializes the basic parts of a helm chart component.

<a name="InitComponentJsonnet"></a>
## func [InitComponentJsonnet](<https://github.com:icebergtech/kr8/blob/main/pkg/kr8_init/component.go#L111-L115>)

```go
func InitComponentJsonnet(compJson kr8_types.Kr8ComponentJsonnet, dstDir string, componentOptions Kr8InitOptions) error
```

Initializes the basic parts of a jsonnet\-based component.

<a name="InitComponentTemplate"></a>
## func [InitComponentTemplate](<https://github.com:icebergtech/kr8/blob/main/pkg/kr8_init/component.go#L77-L81>)

```go
func InitComponentTemplate(compJson kr8_types.Kr8ComponentJsonnet, dstDir string, componentOptions Kr8InitOptions) error
```

Initializes the based parts of a template\-based component.

<a name="InitComponentYaml"></a>
## func [InitComponentYaml](<https://github.com:icebergtech/kr8/blob/main/pkg/kr8_init/component.go#L96>)

```go
func InitComponentYaml(compJson kr8_types.Kr8ComponentJsonnet, dstDir string, componentOptions Kr8InitOptions) error
```

Initializes the basic parts of a yaml\-based component.

<a name="Kr8InitOptions"></a>
## type [Kr8InitOptions](<https://github.com:icebergtech/kr8/blob/main/pkg/kr8_init/kr8_init.go#L6-L19>)

Kr8InitOptions defines the options used by the init subcommands.

```go
type Kr8InitOptions struct {
    // URL to fetch the skeleton directory from
    InitUrl string
    // Name of the cluster to initialize
    ClusterName string
    // Name of the component to initialize
    ComponentName string
    // Type of component to initialize (e.g. jsonnet, yml, chart, compose)
    ComponentType string
    // Determines whether to run in interactive mode
    Interactive bool
    // Determines whether to fetch remote resources
    Fetch bool
}
```