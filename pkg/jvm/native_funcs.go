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
)

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
		Func:   nativeTemplate,
	})

	// Register the escapeStringRegex function
	// Escapes a string for use in regex
	vm.NativeFunction(&jsonnet.NativeFunction{
		Name:   "escapeStringRegex",
		Params: []jsonnetAst.Identifier{"str"},
		Func:   nativeRegexEscape,
	})

	// Register the regexMatch function
	// Matches a string against a regex pattern
	vm.NativeFunction(&jsonnet.NativeFunction{
		Name:   "regexMatch",
		Params: []jsonnetAst.Identifier{"regex", "string"},
		Func:   nativeRegexMatch,
	})

	// Register the regexSubst function
	// Substitutes a regex pattern in a string with another string
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
	// Allows converting a docker-compose file into kubernetes resources using kompose
	// Source: https://github.com/kubernetes/kompose/blob/main/cmd/convert.go
	vm.NativeFunction(&jsonnet.NativeFunction{
		Name:   "kompose",
		Params: []jsonnetAst.Identifier{"input", "komposeOpts"},
		Func:   nativeKompose,
	})

}

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

// Inputs: "str"
func nativeRegexEscape(args []interface{}) (res interface{}, err error) {
	return regexp.QuoteMeta(args[0].(string)), nil
}

// Inputs: "regex", "string"
func nativeRegexMatch(args []interface{}) (res interface{}, err error) {
	return regexp.MatchString(args[0].(string), args[1].(string))
}

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

// Inputs: "input", "komposeOpts"
func nativeKompose(args []interface{}) (res interface{}, err error) {
	//input := args[0].(string)
	// Set the output controller ("deployment"|"daemonSet"|"replicationController")
	// depType := "deployment"
	// options :=

	// 	kompose.ValidateComposeFile(&options)
	// kompose.Convert(options)
	return "", nil
}
