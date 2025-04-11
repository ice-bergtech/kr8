# kr8p\_init

```go
import "github.com/ice-bergtech/kr8p/pkg/kr8p_init"
```

## Index

- [func GenerateClusterJsonnet\(cSpec types.Kr8pClusterSpec, dstDir string\) error](<#GenerateClusterJsonnet>)
- [func GenerateComponentJsonnet\(componentOptions Kr8pInitOptions, dstDir string\) error](<#GenerateComponentJsonnet>)
- [func GenerateLib\(fetch bool, dstDir string\) error](<#GenerateLib>)
- [func GenerateReadme\(dstDir string, cmdOptions Kr8pInitOptions, clusterSpec types.Kr8pClusterSpec\) error](<#GenerateReadme>)
- [type Kr8pInitOptions](<#Kr8pInitOptions>)


<a name="GenerateClusterJsonnet"></a>
## func [GenerateClusterJsonnet](<https://github.com:icebergtech/kr8p/blob/main/pkg/kr8p_init/init.go#L28>)

```go
func GenerateClusterJsonnet(cSpec types.Kr8pClusterSpec, dstDir string) error
```

Generate a cluster.jsonnet file based on the provided Kr8ClusterSpec and store it in the specified directory.

<a name="GenerateComponentJsonnet"></a>
## func [GenerateComponentJsonnet](<https://github.com:icebergtech/kr8p/blob/main/pkg/kr8p_init/init.go#L49>)

```go
func GenerateComponentJsonnet(componentOptions Kr8pInitOptions, dstDir string) error
```

Generate default component kr8\_spec values and store in params.jsonnet. Based on the type:

jsonnet: create a component.jsonnet file and reference it from the params.jsonnet file

yml: leave a note in the params.jsonnet file about where and how the yml files can be referenced

chart: generate a simple taskfile that handles vendoring the chart data

<a name="GenerateLib"></a>
## func [GenerateLib](<https://github.com:icebergtech/kr8p/blob/main/pkg/kr8p_init/init.go#L104>)

```go
func GenerateLib(fetch bool, dstDir string) error
```

Downloads a starter kr8p jsonnet lib from github. If fetch is true, it will download the repo in the /lib directory. If false, it will print the git commands to run. Repo: https://github.com/ice-bergtech/kr8-libsonnet . return util.FetchRepoUrl\("https://github.com/ice-bergtech/kr8-libsonnet", dstDir\+"/kr8\-lib", \!fetch\).

<a name="GenerateReadme"></a>
## func [GenerateReadme](<https://github.com:icebergtech/kr8p/blob/main/pkg/kr8p_init/init.go#L113>)

```go
func GenerateReadme(dstDir string, cmdOptions Kr8pInitOptions, clusterSpec types.Kr8pClusterSpec) error
```

Generates a starter readme for the repo, and writes it to the destination directory.

<a name="Kr8pInitOptions"></a>
## type [Kr8pInitOptions](<https://github.com:icebergtech/kr8p/blob/main/pkg/kr8p_init/init.go#L12-L25>)

Kr8pInitOptions defines the options used by the init subcommands.

```go
type Kr8pInitOptions struct {
    // URL to fetch the skeleton directory from
    InitUrl string
    // Name of the cluster to initialize
    ClusterName string
    // Name of the component to initialize
    ComponentName string
    // Type of component to initialize (e.g. jsonnet, yml, chart, compose)
    ComponentType string
    // Whether to run in interactive mode or not
    Interactive bool
    // Whether to fetch remote resources or not
    Fetch bool
}
```