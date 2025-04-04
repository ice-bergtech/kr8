package cmd

import (
	goyaml "github.com/ghodss/yaml"
	"github.com/spf13/cobra"

	"bufio"
	"fmt"
	"io"
	"os"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/yaml"

	jvm "github.com/ice-bergtech/kr8/pkg/jvm"
	types "github.com/ice-bergtech/kr8/pkg/types"
	util "github.com/ice-bergtech/kr8/pkg/util"
)

// Contains parameters for the kr8 render command
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
}

// Stores the render command options
var cmdRenderFlags CmdRenderOptions

func init() {
	RootCmd.AddCommand(RenderCmd)

	RenderCmd.AddCommand(RenderJsonnetCmd)
	RenderJsonnetCmd.PersistentFlags().BoolVarP(&cmdRenderFlags.Prune, "prune", "", true, "Prune null and empty objects from rendered json")
	RenderJsonnetCmd.PersistentFlags().StringVarP(&cmdRenderFlags.ClusterParams, "clusterparams", "p", "", "provide cluster params as single file - can be combined with --cluster to override cluster")
	RenderJsonnetCmd.PersistentFlags().StringVarP(&cmdRenderFlags.ComponentName, "component", "c", "", "component to render params for")
	RenderJsonnetCmd.PersistentFlags().StringVarP(&cmdRenderFlags.Cluster, "cluster", "C", "", "cluster to render params for")
	RenderJsonnetCmd.PersistentFlags().StringVarP(&cmdRenderFlags.Format, "format", "F", "json", "Output format: json, yaml, stream")

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
		for _, f := range args {
			jvm.JsonnetRender(
				types.CmdJsonnetOptions{
					Prune:         cmdRenderFlags.Prune,
					ClusterParams: cmdRenderFlags.ClusterParams,
					Cluster:       cmdRenderFlags.Cluster,
					Component:     cmdRenderFlags.ComponentName,
					Format:        cmdRenderFlags.Format,
				}, f, RootConfig.VMConfig)
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
			if err == io.EOF {
				break
			} else if err != nil {
				util.FatalErrorCheck(err, "Error decoding decoding yaml stream")
			}
			if len(bytes) == 0 {
				continue
			}
			jsonData, err := yaml.ToJSON(bytes)
			util.FatalErrorCheck(err, "Error converting yaml to JSON")
			if string(jsonData) == "null" {
				// skip empty json
				continue
			}
			_, _, err = unstructured.UnstructuredJSONScheme.Decode(jsonData, nil, nil)
			util.FatalErrorCheck(err, "Error handling unstructured JSON")
			jsa = append(jsa, jsonData)
		}
		for _, j := range jsa {
			out, err := goyaml.JSONToYAML(j)
			util.FatalErrorCheck(err, "Error encoding JSON to YAML")
			fmt.Println("---")
			fmt.Println(string(out))
		}
	},
}
