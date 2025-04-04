# jvm

```go
import "github.com/ice-bergtech/kr8/pkg/jvm"
```

## Index

- [func JsonnetRender\(cmdFlagsJsonnet types.CmdJsonnetOptions, filename string, vmConfig types.VMConfig\)](<#JsonnetRender>)
- [func JsonnetVM\(vmconfig types.VMConfig\) \(\*jsonnet.VM, error\)](<#JsonnetVM>)
- [func RegisterNativeFuncs\(vm \*jsonnet.VM\)](<#RegisterNativeFuncs>)
- [func RenderClusterParams\(vmconfig types.VMConfig, clusterName string, componentNames \[\]string, clusterParams string, prune bool\) string](<#RenderClusterParams>)
- [func RenderClusterParamsOnly\(vmconfig types.VMConfig, clusterName string, clusterParams string, prune bool\) string](<#RenderClusterParamsOnly>)
- [func RenderJsonnet\(vmConfig types.VMConfig, files \[\]string, param string, prune bool, prepend string, source string\) string](<#RenderJsonnet>)


<a name="JsonnetRender"></a>
## func [JsonnetRender](<https://github.com/ice-bergtech/kr8/blob/main/pkg/jvm/jsonnet.go#L112>)

```go
func JsonnetRender(cmdFlagsJsonnet types.CmdJsonnetOptions, filename string, vmConfig types.VMConfig)
```

Renders a jsonnet file with the specified options.

<a name="JsonnetVM"></a>
## func [JsonnetVM](<https://github.com/ice-bergtech/kr8/blob/main/pkg/jvm/jsonnet.go#L41>)

```go
func JsonnetVM(vmconfig types.VMConfig) (*jsonnet.VM, error)
```



<a name="RegisterNativeFuncs"></a>
## func [RegisterNativeFuncs](<https://github.com/ice-bergtech/kr8/blob/main/pkg/jvm/native_funcs.go#L42>)

```go
func RegisterNativeFuncs(vm *jsonnet.VM)
```

Registers additional native functions in the jsonnet VM These functions are used to extend the functionality of jsonnet Adds on to functions part of the jsonnet stdlib: https://jsonnet.org/ref/stdlib.html

<a name="RenderClusterParams"></a>
## func [RenderClusterParams](<https://github.com/ice-bergtech/kr8/blob/main/pkg/jvm/jsonnet.go#L164>)

```go
func RenderClusterParams(vmconfig types.VMConfig, clusterName string, componentNames []string, clusterParams string, prune bool) string
```

render cluster params, merged with one or more component's parameters. Empty componentName list renders all component parameters

<a name="RenderClusterParamsOnly"></a>
## func [RenderClusterParamsOnly](<https://github.com/ice-bergtech/kr8/blob/main/pkg/jvm/jsonnet.go#L149>)

```go
func RenderClusterParamsOnly(vmconfig types.VMConfig, clusterName string, clusterParams string, prune bool) string
```

only render cluster params \(\_cluster\), without components

<a name="RenderJsonnet"></a>
## func [RenderJsonnet](<https://github.com/ice-bergtech/kr8/blob/main/pkg/jvm/jsonnet.go#L71>)

```go
func RenderJsonnet(vmConfig types.VMConfig, files []string, param string, prune bool, prepend string, source string) string
```

Takes a list of jsonnet files and imports each one and mixes them with "\+"