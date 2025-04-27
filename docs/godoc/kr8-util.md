# util

```go
import "github.com/ice-bergtech/kr8/pkg/util"
```

Package util contains various utility functions for directories and files. It includes functions for formatting JSON, writing to files, directory management, and go control\-flow helpers

## Index

- [func CalculateClusterIncludesExcludes\(input map\[string\]string, filters PathFilterOptions\) \[\]string](<#CalculateClusterIncludesExcludes>)
- [func CheckObjectMatch\(input gjson.Result, filterString string\) bool](<#CheckObjectMatch>)
- [func CleanOutputDir\(outputFileMap map\[string\]bool, componentOutputDir string\) error](<#CleanOutputDir>)
- [func Colorize\(input interface\{\}, colorNum int, disabled bool\) string](<#Colorize>)
- [func ErrorIfCheck\(message string, err error\) error](<#ErrorIfCheck>)
- [func FatalErrorCheck\(message string, err error, logger zerolog.Logger\)](<#FatalErrorCheck>)
- [func FetchRepoUrl\(url string, destination string, noop bool\) error](<#FetchRepoUrl>)
- [func Filter\(vs \[\]string, f func\(string\) bool\) \[\]string](<#Filter>)
- [func FilterItems\(input map\[string\]string, pfilter PathFilterOptions\) \[\]string](<#FilterItems>)
- [func FormatJsonnetString\(input string\) \(string, error\)](<#FormatJsonnetString>)
- [func FormatJsonnetStringCustom\(input string, opts formatter.Options\) \(string, error\)](<#FormatJsonnetStringCustom>)
- [func GetClusterFilenames\(searchDir string\) \(\[\]types.Kr8Cluster, error\)](<#GetClusterFilenames>)
- [func GetClusterParamsFilenames\(basePath string, targetPath string\) \[\]string](<#GetClusterParamsFilenames>)
- [func GetClusterPaths\(searchDir string, clusterName string\) \(string, error\)](<#GetClusterPaths>)
- [func GetDefaultFormatOptions\(\) formatter.Options](<#GetDefaultFormatOptions>)
- [func JsonnetPrint\(output string, format string, color bool\) error](<#JsonnetPrint>)
- [func LogErrorIfCheck\(message string, err error, logger zerolog.Logger\) error](<#LogErrorIfCheck>)
- [func Pretty\(inputJson string, colorOutput bool\) \(string, error\)](<#Pretty>)
- [func SetupLogger\(enableColor bool\) zerolog.Logger](<#SetupLogger>)
- [func WriteObjToJsonFile\(filename string, path string, objStruct interface\{\}\) \(string, error\)](<#WriteObjToJsonFile>)
- [type ClusterTreeNode](<#ClusterTreeNode>)
- [type PathFilterOptions](<#PathFilterOptions>)


<a name="CalculateClusterIncludesExcludes"></a>
## func [CalculateClusterIncludesExcludes](<https://github.com:icebergtech/kr8/blob/main/pkg/util/util.go#L160>)

```go
func CalculateClusterIncludesExcludes(input map[string]string, filters PathFilterOptions) []string
```

Using the allClusterParams variable and command flags to create a list of clusters to generate. Clusters can be filtered with "=" for equality or "\~" for regex match.

<a name="CheckObjectMatch"></a>
## func [CheckObjectMatch](<https://github.com:icebergtech/kr8/blob/main/pkg/util/util.go#L80>)

```go
func CheckObjectMatch(input gjson.Result, filterString string) bool
```

Checks if a input object matches a filter string. The filter string can be an equality match or a regex match.

<a name="CleanOutputDir"></a>
## func [CleanOutputDir](<https://github.com:icebergtech/kr8/blob/main/pkg/util/directories.go#L117>)

```go
func CleanOutputDir(outputFileMap map[string]bool, componentOutputDir string) error
```

Given a map of filenames, prunes all \*.yaml files that are NOT in the map from the directory.

<a name="Colorize"></a>
## func [Colorize](<https://github.com:icebergtech/kr8/blob/main/pkg/util/json.go#L42>)

```go
func Colorize(input interface{}, colorNum int, disabled bool) string
```

Colorize function from zerolog console.go file to replicate their coloring functionality. Source: https://github.com/rs/zerolog/blob/a21d6107dcda23e36bc5cfd00ce8fdbe8f3ddc23/console.go#L389 Replicated here because it's a private function.

<a name="ErrorIfCheck"></a>
## func [ErrorIfCheck](<https://github.com:icebergtech/kr8/blob/main/pkg/util/util.go#L140>)

```go
func ErrorIfCheck(message string, err error) error
```



<a name="FatalErrorCheck"></a>
## func [FatalErrorCheck](<https://github.com:icebergtech/kr8/blob/main/pkg/util/util.go#L134>)

```go
func FatalErrorCheck(message string, err error, logger zerolog.Logger)
```

Logs an error and exits the program if the error is not nil. Saves 3 lines per use and centralizes fatal errors for rewriting.

<a name="FetchRepoUrl"></a>
## func [FetchRepoUrl](<https://github.com:icebergtech/kr8/blob/main/pkg/util/remote.go#L13>)

```go
func FetchRepoUrl(url string, destination string, noop bool) error
```

Fetch a git repo from a url and clone it to a destination directory. If the noop flag is true, it print commands to fetch manually without doing anything.

<a name="Filter"></a>
## func [Filter](<https://github.com:icebergtech/kr8/blob/main/pkg/util/util.go#L23>)

```go
func Filter(vs []string, f func(string) bool) []string
```

Filter returns a new slice containing only the elements that satisfy the predicate function. From https://gobyexample.com/collection-functions

<a name="FilterItems"></a>
## func [FilterItems](<https://github.com:icebergtech/kr8/blob/main/pkg/util/util.go#L100>)

```go
func FilterItems(input map[string]string, pfilter PathFilterOptions) []string
```

Given a map of string, filter them based on the provided options. The map value is parsed as a gjson result and then checked against the provided options.

<a name="FormatJsonnetString"></a>
## func [FormatJsonnetString](<https://github.com:icebergtech/kr8/blob/main/pkg/util/json.go#L108>)

```go
func FormatJsonnetString(input string) (string, error)
```

Formats a jsonnet string using the default options.

<a name="FormatJsonnetStringCustom"></a>
## func [FormatJsonnetStringCustom](<https://github.com:icebergtech/kr8/blob/main/pkg/util/json.go#L113>)

```go
func FormatJsonnetStringCustom(input string, opts formatter.Options) (string, error)
```

Formats a jsonnet string using custom options.

<a name="GetClusterFilenames"></a>
## func [GetClusterFilenames](<https://github.com:icebergtech/kr8/blob/main/pkg/util/directories.go#L22>)

```go
func GetClusterFilenames(searchDir string) ([]types.Kr8Cluster, error)
```

Get a list of cluster from within a directory. Walks the directory tree, creating a types.Kr8Cluster for each cluster.jsonnet file found.

<a name="GetClusterParamsFilenames"></a>
## func [GetClusterParamsFilenames](<https://github.com:icebergtech/kr8/blob/main/pkg/util/directories.go#L78>)

```go
func GetClusterParamsFilenames(basePath string, targetPath string) []string
```

Get all cluster parameters within a directory. Walks through the directory hierarchy and returns all paths to \`params.jsonnet\` files.

<a name="GetClusterPaths"></a>
## func [GetClusterPaths](<https://github.com:icebergtech/kr8/blob/main/pkg/util/directories.go#L52>)

```go
func GetClusterPaths(searchDir string, clusterName string) (string, error)
```

Get a specific cluster within a directory by name. Walks the cluster directory searching for the given clusterName. Returns the path to the cluster.

<a name="GetDefaultFormatOptions"></a>
## func [GetDefaultFormatOptions](<https://github.com:icebergtech/kr8/blob/main/pkg/util/json.go#L90>)

```go
func GetDefaultFormatOptions() formatter.Options
```

Configures the default options for the jsonnet formatter.

<a name="JsonnetPrint"></a>
## func [JsonnetPrint](<https://github.com:icebergtech/kr8/blob/main/pkg/util/json.go#L55>)

```go
func JsonnetPrint(output string, format string, color bool) error
```

Print the jsonnet in the specified format. Acceptable formats are: yaml, stream, json.

<a name="LogErrorIfCheck"></a>
## func [LogErrorIfCheck](<https://github.com:icebergtech/kr8/blob/main/pkg/util/util.go#L148>)

```go
func LogErrorIfCheck(message string, err error, logger zerolog.Logger) error
```



<a name="Pretty"></a>
## func [Pretty](<https://github.com:icebergtech/kr8/blob/main/pkg/util/json.go#L18>)

```go
func Pretty(inputJson string, colorOutput bool) (string, error)
```

Pretty formats the input jsonnet string with indentation and optional color output. Returns an error when the input can't properly format the json string input.

<a name="SetupLogger"></a>
## func [SetupLogger](<https://github.com:icebergtech/kr8/blob/main/pkg/util/util.go#L34>)

```go
func SetupLogger(enableColor bool) zerolog.Logger
```



<a name="WriteObjToJsonFile"></a>
## func [WriteObjToJsonFile](<https://github.com:icebergtech/kr8/blob/main/pkg/util/json.go#L120>)

```go
func WriteObjToJsonFile(filename string, path string, objStruct interface{}) (string, error)
```

Write out a struct to a specified path and file. Marshals the given interface and generates a formatted json string. All parent directories needed are created.

<a name="ClusterTreeNode"></a>
## type [ClusterTreeNode](<https://github.com:icebergtech/kr8/blob/main/pkg/util/directories.go#L14-L18>)



```go
type ClusterTreeNode struct {
    DirName string
    VM      *jsonnet.VM
    Parent  *ClusterTreeNode
}
```

<a name="PathFilterOptions"></a>
## type [PathFilterOptions](<https://github.com:icebergtech/kr8/blob/main/pkg/util/util.go#L52-L76>)

Fill with string to include and exclude, using kr8's special parsing.

```go
type PathFilterOptions struct {
    // Comma-separated list of include filters
    // Filters can include:
    //
    // regex filters using the "~" operator. For example, "name~^myregex$"
    // equality matches using the "=" operator. For example, "name=myvalue"
    // substring matches using the "=" operator. For example, "name=myvalue"
    //
    // If no operator is provided, it is treated as a substring match against the "name" field.
    Includes string
    // Comma-separated list of exclude filters.
    // Filters can include:
    //
    // regex filters using the "~" operator. For example, "name~^myregex$"
    // equality matches using the "=" operator. For example, "name=myvalue"
    // substring matches using the "=" operator. For example, "name=myvalue"
    //
    // If no operator is provided, it is treated as a substring match against the "name" field.
    Excludes string
    // Comma separated cluster names.
    // Filters keys on exact match.
    Clusters string
    // Comma separated component names.
    Components string
}
```