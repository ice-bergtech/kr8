package cmd

import (
	"github.com/AlecAivazis/survey/v2"
	kr8init "github.com/ice-bergtech/kr8/pkg/kr8_init"
	types "github.com/ice-bergtech/kr8/pkg/types"
	util "github.com/ice-bergtech/kr8/pkg/util"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var cmdInitFlags kr8init.Kr8InitOptions

func init() {
	RootCmd.AddCommand(InitCmd)
	InitCmd.PersistentFlags().BoolVarP(&cmdInitFlags.Interactive,
		"interactive", "i", false,
		"Initialize a resource interactivly")

	InitCmd.AddCommand(InitRepoCmd)
	InitRepoCmd.PersistentFlags().StringVar(&cmdInitFlags.InitUrl,
		"url", "",
		"Source of skeleton directory to create repo from")
	InitRepoCmd.Flags().StringVarP(&cmdInitFlags.ClusterName,
		"name", "o", "cluster-tpl",
		"Cluster name")
	InitRepoCmd.PersistentFlags().BoolVarP(&cmdInitFlags.Fetch,
		"fetch", "f", false,
		"Fetch remote resources")

	InitCmd.AddCommand(InitClusterCmd)
	InitClusterCmd.Flags().StringVarP(&cmdInitFlags.ClusterName,
		"name", "o", "cluster-tpl",
		"Cluster name")

	InitCmd.AddCommand(InitComponentCmd)
	InitComponentCmd.Flags().StringVarP(&cmdInitFlags.ComponentName,
		"name", "o", "component-tpl",
		"Component name")
	InitComponentCmd.Flags().StringVarP(&cmdInitFlags.ComponentType,
		"type", "t", "jsonnet",
		"Component type, one of: [`jsonnet`, `yml`, `chart`]")
}

// InitCmd represents the command.
// Various subcommands are available to initialize different components of kr8.
var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize kr8 config repos, components and clusters",
	Long: `kr8 requires specific directories and exists for its config to work.
This init command helps in creating directory structure for repos, clusters and 
components`,
}

var InitClusterCmd = &cobra.Command{
	Use:   "cluster [flags]",
	Short: "Init a new cluster config file",
	Long:  "Initialize a new cluster configuration file",
	Run: func(cmd *cobra.Command, args []string) {
		cSpec := types.Kr8ClusterSpec{
			Name:               cmdInitFlags.ClusterName,
			ClusterDir:         RootConfig.ClusterDir,
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
				Default: RootConfig.ClusterDir,
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
		util.FatalErrorCheck("Error generating cluster jsonnet file", kr8init.GenerateClusterJsonnet(cSpec, cSpec.ClusterDir))
	},
}

// Initializes a new kr8 configuration repository
//
// Directory tree:
//
//	components/
//
//	clusters/
//
//	lib/
//
//	generated/
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
		// struct for component and clusters
		cmdInitOptions := kr8init.Kr8InitOptions{
			InitUrl:       cmdInitFlags.InitUrl,
			ClusterName:   cmdInitFlags.ClusterName,
			ComponentName: "example-component",
			ComponentType: "jsonnet",
			Interactive:   false,
		}
		clusterOptions := types.Kr8ClusterSpec{
			PostProcessor:      "",
			GenerateDir:        "generated",
			GenerateShortNames: false,
			PruneParams:        false,
			ClusterDir:         "clusters",
			Name:               cmdInitFlags.ClusterName,
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

var InitComponentCmd = &cobra.Command{
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
			util.FatalErrorCheck("Invalid component name", survey.AskOne(prompt, &cmdInitFlags.ComponentName))

			prompt = &survey.Input{
				Message: "Enter component directory",
				Default: RootConfig.ComponentDir,
			}
			util.FatalErrorCheck("Invalid component directory", survey.AskOne(prompt, &RootConfig.ComponentDir))

			promptS := &survey.Select{
				Message: "Select component type",
				Options: []string{"jsonnet", "yml", "tpl", "chart"},
			}
			util.FatalErrorCheck("Invalid component type", survey.AskOne(promptS, &cmdInitFlags.ComponentType))
		}
		kr8init.GenerateComponentJsonnet(cmdInitFlags, RootConfig.ComponentDir)
	},
}
