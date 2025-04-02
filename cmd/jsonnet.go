package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	goyaml "github.com/ghodss/yaml"
	jsonnet "github.com/google/go-jsonnet"
	jsonnetAst "github.com/google/go-jsonnet/ast"
	"github.com/grafana/tanka/pkg/helm"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// Create Jsonnet VM. Configure with env vars and command line flags
/*

This code was originally copied almost verbatim from the kubecfg project: https://github.com/ksonnet/kubecfg

Copyright 2018 ksonnet

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.

*/

type VMConfig struct {
	// VMConfig is a configuration for the Jsonnet VM
	Jpaths  []string `json:"jpath" yaml:"jpath"`
	ExtVars []string `json:"ext_str_file" yaml:"ext_str_files"`
}

func JsonnetVM(vmconfig VMConfig) (*jsonnet.VM, error) {
	vm := jsonnet.MakeVM()
	RegisterNativeFuncs(vm)

	// always add lib directory in base directory to path
	jpath := []string{rootConfig.BaseDir + "/lib"}

	jpath = append(jpath, filepath.SplitList(os.Getenv("KR8_JPATH"))...)
	jpathArgs := vmconfig.Jpaths
	jpath = append(jpath, jpathArgs...)

	vm.Importer(&jsonnet.FileImporter{
		JPaths: jpath,
	})

	for _, extvar := range vmconfig.ExtVars {
		kv := strings.SplitN(extvar, "=", 2)
		if len(kv) != 2 {
			log.Fatal().Str("ext-str-file", extvar).Msg("Failed to parse. Missing '=' in parameter`")
		}
		v, err := os.ReadFile(kv[1])
		if err != nil {
			panic(err)
		}
		vm.ExtVar(kv[0], string(v))
	}
	return vm, nil
}

// Takes a list of jsonnet files and imports each one and mixes them with "+"
func renderJsonnet(vmConfig VMConfig, files []string, param string, prune bool, prepend string, source string) string {

	// copy the slice so that we don't unitentionally modify the original
	jsonnetPaths := make([]string, len(files[:0]))
	copy(jsonnetPaths, files[:0])

	// range through the files
	for _, s := range files {
		jsonnetPaths = append(jsonnetPaths, fmt.Sprintf("(import '%s')", s))
	}

	// Create a JSonnet VM
	vm, err := JsonnetVM(vmConfig)
	fatalErrorCheck(err, "Error creating jsonnet VM")

	// Join the slices into a jsonnet compat string. Prepend code from "prepend" variable, if set.
	var jsonnetImport string
	if prepend != "" {
		jsonnetImport = prepend + "+" + strings.Join(jsonnetPaths, "+")
	} else {
		jsonnetImport = strings.Join(jsonnetPaths, "+")
	}

	if param != "" {
		jsonnetImport = "(" + jsonnetImport + ")" + param
	}

	if prune {
		// wrap in std.prune, to remove nulls, empty arrays and hashes
		jsonnetImport = "std.prune(" + jsonnetImport + ")"
	}

	// render the jsonnet
	out, err := vm.EvaluateAnonymousSnippet(source, jsonnetImport)
	fatalErrorCheck(err, "Error evaluating jsonnet snippet")

	return out

}

// Native Jsonnet funcs to add
/*

This code is copied almost verbatim from the kubecfg project: https://github.com/ksonnet/kubecfg

Copyright 2018 ksonnet

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.

*/

func RegisterNativeFuncs(vm *jsonnet.VM) {
	// Adds on to functions described here: https://jsonnet.org/ref/stdlib.html

	vm.NativeFunction(&jsonnet.NativeFunction{
		Name:   "template",
		Params: []jsonnetAst.Identifier{"config", "str"},
		Func: func(args []interface{}) (res interface{}, err error) {
			var config any
			err = json.Unmarshal([]byte(args[0].(string)), config)
			if err != nil {
				return "", err
			}

			input := []byte(args[1].(string))
			tmpl, err := template.New("file").Parse(string(input))
			if err != nil {
				return "", err
			}

			var buff bytes.Buffer
			err = tmpl.Execute(&buff, config)
			return buff.String(), err
		},
	})

	vm.NativeFunction(&jsonnet.NativeFunction{
		Name:   "escapeStringRegex",
		Params: []jsonnetAst.Identifier{"str"},
		Func: func(args []interface{}) (res interface{}, err error) {
			return regexp.QuoteMeta(args[0].(string)), nil
		},
	})

	vm.NativeFunction(&jsonnet.NativeFunction{
		Name:   "regexMatch",
		Params: []jsonnetAst.Identifier{"regex", "string"},
		Func: func(args []interface{}) (res interface{}, err error) {
			return regexp.MatchString(args[0].(string), args[1].(string))
		},
	})

	vm.NativeFunction(&jsonnet.NativeFunction{
		Name:   "regexSubst",
		Params: []jsonnetAst.Identifier{"regex", "src", "repl"},
		Func: func(args []interface{}) (res interface{}, err error) {
			regex := args[0].(string)
			src := args[1].(string)
			repl := args[2].(string)

			r, err := regexp.Compile(regex)
			if err != nil {
				return "", err
			}
			return r.ReplaceAllString(src, repl), nil
		},
	})

	vm.NativeFunction(helm.NativeFunc(helm.ExecHelm{}))

}

var jsonnetCmd = &cobra.Command{
	Use:   "jsonnet",
	Short: "Jsonnet utilities",
	Long:  `Utility commands to process jsonnet`,
}

func jsonnetRender(cmdFlagsJsonnet CmdJsonnetOptions, filename string, vmConfig VMConfig) {

	if cmdFlagsJsonnet.Cluster == "" && cmdFlagsJsonnet.ClusterParams == "" {
		log.Fatal().Msg("Please specify a --cluster name and/or --clusterparams")
	}

	config := renderClusterParams(vmConfig, cmdFlagsJsonnet.Cluster, []string{cmdFlagsJsonnet.Component}, cmdFlagsJsonnet.ClusterParams, false)

	// VM
	vm, _ := JsonnetVM(vmConfig)

	var input string
	// pass component, _cluster and _components as extvars

	vm.ExtCode("kr8_cluster", "std.prune("+config+"._cluster)")
	vm.ExtCode("kr8_components", "std.prune("+config+"._components)")
	vm.ExtCode("kr8", "std.prune("+config+"."+cmdFlagsJsonnet.Component+")")
	vm.ExtCode("kr8_unpruned", config+"."+cmdFlagsJsonnet.Component)

	if cmdFlagsJsonnet.Prune {
		input = "std.prune(import '" + filename + "')"
	} else {
		input = "( import '" + filename + "')"
	}
	j, err := vm.EvaluateAnonymousSnippet("file", input)

	if err != nil {
		log.Fatal().Err(err).Msg("Error evaluating jsonnet snippet")
	}
	switch cmdFlagsJsonnet.Format {
	case "yaml":
		yaml, err := goyaml.JSONToYAML([]byte(j))
		if err != nil {
			log.Fatal().Err(err).Msg("Error converting JSON to YAML")
		}
		fmt.Println(string(yaml))
	case "stream": // output yaml stream
		var o []interface{}
		if err := json.Unmarshal([]byte(j), &o); err != nil {
			log.Fatal().Err(err).Msg("")
		}
		for _, jobj := range o {
			fmt.Println("---")
			buf, err := goyaml.Marshal(jobj)
			if err != nil {
				log.Fatal().Err(err).Msg("")
			}
			fmt.Println(string(buf))
		}
	case "json":
		formatted := Pretty(j, rootConfig.Color)
		fmt.Println(formatted)
	default:
		log.Fatal().Msg("Output format must be json, yaml or stream")
	}
}

var jsonnetRenderCmd = &cobra.Command{
	Use:   "render [flags] file [file ...]",
	Short: "Render a jsonnet file",
	Long:  `Render a jsonnet file to JSON or YAML`,

	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		for _, f := range args {
			jsonnetRender(cmdFlagsJsonnet, f, rootConfig.VMConfig)
		}
	},
}

type CmdJsonnetOptions struct {
	Prune         bool
	Cluster       string
	ClusterParams string
	Component     string
	Format        string
}

var cmdFlagsJsonnet CmdJsonnetOptions

func init() {
	RootCmd.AddCommand(jsonnetCmd)
	jsonnetCmd.AddCommand(jsonnetRenderCmd)
	jsonnetRenderCmd.PersistentFlags().BoolVarP(&cmdFlagsJsonnet.Prune, "prune", "", true, "Prune null and empty objects from rendered json")
	jsonnetRenderCmd.PersistentFlags().StringVarP(&cmdFlagsJsonnet.ClusterParams, "clusterparams", "p", "", "provide cluster params as single file - can be combined with --cluster to override cluster")
	jsonnetRenderCmd.PersistentFlags().StringVarP(&cmdFlagsJsonnet.Component, "component", "c", "", "component to render params for")
	jsonnetRenderCmd.PersistentFlags().StringVarP(&cmdFlagsJsonnet.Format, "format", "F", "json", "Output format: json, yaml, stream")
	jsonnetRenderCmd.PersistentFlags().StringVarP(&cmdFlagsJsonnet.Cluster, "cluster", "C", "", "cluster to render params for")

}
