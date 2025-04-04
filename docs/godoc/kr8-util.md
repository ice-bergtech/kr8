# util

```go
import "github.com/ice-bergtech/kr8/pkg/util"
```

## Index

- [func CalculateClusterIncludesExcludes\(input map\[string\]string, filters PathFilterOptions\) \[\]string](<#CalculateClusterIncludesExcludes>)
- [func CheckObjectMatch\(input gjson.Result, filterString string\) bool](<#CheckObjectMatch>)
- [func Colorize\(s interface\{\}, c int, disabled bool\) string](<#Colorize>)
- [func FatalErrorCheck\(err error, message string\)](<#FatalErrorCheck>)
- [func FetchRepoUrl\(url string, destination string, performFetch bool\)](<#FetchRepoUrl>)
- [func Filter\(vs \[\]string, f func\(string\) bool\) \[\]string](<#Filter>)
- [func FilterItems\(input map\[string\]string, pf PathFilterOptions\) \[\]string](<#FilterItems>)
- [func FormatJsonnetString\(input string\) \(string, error\)](<#FormatJsonnetString>)
- [func FormatJsonnetStringCustom\(input string, opts formatter.Options\) \(string, error\)](<#FormatJsonnetStringCustom>)
- [func GetCluster\(searchDir string, clusterName string\) string](<#GetCluster>)
- [func GetClusterParams\(basePath string, targetPath string\) \[\]string](<#GetClusterParams>)
- [func GetClusters\(searchDir string\) \(\[\]types.Kr8Cluster, error\)](<#GetClusters>)
- [func GetDefaultFormatOptions\(\) formatter.Options](<#GetDefaultFormatOptions>)
- [func JsonnetPrint\(output string, format string, color bool\)](<#JsonnetPrint>)
- [func Pretty\(input string, colorOutput bool\) string](<#Pretty>)
- [func WriteObjToJsonFile\(filename string, path string, objStruct interface\{\}\) error](<#WriteObjToJsonFile>)
- [type PathFilterOptions](<#PathFilterOptions>)


<a name="CalculateClusterIncludesExcludes"></a>
## func [CalculateClusterIncludesExcludes](<https://github.com/ice-bergtech/kr8/blob/main/pkg/util/util.go#L105>)

```go
func CalculateClusterIncludesExcludes(input map[string]string, filters PathFilterOptions) []string
```

Using the allClusterParams variable and command flags to create a list of clusters to generate Clusters can be filtered with "=" for equality or "\~" for regex match

<a name="CheckObjectMatch"></a>
## func [CheckObjectMatch](<https://github.com/ice-bergtech/kr8/blob/main/pkg/util/util.go#L46>)

```go
func CheckObjectMatch(input gjson.Result, filterString string) bool
```

Checks if a input object matches a filter string. The filter string can be an equality match or a regex match.

<a name="Colorize"></a>
## func [Colorize](<https://github.com/ice-bergtech/kr8/blob/main/pkg/util/json.go#L37>)

```go
func Colorize(s interface{}, c int, disabled bool) string
```

colorize function from zerolog console.go file to replicate their coloring functionality. https://github.com/rs/zerolog/blob/a21d6107dcda23e36bc5cfd00ce8fdbe8f3ddc23/console.go#L389

<a name="FatalErrorCheck"></a>
## func [FatalErrorCheck](<https://github.com/ice-bergtech/kr8/blob/main/pkg/util/util.go#L97>)

```go
func FatalErrorCheck(err error, message string)
```

util.FatalErrorCheck is a helper function that logs an error and exits the program if the error is not nil. Saves 3 lines per use and centralizes fatal errors for rewriting

<a name="FetchRepoUrl"></a>
## func [FetchRepoUrl](<https://github.com/ice-bergtech/kr8/blob/main/pkg/util/remote.go#L12>)

```go
func FetchRepoUrl(url string, destination string, performFetch bool)
```

Fetch a git repo from a url and clone it to a destination directory if the performFetch flag is false, it will log the command that would be run and return without doing anything

<a name="Filter"></a>
## func [Filter](<https://github.com/ice-bergtech/kr8/blob/main/pkg/util/util.go#L13>)

```go
func Filter(vs []string, f func(string) bool) []string
```

Filter returns a new slice containing only the elements that satisfy the predicate function. From https://gobyexample.com/collection-functions

<a name="FilterItems"></a>
## func [FilterItems](<https://github.com/ice-bergtech/kr8/blob/main/pkg/util/util.go#L63>)

```go
func FilterItems(input map[string]string, pf PathFilterOptions) []string
```



<a name="FormatJsonnetString"></a>
## func [FormatJsonnetString](<https://github.com/ice-bergtech/kr8/blob/main/pkg/util/json.go#L92>)

```go
func FormatJsonnetString(input string) (string, error)
```

Formats a jsonnet string using the default options

<a name="FormatJsonnetStringCustom"></a>
## func [FormatJsonnetStringCustom](<https://github.com/ice-bergtech/kr8/blob/main/pkg/util/json.go#L97>)

```go
func FormatJsonnetStringCustom(input string, opts formatter.Options) (string, error)
```

Formats a jsonnet string using custom options

<a name="GetCluster"></a>
## func [GetCluster](<https://github.com/ice-bergtech/kr8/blob/main/pkg/util/directories.go#L43>)

```go
func GetCluster(searchDir string, clusterName string) string
```



<a name="GetClusterParams"></a>
## func [GetClusterParams](<https://github.com/ice-bergtech/kr8/blob/main/pkg/util/directories.go#L67>)

```go
func GetClusterParams(basePath string, targetPath string) []string
```



<a name="GetClusters"></a>
## func [GetClusters](<https://github.com/ice-bergtech/kr8/blob/main/pkg/util/directories.go#L13>)

```go
func GetClusters(searchDir string) ([]types.Kr8Cluster, error)
```



<a name="GetDefaultFormatOptions"></a>
## func [GetDefaultFormatOptions](<https://github.com/ice-bergtech/kr8/blob/main/pkg/util/json.go#L75>)

```go
func GetDefaultFormatOptions() formatter.Options
```

Configures the default options for the jsonnet formatter

<a name="JsonnetPrint"></a>
## func [JsonnetPrint](<https://github.com/ice-bergtech/kr8/blob/main/pkg/util/json.go#L51>)

```go
func JsonnetPrint(output string, format string, color bool)
```

Print the jsonnet output in the specified format allows for: yaml, stream, json

<a name="Pretty"></a>
## func [Pretty](<https://github.com/ice-bergtech/kr8/blob/main/pkg/util/json.go#L17>)

```go
func Pretty(input string, colorOutput bool) string
```

Pretty formats the input jsonnet string with indentation and optional color output.

<a name="WriteObjToJsonFile"></a>
## func [WriteObjToJsonFile](<https://github.com/ice-bergtech/kr8/blob/main/pkg/util/json.go#L102>)

```go
func WriteObjToJsonFile(filename string, path string, objStruct interface{}) error
```

Write out a struct to a specified path and file

<a name="PathFilterOptions"></a>
## type [PathFilterOptions](<https://github.com/ice-bergtech/kr8/blob/main/pkg/util/util.go#L24-L42>)

Fill with string to include and exclude, using kr8's special parsing

```go
type PathFilterOptions struct {
    // Comma-separated list of include filters
    // Filters can include regex filters using the "~" operator. For example, "name~^myregex$"
    // Filters can include equality matches using the "=" operator. For example, "name=myvalue"
    // Filters can include substring matches using the "=" operator. For example, "name=myvalue"
    // If no operator is provided, it is treated as a substring match against the "name" field.
    Includes string
    // Comma-separated list of exclude filters
    // Filters can include regex filters using the "~" operator. For example, "name~^myregex$"
    // Filters can include equality matches using the "=" operator. For example, "name=myvalue"
    // Filters can include substring matches using the "=" operator. For example, "name=myvalue"
    // If no operator is provided, it is treated as a substring match against the "name" field.
    Excludes string
    // Comma separated cluster names
    // Filters keys on exact match
    Clusters string
    // Comma separated component names
    Components string
}
```