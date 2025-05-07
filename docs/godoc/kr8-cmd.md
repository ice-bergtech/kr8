# cmd

```go
import "github.com/ice-bergtech/kr8/cmd"
```

Copyright © 2019 kubecfg Authors

Licensed under the Apache License, Version 2.0 \(the "License"\); you may not use this file except in compliance with the License. You may obtain a copy of the License at

```
http://www.apache.org/licenses/LICENSE-2.0
```

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.

Defines the cli\-interface commands available to the user.

## Index

- [Variables](<#variables>)
- [func ConfigureLogger\(debug bool\)](<#ConfigureLogger>)
- [func Execute\(ver string\)](<#Execute>)
- [func FormatFile\(filename string, logger zerolog.Logger\) error](<#FormatFile>)
- [func GenerateCmdClusterListBuilder\(allClusterParams map\[string\]string\) \[\]string](<#GenerateCmdClusterListBuilder>)
- [func GenerateCommand\(cmd \*cobra.Command, args \[\]string\)](<#GenerateCommand>)
- [func InitConfig\(\)](<#InitConfig>)
- [func ProfilingFinalizer\(\)](<#ProfilingFinalizer>)
- [func ProfilingInitializer\(\)](<#ProfilingInitializer>)
- [type CmdFormatOptions](<#CmdFormatOptions>)
- [type CmdGenerateOptions](<#CmdGenerateOptions>)
- [type CmdGetOptions](<#CmdGetOptions>)
- [type CmdRenderOptions](<#CmdRenderOptions>)
- [type CmdRootOptions](<#CmdRootOptions>)
- [type Stamp](<#Stamp>)


## Variables

<a name="FormatCmd"></a>

```go
var FormatCmd = &cobra.Command{
    Use:     "format [flags] [files or directories]",
    Aliases: []string{"fmt", "f"},
    Short:   "Format jsonnet files in a directory.  Defaults to `./`",
    Long: `Formats jsonnet and libsonnet files.
A list of files and/or directories. Defaults to current directory (./).
If path is a directory, scans directories for files with matching extensions.
Formats files with the following options: ` + prettyPrintFormattingOpts(),
    Args: cobra.MinimumNArgs(0),
    Run: func(cmd *cobra.Command, args []string) {
        log.Debug().Any("options", cmdFormatOptions).Msg("Formatting files...")
        paths := args
        if len(paths) == 0 {
            paths = []string{"./"}
        }

        for _, path := range paths {
            logger := log.With().Str("param", path).Logger()
            err := util.FileFuncInDir(path, cmdFormatOptions.Recursive, FormatFile, logger)
            if err != nil {
                logger.Error().Err(err).Msg("issue formatting path")
            }
        }
    },
}
```

<a name="GenerateCmd"></a>

```go
var GenerateCmd = &cobra.Command{
    Use:     "generate [flags]",
    Aliases: []string{"gen", "g"},
    Short:   "Generate components",
    Long:    `Generate components in clusters`,
    Example: "kr8 generate",

    Args: cobra.MinimumNArgs(0),
    Run:  GenerateCommand,
}
```

<a name="GetClustersCmd"></a>

```go
var GetClustersCmd = &cobra.Command{
    Use:   "clusters [flags]",
    Short: "Get all clusters",
    Long:  "Get all clusters defined in kr8+ config hierarchy",
    Run: func(cmd *cobra.Command, args []string) {

        clusters, err := util.GetClusterFilenames(RootConfig.ClusterDir)
        util.FatalErrorCheck("Error getting clusters", err, log.Logger)

        if cmdGetFlags.NoTable {
            for _, c := range clusters {
                println(c.Name + ": " + c.Path)
            }

            return
        }

        var entry []string
        table := tablewriter.NewWriter(os.Stdout)
        table.SetHeader([]string{"Name", "Path"})

        for _, c := range clusters {
            entry = append(entry, c.Name)
            entry = append(entry, c.Path)
            table.Append(entry)
            entry = entry[:0]
        }
        table.Render()

    },
}
```

<a name="GetCmd"></a>GetCmd represents the get command.

```go
var GetCmd = &cobra.Command{
    Use:   "get",
    Short: "Display one or many kr8+ resources",
    Long:  `Displays information about kr8+ resources such as clusters and components`,
}
```

<a name="GetComponentsCmd"></a>

```go
var GetComponentsCmd = &cobra.Command{
    Use:   "components [flags]",
    Short: "Get all components",
    Long:  "Get all available components defined in the kr8+ config hierarchy",
    Run: func(cmd *cobra.Command, args []string) {

        if cmdGetFlags.Cluster == "" && cmdGetFlags.ClusterParams == "" {
            log.Fatal().Msg("Please specify a --cluster name and/or --clusterparams file")
        }

        var params []string
        if cmdGetFlags.Cluster != "" {
            clusterPath, err := util.GetClusterPath(RootConfig.ClusterDir, cmdGetFlags.Cluster)
            util.FatalErrorCheck("error getting cluster path for "+cmdGetFlags.Cluster, err, log.Logger)
            params = util.GetClusterParamsFilenames(RootConfig.ClusterDir, clusterPath)
        }
        if cmdGetFlags.ClusterParams != "" {
            params = append(params, cmdGetFlags.ClusterParams)
        }

        jvm, err := jnetvm.JsonnetRenderFiles(RootConfig.VMConfig, params, "._components", true, "", "components", false)
        util.FatalErrorCheck("error rendering jsonnet files", err, log.Logger)
        if cmdGetFlags.ParamField != "" {
            value := gjson.Get(jvm, cmdGetFlags.ParamField)
            if value.String() == "" {
                log.Fatal().Msg("Error getting param: " + cmdGetFlags.ParamField)
            } else {
                formatted, err := util.Pretty(jvm, RootConfig.Color)
                util.FatalErrorCheck("error pretty printing jsonnet", err, log.Logger)
                fmt.Println(formatted)
            }
        } else {
            formatted, err := util.Pretty(jvm, RootConfig.Color)
            util.FatalErrorCheck("error pretty printing jsonnet", err, log.Logger)
            fmt.Println(formatted)
        }
    },
}
```

<a name="GetParamsCmd"></a>

```go
var GetParamsCmd = &cobra.Command{
    Use:   "params [flags]",
    Short: "Get parameter for components and clusters",
    Long:  "Get parameters assigned to clusters and components in the kr8+ config hierarchy",
    Run: func(cmd *cobra.Command, args []string) {
        if cmdGetFlags.Cluster == "" {
            log.Fatal().Msg("Please specify a --cluster")
        }

        var cList []string
        if cmdGetFlags.Component != "" {
            cList = append(cList, cmdGetFlags.Component)
        }

        params, err := jnetvm.JsonnetRenderClusterParams(
            RootConfig.VMConfig,
            cmdGetFlags.Cluster,
            cList,
            cmdGetFlags.ClusterParams,
            true,
            false,
        )
        util.FatalErrorCheck("error rendering cluster params", err, log.Logger)

        if cmdGetFlags.ParamField == "" {
            if cmdGetFlags.Component != "" {
                result := gjson.Get(params, cmdGetFlags.Component).String()
                formatted, err := util.Pretty(result, RootConfig.Color)
                util.FatalErrorCheck("error pretty printing jsonnet", err, log.Logger)
                fmt.Println(formatted)
            } else {
                formatted, err := util.Pretty(params, RootConfig.Color)
                util.FatalErrorCheck("error pretty printing jsonnet", err, log.Logger)
                fmt.Println(formatted)
            }

            return
        }

        if cmdGetFlags.ParamField != "" {
            value := gjson.Get(params, cmdGetFlags.ParamField)
            if value.String() == "" {
                log.Fatal().Msg("Error getting param: " + cmdGetFlags.ParamField)
            }

            fmt.Println(value)
        }
    },
}
```

<a name="InitClusterCmd"></a>

```go
var InitClusterCmd = &cobra.Command{
    Use:   "cluster [flags]",
    Short: "Init a new cluster config file",
    Long:  "Initialize a new cluster configuration file",
    Run: func(cmd *cobra.Command, args []string) {
        cSpec := kr8_types.Kr8ClusterSpec{
            Name:               cmdInitFlags.ClusterName,
            PostProcessor:      "function(input) input",
            GenerateDir:        "generated",
            GenerateShortNames: false,
            PruneParams:        false,
            ClusterOutputDir:   RootConfig.ClusterDir,
            EnableCache:        true,
            CompressCache:      true,
        }

        if cmdInitFlags.Interactive {

            prompt := &survey.Input{
                Message: "Set the cluster configuration directory",
                Default: RootConfig.ClusterDir,
                Help:    "Set the root directory to store cluster configurations, optionally including subdirectories",
            }
            util.FatalErrorCheck("Invalid cluster directory", survey.AskOne(prompt, &cSpec.ClusterOutputDir), log.Logger)

            prompt = &survey.Input{
                Message: "Set the cluster name",
                Default: cmdInitFlags.ClusterName,
                Help:    "Distinct name for the cluster",
            }
            util.FatalErrorCheck("Invalid cluster name", survey.AskOne(prompt, &cSpec.Name), log.Logger)

            promptB := &survey.Confirm{
                Message: "Generate short names for output file names?",
                Default: cSpec.GenerateShortNames,
                Help:    "Shortens component names and file structure",
            }
            util.FatalErrorCheck("Invalid option", survey.AskOne(promptB, &cSpec.GenerateShortNames), log.Logger)

            promptB = &survey.Confirm{
                Message: "Prune component parameters?",
                Default: cSpec.PruneParams,
                Help:    "This removes empty and null parameters from configuration",
            }
            util.FatalErrorCheck("Invalid option", survey.AskOne(promptB, &cSpec.PruneParams), log.Logger)
        }

        util.FatalErrorCheck(
            "Error generating cluster jsonnet file",
            kr8init.GenerateClusterJsonnet(cSpec, cSpec.ClusterOutputDir),
            log.Logger,
        )
    },
}
```

<a name="InitCmd"></a>InitCmd represents the command. Various subcommands are available to initialize different components of kr8.

```go
var InitCmd = &cobra.Command{
    Use:   "init",
    Short: "Initialize kr8+ config repos, components and clusters",
    Long: `kr8+ requires specific directories and exists for its config to work.
This init command helps in creating directory structure for repos, clusters and 
components`,
}
```

<a name="InitComponentCmd"></a>

```go
var InitComponentCmd = &cobra.Command{
    Use:   "component [flags]",
    Short: "Init a new component config file",
    Long:  "Initialize a new component configuration file",
    Run: func(cmd *cobra.Command, args []string) {

        if cmdInitFlags.Interactive {

            prompt := &survey.Input{
                Message: "Enter component directory",
                Default: RootConfig.ComponentDir,
                Help:    "Enter the root directory to store components in",
            }
            util.FatalErrorCheck("Invalid component directory", survey.AskOne(prompt, &RootConfig.ComponentDir), log.Logger)

            prompt = &survey.Input{
                Message: "Enter component name",
                Default: cmdInitFlags.ComponentName,
                Help:    "Enter the name of the component you want to create",
            }
            util.FatalErrorCheck("Invalid component name", survey.AskOne(prompt, &cmdInitFlags.ComponentName), log.Logger)

            promptS := &survey.Select{
                Message: "Select component type",
                Options: []string{"jsonnet", "yml", "tpl", "chart"},
                Help:    "Select the type of component you want to create",
                Default: "jsonnet",
                Description: func(value string, index int) string {
                    switch value {
                    case "jsonnet":
                        return "Use a Jsonnet file to describe the component resources"
                    case "chart":
                        return "Use a Helm chart to describe the component resources"
                    case "yml":
                        return "Use a yml (docker-compose) file to describe the component resources"
                    case "tpl":
                        return "Use a template file to describe the component resources"
                    default:
                        return ""
                    }
                },
            }
            util.FatalErrorCheck("Invalid component type", survey.AskOne(promptS, &cmdInitFlags.ComponentType), log.Logger)
        }
        util.FatalErrorCheck(
            "Error generating component jsonnet",
            kr8init.GenerateComponentJsonnet(cmdInitFlags, RootConfig.ComponentDir),
            log.Logger,
        )
    },
}
```

<a name="InitRepoCmd"></a>Initializes a new kr8\+ configuration repository

Directory tree:

- components/
- clusters/
- lib/
- generated/

```go
var InitRepoCmd = &cobra.Command{
    Use:   "repo [flags] dir",
    Args:  cobra.MinimumNArgs(1),
    Short: "Initialize a new kr8+ config repo",
    Long: `Initialize a new kr8+ config repo by downloading the kr8+ config skeleton repo
and initialize a git repo so you can get started`,
    Run: func(cmd *cobra.Command, args []string) {
        if len(args) == 0 {
            log.Fatal().Msg("Error: no directory specified")
        }
        outDir := args[len(args)-1]
        log.Debug().Msg("Initializing kr8+ config repo in " + outDir)
        if cmdInitFlags.InitUrl != "" {
            util.FatalErrorCheck(
                "Issue fetching repo",
                util.FetchRepoUrl(cmdInitFlags.InitUrl, outDir, !cmdInitFlags.Fetch),
                log.Logger,
            )

            return
        }

        cmdInitOptions := kr8init.Kr8InitOptions{
            InitUrl:       cmdInitFlags.InitUrl,
            ClusterName:   cmdInitFlags.ClusterName,
            ComponentName: "example-component",
            ComponentType: "jsonnet",
            Interactive:   false,
            Fetch:         cmdInitFlags.Fetch,
        }
        clusterOptions := kr8_types.Kr8ClusterSpec{
            Name:               cmdInitFlags.ClusterName,
            PostProcessor:      "",
            GenerateDir:        "generated",
            GenerateShortNames: false,
            PruneParams:        false,
            ClusterOutputDir:   "generated" + "/" + cmdInitFlags.ClusterName,
            EnableCache:        true,
            CompressCache:      true,
        }

        util.FatalErrorCheck(
            "Issue creating cluster.jsonnet",
            kr8init.GenerateClusterJsonnet(clusterOptions, outDir+"/clusters"),
            log.Logger,
        )
        util.FatalErrorCheck(
            "Issue creating example component.jsonnet",
            kr8init.GenerateComponentJsonnet(cmdInitOptions, outDir+"/components"),
            log.Logger,
        )
        util.FatalErrorCheck(
            "Issue creating lib folder",
            kr8init.GenerateLib(cmdInitOptions.Fetch, outDir+"/lib"),
            log.Logger,
        )
        util.FatalErrorCheck(
            "Issue creating Readme.md",
            kr8init.GenerateReadme(outDir, cmdInitOptions, clusterOptions),
            log.Logger,
        )
    },
}
```

<a name="JsonnetCmd"></a>

```go
var JsonnetCmd = &cobra.Command{
    Use:   "jsonnet",
    Short: "Jsonnet utilities",
    Long:  `Utility commands to process jsonnet`,
}
```

<a name="JsonnetRenderCmd"></a>

```go
var JsonnetRenderCmd = &cobra.Command{
    Use:   "render [flags] file [file ...]",
    Short: "Render a jsonnet file",
    Long:  `Render a jsonnet file to JSON or YAML`,

    Args: cobra.MinimumNArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        for _, f := range args {
            err := jvm.JsonnetRender(cmdFlagsJsonnet, f, RootConfig.VMConfig, log.Logger)
            if err != nil {
                log.Fatal().Str("file", f).Err(err).Msg("error rendering jsonnet file")
            }
        }
    },
}
```

<a name="RenderCmd"></a>

```go
var RenderCmd = &cobra.Command{
    Use:   "render",
    Short: "Render files",
    Long:  `Render files in jsonnet or YAML`,
}
```

<a name="RenderHelmCmd"></a>

```go
var RenderHelmCmd = &cobra.Command{
    Use:   "helm",
    Short: "Clean YAML stream from Helm Template output - Reads from Stdin",
    Long:  `Removes Null YAML objects from a YAML stream`,
    Run: func(cmd *cobra.Command, args []string) {
        decoder := yaml.NewYAMLReader(bufio.NewReader(os.Stdin))
        jsa := [][]byte{}
        for {
            bytes, err := decoder.Read()
            if errors.Is(err, io.EOF) {
                break
            } else if err != nil {
                util.FatalErrorCheck("Error decoding yaml stream", err, log.Logger)
            }
            if len(bytes) == 0 {
                continue
            }
            jsonData, err := yaml.ToJSON(bytes)
            util.FatalErrorCheck("Error converting yaml to JSON", err, log.Logger)
            if string(jsonData) == "null" {

                continue
            }
            _, _, err = unstructured.UnstructuredJSONScheme.Decode(jsonData, nil, nil)
            util.FatalErrorCheck("Error handling unstructured JSON", err, log.Logger)
            jsa = append(jsa, jsonData)
        }
        for _, j := range jsa {
            out, err := goyaml.JSONToYAML(j)
            util.FatalErrorCheck("Error encoding JSON to YAML", err, log.Logger)
            fmt.Println("---")
            fmt.Println(string(out))
        }
    },
}
```

<a name="RenderJsonnetCmd"></a>

```go
var RenderJsonnetCmd = &cobra.Command{
    Use:   "jsonnet file [file ...]",
    Short: "Render a jsonnet file",
    Long:  `Render a jsonnet file to JSON or YAML`,

    Args: cobra.MinimumNArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        for _, fileName := range args {
            err := jvm.JsonnetRender(
                types.CmdJsonnetOptions{
                    Prune:         cmdRenderFlags.Prune,
                    ClusterParams: cmdRenderFlags.ClusterParams,
                    Cluster:       cmdRenderFlags.Cluster,
                    Component:     cmdRenderFlags.ComponentName,
                    Format:        cmdRenderFlags.Format,
                    Color:         false,
                    Lint:          cmdRenderFlags.Lint,
                }, fileName, RootConfig.VMConfig, log.Logger)
            if err != nil {
                log.Fatal().Str("filename", fileName).Err(err).Msg("error rendering jsonnet")
            }
        }
    },
}
```

<a name="RootCmd"></a>RootCmd represents the base command when called without any subcommands.

```go
var RootCmd = &cobra.Command{
    Use:   "kr8",
    Short: "A jsonnet-powered config management tool",
    Long:  `An opinionated configuration management tool for Kubernetes Clusters powered by jsonnet`,
}
```

<a name="VersionCmd"></a>Print out versions of packages in use. Chore\(\) \- Updated manually.

```go
var VersionCmd = &cobra.Command{
    Use:   "version",
    Short: "Return the current version of kr8+",
    Long:  `Return the current version of kr8+`,
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println(RootCmd.Use + "+ Version: " + version)
        info, ok := debug.ReadBuildInfo()
        if !ok {
            log.Fatal().Msg("could not read build info")
        }
        stamp := retrieveStamp(info)
        fmt.Printf("  Built with %s on %s\n", stamp.InfoGoCompiler, stamp.InfoBuildTime)
        fmt.Printf("  VCS revision: %s\n", stamp.VCSRevision)
        fmt.Printf("  Go version %s, GOOS %s, GOARCH %s\n", info.GoVersion, stamp.InfoGOOS, stamp.InfoGOARCH)
        fmt.Print("  Dependencies:\n")
        for _, mod := range retrieveDepends(info) {
            fmt.Printf("    %s\n", mod)
        }

    },
}
```

<a name="ConfigureLogger"></a>
## func [ConfigureLogger](<https://github.com:icebergtech/kr8/blob/main/cmd/root.go#L109>)

```go
func ConfigureLogger(debug bool)
```



<a name="Execute"></a>
## func [Execute](<https://github.com:icebergtech/kr8/blob/main/cmd/root.go#L35>)

```go
func Execute(ver string)
```

Execute adds all child commands to the root command sets flags appropriately. This is called by main.main\(\). It only needs to happen once to the rootCmd.

<a name="FormatFile"></a>
## func [FormatFile](<https://github.com:icebergtech/kr8/blob/main/cmd/format.go#L33>)

```go
func FormatFile(filename string, logger zerolog.Logger) error
```

Read, format, and write back a file. github.com/google/go\-jsonnet/formatter is used to format files.

<a name="GenerateCmdClusterListBuilder"></a>
## func [GenerateCmdClusterListBuilder](<https://github.com:icebergtech/kr8/blob/main/cmd/generate.go#L127>)

```go
func GenerateCmdClusterListBuilder(allClusterParams map[string]string) []string
```



<a name="GenerateCommand"></a>
## func [GenerateCommand](<https://github.com:icebergtech/kr8/blob/main/cmd/generate.go#L72>)

```go
func GenerateCommand(cmd *cobra.Command, args []string)
```

This function generates the components for each cluster in parallel. It uses a wait group to ensure that all clusters have been processed before exiting.

<a name="InitConfig"></a>
## func [InitConfig](<https://github.com:icebergtech/kr8/blob/main/cmd/root.go#L135>)

```go
func InitConfig()
```

InitConfig reads in config file and ENV variables if set.

<a name="ProfilingFinalizer"></a>
## func [ProfilingFinalizer](<https://github.com:icebergtech/kr8/blob/main/cmd/root.go#L174>)

```go
func ProfilingFinalizer()
```

Stop profiling and write cpu and memory profiling files if configured.

<a name="ProfilingInitializer"></a>
## func [ProfilingInitializer](<https://github.com:icebergtech/kr8/blob/main/cmd/root.go#L199>)

```go
func ProfilingInitializer()
```

Sets up program profiling.

<a name="CmdFormatOptions"></a>
## type [CmdFormatOptions](<https://github.com:icebergtech/kr8/blob/main/cmd/format.go#L18-L20>)

Contains the paths to include and exclude for a format command.

```go
type CmdFormatOptions struct {
    Recursive bool
}
```

<a name="CmdGenerateOptions"></a>
## type [CmdGenerateOptions](<https://github.com:icebergtech/kr8/blob/main/cmd/generate.go#L21-L30>)

Stores the options for the 'generate' command.

```go
type CmdGenerateOptions struct {
    // Stores the path to the cluster params file
    ClusterParamsFile string
    // Stores the output directory for generated files
    GenerateDir string
    // Stores the filters to apply to clusters and components when generating files
    Filters util.PathFilterOptions
    // Lint Files with jsonnet linter before generating output
    Lint bool
}
```

<a name="CmdGetOptions"></a>
## type [CmdGetOptions](<https://github.com:icebergtech/kr8/blob/main/cmd/get.go#L39-L53>)

Holds the options for the get command.

```go
type CmdGetOptions struct {
    // ClusterParams provides a way to provide cluster params as a single file.
    // This can be combined with --cluster to override the cluster.
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

<a name="CmdRenderOptions"></a>
## type [CmdRenderOptions](<https://github.com:icebergtech/kr8/blob/main/cmd/render.go#L25-L38>)

Contains parameters for the kr8\+ render command.

```go
type CmdRenderOptions struct {
    // Prune null and empty objects from rendered json
    Prune bool
    // Filename to read cluster configuration from
    ClusterParams string
    // Name of the component to render
    ComponentName string
    // Name of the cluster to render
    Cluster string
    // Format of the output (yaml, json or stream)
    Format string
    // Lint Files with jsonnet linter before generating output
    Lint bool
}
```

<a name="CmdRootOptions"></a>
## type [CmdRootOptions](<https://github.com:icebergtech/kr8/blob/main/cmd/root.go#L44-L67>)

Default options that are available to all commands.

```go
type CmdRootOptions struct {
    // kr8+ config base directory
    BaseDir string
    // kr8+ cluster directory
    ClusterDir string
    // kr8+ component directory
    ComponentDir string
    // A config file with kr8+ configuration
    ConfigFile string
    // parallelism - defaults to runtime.GOMAXPROCS(0)
    Parallel int
    // log more information about what kr8+ is doing. Overrides --loglevel
    Debug bool
    // set log level
    LogLevel string
    // enable colorized output (default true). Set to false to disable")
    Color bool
    // contains information to configure jsonnet vm
    VMConfig types.VMConfig
    // Profiling output directory.  Only captured if set.
    ProfilingDir string
    // CPU profiling output file handle.
    ProfilingCPUFile *os.File
}
```

<a name="RootConfig"></a>

```go
var RootConfig CmdRootOptions
```

<a name="Stamp"></a>
## type [Stamp](<https://github.com:icebergtech/kr8/blob/main/cmd/version.go#L14-L21>)



```go
type Stamp struct {
    InfoGoVersion  string
    InfoGoCompiler string
    InfoGOARCH     string
    InfoGOOS       string
    InfoBuildTime  string
    VCSRevision    string
}
```