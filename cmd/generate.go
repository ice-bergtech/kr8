//nolint:gochecknoinits,gochecknoglobals
package cmd

import (
	"strconv"
	"sync"

	"github.com/panjf2000/ants/v2"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	//nolint:exptostd
	"golang.org/x/exp/maps"

	"github.com/ice-bergtech/kr8/pkg/generate"
	"github.com/ice-bergtech/kr8/pkg/types"
	util "github.com/ice-bergtech/kr8/pkg/util"
)

// Stores the options for the 'generate' command.
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

var cmdGenerateFlags CmdGenerateOptions

func init() {
	RootCmd.AddCommand(GenerateCmd)
	GenerateCmd.Flags().StringVarP(&cmdGenerateFlags.ClusterParamsFile,
		"clusterparams", "p", "",
		"provide cluster params as single file - can be combined with --cluster to override cluster")
	GenerateCmd.Flags().StringVarP(&cmdGenerateFlags.Filters.Clusters,
		"clusters", "C", "",
		"clusters to generate - comma separated list of cluster names and/or regular expressions ")
	GenerateCmd.Flags().StringVarP(&cmdGenerateFlags.Filters.Components, "components", "c", "",
		"components to generate - comma separated list of component names and/or regular expressions")
	GenerateCmd.Flags().StringVarP(&cmdGenerateFlags.GenerateDir,
		"generate-dir", "o", "generated",
		"output directory")
	GenerateCmd.Flags().StringVarP(&cmdGenerateFlags.Filters.Includes,
		"clincludes", "i", "",
		"filter included cluster by including clusters with matching cluster parameters - "+
			"comma separate list of key/value conditions separated by = or ~ (for regex match)")
	GenerateCmd.Flags().StringVarP(&cmdGenerateFlags.Filters.Excludes,
		"clexcludes", "x", "",
		"filter included cluster by excluding clusters with matching cluster parameters - "+
			"comma separate list of key/value conditions separated by = or ~ (for regex match)")
	GenerateCmd.Flags().BoolVarP(&cmdGenerateFlags.Lint, "lint", "l", true,
		"lint Files with jsonnet linter before generating output")
}

var GenerateCmd = &cobra.Command{
	Use:     "generate [flags]",
	Aliases: []string{"gen", "g"},
	Short:   "Generate components",
	Long:    `Generate components in clusters`,
	Example: "kr8 generate",

	Args: cobra.MinimumNArgs(0),
	Run:  GenerateCommand,
}

// This function generates the components for each cluster in parallel.
// It uses a wait group to ensure that all clusters have been processed before exiting.
func GenerateCommand(cmd *cobra.Command, args []string) {
	// get list of all clusters, render cluster level params for all of them
	allClusterParams, err := generate.GetClusterParams(
		RootConfig.ClusterDir,
		RootConfig.VMConfig,
		cmdGenerateFlags.Lint,
		log.Logger,
	)
	util.FatalErrorCheck("error getting cluster params from "+RootConfig.ClusterDir, err, log.Logger)

	clusterList := GenerateCmdClusterListBuilder(allClusterParams)

	// Setup the threading pools, one for clusters and one for clusters
	var waitGroup sync.WaitGroup
	ants_cp, _ := ants.NewPool(RootConfig.Parallel)
	ants_cl, _ := ants.NewPool(RootConfig.Parallel)

	kr8Opts := types.Kr8Opts{
		BaseDir:      RootConfig.BaseDir,
		ComponentDir: RootConfig.ComponentDir,
		ClusterDir:   RootConfig.ClusterDir,
	}

	// Generate config for each cluster in parallel
	for _, clusterName := range clusterList {
		waitGroup.Add(1)
		_ = ants_cl.Submit(func() {
			defer waitGroup.Done()
			subLogger := log.With().Str("cluster", clusterName).Logger()

			genFlags := generate.GenerateProcessRootConfig{
				ClusterName:       clusterName,
				ClusterDir:        RootConfig.ClusterDir,
				BaseDir:           RootConfig.BaseDir,
				GenerateDir:       cmdGenerateFlags.GenerateDir,
				Kr8Opts:           kr8Opts,
				ClusterParamsFile: cmdGenerateFlags.ClusterParamsFile,
				Filters:           cmdGenerateFlags.Filters,
				VmConfig:          RootConfig.VMConfig,
				Noop:              false,
				Lint:              cmdGenerateFlags.Lint,
			}

			err = generate.GenProcessCluster(
				&genFlags,
				ants_cp,
				subLogger)
			if err != nil {
				subLogger.Fatal().Err(err).Msg("error processing cluster")
			}
		})
	}
	waitGroup.Wait()
}

func GenerateCmdClusterListBuilder(allClusterParams map[string]string) []string {
	var clusterList []string
	// Filter out and cluster or components we don't want to generate
	if cmdGenerateFlags.Filters.Includes != "" || cmdGenerateFlags.Filters.Excludes != "" {
		clusterList = util.CalculateClusterIncludesExcludes(allClusterParams, cmdGenerateFlags.Filters)
		log.Debug().Msg("Have " + strconv.Itoa(len(clusterList)) + " after filtering")
	} else {
		//nolint:exptostd
		clusterList = maps.Keys(allClusterParams)
	}

	return clusterList
}
