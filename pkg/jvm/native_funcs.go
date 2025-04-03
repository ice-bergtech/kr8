package jvm

import (
	"bytes"
	"encoding/json"
	"regexp"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	jsonnet "github.com/google/go-jsonnet"
	jsonnetAst "github.com/google/go-jsonnet/ast"
	"github.com/grafana/tanka/pkg/helm"

	types "github.com/ice-bergtech/kr8/pkg/types"
)

// Native Jsonnet funcs to add
/*

Much of this code is based on the kubecfg project: https://github.com/ksonnet/kubecfg
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
	vm.NativeFunction(&jsonnet.NativeFunction{
		Name:   "template",
		Params: []jsonnetAst.Identifier{"config", "str"},
		Func:   nativeTemplate,
	})

	// Register the escapeStringRegex function
	vm.NativeFunction(&jsonnet.NativeFunction{
		Name:   "escapeStringRegex",
		Params: []jsonnetAst.Identifier{"str"},
		Func:   nativeRegexEscape,
	})

	// Register the regexMatch function
	vm.NativeFunction(&jsonnet.NativeFunction{
		Name:   "regexMatch",
		Params: []jsonnetAst.Identifier{"regex", "string"},
		Func:   nativeRegexMatch,
	})

	// Register the regexSubst function
	vm.NativeFunction(&jsonnet.NativeFunction{
		Name:   "regexSubst",
		Params: []jsonnetAst.Identifier{"regex", "src", "repl"},
		Func:   nativeRegexSubst,
	})

	// Register the helm function
	// Allows executing helm template to process a helm chart and make available to kr8 configuration
	// Source: https://github.com/grafana/tanka/blob/v0.27.1/pkg/helm/template.go#L23
	vm.NativeFunction(helm.NativeFunc(helm.ExecHelm{}))

	// Register the kompose function
	vm.NativeFunction(&jsonnet.NativeFunction{
		Name:   "kompose",
		Params: []jsonnetAst.Identifier{"input", "komposeOpts"},
		Func:   nativeKompose,
	})

}

// Uses sprig to process passed in config data and template
// Inputs: "config" "str"
func nativeTemplate(args []interface{}) (res interface{}, err error) {
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
}

// Escapes a string for use in regex
// Inputs: "str"
func nativeRegexEscape(args []interface{}) (res interface{}, err error) {
	return regexp.QuoteMeta(args[0].(string)), nil
}

// Matches a string against a regex pattern
// Inputs: "regex", "string"
func nativeRegexMatch(args []interface{}) (res interface{}, err error) {
	return regexp.MatchString(args[0].(string), args[1].(string))
}

// Substitutes a regex pattern in a string with another string
// Inputs: "regex", "src", "repl"
func nativeRegexSubst(args []interface{}) (res interface{}, err error) {
	regex := args[0].(string)
	src := args[1].(string)
	repl := args[2].(string)

	r, err := regexp.Compile(regex)
	if err != nil {
		return "", err
	}
	return r.ReplaceAllString(src, repl), nil
}

// Allows converting a docker-compose file string into kubernetes resources using kompose
// Source: https://github.com/kubernetes/kompose/blob/main/cmd/convert.go
// Inputs: "input filename", "outdir", "componentConfig"
// Filename must be one of: “
func nativeKompose(args []interface{}) (res interface{}, err error) {
	input := args[0].(string)
	outDir := args[1].(string)
	var componentSpec types.Kr8ComponentJsonnet
	if err := json.Unmarshal([]byte(args[2].(string)), componentSpec); err != nil {
		return "", err
	}
	// Set the output controller ("deployment"|"daemonSet"|"replicationController")
	// depType := "deployment"

	options := types.Create([]string{input}, outDir, componentSpec)
	if err := options.Validate(); err != nil {
		return "", err
	}
	return options.Convert()
}
