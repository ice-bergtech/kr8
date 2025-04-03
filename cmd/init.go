package cmd

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/hashicorp/go-getter"
	util "github.com/ice-bergtech/kr8/pkg/util"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// cmdInitOptions defines the options used by the init subcommands
type cmdInitOptions struct {
	// URL to fetch the skeleton directory from
	InitUrl string
	// Name of the cluster to initialize
	ClusterName string
	// Name of the component to initialize
	ComponentName string
	// Type of component to initialize (e.g. jsonnet, yml, chart, compose)
	ComponentType string
	// Whether to run in interactive mode or not
	Interactive bool
	// Whether to fetch remote resources or not
	Fetch bool
}

var cmdInitFlags cmdInitOptions

func init() {
	RootCmd.AddCommand(initCmd)
	initCmd.PersistentFlags().BoolVarP(&cmdInitFlags.Interactive, "interactive", "i", false, "Initialize a resource interactivly")

	initCmd.AddCommand(repoCmd)
	repoCmd.PersistentFlags().StringVar(&cmdInitFlags.InitUrl, "url", "", "Source of skeleton directory to create repo from")
	repoCmd.Flags().StringVarP(&cmdInitFlags.ClusterName, "name", "o", "cluster-tpl", "Cluster name")
	repoCmd.PersistentFlags().BoolVarP(&cmdInitFlags.Fetch, "fetch", "f", false, "Fetch remote resources")

	initCmd.AddCommand(initCluster)
	initCluster.Flags().StringVarP(&cmdInitFlags.ClusterName, "name", "o", "cluster-tpl", "Cluster name")

	initCmd.AddCommand(initComponent)
	initComponent.Flags().StringVarP(&cmdInitFlags.ComponentName, "name", "o", "component-tpl", "Component name")
	initComponent.Flags().StringVarP(&cmdInitFlags.ComponentType, "type", "t", "jsonnet", "Component type, one of: [`jsonnet`, `yml`, `chart`]")

}

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize kr8 config repos, components and clusters",
	Long: `kr8 requires specific directories and exists for its config to work.
This init command helps in creating directory structure for repos, clusters and 
components`,
	//Run: func(cmd *cobra.Command, args []string) {},
	// Directory tree:
	//   components/
	//   clusters/
	//   lib/
	//   generated/
}

var initCluster = &cobra.Command{
	Use:   "cluster [flags]",
	Short: "Init a new cluster config file",
	Long:  "Initialize a new cluster configuration file",
	Run: func(cmd *cobra.Command, args []string) {
		cSpec := Kr8ClusterSpec{
			Name:               cmdInitFlags.ClusterName,
			ClusterDir:         rootConfig.ClusterDir,
			PostProcessor:      "function(input) input",
			GenerateDir:        "generated",
			GenerateShortNames: false,
			PruneParams:        false,
		}
		// Get cluster name, path from user if not set
		if cmdInitFlags.Interactive {
			prompt := &survey.Input{
				Message: "Set the cluster name",
				Default: cmdInitFlags.ClusterName,
			}
			survey.AskOne(prompt, &cSpec.Name)

			prompt = &survey.Input{
				Message: "Set the cluster configuration directory",
				Default: rootConfig.ClusterDir,
			}
			survey.AskOne(prompt, &cSpec.ClusterDir)

			promptB := &survey.Confirm{
				Message: "Generate short names for output file names?",
				Default: cSpec.GenerateShortNames,
			}
			survey.AskOne(promptB, &cSpec.GenerateShortNames)

			promptB = &survey.Confirm{
				Message: "Prune component parameters?",
				Default: cSpec.PruneParams,
			}
			survey.AskOne(promptB, &cSpec.PruneParams)

			prompt = &survey.Input{
				Message: "Set the cluster spec post-processor",
				Default: cSpec.PostProcessor,
			}
			survey.AskOne(prompt, &cSpec.PostProcessor)
		}
		// Generate the jsonnet file based on the config
		util.FatalErrorCheck(GenerateClusterJsonnet(cSpec, cSpec.ClusterDir), "Error generating cluster jsonnet file")
	},
}

var repoCmd = &cobra.Command{
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
			fetchRepoUrl(cmdInitFlags.InitUrl, outDir, cmdInitFlags.Fetch)
			return
		}
		// struct for component and clusters
		cmdInitOptions := cmdInitOptions{
			InitUrl:       cmdInitFlags.InitUrl,
			ClusterName:   cmdInitFlags.ClusterName,
			ComponentName: "example-component",
			ComponentType: "jsonnet",
			Interactive:   false,
		}
		clusterOptions := Kr8ClusterSpec{
			PostProcessor:      "",
			GenerateDir:        "generated",
			GenerateShortNames: false,
			PruneParams:        false,
			ClusterDir:         "clusters",
			Name:               cmdInitFlags.ClusterName,
		}
		GenerateClusterJsonnet(clusterOptions, outDir+"/clusters")
		GenerateComponentJsonnet(cmdInitOptions, outDir+"/components")
		GenerateLib(cmdInitFlags.Fetch, outDir+"/lib")
		GenerateReadme(outDir, cmdInitOptions, clusterOptions)
	},
}

var initComponent = &cobra.Command{
	Use:   "component [flags]",
	Short: "Init a new component config file",
	Long:  "Initialize a new component configuration file",
	Run: func(cmd *cobra.Command, args []string) {
		// Get component name, path and type from user if not set
		if cmdInitFlags.Interactive {
			prompt := &survey.Input{
				Message: "Enter component name",
				Default: cmdInitFlags.ComponentName,
			}
			survey.AskOne(prompt, &cmdInitFlags.ComponentName)

			prompt = &survey.Input{
				Message: "Enter component directory",
				Default: rootConfig.ComponentDir,
			}
			survey.AskOne(prompt, &rootConfig.ComponentDir)

			promptS := &survey.Select{
				Message: "Select component type",
				Options: []string{"jsonnet", "yml", "tpl", "chart"},
			}
			survey.AskOne(promptS, &cmdInitFlags.ComponentType)
		}
		GenerateComponentJsonnet(cmdInitFlags, rootConfig.ComponentDir)
	},
}

// Write out a struct to a specified path and file
func writeObjToJsonFile(filename string, path string, objStruct interface{}) error {
	util.FatalErrorCheck(os.MkdirAll(path, 0755), "error creating resource directory")

	jsonStr, errJ := json.MarshalIndent(objStruct, "", "  ")
	util.FatalErrorCheck(errJ, "error marshalling component resource to json")

	jsonStrFormatted, errF := formatJsonnetString(string(jsonStr))
	util.FatalErrorCheck(errF, "error formatting component resource to json")

	return (os.WriteFile(path+"/"+filename, []byte(jsonStrFormatted), 0644))
}

// Generate a cluster.jsonnet file based on the provided Kr8ClusterSpec and store it in the specified directory.
func GenerateClusterJsonnet(cSpec Kr8ClusterSpec, dstDir string) error {
	filename := "cluster.jsonnet"
	clusterJson := Kr8ClusterJsonnet{
		ClusterSpec: cSpec,
		Cluster:     kr8Cluster{Name: cSpec.Name},
		Components:  map[string]Kr8ClusterComponentRef{},
	}
	return writeObjToJsonFile(filename, dstDir+"/"+cSpec.Name, clusterJson)
}

// Generate default component kr8_spec values and store in params.jsonnet
// Based on the type:
// jsonnet: create a component.jsonnet file and reference it from the params.jsonnet file
// yml: leave a note in the params.jsonnet file about where and how the yml files can be referenced
// chart: generate a simple taskfile that handles vendoring the chart data
func GenerateComponentJsonnet(componentOptions cmdInitOptions, dstDir string) error {

	compJson := Kr8ComponentJsonnet{
		Kr8Spec: Kr8ComponentSpec{
			Kr8_allparams:         false,
			Kr8_allclusters:       false,
			DisableOutputDirClean: false,
			Includes:              []interface{}{},
			ExtFiles:              map[string]string{},
			JPaths:                []string{},
		},
		ReleaseName: strings.ReplaceAll(componentOptions.ComponentName, "_", "-"),
		Namespace:   "Default",
		Version:     "1.0.0",
	}
	switch componentOptions.ComponentType {
	case "jsonnet":
		compJson.Kr8Spec.Includes = append(compJson.Kr8Spec.Includes, "component.jsonnet")
	case "yml":
		compJson.Kr8Spec.Includes = append(compJson.Kr8Spec.Includes, Kr8ComponentSpecIncludeObject{File: "input.yml", DestName: "glhf"})
	case "tpl":
		compJson.Kr8Spec.Includes = append(compJson.Kr8Spec.Includes, Kr8ComponentSpecIncludeObject{File: "README.tpl", DestDir: "docs", DestExt: "md"})
	case "chart":
		break
	default:
		break
	}

	return writeObjToJsonFile("params.jsonnet", dstDir+"/"+componentOptions.ComponentName, compJson)
}

// Fetch a git repo from a url and clone it to a destination directory
// if the performFetch flag is false, it will log the command that would be run and return without doing anything
func fetchRepoUrl(url string, destination string, performFetch bool) {
	if !performFetch {
		gitCommand := "git clone -- " + url + " " + destination
		cleanupCmd := "rm -rf \"" + destination + "/.git\""
		log.Info().Msg("Fetch disabled. Would have ran: ")
		log.Info().Msg("`" + gitCommand + "`")
		log.Info().Msg("`" + cleanupCmd + "`")
		return
	}

	// Get the current working directory
	pwd, err := os.Getwd()
	util.FatalErrorCheck(err, "Error getting working directory")

	// Download the skeletion directory
	log.Debug().Msg("Downloading skeleton repo from git::" + url)
	client := &getter.Client{
		Src:  "git::" + url,
		Dst:  destination,
		Pwd:  pwd,
		Mode: getter.ClientModeAny,
	}

	util.FatalErrorCheck(client.Get(), "Error getting repo")

	// Check for .git folder
	if _, err := os.Stat(destination + "/.git"); !os.IsNotExist(err) {
		log.Debug().Msg("Removing .git directory")
		os.RemoveAll(destination + "/.git")
	}
}

func GenerateLib(fetch bool, dstDir string) {
	util.FatalErrorCheck(os.MkdirAll(dstDir, 0755), "error creating lib directory")
	fetchRepoUrl("https://github.com/kube-libsonnet/kube-libsonnet.git", dstDir+"/klib", fetch)
}

func GenerateReadme(dstDir string, cmdOptions cmdInitOptions, clusterSpec Kr8ClusterSpec) {
	type templateVars struct {
		Cmd     cmdInitOptions
		Cluster Kr8ClusterSpec
	}
	var fetch string
	if cmdOptions.Fetch {
		fetch = "true"
	} else {
		fetch = "false"
	}

	readmeTemplate := strings.Join([]string{
		"# Stack " + cmdOptions.ClusterName + " Readme",
		"",
		"## Project Overview",
		"",
		"This project is a cluster stack initialized by kr8+",
		"",
		"* Generate and customize component configuration for Kubernetes clusters across environments, regions and platforms",
		"* Opinionated config, flexible deployment. kr8+ simply generates manifests for you, you decide how to deploy them",
		"* Render and override component config from multiple sources",
		"  * Helm, Kustomize, Static manifests, raw configuration",
		"* Generate static configuration across clusters that is CI/CD friendly",
		"  * Kubernetes manifests, Helm charts, Kustomize overlays, Documentation, text files",
		"",
		"## Usage",
		"",
		"1. Define components in the `components` directory.",
		"2. Define tiered cluster configuration in the `" + clusterSpec.ClusterDir + "` directory.",
		"3. Run `kr8 generate` to generate component configuration files.",
		"",
		"## Info ",
		"",
		"This project is initialized with the following parameters:",
		"",
		"	* ClusterName: `" + cmdOptions.ClusterName + "`",
		"	* Fetch External Libs: " + fetch,
		"   * Cluster config root directory: `" + clusterSpec.ClusterDir + "`",
		"   * Component root directory: `components`",
		"   * Cluster config root directory: `" + clusterSpec.ClusterDir + "`",
		"   * Generated config outpu directory: `" + clusterSpec.GenerateDir + "`",
		"",
		"Generated using [kr8+](https://github.com/ice-bergtech/kr8) V `" + Version + "`",
	}, "\n")

	os.WriteFile(dstDir+"/Readme.md", []byte(readmeTemplate), 0644)
}
