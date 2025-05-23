# jnetvm

```go
import "github.com/ice-bergtech/kr8/pkg/jnetvm"
```

Package jvm contains the jsonnet rendering logic.

## Index

- [func JsonnetRender\(cmdFlagsJsonnet types.CmdJsonnetOptions, filename string, vmConfig types.VMConfig, logger zerolog.Logger\) error](<#JsonnetRender>)
- [func JsonnetRenderClusterParams\(vmConfig types.VMConfig, clusterName string, componentNames \[\]string, clusterParams string, prune bool, lint bool\) \(string, error\)](<#JsonnetRenderClusterParams>)
- [func JsonnetRenderClusterParamsOnly\(vmConfig types.VMConfig, clusterName string, clusterParams string, prune bool, lint bool\) \(string, error\)](<#JsonnetRenderClusterParamsOnly>)
- [func JsonnetRenderFiles\(vmConfig types.VMConfig, files \[\]string, param string, prune bool, prepend string, source string, lint bool\) \(string, error\)](<#JsonnetRenderFiles>)
- [func JsonnetVM\(vmConfig types.VMConfig\) \(\*jsonnet.VM, error\)](<#JsonnetVM>)
- [func MergeComponentDefaults\(componentMap map\[string\]kr8\_types.Kr8ClusterComponentRef, componentNames \[\]string, vmConfig types.VMConfig\) \(string, error\)](<#MergeComponentDefaults>)


<a name="JsonnetRender"></a>
## func [JsonnetRender](<https://github.com:icebergtech/kr8/blob/main/pkg/jnetvm/jsonnet.go#L142-L147>)

```go
func JsonnetRender(cmdFlagsJsonnet types.CmdJsonnetOptions, filename string, vmConfig types.VMConfig, logger zerolog.Logger) error
```

Renders a jsonnet file with the specified options.

<a name="JsonnetRenderClusterParams"></a>
## func [JsonnetRenderClusterParams](<https://github.com:icebergtech/kr8/blob/main/pkg/jnetvm/jsonnet.go#L222-L229>)

```go
func JsonnetRenderClusterParams(vmConfig types.VMConfig, clusterName string, componentNames []string, clusterParams string, prune bool, lint bool) (string, error)
```

Render cluster params, merged with one or more component's parameters. Empty componentName list renders all component parameters.

<a name="JsonnetRenderClusterParamsOnly"></a>
## func [JsonnetRenderClusterParamsOnly](<https://github.com:icebergtech/kr8/blob/main/pkg/jnetvm/jsonnet.go#L198-L204>)

```go
func JsonnetRenderClusterParamsOnly(vmConfig types.VMConfig, clusterName string, clusterParams string, prune bool, lint bool) (string, error)
```

Only render cluster params \(\_cluster\), without components.

<a name="JsonnetRenderFiles"></a>
## func [JsonnetRenderFiles](<https://github.com:icebergtech/kr8/blob/main/pkg/jnetvm/jsonnet.go#L84-L92>)

```go
func JsonnetRenderFiles(vmConfig types.VMConfig, files []string, param string, prune bool, prepend string, source string, lint bool) (string, error)
```

Takes a list of jsonnet files and imports each one. Formats the string for jsonnet using "\+". source is only used for error messages.

<a name="JsonnetVM"></a>
## func [JsonnetVM](<https://github.com:icebergtech/kr8/blob/main/pkg/jnetvm/jsonnet.go#L51>)

```go
func JsonnetVM(vmConfig types.VMConfig) (*jsonnet.VM, error)
```

Create a Jsonnet VM to run commands in. It:

- creates a jsonnet VM
- registers kr8\+ native functions
- Add jsonnet library directories
- loads external files into extVars

<a name="MergeComponentDefaults"></a>
## func [MergeComponentDefaults](<https://github.com:icebergtech/kr8/blob/main/pkg/jnetvm/jsonnet.go#L268-L272>)

```go
func MergeComponentDefaults(componentMap map[string]kr8_types.Kr8ClusterComponentRef, componentNames []string, vmConfig types.VMConfig) (string, error)
```

