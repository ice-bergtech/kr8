//nolint:gochecknoinits,gochecknoglobals
package cmd

import (
	"path/filepath"

	jvm "github.com/ice-bergtech/kr8/pkg/jnetvm"
	types "github.com/ice-bergtech/kr8/pkg/types"
	util "github.com/ice-bergtech/kr8/pkg/util"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// Create Jsonnet VM. Configure with env vars and command line flags

var JsonnetCmd = &cobra.Command{
	Use:     "jsonnet",
	Short:   "Jsonnet utilities",
	Aliases: []string{"j"},
	Long:    `Utility commands to process jsonnet`,
}

var JsonnetRenderCmd = &cobra.Command{
	Use:     "render [flags] file [file ...]",
	Aliases: []string{"r"},
	Short:   "Render a jsonnet file",
	Long:    `Render a jsonnet file to JSON or YAML`,

	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cleanedArgs := make([]string, len(args))
		for i, arg := range args {
			cleanedArgs[i] = filepath.Clean(arg)
		}

		vmConfig := types.VMConfig{}
		output, err := jvm.JsonnetRenderFiles(vmConfig, cleanedArgs, "", cmdFlagsJsonnet.Prune, "", "cli", cmdFlagsJsonnet.Lint)
		if err != nil {
			log.Fatal().Err(err).Msg("error rendering jsonnet file")
		}

		if cmdFlagsJsonnet.Output != "" {
			err := util.WriteFile([]byte(output), cmdFlagsJsonnet.Output)
			if err != nil {
				log.Fatal().Err(err).Msg("error writing output file")
			}
			log.Info().Str("file", cmdFlagsJsonnet.Output).Msg("output written to file")
		} else {
			log.Info().Msg(output)
		}
	},
}

type CmdJsonnetRenderOptions struct {
	Prune  bool
	Format string
	Lint   bool
	Output string
}

var cmdFlagsJsonnet CmdJsonnetRenderOptions

func init() {
	RootCmd.AddCommand(JsonnetCmd)
	JsonnetCmd.AddCommand(JsonnetRenderCmd)
	JsonnetRenderCmd.PersistentFlags().BoolVarP(&cmdFlagsJsonnet.Prune,
		"prune", "", true,
		"removes null and empty objects from ingested jsonnet files")
	JsonnetRenderCmd.PersistentFlags().StringVarP(&cmdFlagsJsonnet.Format,
		"format", "F", "json",
		"output format: json, yaml, stream")
	JsonnetRenderCmd.Flags().BoolVarP(&cmdFlagsJsonnet.Lint, "lint", "l", true,
		"lint Files with jsonnet linter before generating output")
	JsonnetRenderCmd.Flags().StringVarP(&cmdFlagsJsonnet.Output, "output", "o", "",
		"write output to file instead of stdout")
}
