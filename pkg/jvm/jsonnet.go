package jvm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"golang.org/x/exp/maps"

	"github.com/Masterminds/sprig/v3"
	jsonnet "github.com/google/go-jsonnet"
	jsonnetAst "github.com/google/go-jsonnet/ast"
	"github.com/grafana/tanka/pkg/helm"
	"github.com/rs/zerolog/log"
	"github.com/tidwall/gjson"

	types "github.com/ice-bergtech/kr8/pkg/types"
	util "github.com/ice-bergtech/kr8/pkg/util"
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

func JsonnetVM(vmconfig types.VMConfig) (*jsonnet.VM, error) {
	vm := jsonnet.MakeVM()
	RegisterNativeFuncs(vm)

	// always add lib directory in base directory to path
	jpath := []string{vmconfig.BaseDir + "/lib"}

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
func RenderJsonnet(vmConfig types.VMConfig, files []string, param string, prune bool, prepend string, source string) string {

	// copy the slice so that we don't unitentionally modify the original
	jsonnetPaths := make([]string, len(files[:0]))
	copy(jsonnetPaths, files[:0])

	// range through the files
	for _, s := range files {
		jsonnetPaths = append(jsonnetPaths, fmt.Sprintf("(import '%s')", s))
	}

	// Create a JSonnet VM
	vm, err := JsonnetVM(vmConfig)
	util.FatalErrorCheck(err, "Error creating jsonnet VM")

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
	util.FatalErrorCheck(err, "Error evaluating jsonnet snippet")

	return out

}

// Native Jsonnet funcs to add
/*

This code is copied almost verbatim from the kubecfg project: https://github.com/ksonnet/kubecfg
Native funcs: https://github.com/kubecfg/kubecfg/blob/main/utils/nativefuncs.go

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

// Registers additional native functions in the jsonnet VM
// These functions are used to extend the functionality of jsonnet
// Adds on to functions part of the jsonnet stdlib: https://jsonnet.org/ref/stdlib.html
func RegisterNativeFuncs(vm *jsonnet.VM) {
	// Register the template function
	// Uses sprig to process as passed in template and config
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
			tmpl, err := template.New("file").Funcs(sprig.FuncMap()).Parse(string(input))
			if err != nil {
				return "", err
			}

			var buff bytes.Buffer
			err = tmpl.Execute(&buff, config)
			return buff.String(), err
		},
	})

	// Register the escapeStringRegex function
	// Escapes a string for use in regex
	vm.NativeFunction(&jsonnet.NativeFunction{
		Name:   "escapeStringRegex",
		Params: []jsonnetAst.Identifier{"str"},
		Func: func(args []interface{}) (res interface{}, err error) {
			return regexp.QuoteMeta(args[0].(string)), nil
		},
	})

	// Register the regexMatch function
	// Matches a string against a regex pattern
	vm.NativeFunction(&jsonnet.NativeFunction{
		Name:   "regexMatch",
		Params: []jsonnetAst.Identifier{"regex", "string"},
		Func: func(args []interface{}) (res interface{}, err error) {
			return regexp.MatchString(args[0].(string), args[1].(string))
		},
	})

	// Register the regexSubst function
	// Substitutes a regex pattern in a string with another string
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

	// Register the helm function
	// Allows executing helm template to process a helm chart and make available to kr8 configuration
	// Source: https://github.com/grafana/tanka/blob/v0.27.1/pkg/helm/template.go#L23
	vm.NativeFunction(helm.NativeFunc(helm.ExecHelm{}))

	// Register the kompose function
	// Allows converting a docker-compose file into kubernetes resources using kompose
	// Source: https://github.com/kubernetes/kompose/blob/main/cmd/convert.go
	vm.NativeFunction(&jsonnet.NativeFunction{
		Name:   "kompose",
		Params: []jsonnetAst.Identifier{"input", "komposeOpts"},
		Func: func(args []interface{}) (res interface{}, err error) {
			//input := args[0].(string)
			// Set the output controller ("deployment"|"daemonSet"|"replicationController")
			// depType := "deployment"
			// options :=

			// 	kompose.ValidateComposeFile(&options)
			// kompose.Convert(options)
			return "", nil
		},
	})

}

// Renders a jsonnet file with the specified options.
func JsonnetRender(cmdFlagsJsonnet types.CmdJsonnetOptions, filename string, vmConfig types.VMConfig) {
	// Check if cluster and/or clusterparams are specified
	if cmdFlagsJsonnet.Cluster == "" && cmdFlagsJsonnet.ClusterParams == "" {
		log.Fatal().Msg("Please specify a --cluster name and/or --clusterparams")
	}

	// Render the cluster parameters
	config := RenderClusterParams(vmConfig, cmdFlagsJsonnet.Cluster, []string{cmdFlagsJsonnet.Component}, cmdFlagsJsonnet.ClusterParams, false)

	// Create a new VM instance
	vm, _ := JsonnetVM(vmConfig)
	// Setup kr8 config as external vars
	vm.ExtCode("kr8_cluster", "std.prune("+config+"._cluster)")
	vm.ExtCode("kr8_components", "std.prune("+config+"._components)")
	vm.ExtCode("kr8", "std.prune("+config+"."+cmdFlagsJsonnet.Component+")")
	vm.ExtCode("kr8_unpruned", config+"."+cmdFlagsJsonnet.Component)

	var input string
	// If pruning is enabled, prune the input before rendering
	// This removes all null and empty fields from the imported file
	if cmdFlagsJsonnet.Prune {
		input = "std.prune(import '" + filename + "')"
	} else {
		input = "( import '" + filename + "')"
	}

	//
	// Evaluate the jsonnet snippet and print the result
	// This is where the magic happens! The jsonnet code is evaluated and the result is stored
	//
	j, err := vm.EvaluateAnonymousSnippet("file", input)
	util.FatalErrorCheck(err, "Error evaluating jsonnet snippet")

	util.JsonnetPrint(j, cmdFlagsJsonnet.Format, cmdFlagsJsonnet.Color)
}

// only render cluster params (_cluster), without components
func RenderClusterParamsOnly(vmconfig types.VMConfig, clusterName string, clusterParams string, prune bool) string {
	var params []string
	if clusterName != "" {
		clusterPath := util.GetCluster(vmconfig.BaseDir, clusterName)
		params = util.GetClusterParams(vmconfig.BaseDir, clusterPath)
	}
	if clusterParams != "" {
		params = append(params, clusterParams)
	}
	renderedParams := RenderJsonnet(vmconfig, params, "._cluster", prune, "", "clusterparams")

	return renderedParams
}

// render cluster params, merged with one or more component's parameters. Empty componentName list renders all component parameters
func RenderClusterParams(vmconfig types.VMConfig, clusterName string, componentNames []string, clusterParams string, prune bool) string {
	if clusterName == "" && clusterParams == "" {
		log.Fatal().Msg("Please specify a --cluster name and/or --clusterparams")
	}

	var params []string
	var componentMap map[string]types.Kr8ClusterComponentRef

	if clusterName != "" {
		clusterPath := util.GetCluster(vmconfig.BaseDir, clusterName)
		params = util.GetClusterParams(vmconfig.BaseDir, clusterPath)
	}
	if clusterParams != "" {
		params = append(params, clusterParams)
	}

	compParams := RenderJsonnet(vmconfig, params, "", true, "", "clusterparams")

	compString := gjson.Get(compParams, "_components")
	err := json.Unmarshal([]byte(compString.String()), &componentMap)
	util.FatalErrorCheck(err, "failed to parse component map")
	componentDefaultsMerged := "{"

	listComponentKeys := maps.Keys(componentMap)
	if len(componentNames) > 0 {
		listComponentKeys = componentNames
	}

	// all components
	for _, key := range listComponentKeys {
		if value, ok := componentMap[key]; ok {
			path := vmconfig.BaseDir + "/" + value.Path + "/params.jsonnet"
			fileC, err := os.ReadFile(path)
			util.FatalErrorCheck(err, "Error reading file "+path)
			componentDefaultsMerged = componentDefaultsMerged + fmt.Sprintf("'%s': %s,", key, string(fileC))
		}
	}
	componentDefaultsMerged = componentDefaultsMerged + "}"

	compParams = RenderJsonnet(vmconfig, params, "", prune, componentDefaultsMerged, "componentparams")

	return compParams
}
