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
)

// Contains parameters for the kr8 render command
type cmdRenderOptions struct {
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
var cmdRenderFlags cmdRenderOptions

func init() {
	RootCmd.AddCommand(renderCmd)

	renderCmd.AddCommand(renderJsonnetCmd)
	renderJsonnetCmd.PersistentFlags().BoolVarP(&cmdRenderFlags.Prune, "prune", "", true, "Prune null and empty objects from rendered json")
	renderJsonnetCmd.PersistentFlags().StringVarP(&cmdRenderFlags.ClusterParams, "clusterparams", "p", "", "provide cluster params as single file - can be combined with --cluster to override cluster")
	renderJsonnetCmd.PersistentFlags().StringVarP(&cmdRenderFlags.ComponentName, "component", "c", "", "component to render params for")
	renderJsonnetCmd.PersistentFlags().StringVarP(&cmdRenderFlags.Cluster, "cluster", "C", "", "cluster to render params for")
	renderJsonnetCmd.PersistentFlags().StringVarP(&cmdRenderFlags.Format, "format", "F", "json", "Output format: json, yaml, stream")

	renderCmd.AddCommand(helmCleanCmd)
}

var renderCmd = &cobra.Command{
	Use:   "render",
	Short: "Render files",
	Long:  `Render files in jsonnet or YAML`,
}

var renderJsonnetCmd = &cobra.Command{
	Use:   "jsonnet file [file ...]",
	Short: "Render a jsonnet file",
	Long:  `Render a jsonnet file to JSON or YAML`,

	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		for _, f := range args {
			jsonnetRender(
				CmdJsonnetOptions{
					Prune:         cmdRenderFlags.Prune,
					ClusterParams: cmdRenderFlags.ClusterParams,
					Cluster:       cmdRenderFlags.Cluster,
					Component:     cmdRenderFlags.ComponentName,
					Format:        cmdRenderFlags.Format,
				}, f, rootConfig.VMConfig)
		}
	},
}

var helmCleanCmd = &cobra.Command{
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
				fatalErrorCheck(err, "Error decoding decoding yaml stream")
			}
			if len(bytes) == 0 {
				continue
			}
			jsonData, err := yaml.ToJSON(bytes)
			fatalErrorCheck(err, "Error converting yaml to JSON")
			if string(jsonData) == "null" {
				// skip empty json
				continue
			}
			_, _, err = unstructured.UnstructuredJSONScheme.Decode(jsonData, nil, nil)
			fatalErrorCheck(err, "Error handling unstructured JSON")
			jsa = append(jsa, jsonData)
		}
		for _, j := range jsa {
			out, err := goyaml.JSONToYAML(j)
			fatalErrorCheck(err, "Error encoding JSON to YAML")
			fmt.Println("---")
			fmt.Println(string(out))
		}
	},
}
