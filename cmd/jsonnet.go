//nolint:gochecknoinits,gochecknoglobals
package cmd

import (
	jvm "github.com/ice-bergtech/kr8/pkg/jnetvm"
	types "github.com/ice-bergtech/kr8/pkg/types"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// Create Jsonnet VM. Configure with env vars and command line flags

var JsonnetCmd = &cobra.Command{
	Use:   "jsonnet",
	Short: "Jsonnet utilities",
	Long:  `Utility commands to process jsonnet`,
}

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

var cmdFlagsJsonnet types.CmdJsonnetOptions

func init() {
	RootCmd.AddCommand(JsonnetCmd)
	JsonnetCmd.AddCommand(JsonnetRenderCmd)
	JsonnetRenderCmd.PersistentFlags().BoolVarP(&cmdFlagsJsonnet.Prune,
		"prune", "", true,
		"removes null and empty objects from ingested jsonnet files")
	JsonnetRenderCmd.PersistentFlags().StringVarP(&cmdFlagsJsonnet.ClusterParams,
		"clusterparams", "p", "",
		"provide cluster params as single file - can be combined with --cluster to override cluster")
	JsonnetRenderCmd.PersistentFlags().StringVarP(&cmdFlagsJsonnet.Component,
		"component", "c", "",
		"component to render params for")
	JsonnetRenderCmd.PersistentFlags().StringVarP(&cmdFlagsJsonnet.Format,
		"format", "F", "json",
		"output format: json, yaml, stream")
	JsonnetRenderCmd.PersistentFlags().StringVarP(&cmdFlagsJsonnet.Cluster,
		"cluster", "C", "",
		"cluster to render params for")
}
