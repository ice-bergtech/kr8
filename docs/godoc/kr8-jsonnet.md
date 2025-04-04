# jvm

```go
import "github.com/ice-bergtech/kr8/pkg/jvm"
```

Package jvm contains the jsonnet rendering logic.

## Index

- [func JsonnetRender\(cmdFlagsJsonnet types.CmdJsonnetOptions, filename string, vmConfig types.VMConfig\)](<#JsonnetRender>)
- [func JsonnetRenderClusterParams\(vmconfig types.VMConfig, clusterName string, componentNames \[\]string, clusterParams string, prune bool\) string](<#JsonnetRenderClusterParams>)
- [func JsonnetRenderClusterParamsOnly\(vmconfig types.VMConfig, clusterName string, clusterParams string, prune bool\) string](<#JsonnetRenderClusterParamsOnly>)
- [func JsonnetRenderFiles\(vmConfig types.VMConfig, files \[\]string, param string, prune bool, prepend string, source string\) string](<#JsonnetRenderFiles>)
- [func JsonnetVM\(vmconfig types.VMConfig\) \(\*jsonnet.VM, error\)](<#JsonnetVM>)
- [func NativeHelmTemplate\(\) \*jsonnet.NativeFunction](<#NativeHelmTemplate>)
- [func NativeKompose\(\) \*jsonnet.NativeFunction](<#NativeKompose>)
- [func NativeRegexEscape\(\) \*jsonnet.NativeFunction](<#NativeRegexEscape>)
- [func NativeRegexMatch\(\) \*jsonnet.NativeFunction](<#NativeRegexMatch>)
- [func NativeRegexSubst\(\) \*jsonnet.NativeFunction](<#NativeRegexSubst>)
- [func NativeSprigTemplate\(\) \*jsonnet.NativeFunction](<#NativeSprigTemplate>)
- [func RegisterNativeFuncs\(vm \*jsonnet.VM\)](<#RegisterNativeFuncs>)


<a name="JsonnetRender"></a>
## func [JsonnetRender](<https://github.com/ice-bergtech/kr8/blob/main/pkg/jvm/jsonnet.go#L111>)

```go
func JsonnetRender(cmdFlagsJsonnet types.CmdJsonnetOptions, filename string, vmConfig types.VMConfig)
```

Renders a jsonnet file with the specified options.

<a name="JsonnetRenderClusterParams"></a>
## func [JsonnetRenderClusterParams](<https://github.com/ice-bergtech/kr8/blob/main/pkg/jvm/jsonnet.go#L163>)

```go
func JsonnetRenderClusterParams(vmconfig types.VMConfig, clusterName string, componentNames []string, clusterParams string, prune bool) string
```

Render cluster params, merged with one or more component's parameters. Empty componentName list renders all component parameters

<a name="JsonnetRenderClusterParamsOnly"></a>
## func [JsonnetRenderClusterParamsOnly](<https://github.com/ice-bergtech/kr8/blob/main/pkg/jvm/jsonnet.go#L148>)

```go
func JsonnetRenderClusterParamsOnly(vmconfig types.VMConfig, clusterName string, clusterParams string, prune bool) string
```

Only render cluster params \(\_cluster\), without components

<a name="JsonnetRenderFiles"></a>
## func [JsonnetRenderFiles](<https://github.com/ice-bergtech/kr8/blob/main/pkg/jvm/jsonnet.go#L70>)

```go
func JsonnetRenderFiles(vmConfig types.VMConfig, files []string, param string, prune bool, prepend string, source string) string
```

Takes a list of jsonnet files and imports each one and mixes them with "\+"

<a name="JsonnetVM"></a>
## func [JsonnetVM](<https://github.com/ice-bergtech/kr8/blob/main/pkg/jvm/jsonnet.go#L40>)

```go
func JsonnetVM(vmconfig types.VMConfig) (*jsonnet.VM, error)
```

Create a Jsonnet VM to run commands in

<a name="NativeHelmTemplate"></a>
## func [NativeHelmTemplate](<https://github.com/ice-bergtech/kr8/blob/main/pkg/jvm/native_funcs.go#L69>)

```go
func NativeHelmTemplate() *jsonnet.NativeFunction
```

Allows executing helm template to process a helm chart and make available to kr8 configuration.

Source: https://github.com/grafana/tanka/blob/v0.27.1/pkg/helm/template.go#L23

<a name="NativeKompose"></a>
## func [NativeKompose](<https://github.com/ice-bergtech/kr8/blob/main/pkg/jvm/native_funcs.go#L152>)

```go
func NativeKompose() *jsonnet.NativeFunction
```

Allows converting a docker\-compose file string into kubernetes resources using kompose

Source: https://github.com/kubernetes/kompose/blob/main/cmd/convert.go

Files in the directory must be in the format \`\[docker\-\]compose.ym\[a\]l\`

Inputs: \`inFile\`, \`outPath\`, \`opts\`

<a name="NativeRegexEscape"></a>
## func [NativeRegexEscape](<https://github.com/ice-bergtech/kr8/blob/main/pkg/jvm/native_funcs.go#L104>)

```go
func NativeRegexEscape() *jsonnet.NativeFunction
```

Escapes a string for use in regex

Inputs: "str"

<a name="NativeRegexMatch"></a>
## func [NativeRegexMatch](<https://github.com/ice-bergtech/kr8/blob/main/pkg/jvm/native_funcs.go#L116>)

```go
func NativeRegexMatch() *jsonnet.NativeFunction
```

Matches a string against a regex pattern

Inputs: "regex", "string"

<a name="NativeRegexSubst"></a>
## func [NativeRegexSubst](<https://github.com/ice-bergtech/kr8/blob/main/pkg/jvm/native_funcs.go#L128>)

```go
func NativeRegexSubst() *jsonnet.NativeFunction
```

Substitutes a regex pattern in a string with another string

Inputs: "regex", "src", "repl"

<a name="NativeSprigTemplate"></a>
## func [NativeSprigTemplate](<https://github.com/ice-bergtech/kr8/blob/main/pkg/jvm/native_funcs.go#L78>)

```go
func NativeSprigTemplate() *jsonnet.NativeFunction
```

Uses sprig to process passed in config data and template.

Sprig template guide: https://masterminds.github.io/sprig/

Inputs: "config" "str"

<a name="RegisterNativeFuncs"></a>
## func [RegisterNativeFuncs](<https://github.com/ice-bergtech/kr8/blob/main/pkg/jvm/native_funcs.go#L45>)

```go
func RegisterNativeFuncs(vm *jsonnet.VM)
```

Registers additional native functions in the jsonnet VM. These functions are used to extend the functionality of jsonnet. Adds on to functions part of the jsonnet stdlib: https://jsonnet.org/ref/stdlib.html