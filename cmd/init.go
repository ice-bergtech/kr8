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

var (
	dl_url          string
	real_url        string
	initClName      string
	initClPath      string
	initCoName      string
	initCoPath      string
	initCoType      string
	initInteractive bool
	//initSkipDocs    bool
)

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
	Use:   "cluster",
	Short: "Init a cluster config file",
	Long:  "Initialize a new cluster configuration file",
	Run: func(cmd *cobra.Command, args []string) {
		cSpec := ClusterSpec{
			Name:               initClName,
			ClusterDir:         initClPath,
			PostProcessor:      "function(input) input",
			GenerateDir:        "generated",
			GenerateShortNames: false,
			PruneParams:        false,
		}
		// Get cluster name, path from user if not set
		if initInteractive {
			prompt := &survey.Input{
				Message: "Set the cluster name",
				Default: initClName,
			}
			survey.AskOne(prompt, &cSpec.Name)

			prompt = &survey.Input{
				Message: "Set the cluster path",
				Default: initClPath,
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
		fatalErrorCheck(generateClusterJsonnet(cSpec), "Error generating cluster jsonnet file")
	},
}

// Write out a struct to a specified path and file
func writeInitializedStruct(filename string, path string, objStruct interface{}) error {
	fatalErrorCheck(os.MkdirAll(componentDir, 0755), "error creating component directory")

	jsonStr, errJ := json.MarshalIndent(objStruct, "", "  ")
	fatalErrorCheck(errJ, "error marshalling component jsonnet to json")

	jsonStrFormatted, errF := formatJsonnetString(string(jsonStr))
	fatalErrorCheck(errF, "error formatting component jsonnet to json")

	return (os.WriteFile(filename, []byte(jsonStrFormatted), 0644))
}

func generateClusterJsonnet(cSpec ClusterSpec) error {
	filename := "cluster.jsonnet"
	clusterJson := ClusterJsonnet{
		ClusterSpec: cSpec,
		Cluster:     Cluster{Name: cSpec.Name},
		Components:  map[string]ClusterComponent{},
	}
	clOutDir := clusterDir + "/" + initClName
	if initClPath != "" {
		clOutDir = initClPath
	}
	return writeInitializedStruct(filename, clOutDir, clusterJson)
}

var initComponent = &cobra.Command{
	Use:   "component",
	Short: "Init a component config file",
	Long:  "Initialize a new component configuration file",
	Run: func(cmd *cobra.Command, args []string) {
		// Get component name, path and type from user if not set
		if initInteractive {
			prompt := &survey.Input{
				Message: "Enter component name",
				Default: initCoName,
			}
			survey.AskOne(prompt, &initCoName)

			prompt = &survey.Input{
				Message: "Enter component path",
				Default: initCoPath,
			}
			survey.AskOne(prompt, &initCoPath)

			promptS := &survey.Select{
				Message: "Select component type",
				Options: []string{"jsonnet", "yml", "tpl", "chart"},
			}
			survey.AskOne(promptS, &initCoType)
		}
		generateComponentJsonnet()
	},
}

// Generate default component kr8_spec values and store in params.jsonnet
// Based on the type:
// jsonnet: create a component.jsonnet file and reference it from the params.jsonnet file
// yml: leave a note in the params.jsonnet file about where and how the yml files can be referenced
// chart: generate a simple taskfile that handles vendoring the chart data
func generateComponentJsonnet() error {

	compJson := ComponentJsonnet{
		Kr8Spec: ComponentSpec{
			Kr8_allparams:         false,
			Kr8_allclusters:       false,
			DisableOutputDirClean: false,
			Includes:              []interface{}{},
			ExtFiles:              map[string]string{},
			JPaths:                []string{},
		},
		ReleaseName: strings.ReplaceAll(initCoName, "_", "-"),
		Namespace:   "Default",
		Version:     "1.0.0",
	}
	switch initCoType {
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

	filename := "params.jsonnet"
	componentDir := clusterDir + "/" + initCoName
	if initCoPath != "" {
		componentDir = initCoPath
	}

	return writeInitializedStruct(filename, componentDir, compJson)
}

var repoCmd = &cobra.Command{
	Use:   "repo dir",
	Args:  cobra.MinimumNArgs(1),
	Short: "Initialize a new kr8 config repo",
	Long: `Initialize a new kr8 config repo by downloading the kr8 config skeleton repo
and initialize a git repo so you can get started`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) < 1 {
			log.Fatal().Msg("Must specify a destination")
		}

		if dl_url != "" {
			real_url = dl_url
		} else {
			log.Fatal().Msg("Must specify a URL")
		}
		// Get the current working directory
		pwd, err := os.Getwd()
		if err != nil {
			log.Fatal().Err(err).Msg("Error getting working directory")
		}

		// Download the skeletion directory
		log.Debug().Msg("Downloading skeleton repo from " + real_url)
		client := &getter.Client{
			Src:  real_url,
			Dst:  args[0],
			Pwd:  pwd,
			Mode: getter.ClientModeAny,
		}

		if err := client.Get(); err != nil {
			log.Fatal().Err(err).Msg("")
			os.Exit(1)
		}

		// Check for .git folder
		if _, err := os.Stat(args[0] + "/.git"); !os.IsNotExist(err) {
			log.Debug().Msg("Removing .git directory")
			os.RemoveAll(args[0] + "/.git")
		}
	},
}

func init() {
	RootCmd.AddCommand(initCmd)
	initCmd.AddCommand(repoCmd)
	initCmd.AddCommand(initCluster)
	initCmd.AddCommand(initComponent)

	initCmd.PersistentFlags().BoolVarP(&initInteractive, "interactive", "i", false, "Initialize a resource interactivly")
	//initCmd.PersistentFlags().BoolVarP(&initSkipDocs, "skip-docs", "s", false, "Skip config doc lines")

	repoCmd.PersistentFlags().StringVar(&dl_url, "url", "", "Source of skeleton directory to create repo from")

	initCluster.Flags().StringVarP(&initClName, "name", "o", "cluster-tpl", "Cluster name")
	initCluster.Flags().StringVarP(&initClPath, "path", "p", "clusters", "Cluster path")

	initComponent.Flags().StringVarP(&initCoName, "name", "o", "component-tpl", "Component name")
	initComponent.Flags().StringVarP(&initCoPath, "path", "p", "components", "Component path")
	initComponent.Flags().StringVarP(&initCoType, "type", "t", "jsonnet", "Component type, one of: [`jsonnet`, `yml`, `chart`]")

}
