package jnetvm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"path/filepath"
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

	jvm.NativeFunction(NativeNetUrl())

	// IP helpers
	jvm.NativeFunction(NativeNetIPInfo())
	jvm.NativeFunction(NativeNetAddressCompare())
	jvm.NativeFunction(NativeNetAddressDelta())
	jvm.NativeFunction(NativeNetAddressSort())
	jvm.NativeFunction(NativeNetAddressInc())
	jvm.NativeFunction(NativeNetAddressIncBy())
	jvm.NativeFunction(NativeNetAddressDec())
	jvm.NativeFunction(NativeNetAddressDecBy())
}

// Allows executing helm template to process a helm chart and make available to kr8 configuration.
//
// Source: https://github.com/grafana/tanka/blob/v0.27.1/pkg/helm/template.go#L23
func NativeHelmTemplate() *jsonnet.NativeFunction {
	return helm.NativeFunc(helm.ExecHelm{})
}

// Uses sprig to process passed in config data and template.
// Sprig template guide: https://masterminds.github.io/sprig/ .
//
// Inputs: "config" "str".
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

// Allows converting a docker-compose file string into kubernetes resources using kompose.
// Files in the directory must be in the format `[docker-]compose.ym[a]l`.
//
// Source: https://github.com/kubernetes/kompose/blob/main/cmd/convert.go
//
// Inputs: `inFile`, `outPath`, `opts`.
func NativeKompose() *jsonnet.NativeFunction {
	return &jsonnet.NativeFunction{
		Name:   "komposeFile",
		Params: jsonnetAst.Identifiers{"inFile", "outPath", "opts"},
		Func: func(args []interface{}) (interface{}, error) {
			inFile, argOk := args[0].(string)
			log.Debug().Msg("inFile: " + inFile)
			if !argOk {
				return nil, jsonnet.RuntimeError{
					Msg:        "first argument 'inFile' must be of 'string' type, got " + fmt.Sprintf("%T", args[0]),
					StackTrace: nil,
				}
			}

			outPath, argOk := args[1].(string)
			log.Debug().Msg("outPath: " + outPath)
			if !argOk {
				return nil, jsonnet.RuntimeError{
					Msg:        "second argument 'outPath' must be of 'string' type, got " + fmt.Sprintf("%T", args[1]),
					StackTrace: nil,
				}
			}

			opts, err := parseOpts(args[2])
			if err != nil {
				return "", err
			}

			root := filepath.Dir(opts.CalledFrom)

			options := types.Create([]string{filepath.Join(root, inFile)}, root+"/"+outPath, *opts)
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

	// Charts are only allowed at relative paths. Use conf.CalledFrom to find the callers directory.
	if opts.Namespace == "" {
		return nil, jsonnet.RuntimeError{
			Msg:        "kompose: 'opts.namespace' is unset or empty.",
			StackTrace: nil,
		}
	}
	if opts.CalledFrom == "" {
		return nil, jsonnet.RuntimeError{
			Msg:        "kompose: 'opts.called_from' is unset or empty.",
			StackTrace: nil,
		}
	}

	return &opts, nil
}
