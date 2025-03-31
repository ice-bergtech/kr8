package cmd

import (
	goyaml "github.com/ghodss/yaml"
	"github.com/spf13/cobra"

	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/rs/zerolog/log"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/yaml"
)

var renderCmd = &cobra.Command{
	Use:   "render",
	Short: "Render files",
	Long:  `Render files in jsonnet or UAML`,
}

var renderJsonnetCmd = &cobra.Command{
	Use:   "jsonnet file [file ...]",
	Short: "Render a jsonnet file",
	Long:  `Render a jsonnet file to JSON or YAML`,

	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		jsonnetRenderCmd.Run(cmd, args)
	},
}

var helmCleanCmd = &cobra.Command{
	Use:   "helmclean",
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
				log.Fatal().Err(err).Msg("Error decoding decoding yaml stream")
			}
			if len(bytes) == 0 {
				continue
			}
			jsonData, err := yaml.ToJSON(bytes)
			if err != nil {
				log.Fatal().Err(err).Msg("Error encoding yaml to JSON")
			}
			if string(jsonData) == "null" {
				// skip empty json
				continue
			}
			_, _, err = unstructured.UnstructuredJSONScheme.Decode(jsonData, nil, nil)
			if err != nil {
				log.Fatal().Err(err).Msg("Error handling unstructured JSON")
			}
			jsa = append(jsa, jsonData)
		}
		for _, j := range jsa {
			out, err := goyaml.JSONToYAML(j)
			if err != nil {
				log.Fatal().Err(err).Msg("Error encoding JSON to YAML")
			}
			fmt.Println("---")
			fmt.Println(string(out))
		}
	},
}

func init() {
	RootCmd.AddCommand(renderCmd)
	renderCmd.AddCommand(renderJsonnetCmd)
	renderJsonnetCmd.PersistentFlags().BoolVarP(&pruneFlag, "prune", "", true, "Prune null and empty objects from rendered json")
	renderJsonnetCmd.PersistentFlags().StringVarP(&clusterParams, "clusterparams", "", "", "provide cluster params as single file - can be combined with --cluster to override cluster")
	renderJsonnetCmd.PersistentFlags().StringVarP(&componentName, "component", "C", "", "component to render params for")
	renderJsonnetCmd.PersistentFlags().StringVarP(&outputFormat, "format", "F", "json", "Output format: json, yaml, stream")
	renderJsonnetCmd.PersistentFlags().StringVarP(&cluster, "cluster", "c", "", "cluster to render params for")
	renderCmd.AddCommand(helmCleanCmd)
}
