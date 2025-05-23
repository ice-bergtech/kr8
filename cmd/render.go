//nolint:gochecknoinits,gochecknoglobals
package cmd

import (
	"errors"

	goyaml "github.com/ghodss/yaml"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"bufio"
	"fmt"
	"io"
	"os"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/yaml"

	jvm "github.com/ice-bergtech/kr8/pkg/jnetvm"
	types "github.com/ice-bergtech/kr8/pkg/types"
	util "github.com/ice-bergtech/kr8/pkg/util"
)

// Contains parameters for the kr8+ render command.
type CmdRenderOptions struct {
	// Prune null and empty objects from rendered json
	Prune bool
	// Filename to read cluster configuration from
	ClusterParams string
	// Name of the component to render
	ComponentName string
	// Name of the cluster to render
	Cluster string
	// Format of the output (yaml, json or stream)
	Format string
	// Lint Files with jsonnet linter before generating output
	Lint bool
}

// Stores the render command options.
var cmdRenderFlags CmdRenderOptions

func init() {
	RootCmd.AddCommand(RenderCmd)

	RenderCmd.AddCommand(RenderJsonnetCmd)
	RenderJsonnetCmd.PersistentFlags().BoolVarP(&cmdRenderFlags.Prune,
		"prune", "", true,
		"prune null and empty objects from rendered json")
	RenderJsonnetCmd.PersistentFlags().StringVarP(&cmdRenderFlags.ClusterParams,
		"clusterparams", "p", "",
		"provide cluster params as single file - can be combined with --cluster to override cluster")
	RenderJsonnetCmd.PersistentFlags().StringVarP(&cmdRenderFlags.ComponentName,
		"component", "c", "",
		"component to render params for")
	RenderJsonnetCmd.PersistentFlags().StringVarP(&cmdRenderFlags.Cluster,
		"cluster", "C", "",
		"cluster to render params for")
	RenderJsonnetCmd.PersistentFlags().StringVarP(&cmdRenderFlags.Format,
		"format", "F", "json",
		"output format: json, yaml, stream")
	RenderJsonnetCmd.Flags().BoolVarP(&cmdRenderFlags.Lint, "lint", "l", false,
		"lint Files with jsonnet linter before generating output")

	RenderCmd.AddCommand(RenderHelmCmd)
}

var RenderCmd = &cobra.Command{
	Use:   "render",
	Short: "Render files",
	Long:  `Render files in jsonnet or YAML`,
}

var RenderJsonnetCmd = &cobra.Command{
	Use:   "jsonnet file [file ...]",
	Short: "Render a jsonnet file",
	Long:  `Render a jsonnet file to JSON or YAML`,

	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		for _, fileName := range args {
			err := jvm.JsonnetRender(
				types.CmdJsonnetOptions{
					Prune:         cmdRenderFlags.Prune,
					ClusterParams: cmdRenderFlags.ClusterParams,
					Cluster:       cmdRenderFlags.Cluster,
					Component:     cmdRenderFlags.ComponentName,
					Format:        cmdRenderFlags.Format,
					Color:         false,
					Lint:          cmdRenderFlags.Lint,
				}, fileName, RootConfig.VMConfig, log.Logger)
			if err != nil {
				log.Fatal().Str("filename", fileName).Err(err).Msg("error rendering jsonnet")
			}
		}
	},
}

var RenderHelmCmd = &cobra.Command{
	Use:   "helm",
	Short: "Clean YAML stream from Helm Template output - Reads from Stdin",
	Long:  `Removes Null YAML objects from a YAML stream`,
	Run: func(cmd *cobra.Command, args []string) {
		decoder := yaml.NewYAMLReader(bufio.NewReader(os.Stdin))
		jsa := [][]byte{}
		for {
			bytes, err := decoder.Read()
			if errors.Is(err, io.EOF) {
				break
			} else if err != nil {
				util.FatalErrorCheck("Error decoding yaml stream", err, log.Logger)
			}
			if len(bytes) == 0 {
				continue
			}
			jsonData, err := yaml.ToJSON(bytes)
			util.FatalErrorCheck("Error converting yaml to JSON", err, log.Logger)
			if string(jsonData) == "null" {
				// skip empty json
				continue
			}
			_, _, err = unstructured.UnstructuredJSONScheme.Decode(jsonData, nil, nil)
			util.FatalErrorCheck("Error handling unstructured JSON", err, log.Logger)
			jsa = append(jsa, jsonData)
		}
		for _, j := range jsa {
			out, err := goyaml.JSONToYAML(j)
			util.FatalErrorCheck("Error encoding JSON to YAML", err, log.Logger)
			fmt.Println("---")
			fmt.Println(string(out))
		}
	},
}
