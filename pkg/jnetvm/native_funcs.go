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

package jnetvm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"path/filepath"
	"regexp"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	jsonnet "github.com/google/go-jsonnet"
	jsonnetAst "github.com/google/go-jsonnet/ast"
	"github.com/grafana/tanka/pkg/helm"

	"github.com/rs/zerolog/log"

	types "github.com/ice-bergtech/kr8/pkg/types"
)

// Registers additional native functions in the jsonnet VM.
// These functions are used to extend the functionality of jsonnet.
// Adds on to functions part of the jsonnet stdlib: https://jsonnet.org/ref/stdlib.html
func RegisterNativeFuncs(jvm *jsonnet.VM) {
	// Register the template function
	jvm.NativeFunction(NativeSprigTemplate())

	// Register the escapeStringRegex function
	jvm.NativeFunction(NativeRegexEscape())

	// Register the regexMatch function
	jvm.NativeFunction(NativeRegexMatch())

	// Register the regexSubst function
	jvm.NativeFunction(NativeRegexSubst())

	// Register the helm function
	jvm.NativeFunction(NativeHelmTemplate())

	// Register the kompose function
	jvm.NativeFunction(NativeKompose())
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
		Func: func(args []interface{}) (interface{}, error) {
			var config any
			err := json.Unmarshal([]byte(args[0].(string)), config)
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
		Func: func(args []interface{}) (interface{}, error) {
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
		Func: func(args []interface{}) (interface{}, error) {
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
		Func: func(args []interface{}) (interface{}, error) {
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
// Inputs: `inFile`, `outPath`, `opts`
func NativeKompose() *jsonnet.NativeFunction {
	return &jsonnet.NativeFunction{
		Name:   "komposeFile",
		Params: jsonnetAst.Identifiers{"inFile", "outPath", "opts"},
		Func: func(args []interface{}) (interface{}, error) {
			inFile, argOk := args[0].(string)
			log.Debug().Msg("inFile: " + inFile)
			if !argOk {
				return nil, fmt.Errorf("first argument 'inFile' must be of 'string' type, got '%T' instead", args[0])
			}

			outPath, argOk := args[1].(string)
			log.Debug().Msg("outPath: " + outPath)
			if !argOk {
				return nil, fmt.Errorf("second argument 'outPath' must be of 'string' type, got '%T' instead", args[1])
			}

			opts, err := parseOpts(args[2])
			if err != nil {
				return "", err
			}

			root := filepath.Dir(opts.CalledFrom)

			options := types.Create([]string{root + "/" + inFile}, root+"/"+outPath, *opts)
			if err := options.Validate(); err != nil {
				return "", err
			}

			return options.Convert()
		},
	}
}

func parseOpts(data interface{}) (*types.Kr8ComponentJsonnet, error) {
	component, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	opts := types.Kr8ComponentJsonnet{}

	if err := json.Unmarshal(component, &opts); err != nil {
		return nil, err
	}

	// Charts are only allowed at relative paths. Use conf.CalledFrom to find the callers directory
	if opts.Namespace == "" {
		return nil, fmt.Errorf("kompose: 'opts.Namespace' is unset or empty.")
	}
	if opts.CalledFrom == "" {
		return nil, fmt.Errorf("kompose: 'opts.CalledFrom' is unset or empty.")
	}

	return &opts, nil
}
