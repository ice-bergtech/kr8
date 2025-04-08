# cmd

```go
import "github.com/ice-bergtech/kr8/cmd"
```

## Index

- [Variables](<#variables>)
- [func ConfigureLogger\(debug bool\)](<#ConfigureLogger>)
- [func Execute\(ver string\)](<#Execute>)
- [func GenerateCommand\(cmd \*cobra.Command, args \[\]string\)](<#GenerateCommand>)
- [func InitConfig\(\)](<#InitConfig>)
- [type CmdGenerateOptions](<#CmdGenerateOptions>)
- [type CmdGetOptions](<#CmdGetOptions>)
- [type CmdRenderOptions](<#CmdRenderOptions>)
- [type CmdRootOptions](<#CmdRootOptions>)
- [type Stamp](<#Stamp>)


## Variables

<a name="FormatCmd"></a>

```go
var FormatCmd = &cobra.Command{
    Use:   "format [flags]",
    Short: "Format jsonnet files",
    Long:  `Format jsonnet configuration files`,

    Args: cobra.MinimumNArgs(0),
    Run: func(cmd *cobra.Command, args []string) {
        // First get a list of all files in the base directory and subdirectories. Ignore .git directories.
        var fileList []string
        err := filepath.Walk(RootConfig.BaseDir, func(path string, info fs.FileInfo, err error) error {
            if info.IsDir() {
                if info.Name() == ".git" {
                    return filepath.SkipDir
                }

                return nil
            }
            fileList = append(fileList, path)

            return nil
        })
        util.FatalErrorCheck("Error walking the path "+RootConfig.BaseDir, err)

        fileList = util.Filter(fileList, func(s string) bool {
            var result bool
            for _, f := range strings.Split(cmdFormatFlags.Includes, ",") {
                t, _ := filepath.Match(f, s)
                if t {
                    return t
                }
                result = result || t
            }

            return result
        })

        fileList = util.Filter(fileList, func(s string) bool {
            var result bool
            for _, f := range strings.Split(cmdFormatFlags.Excludes, ",") {
                t, _ := filepath.Match(f, s)
                if t {
                    return !t
                }
                result = result || t
            }

            return !result
        })
        log.Debug().Msg("Filtered file list: " + fmt.Sprintf("%v", fileList))
        log.Debug().Msg("Formatting files...")

        var waitGroup sync.WaitGroup
        parallel, err := cmd.Flags().GetInt("parallel")
        util.FatalErrorCheck("Error getting parallel flag", err)
        log.Debug().Msg("Parallel set to " + strconv.Itoa(parallel))
        ants_file, _ := ants.NewPool(parallel)
        for _, filename := range fileList {
            waitGroup.Add(1)
            _ = ants_file.Submit(func() {
                defer waitGroup.Done()
                var bytes []byte
                bytes, err = os.ReadFile(filepath.Clean(filename))
                output, err := formatter.Format(filename, string(bytes), util.GetDefaultFormatOptions())
                if err != nil {
                    fmt.Fprintln(os.Stderr, err.Error())

                    return
                }
                err = os.WriteFile(filepath.Clean(filename), []byte(output), 0600)
                if err != nil {
                    fmt.Fprintln(os.Stderr, err.Error())

                    return
                }
            })
        }
        waitGroup.Wait()
    },
}
```

<a name="GenerateCmd"></a>

```go
var GenerateCmd = &cobra.Command{
    Use:   "generate [flags]",
    Short: "Generate components",
    Long:  `Generate components in clusters`,

    Args: cobra.MinimumNArgs(0),
    Run:  GenerateCommand,
}
```

<a name="GetClustersCmd"></a>

```go
var GetClustersCmd = &cobra.Command{
    Use:   "clusters [flags]",
    Short: "Get all clusters",
    Long:  "Get all clusters defined in kr8 config hierarchy",
    Run: func(cmd *cobra.Command, args []string) {

        clusters, err := util.GetClusterFilenames(RootConfig.ClusterDir)
        util.FatalErrorCheck("Error getting clusters", err)

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
    Short: "Display one or many kr8 resources",
    Long:  `Displays information about kr8 resources such as clusters and components`,
}
```

<a name="GetComponentsCmd"></a>

```go
var GetComponentsCmd = &cobra.Command{
    Use:   "components [flags]",
    Short: "Get all components",
    Long:  "Get all available components defined in the kr8 config hierarchy",
    Run: func(cmd *cobra.Command, args []string) {

        if cmdGetFlags.Cluster == "" && cmdGetFlags.ClusterParams == "" {
            log.Fatal().Msg("Please specify a --cluster name and/or --clusterparams file")
        }

        var params []string
        if cmdGetFlags.Cluster != "" {
            clusterPath := util.GetClusterPaths(RootConfig.ClusterDir, cmdGetFlags.Cluster)
            params = util.GetClusterParamsFilenames(RootConfig.ClusterDir, clusterPath)
        }
        if cmdGetFlags.ClusterParams != "" {
            params = append(params, cmdGetFlags.ClusterParams)
        }

        jvm := jnetvm.JsonnetRenderFiles(RootConfig.VMConfig, params, "._components", true, "", "components")
        if cmdGetFlags.ParamField != "" {
            value := gjson.Get(jvm, cmdGetFlags.ParamField)
            if value.String() == "" {
                log.Fatal().Msg("Error getting param: " + cmdGetFlags.ParamField)
            } else {
                formatted := util.Pretty(jvm, RootConfig.Color)
                fmt.Println(formatted)
            }
        } else {
            formatted := util.Pretty(jvm, RootConfig.Color)
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
    Long:  "Get parameters assigned to clusters and components in the kr8 config hierarchy",
    Run: func(cmd *cobra.Command, args []string) {
        if cmdGetFlags.Cluster == "" {
            log.Fatal().Msg("Please specify a --cluster")
        }

        var cList []string
        if cmdGetFlags.Component != "" {
            cList = append(cList, cmdGetFlags.Component)
        }

        params := jnetvm.JsonnetRenderClusterParams(
            RootConfig.VMConfig,
            cmdGetFlags.Cluster,
            cList,
            cmdGetFlags.ClusterParams,
            true,
        )

        if cmdGetFlags.ParamField == "" {
            if cmdGetFlags.Component != "" {
                result := gjson.Get(params, cmdGetFlags.Component).String()
                fmt.Println(util.Pretty(result, RootConfig.Color))
            } else {
                fmt.Println(util.Pretty(params, RootConfig.Color))
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
        cSpec := types.Kr8ClusterSpec{
            Name:               cmdInitFlags.ClusterName,
            PostProcessor:      "function(input) input",
            GenerateDir:        "generated",
            GenerateShortNames: false,
            PruneParams:        false,
            ClusterDir:         RootConfig.ClusterDir,
        }

        if cmdInitFlags.Interactive {
            prompt := &survey.Input{
                Message: "Set the cluster name",
                Default: cmdInitFlags.ClusterName,
                Help:    "Distinct name for the cluster",
            }
            util.FatalErrorCheck("Invalid cluster name", survey.AskOne(prompt, &cSpec.Name))

            prompt = &survey.Input{
                Message: "Set the cluster configuration directory",
                Default: RootConfig.ClusterDir,
                Help:    "Set the root directory for the new cluster",
            }
            util.FatalErrorCheck("Invalid cluster directory", survey.AskOne(prompt, &cSpec.ClusterDir))

            promptB := &survey.Confirm{
                Message: "Generate short names for output file names?",
                Default: cSpec.GenerateShortNames,
                Help:    "Shortens component names and file structure",
            }
            util.FatalErrorCheck("Invalid option", survey.AskOne(promptB, &cSpec.GenerateShortNames))

            promptB = &survey.Confirm{
                Message: "Prune component parameters?",
                Default: cSpec.PruneParams,
                Help:    "This removes empty and null parameters from configuration",
            }
            util.FatalErrorCheck("Invalid option", survey.AskOne(promptB, &cSpec.PruneParams))
        }

        util.FatalErrorCheck("Error generating cluster jsonnet file", kr8init.GenerateClusterJsonnet(cSpec, cSpec.ClusterDir))
    },
}
```

<a name="InitCmd"></a>InitCmd represents the command. Various subcommands are available to initialize different components of kr8.

```go
var InitCmd = &cobra.Command{
    Use:   "init",
    Short: "Initialize kr8 config repos, components and clusters",
    Long: `kr8 requires specific directories and exists for its config to work.
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
                Message: "Enter component name",
                Default: cmdInitFlags.ComponentName,
                Help:    "Enter the name of the component you want to create",
            }
            util.FatalErrorCheck("Invalid component name", survey.AskOne(prompt, &cmdInitFlags.ComponentName))

            prompt = &survey.Input{
                Message: "Enter component directory",
                Default: RootConfig.ComponentDir,
                Help:    "Enter the directory where you want to create the component",
            }
            util.FatalErrorCheck("Invalid component directory", survey.AskOne(prompt, &RootConfig.ComponentDir))

            promptS := &survey.Select{
                Message: "Select component type",
                Options: []string{"jsonnet", "yml", "tpl", "chart"},
                Help:    "Select the type of component you want to create",
                Default: "jsonnet",
                Description: func(value string, index int) string {
                    switch value {
                    case "jsonnet":
                        return "Use a Jsonnet file to describe the component resources"
                    case "yml":
                        return "Use a yml (docker-compose) file to describe the component resources"
                    case "tpl":
                        return "Use a template file to describe the component resources"
                    case "chart":
                        return "Use a Helm chart to describe the component resources"
                    default:
                        return ""
                    }
                },
            }
            util.FatalErrorCheck("Invalid component type", survey.AskOne(promptS, &cmdInitFlags.ComponentType))
        }
        util.FatalErrorCheck(
            "Error generating component jsonnet",
            kr8init.GenerateComponentJsonnet(cmdInitFlags, RootConfig.ComponentDir),
        )
    },
}
```

<a name="InitRepoCmd"></a>Initializes a new kr8 configuration repository

Directory tree:

```
components/

clusters/

lib/

generated/
```

```go
var InitRepoCmd = &cobra.Command{
    Use:   "repo [flags] dir",
    Args:  cobra.MinimumNArgs(1),
    Short: "Initialize a new kr8 config repo",
    Long: `Initialize a new kr8 config repo by downloading the kr8 config skeleton repo
and initialize a git repo so you can get started`,
    Run: func(cmd *cobra.Command, args []string) {
        if len(args) == 0 {
            log.Fatal().Msg("Error: no directory specified")
        }
        outDir := args[len(args)-1]
        log.Debug().Msg("Initializing kr8 config repo in " + outDir)
        if cmdInitFlags.InitUrl != "" {
            util.FatalErrorCheck(
                "Issue fetching repo",
                util.FetchRepoUrl(cmdInitFlags.InitUrl, outDir, cmdInitFlags.Fetch),
            )

            return
        }

        cmdInitOptions := kr8init.Kr8InitOptions{
            InitUrl:       cmdInitFlags.InitUrl,
            ClusterName:   cmdInitFlags.ClusterName,
            ComponentName: "example-component",
            ComponentType: "jsonnet",
            Interactive:   false,
            Fetch:         false,
        }
        clusterOptions := types.Kr8ClusterSpec{
            Name:               cmdInitFlags.ClusterName,
            PostProcessor:      "",
            GenerateDir:        "generated",
            GenerateShortNames: false,
            PruneParams:        false,
            ClusterDir:         "clusters",
        }

        util.FatalErrorCheck(
            "Issue creating cluster.jsonnet",
            kr8init.GenerateClusterJsonnet(clusterOptions, outDir+"/clusters"),
        )
        util.FatalErrorCheck(
            "Issue creating example component.jsonnet",
            kr8init.GenerateComponentJsonnet(cmdInitOptions, outDir+"/components"),
        )
        util.FatalErrorCheck(
            "Issue creating lib folder",
            kr8init.GenerateLib(cmdInitFlags.Fetch, outDir+"/lib"),
        )
        util.FatalErrorCheck(
            "Issue creating Readme.md",
            kr8init.GenerateReadme(outDir, cmdInitOptions, clusterOptions),
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
            jvm.JsonnetRender(cmdFlagsJsonnet, f, RootConfig.VMConfig)
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
                util.FatalErrorCheck("Error decoding yaml stream", err)
            }
            if len(bytes) == 0 {
                continue
            }
            jsonData, err := yaml.ToJSON(bytes)
            util.FatalErrorCheck("Error converting yaml to JSON", err)
            if string(jsonData) == "null" {

                continue
            }
            _, _, err = unstructured.UnstructuredJSONScheme.Decode(jsonData, nil, nil)
            util.FatalErrorCheck("Error handling unstructured JSON", err)
            jsa = append(jsa, jsonData)
        }
        for _, j := range jsa {
            out, err := goyaml.JSONToYAML(j)
            util.FatalErrorCheck("Error encoding JSON to YAML", err)
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
            jvm.JsonnetRender(
                types.CmdJsonnetOptions{
                    Prune:         cmdRenderFlags.Prune,
                    ClusterParams: cmdRenderFlags.ClusterParams,
                    Cluster:       cmdRenderFlags.Cluster,
                    Component:     cmdRenderFlags.ComponentName,
                    Format:        cmdRenderFlags.Format,
                    Color:         false,
                }, fileName, RootConfig.VMConfig)
        }
    },
}
```

<a name="RootCmd"></a>RootCmd represents the base command when called without any subcommands.

```go
var RootCmd = &cobra.Command{
    Use:   "kr8",
    Short: "Kubernetes config parameter framework",
    Long: `A tool to generate Kubernetes configuration from a hierarchy
	of jsonnet files`,
}
```

<a name="VersionCmd"></a>Print out versions of packages in use. Chore\(\) \- Updated manually.

```go
var VersionCmd = &cobra.Command{
    Use:   "version",
    Short: "Return the current version of kr8+",
    Long:  `return the current version of kr8+`,
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println(RootCmd.Use + "+ Version: " + version)
        info, ok := debug.ReadBuildInfo()
        if !ok {
            panic("Could not read build info")
        }
        stamp := retrieveStamp(info)
        fmt.Printf("  Built with %s on %s\n", stamp.InfoGoCompiler, stamp.InfoBuildTime)
        fmt.Printf("  VCS revision: %s\n", stamp.VCSRevision)
        fmt.Printf("  Go version %s, GOOS %s, GOARCH %s\n", info.GoVersion, stamp.InfoGOOS, stamp.InfoGOARCH)
        fmt.Print("  Dependencies:\n")
        for _, mod := range retrieveDepends(info) {
            fmt.Println("    %s\n", mod)
        }

    },
}
```

<a name="ConfigureLogger"></a>
## func [ConfigureLogger](<https://github.com/ice-bergtech/kr8/blob/main/cmd/root.go#L98>)

```go
func ConfigureLogger(debug bool)
```



<a name="Execute"></a>
## func [Execute](<https://github.com/ice-bergtech/kr8/blob/main/cmd/root.go#L33>)

```go
func Execute(ver string)
```

Execute adds all child commands to the root command sets flags appropriately. This is called by main.main\(\). It only needs to happen once to the rootCmd.

<a name="GenerateCommand"></a>
## func [GenerateCommand](<https://github.com/ice-bergtech/kr8/blob/main/cmd/generate.go#L63>)

```go
func GenerateCommand(cmd *cobra.Command, args []string)
```

This function will generate the components for each cluster in parallel. It uses a wait group to ensure that all clusters have been processed before exiting.

<a name="InitConfig"></a>
## func [InitConfig](<https://github.com/ice-bergtech/kr8/blob/main/cmd/root.go#L137>)

```go
func InitConfig()
```

InitConfig reads in config file and ENV variables if set.

<a name="CmdGenerateOptions"></a>
## type [CmdGenerateOptions](<https://github.com/ice-bergtech/kr8/blob/main/cmd/generate.go#L18-L25>)

Stores the options for the 'generate' command.

```go
type CmdGenerateOptions struct {
    // Stores the path to the cluster params file
    ClusterParamsFile string
    // Stores the output directory for generated files
    GenerateDir string
    // Stores the filters to apply to clusters and components when generating files
    Filters util.PathFilterOptions
}
```

<a name="CmdGetOptions"></a>
## type [CmdGetOptions](<https://github.com/ice-bergtech/kr8/blob/main/cmd/get.go#L38-L52>)

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
## type [CmdRenderOptions](<https://github.com/ice-bergtech/kr8/blob/main/cmd/render.go#L23-L34>)

Contains parameters for the kr8 render command.

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
}
```

<a name="CmdRootOptions"></a>
## type [CmdRootOptions](<https://github.com/ice-bergtech/kr8/blob/main/cmd/root.go#L42-L61>)

Default options that are available to all commands.

```go
type CmdRootOptions struct {
    // kr8 config base directory
    BaseDir string
    // kr8 cluster directory
    ClusterDir string
    // kr8 component directory
    ComponentDir string
    // A config file with kr8 configuration
    ConfigFile string
    // parallelism - defaults to runtime.GOMAXPROCS(0)
    Parallel int
    // log more information about what kr8 is doing. Overrides --loglevel
    Debug bool
    // set log level
    LogLevel string
    // enable colorized output (default true). Set to false to disable")
    Color bool
    // contains ingormation to configure jsonnet vm
    VMConfig types.VMConfig
}
```

<a name="RootConfig"></a>

```go
var RootConfig CmdRootOptions
```

<a name="Stamp"></a>
## type [Stamp](<https://github.com/ice-bergtech/kr8/blob/main/cmd/version.go#L11-L18>)



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