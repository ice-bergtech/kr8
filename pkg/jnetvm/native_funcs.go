package jnetvm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	jsonnet "github.com/google/go-jsonnet"
	jsonnetAst "github.com/google/go-jsonnet/ast"
	"github.com/grafana/tanka/pkg/helm"
	types "github.com/ice-bergtech/kr8/pkg/types"
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
				return nil, types.Kr8Error{
					Message: "second argument 'templateStr' must be of 'string' type, got " + fmt.Sprintf("%T", args[1]),
					Value:   nil,
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
