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

package jvm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	jsonnet "github.com/google/go-jsonnet"
	jsonnetAst "github.com/google/go-jsonnet/ast"
	"github.com/grafana/tanka/pkg/helm"
	"github.com/rs/zerolog/log"
	"k8s.io/apimachinery/pkg/runtime"

	types "github.com/ice-bergtech/kr8/pkg/types"
)

// Registers additional native functions in the jsonnet VM.
// These functions are used to extend the functionality of jsonnet.
// Adds on to functions part of the jsonnet stdlib: https://jsonnet.org/ref/stdlib.html
func RegisterNativeFuncs(vm *jsonnet.VM) {
	// Register the template function
	vm.NativeFunction(NativeSprigTemplate())

	// Register the escapeStringRegex function
	vm.NativeFunction(NativeRegexEscape())

	// Register the regexMatch function
	vm.NativeFunction(NativeRegexMatch())

	// Register the regexSubst function
	vm.NativeFunction(NativeRegexSubst())

	// Register the helm function
	vm.NativeFunction(NativeHelmTemplate())

	// Register the kompose function
	vm.NativeFunction(NativeKompose())

}

// Allows executing helm template to process a helm chart and make available to kr8 configuration.
//
// Source: https://github.com/grafana/tanka/blob/v0.27.1/pkg/helm/template.go#L23
func NativeHelmTemplate() *jsonnet.NativeFunction {
	return helm.NativeFunc(helm.ExecHelm{})
}

// Uses sprig to process passed in config data and template.
//
// Sprig template guide: https://masterminds.github.io/sprig/
//
// Inputs: "config" "str"
func NativeSprigTemplate() *jsonnet.NativeFunction {
	return &jsonnet.NativeFunction{
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
		}}
}

// Escapes a string for use in regex
//
// Inputs: "str"
func NativeRegexEscape() *jsonnet.NativeFunction {
	return &jsonnet.NativeFunction{
		Name:   "escapeStringRegex",
		Params: []jsonnetAst.Identifier{"str"},
		Func: func(args []interface{}) (res interface{}, err error) {
			return regexp.QuoteMeta(args[0].(string)), nil
		}}
}

// Matches a string against a regex pattern
//
// Inputs: "regex", "string"
func NativeRegexMatch() *jsonnet.NativeFunction {
	return &jsonnet.NativeFunction{
		Name:   "regexMatch",
		Params: []jsonnetAst.Identifier{"regex", "string"},
		Func: func(args []interface{}) (res interface{}, err error) {
			return regexp.MatchString(args[0].(string), args[1].(string))
		}}
}

// Substitutes a regex pattern in a string with another string
//
// Inputs: "regex", "src", "repl"
func NativeRegexSubst() *jsonnet.NativeFunction {
	return &jsonnet.NativeFunction{
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
		}}
}

// Allows converting a docker-compose file string into kubernetes resources using kompose
//
// Source: https://github.com/kubernetes/kompose/blob/main/cmd/convert.go
//
// Files in the directory must be in the format `[docker-]compose.ym[a]l`
//
// Inputs: `inPath`, `outPath`, `opts`
func NativeKompose() *jsonnet.NativeFunction {
	return &jsonnet.NativeFunction{
		Name:   "kompose",
		Params: []jsonnetAst.Identifier{"input", "komposeOpts"},
		Func: func(args []interface{}) (res interface{}, err error) {
			input := args[0].(string)
			outDir := args[1].(string)
			var componentSpec types.Kr8ComponentJsonnet
			if err := json.Unmarshal([]byte(args[2].(string)), &componentSpec); err != nil {
				return "", err
			}

			options := types.Create([]string{input}, outDir, componentSpec)
			if err := options.Validate(); err != nil {
				return "", err
			}
			return options.Convert()
		}}
}
