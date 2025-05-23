# types

```go
import "github.com/ice-bergtech/kr8/pkg/types"
```

Package types contains shared types used across kr8\+ packages.

## Index

- [type CmdJsonnetOptions](<#CmdJsonnetOptions>)
- [type Kr8Cluster](<#Kr8Cluster>)
- [type Kr8Error](<#Kr8Error>)
  - [func \(e Kr8Error\) Error\(\) string](<#Kr8Error.Error>)
- [type Kr8Opts](<#Kr8Opts>)
- [type VMConfig](<#VMConfig>)


<a name="CmdJsonnetOptions"></a>
## type [CmdJsonnetOptions](<https://github.com:icebergtech/kr8/blob/main/pkg/types/types.go#L26-L34>)

Options for running the jsonnet command. Used by a few packages and commands.

```go
type CmdJsonnetOptions struct {
    Prune         bool
    Cluster       string
    ClusterParams string
    Component     string
    Format        string
    Color         bool
    Lint          bool
}
```

<a name="Kr8Cluster"></a>
## type [Kr8Cluster](<https://github.com:icebergtech/kr8/blob/main/pkg/types/types.go#L9-L12>)

An object that stores variables that can be referenced by components.

```go
type Kr8Cluster struct {
    Name string `json:"name"`
    Path string `json:"-"`
}
```

<a name="Kr8Error"></a>
## type [Kr8Error](<https://github.com:icebergtech/kr8/blob/main/pkg/types/types.go#L47-L52>)

Shared kr8\+ error struct.

```go
type Kr8Error struct {
    // Message to show the user.
    Message string
    // Value to include with message
    Value interface{}
}
```

<a name="Kr8Error.Error"></a>
### func \(Kr8Error\) [Error](<https://github.com:icebergtech/kr8/blob/main/pkg/types/types.go#L55>)

```go
func (e Kr8Error) Error() string
```

Error implements error.

<a name="Kr8Opts"></a>
## type [Kr8Opts](<https://github.com:icebergtech/kr8/blob/main/pkg/types/types.go#L15-L22>)

Options that configure where kr8\+ looks for files.

```go
type Kr8Opts struct {
    // Base directory of kr8+ configuration
    BaseDir string
    // Directory where component definitions are stored
    ComponentDir string
    // Directory where cluster configurations are stored
    ClusterDir string
}
```

<a name="VMConfig"></a>
## type [VMConfig](<https://github.com:icebergtech/kr8/blob/main/pkg/types/types.go#L37-L44>)

VMConfig describes configuration to initialize the Jsonnet VM with.

```go
type VMConfig struct {
    // JPaths is a list of paths to search for Jsonnet libraries (libsonnet files)
    JPaths []string `json:"jpath" yaml:"jpath"`
    // ExtVars is a list of external variables to pass to Jsonnet VM
    ExtVars []string `json:"ext_str_file" yaml:"ext_str_files"`
    // base directory for the project
    BaseDir string `json:"base_dir" yaml:"base_dir"`
}
```