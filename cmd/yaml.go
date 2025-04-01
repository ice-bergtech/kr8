package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"

	gyaml "github.com/ghodss/yaml"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/yaml"
)

var yamlCmd = &cobra.Command{
	Use:   "yaml",
	Short: "YAML utilities",
	Long:  `Utility commands to process YAML`,
}

var yamlHelmCleanCmd = &cobra.Command{
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
			out, err := gyaml.JSONToYAML(j)
			if err != nil {
				log.Fatal().Err(err).Msg("Error encoding JSON to YAML")
			}
			fmt.Println("---")
			fmt.Println(string(out))
		}
	},
}

func init() {
	RootCmd.AddCommand(yamlCmd)
	yamlCmd.AddCommand(yamlHelmCleanCmd)

}
