package util

import (
	"encoding/json"
	"fmt"
	"os"

	formatter "github.com/google/go-jsonnet/formatter"

	"github.com/fatih/color"
	goyaml "github.com/ghodss/yaml"
	"github.com/hokaccha/go-prettyjson"
	"github.com/rs/zerolog/log"
)

// Pretty formats the input jsonnet string with indentation and optional color output.
func Pretty(input string, colorOutput bool) string {
	if input == "" {
		// escape hatch for empty input
		return ""
	}
	fmtr := prettyjson.NewFormatter()
	fmtr.Indent = 4
	if !colorOutput {
		fmtr.DisabledColor = true
	}
	fmtr.KeyColor = color.New(color.FgRed)
	fmtr.NullColor = color.New(color.Underline)

	formatted, err := fmtr.Format([]byte(input))
	FatalErrorCheck("Error formatting JSON", err)

	return string(formatted)
}

// Colorize function from zerolog console.go file to replicate their coloring functionality.
// Source: https://github.com/rs/zerolog/blob/a21d6107dcda23e36bc5cfd00ce8fdbe8f3ddc23/console.go#L389
func Colorize(input interface{}, colorNum int, disabled bool) string {
	e := os.Getenv("NO_COLOR")

	if disabled || (e != "" || colorNum == 0) {
		// escape hatch for disabled coloring or empty input
		return fmt.Sprintf("%s", input)
	}

	return fmt.Sprintf("\x1b[%dm%v\x1b[0m", colorNum, input)
}

// Print the jsonnet output in the specified format.
// Acceptable formats are: yaml, stream, json.
func JsonnetPrint(output string, format string, color bool) {
	switch format {
	case "yaml":
		yaml, err := goyaml.JSONToYAML([]byte(output))
		FatalErrorCheck("Error converting output JSON to YAML", err)
		fmt.Println(string(yaml))
	case "stream": // output yaml stream
		var o []interface{}
		FatalErrorCheck("Error unmarshalling output JSON", json.Unmarshal([]byte(output), &o))
		for _, jobj := range o {
			fmt.Println("---")
			buf, err := goyaml.Marshal(jobj)
			FatalErrorCheck("Error marshalling output JSON to YAML", err)
			fmt.Println(string(buf))
		}
	case "json":
		formatted := Pretty(output, color)
		fmt.Println(formatted)
	default:
		log.Fatal().Msg("Output format must be json, yaml or stream")
	}
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
// If successful, returns what was written. If not successful, returns an error.
func WriteObjToJsonFile(filename string, path string, objStruct interface{}) (string, error) {
	if err := os.MkdirAll(path, 0750); err != nil {
		return "error creating resource directory", err
	}

	jsonStr, err := json.MarshalIndent(objStruct, "", "  ")
	if err != nil {
		return "error marshalling component resource to json", err
	}

	jsonStrFormatted, err := FormatJsonnetString(string(jsonStr))
	if err != nil {
		return "error formatting component resource to json", err
	}

	return jsonStrFormatted, (os.WriteFile(path+"/"+filename, []byte(jsonStrFormatted), 0600))
}
