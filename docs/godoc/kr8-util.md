# util

```go
import "github.com/ice-bergtech/kr8/pkg/util"
```

Utility functions for directories and files

## Index

- [func CalculateClusterIncludesExcludes\(input map\[string\]string, filters PathFilterOptions\) \[\]string](<#CalculateClusterIncludesExcludes>)
- [func CheckObjectMatch\(input gjson.Result, filterString string\) bool](<#CheckObjectMatch>)
- [func Colorize\(input interface\{\}, colorNum int, disabled bool\) string](<#Colorize>)
- [func FatalErrorCheck\(message string, err error\)](<#FatalErrorCheck>)
- [func FetchRepoUrl\(url string, destination string, performFetch bool\) error](<#FetchRepoUrl>)
- [func Filter\(vs \[\]string, f func\(string\) bool\) \[\]string](<#Filter>)
- [func FilterItems\(input map\[string\]string, pf PathFilterOptions\) \[\]string](<#FilterItems>)
- [func FormatJsonnetString\(input string\) \(string, error\)](<#FormatJsonnetString>)
- [func FormatJsonnetStringCustom\(input string, opts formatter.Options\) \(string, error\)](<#FormatJsonnetStringCustom>)
- [func GetClusterFilenames\(searchDir string\) \(\[\]types.Kr8Cluster, error\)](<#GetClusterFilenames>)
- [func GetClusterParamsFilenames\(basePath string, targetPath string\) \[\]string](<#GetClusterParamsFilenames>)
- [func GetClusterPaths\(searchDir string, clusterName string\) string](<#GetClusterPaths>)
- [func GetDefaultFormatOptions\(\) formatter.Options](<#GetDefaultFormatOptions>)
- [func JsonnetPrint\(output string, format string, color bool\)](<#JsonnetPrint>)
- [func Pretty\(input string, colorOutput bool\) string](<#Pretty>)
- [func WriteObjToJsonFile\(filename string, path string, objStruct interface\{\}\) \(string, error\)](<#WriteObjToJsonFile>)
- [type PathFilterOptions](<#PathFilterOptions>)


<a name="CalculateClusterIncludesExcludes"></a>
## func [CalculateClusterIncludesExcludes](<https://github.com/ice-bergtech/kr8/blob/main/pkg/util/util.go#L117>)

```go
func CalculateClusterIncludesExcludes(input map[string]string, filters PathFilterOptions) []string
```

Using the allClusterParams variable and command flags to create a list of clusters to generate. Clusters can be filtered with "=" for equality or "\~" for regex match.

<a name="CheckObjectMatch"></a>
## func [CheckObjectMatch](<https://github.com/ice-bergtech/kr8/blob/main/pkg/util/util.go#L53>)

```go
func CheckObjectMatch(input gjson.Result, filterString string) bool
```

Checks if a input object matches a filter string. The filter string can be an equality match or a regex match.

<a name="Colorize"></a>
## func [Colorize](<https://github.com/ice-bergtech/kr8/blob/main/pkg/util/json.go#L38>)

```go
func Colorize(input interface{}, colorNum int, disabled bool) string
```

Colorize function from zerolog console.go file to replicate their coloring functionality. Source: https://github.com/rs/zerolog/blob/a21d6107dcda23e36bc5cfd00ce8fdbe8f3ddc23/console.go#L389

<a name="FatalErrorCheck"></a>
## func [FatalErrorCheck](<https://github.com/ice-bergtech/kr8/blob/main/pkg/util/util.go#L109>)

```go
func FatalErrorCheck(message string, err error)
```

Logs an error and exits the program if the error is not nil. Saves 3 lines per use and centralizes fatal errors for rewriting.

<a name="FetchRepoUrl"></a>
## func [FetchRepoUrl](<https://github.com/ice-bergtech/kr8/blob/main/pkg/util/remote.go#L13>)

```go
func FetchRepoUrl(url string, destination string, performFetch bool) error
```

Fetch a git repo from a url and clone it to a destination directory. If the performFetch flag is false, it will log the command that would be run and return without doing anything.

<a name="Filter"></a>
## func [Filter](<https://github.com/ice-bergtech/kr8/blob/main/pkg/util/util.go#L13>)

```go
func Filter(vs []string, f func(string) bool) []string
```

Filter returns a new slice containing only the elements that satisfy the predicate function. From https://gobyexample.com/collection-functions

<a name="FilterItems"></a>
## func [FilterItems](<https://github.com/ice-bergtech/kr8/blob/main/pkg/util/util.go#L73>)

```go
func FilterItems(input map[string]string, pf PathFilterOptions) []string
```

Given a map of string, filter them based on the provided options. The map value is parsed as a gjson result and then checked against the provided options.

<a name="FormatJsonnetString"></a>
## func [FormatJsonnetString](<https://github.com/ice-bergtech/kr8/blob/main/pkg/util/json.go#L93>)

```go
func FormatJsonnetString(input string) (string, error)
```

Formats a jsonnet string using the default options.

<a name="FormatJsonnetStringCustom"></a>
## func [FormatJsonnetStringCustom](<https://github.com/ice-bergtech/kr8/blob/main/pkg/util/json.go#L98>)

```go
func FormatJsonnetStringCustom(input string, opts formatter.Options) (string, error)
```

Formats a jsonnet string using custom options.

<a name="GetClusterFilenames"></a>
## func [GetClusterFilenames](<https://github.com/ice-bergtech/kr8/blob/main/pkg/util/directories.go#L16>)

```go
func GetClusterFilenames(searchDir string) ([]types.Kr8Cluster, error)
```

Get a list of cluster from within a directory. Walks the directory tree, creating a types.Kr8Cluster for each cluster.jsonnet file found.

<a name="GetClusterParamsFilenames"></a>
## func [GetClusterParamsFilenames](<https://github.com/ice-bergtech/kr8/blob/main/pkg/util/directories.go#L73>)

```go
func GetClusterParamsFilenames(basePath string, targetPath string) []string
```

Get all cluster parameters within a directory. Walks through the directory hierarchy and returns all paths to \`params.jsonnet\` files.

<a name="GetClusterPaths"></a>
## func [GetClusterPaths](<https://github.com/ice-bergtech/kr8/blob/main/pkg/util/directories.go#L46>)

```go
func GetClusterPaths(searchDir string, clusterName string) string
```

Get a specific cluster within a directory by name. Returns the path to the cluster.

<a name="GetDefaultFormatOptions"></a>
## func [GetDefaultFormatOptions](<https://github.com/ice-bergtech/kr8/blob/main/pkg/util/json.go#L75>)

```go
func GetDefaultFormatOptions() formatter.Options
```

Configures the default options for the jsonnet formatter.

<a name="JsonnetPrint"></a>
## func [JsonnetPrint](<https://github.com/ice-bergtech/kr8/blob/main/pkg/util/json.go#L51>)

```go
func JsonnetPrint(output string, format string, color bool)
```

Print the jsonnet output in the specified format. Acceptable formats are: yaml, stream, json.

<a name="Pretty"></a>
## func [Pretty](<https://github.com/ice-bergtech/kr8/blob/main/pkg/util/json.go#L17>)

```go
func Pretty(input string, colorOutput bool) string
```

Pretty formats the input jsonnet string with indentation and optional color output.

<a name="WriteObjToJsonFile"></a>
## func [WriteObjToJsonFile](<https://github.com/ice-bergtech/kr8/blob/main/pkg/util/json.go#L104>)

```go
func WriteObjToJsonFile(filename string, path string, objStruct interface{}) (string, error)
```

Write out a struct to a specified path and file. If successful, returns what was written. If not successful, returns an error.

<a name="PathFilterOptions"></a>
## type [PathFilterOptions](<https://github.com/ice-bergtech/kr8/blob/main/pkg/util/util.go#L25-L49>)

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