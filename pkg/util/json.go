package util

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/fatih/color"
	goyaml "github.com/ghodss/yaml"
	formatter "github.com/google/go-jsonnet/formatter"
	"github.com/hokaccha/go-prettyjson"

	types "github.com/ice-bergtech/kr8/pkg/types"
)

// Pretty formats the input jsonnet string with indentation and optional color output.
// Returns an error when the input can't properly format the json string input.
func Pretty(inputJson string, colorOutput bool) (string, error) {
	if inputJson == "" {
		// escape hatch for empty input
		return "", nil
	}
	fmtr := prettyjson.NewFormatter()
	fmtr.Indent = 4
	if !colorOutput {
		fmtr.DisabledColor = true
	}
	fmtr.KeyColor = color.New(color.FgRed)
	fmtr.NullColor = color.New(color.Underline)

	formatted, err := fmtr.Format([]byte(inputJson))
	if err := GenErrorIfCheck("error formatting JSON", err); err != nil {
		return "", err
	}

	return string(formatted), nil
}

// Colorize function from zerolog console.go file to replicate their coloring functionality.
// Source: https://github.com/rs/zerolog/blob/a21d6107dcda23e36bc5cfd00ce8fdbe8f3ddc23/console.go#L389
// Replicated here because it's a private function.
func Colorize(input interface{}, colorNum int, disabled bool) string {
	e := os.Getenv("NO_COLOR")

	if disabled || (e != "" || colorNum == 0) {
		// escape hatch for disabled coloring or empty input
		return fmt.Sprintf("%s", input)
	}

	return fmt.Sprintf("\x1b[%dm%v\x1b[0m", colorNum, input)
}

// Print the jsonnet in the specified format.
// Acceptable formats are: yaml, stream, json.
func JsonnetPrint(output string, format string, color bool) error {
	switch format {
	case "yaml":
		yaml, err := goyaml.JSONToYAML([]byte(output))
		if err := GenErrorIfCheck("error converting output JSON to YAML", err); err != nil {
			return err
		}
		fmt.Println(string(yaml))
	case "stream": // output yaml stream
		var o []interface{}
		if err := GenErrorIfCheck("error unmarshalling output JSON", json.Unmarshal([]byte(output), &o)); err != nil {
			return err
		}
		for _, jobj := range o {
			fmt.Println("---")
			buf, err := goyaml.Marshal(jobj)
			if err := GenErrorIfCheck("error marshalling output JSON to YAML", err); err != nil {
				return err
			}
			fmt.Println(string(buf))
		}
	case "json":
		formatted, err := Pretty(output, color)
		if err := GenErrorIfCheck("error formatting output JSON", err); err != nil {
			return err
		}
		fmt.Println(formatted)
	default:
		return types.Kr8Error{Message: "error: output format must be json, yaml or stream", Value: format}
	}

	return nil
}

// Configures the default options for the jsonnet formatter.
func GetDefaultFormatOptions() formatter.Options {
	return formatter.Options{
		Indent:              2,
		MaxBlankLines:       2,
		StringStyle:         formatter.StringStyleLeave,
		CommentStyle:        formatter.CommentStyleLeave,
		UseImplicitPlus:     false,
		PrettyFieldNames:    true,
		PadArrays:           false,
		PadObjects:          true,
		SortImports:         true,
		StripEverything:     false,
		StripComments:       false,
		StripAllButComments: false,
	}
}

// Formats a jsonnet string using the default options.
func FormatJsonnetString(input string) (string, error) {
	return FormatJsonnetStringCustom(input, GetDefaultFormatOptions())
}

// Formats a jsonnet string using custom options.
func FormatJsonnetStringCustom(input string, opts formatter.Options) (string, error) {
	return formatter.Format("", input, opts)
}

// Write out a struct to a specified path and file.
// Marshals the given interface and generates a formatted json string.
// All parent directories needed are created.
func WriteObjToJsonFile(filename string, path string, objStruct interface{}) (string, error) {
	jsonStr, err := json.MarshalIndent(objStruct, "", "  ")
	if err != nil {
		return "", GenErrorIfCheck("error marshalling component resource to json", err)
	}

	jsonStrFormatted, err := FormatJsonnetString(string(jsonStr))
	if err != nil {
		return "", GenErrorIfCheck("error formatting component resource to json", err)
	}
	// Create directories after we marshal and format the json
	if err := os.MkdirAll(path, 0750); err != nil {
		return "", GenErrorIfCheck("error creating resource directory", err)
	}

	return jsonStrFormatted, GenErrorIfCheck("error writing file to disk",
		os.WriteFile(path+"/"+filename, []byte(jsonStrFormatted), 0600),
	)
}
