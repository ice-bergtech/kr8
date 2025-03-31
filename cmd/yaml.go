// Copyright © 2018 Lee Briggs <lee@leebriggs.co.uk>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

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
