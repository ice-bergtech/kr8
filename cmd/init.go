package cmd

import (
	"github.com/AlecAivazis/survey/v2"
	kr8init "github.com/ice-bergtech/kr8/pkg/kr8_init"
	"github.com/ice-bergtech/kr8/pkg/kr8_types"
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
		cSpec := kr8_types.Kr8ClusterSpec{
			Name:               cmdInitFlags.ClusterName,
			PostProcessor:      "function(input) input",
			GenerateDir:        "generated",
			GenerateShortNames: false,
			PruneParams:        false,
			ClusterOutputDir:   RootConfig.ClusterDir,
		}

		if cmdInitFlags.Interactive {
			prompt := &survey.Input{
				Message: "Set the cluster configuration directory",
				Default: RootConfig.ClusterDir,
				Help:    "Set the root directory to store cluster configurations, optionally including subdirectories",
			}
			util.FatalErrorCheck("Invalid cluster directory", survey.AskOne(prompt, &cSpec.ClusterOutputDir), log.Logger)

			// Get cluster name, path from user if not set
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
		// Generate the jsonnet file based on the config
		util.FatalErrorCheck(
			"Error generating cluster jsonnet file",
			kr8init.GenerateClusterJsonnet(cSpec, cSpec.ClusterOutputDir),
			log.Logger,
		)
	},
}

// Initializes a new kr8 configuration repository
//
// Directory tree:
//   - components/
//   - clusters/
//   - lib/
//   - generated/
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
				util.FetchRepoUrl(cmdInitFlags.InitUrl, outDir, !cmdInitFlags.Fetch),
				log.Logger,
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
			Fetch:         false,
		}
		clusterOptions := kr8_types.Kr8ClusterSpec{
			Name:               cmdInitFlags.ClusterName,
			PostProcessor:      "",
			GenerateDir:        "generated",
			GenerateShortNames: false,
			PruneParams:        false,
			ClusterOutputDir:   "generated" + "/" + cmdInitFlags.ClusterName,
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
			kr8init.GenerateLib(cmdInitFlags.Fetch, outDir+"/lib"),
			log.Logger,
		)
		util.FatalErrorCheck(
			"Issue creating Readme.md",
			kr8init.GenerateReadme(outDir, cmdInitOptions, clusterOptions),
			log.Logger,
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
