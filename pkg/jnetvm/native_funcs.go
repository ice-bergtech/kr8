package jnetvm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	jsonnet "github.com/google/go-jsonnet"
	jsonnetAst "github.com/google/go-jsonnet/ast"
	"github.com/grafana/tanka/pkg/helm"
	kompose_logger "github.com/sirupsen/logrus"

	"github.com/ice-bergtech/kr8/pkg/kr8_types"
)

// Registers additional native functions in the jsonnet VM.
// These functions are used to extend the functionality of jsonnet.
// Adds on to functions part of the jsonnet stdlib: https://jsonnet.org/ref/stdlib.html
func RegisterNativeFuncs(jvm *jsonnet.VM) {
	listFuncs := []*jsonnet.NativeFunction{
		// Process a helm template directory
		NativeHelmTemplate(),
		// Process a template file
		NativeSprigTemplate(),
		// Process a docker-compose file with kompose
		NativeKompose(),
		// Regex Functions
		// Register the escapeStringRegex function
		NativeRegexEscape(),
		// Register the regexMatch function
		NativeRegexMatch(),
		// Register the regexSubst function
		NativeRegexSubst(),
		// IP helpers
		NativeNetUrl(),
		NativeNetIPInfo(),
		NativeNetAddressCompare(),
		NativeNetAddressDelta(),
		NativeNetAddressSort(),
		NativeNetAddressInc(),
		NativeNetAddressIncBy(),
		NativeNetAddressDec(),
		NativeNetAddressDecBy(),
		NativeNetAddressARPA(),
		NativeNetAddressHex(),
		NativeNetAddressBinary(),
		NativeNetAddressNetsBetween(),
		NativeNetAddressCalcSubnetsV4(),
		NativeNetAddressCalcSubnetsV6(),
	}

	for _, nFunc := range listFuncs {
		jvm.NativeFunction(nFunc)
	}
	// Add help function separately
	jvm.NativeFunction(NativeHelp(listFuncs))
}

func NativeHelp(allFuncs []*jsonnet.NativeFunction) *jsonnet.NativeFunction {
	return &jsonnet.NativeFunction{
		Name:   "help",
		Params: []jsonnetAst.Identifier{},
		Func: func(args []interface{}) (interface{}, error) {
			result := "help: " + strings.Join(
				[]string{
					"Print out kr8 native funcion names and parameters.",
					"Functions are called in the format:",
					"`std.native('<function>')(<param1>, <param2>)`",
				},
				"\n",
			) + "\n"
			result += "\n" + "Available functions:\n"
			result += "\n" + "------------------------\n"

			for _, val := range allFuncs {
				params := []string{}
				for _, id := range val.Params {
					// Convert Identifier to string
					params = append(params, fmt.Sprint(id))
				}
				result += val.Name + ": ['" + strings.Join(params, "', '") + "']\n"
			}

			return result, nil
		},
	}
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
// Inputs: "config" "templateStr".
func NativeSprigTemplate() *jsonnet.NativeFunction {
	return &jsonnet.NativeFunction{
		Name:   "template",
		Params: []jsonnetAst.Identifier{"config", "templateStr"},
		Func: func(args []interface{}) (interface{}, error) {
			var config any
			err := json.Unmarshal([]byte(args[0].(string)), config)
			if err != nil {
				return "", err
			}

			// templateStr is a string that contains the sprig template.
			input, argOk := args[1].(string)
			if !argOk {
				return nil, jsonnet.RuntimeError{
					Msg:        "second argument 'templateStr' must be of 'string' type, got " + fmt.Sprintf("%T", args[1]),
					StackTrace: nil,
				}
			}

			tmpl, err := template.New("templateStr").Funcs(sprig.FuncMap()).Parse(input)
			if err != nil {
				return "", err
			}

			var buff bytes.Buffer
			err = tmpl.Execute(&buff, config)

			return buff.String(), err
		}}
}

// Allows converting a docker-compose file string into kubernetes resources using kompose.
// Files in the directory must be in the format `[docker-]compose.y[a]ml`.
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
			if !argOk {
				return nil, jsonnet.RuntimeError{
					Msg:        "first argument 'inFile' must be of 'string' type, got " + fmt.Sprintf("%T", args[0]),
					StackTrace: nil,
				}
			}

			outPath, argOk := args[1].(string)
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

			// ensure that the logger that kopmose uses is set to warn and above
			kompose_logger.SetLevel(kompose_logger.WarnLevel)
			// TODO: add logrus hook to capture and convert events to zerolog

			options := kr8_types.Create([]string{filepath.Join(root, inFile)}, root+"/"+outPath, *opts)
			if err := options.Validate(); err != nil {
				return "", err
			}

			return options.Convert()
		},
	}
}

func parseOpts(data interface{}) (*kr8_types.Kr8ComponentJsonnet, error) {
	component, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	//nolint:exhaustruct
	opts := kr8_types.Kr8ComponentJsonnet{}
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
