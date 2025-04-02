package cmd

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/hashicorp/go-getter"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

type cmdInitOptions struct {
	InitUrl       string
	ClusterName   string
	ComponentName string
	ComponentType string
	Interactive   bool
	Fetch         bool
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
		cSpec := ClusterSpec{
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
		fatalErrorCheck(generateClusterJsonnet(cSpec, cSpec.ClusterDir), "Error generating cluster jsonnet file")
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
		clusterOptions := ClusterSpec{
			PostProcessor:      "",
			GenerateDir:        "generated",
			GenerateShortNames: false,
			PruneParams:        false,
			ClusterDir:         "clusters",
			Name:               cmdInitFlags.ClusterName,
		}
		generateClusterJsonnet(clusterOptions, outDir+"/clusters")
		generateComponentJsonnet(cmdInitOptions, outDir+"/components")
		generateLib(cmdInitFlags.Fetch, outDir+"/lib")
		generateReadme(outDir, cmdInitOptions, clusterOptions)
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
		generateComponentJsonnet(cmdInitFlags, rootConfig.ComponentDir)
	},
}

// Write out a struct to a specified path and file
func writeInitializedStruct(filename string, path string, objStruct interface{}) error {
	fatalErrorCheck(os.MkdirAll(path, 0755), "error creating resource directory")

	jsonStr, errJ := json.MarshalIndent(objStruct, "", "  ")
	fatalErrorCheck(errJ, "error marshalling component resource to json")

	jsonStrFormatted, errF := formatJsonnetString(string(jsonStr))
	fatalErrorCheck(errF, "error formatting component resource to json")

	return (os.WriteFile(path+"/"+filename, []byte(jsonStrFormatted), 0644))
}

func generateClusterJsonnet(cSpec ClusterSpec, dstDir string) error {
	filename := "cluster.jsonnet"
	clusterJson := ClusterJsonnet{
		ClusterSpec: cSpec,
		Cluster:     Cluster{Name: cSpec.Name},
		Components:  map[string]ClusterComponent{},
	}
	return writeInitializedStruct(filename, dstDir+"/"+cSpec.Name, clusterJson)
}

// Generate default component kr8_spec values and store in params.jsonnet
// Based on the type:
// jsonnet: create a component.jsonnet file and reference it from the params.jsonnet file
// yml: leave a note in the params.jsonnet file about where and how the yml files can be referenced
// chart: generate a simple taskfile that handles vendoring the chart data
func generateComponentJsonnet(componentOptions cmdInitOptions, dstDir string) error {

	compJson := ComponentJsonnet{
		Kr8Spec: ComponentSpec{
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
		compJson.Kr8Spec.Includes = append(compJson.Kr8Spec.Includes, IncludeFileEntryStruct{File: "input.yml", DestName: "glhf"})
	case "tpl":
		compJson.Kr8Spec.Includes = append(compJson.Kr8Spec.Includes, IncludeFileEntryStruct{File: "README.tpl", DestDir: "docs", DestExt: "md"})
	case "chart":
		break
	default:
		break
	}

	return writeInitializedStruct("params.jsonnet", dstDir+"/"+componentOptions.ComponentName, compJson)
}

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
	fatalErrorCheck(err, "Error getting working directory")

	// Download the skeletion directory
	log.Debug().Msg("Downloading skeleton repo from " + url)
	client := &getter.Client{
		Src:  url,
		Dst:  destination,
		Pwd:  pwd,
		Mode: getter.ClientModeAny,
	}

	fatalErrorCheck(client.Get(), "Error getting repo")

	// Check for .git folder
	if _, err := os.Stat(destination + "/.git"); !os.IsNotExist(err) {
		log.Debug().Msg("Removing .git directory")
		os.RemoveAll(destination + "/.git")
	}
}

func generateLib(fetch bool, dstDir string) {
	fatalErrorCheck(os.MkdirAll(dstDir, 0755), "error creating lib directory")
	fetchRepoUrl("git::https://github.com/kube-libsonnet/kube-libsonnet.git", dstDir+"/klib", fetch)
}

func generateReadme(dstDir string, cmdOptions cmdInitOptions, clusterSpec ClusterSpec) {

}
