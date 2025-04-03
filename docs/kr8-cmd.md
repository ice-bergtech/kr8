# cmd

```go
import "github.com/ice-bergtech/kr8/cmd"
```

## Index

- [Variables](<#variables>)
- [func Execute\(version string\)](<#Execute>)
- [func GenerateClusterJsonnet\(cSpec types.Kr8ClusterSpec, dstDir string\) error](<#GenerateClusterJsonnet>)
- [func GenerateComponentJsonnet\(componentOptions cmdInitOptions, dstDir string\) error](<#GenerateComponentJsonnet>)
- [func GenerateLib\(fetch bool, dstDir string\)](<#GenerateLib>)
- [func GenerateReadme\(dstDir string, cmdOptions cmdInitOptions, clusterSpec types.Kr8ClusterSpec\)](<#GenerateReadme>)
- [func GetDefaultFormatOptions\(\) formatter.Options](<#GetDefaultFormatOptions>)
- [type CmdGetOptions](<#CmdGetOptions>)


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
## func [Execute](<https://github.com/ice-bergtech/kr8/blob/main/cmd/root.go#L32>)

```go
func Execute(version string)
```

Execute adds all child commands to the root command sets flags appropriately. This is called by main.main\(\). It only needs to happen once to the rootCmd.

<a name="GenerateClusterJsonnet"></a>
## func [GenerateClusterJsonnet](<https://github.com/ice-bergtech/kr8/blob/main/cmd/init.go#L199>)

```go
func GenerateClusterJsonnet(cSpec types.Kr8ClusterSpec, dstDir string) error
```

Generate a cluster.jsonnet file based on the provided Kr8ClusterSpec and store it in the specified directory.

<a name="GenerateComponentJsonnet"></a>
## func [GenerateComponentJsonnet](<https://github.com/ice-bergtech/kr8/blob/main/cmd/init.go#L214>)

```go
func GenerateComponentJsonnet(componentOptions cmdInitOptions, dstDir string) error
```

Generate default component kr8\_spec values and store in params.jsonnet Based on the type: jsonnet: create a component.jsonnet file and reference it from the params.jsonnet file yml: leave a note in the params.jsonnet file about where and how the yml files can be referenced chart: generate a simple taskfile that handles vendoring the chart data

<a name="GenerateLib"></a>
## func [GenerateLib](<https://github.com/ice-bergtech/kr8/blob/main/cmd/init.go#L279>)

```go
func GenerateLib(fetch bool, dstDir string)
```



<a name="GenerateReadme"></a>
## func [GenerateReadme](<https://github.com/ice-bergtech/kr8/blob/main/cmd/init.go#L284>)

```go
func GenerateReadme(dstDir string, cmdOptions cmdInitOptions, clusterSpec types.Kr8ClusterSpec)
```



<a name="GetDefaultFormatOptions"></a>
## func [GetDefaultFormatOptions](<https://github.com/ice-bergtech/kr8/blob/main/cmd/format.go#L21>)

```go
func GetDefaultFormatOptions() formatter.Options
```

Configures the default options for the jsonnet formatter

<a name="CmdGetOptions"></a>
## type [CmdGetOptions](<https://github.com/ice-bergtech/kr8/blob/main/cmd/get.go#L38-L51>)

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