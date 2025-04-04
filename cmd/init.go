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
		cSpec := types.Kr8ClusterSpec{
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
		util.FatalErrorCheck(kr8init.GenerateClusterJsonnet(cSpec, cSpec.ClusterDir), "Error generating cluster jsonnet file")
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
			util.FetchRepoUrl(cmdInitFlags.InitUrl, outDir, cmdInitFlags.Fetch)
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
		kr8init.GenerateClusterJsonnet(clusterOptions, outDir+"/clusters")
		kr8init.GenerateComponentJsonnet(cmdInitOptions, outDir+"/components")
		kr8init.GenerateLib(cmdInitFlags.Fetch, outDir+"/lib")
		kr8init.GenerateReadme(outDir, cmdInitOptions, clusterOptions)
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
		kr8init.GenerateComponentJsonnet(cmdInitFlags, rootConfig.ComponentDir)
	},
}
