package cmd

import (
	jvm "github.com/ice-bergtech/kr8/pkg/jvm"
	types "github.com/ice-bergtech/kr8/pkg/types"
	"github.com/spf13/cobra"
)

// Create Jsonnet VM. Configure with env vars and command line flags

var jsonnetCmd = &cobra.Command{
	Use:   "jsonnet",
	Short: "Jsonnet utilities",
	Long:  `Utility commands to process jsonnet`,
}

var jsonnetRenderCmd = &cobra.Command{
	Use:   "render [flags] file [file ...]",
	Short: "Render a jsonnet file",
	Long:  `Render a jsonnet file to JSON or YAML`,

	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		for _, f := range args {
			jvm.JsonnetRender(cmdFlagsJsonnet, f, rootConfig.VMConfig)
		}
	},
}

var cmdFlagsJsonnet types.CmdJsonnetOptions

func init() {
	RootCmd.AddCommand(jsonnetCmd)
	jsonnetCmd.AddCommand(jsonnetRenderCmd)
	jsonnetRenderCmd.PersistentFlags().BoolVarP(&cmdFlagsJsonnet.Prune, "prune", "", true, "Prune removes null and empty objects from ingested jsonnet files")
	jsonnetRenderCmd.PersistentFlags().StringVarP(&cmdFlagsJsonnet.ClusterParams, "clusterparams", "p", "", "provide cluster params as single file - can be combined with --cluster to override cluster")
	jsonnetRenderCmd.PersistentFlags().StringVarP(&cmdFlagsJsonnet.Component, "component", "c", "", "component to render params for")
	jsonnetRenderCmd.PersistentFlags().StringVarP(&cmdFlagsJsonnet.Format, "format", "F", "json", "Output format: json, yaml, stream")
	jsonnetRenderCmd.PersistentFlags().StringVarP(&cmdFlagsJsonnet.Cluster, "cluster", "C", "", "cluster to render params for")

}
