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
		return ""
	}
	f := prettyjson.NewFormatter()
	f.Indent = 4
	if !colorOutput {
		f.DisabledColor = true
	}
	f.KeyColor = color.New(color.FgRed)
	f.NullColor = color.New(color.Underline)

	formatted, err := f.Format([]byte(input))
	FatalErrorCheck(err, "Error formatting JSON")

	return string(formatted)
}

// Colorize function from zerolog console.go file to replicate their coloring functionality.
// Source: https://github.com/rs/zerolog/blob/a21d6107dcda23e36bc5cfd00ce8fdbe8f3ddc23/console.go#L389
func Colorize(s interface{}, c int, disabled bool) string {
	e := os.Getenv("NO_COLOR")
	if e != "" || c == 0 {
		disabled = true
	}

	if disabled {
		return fmt.Sprintf("%s", s)
	}
	return fmt.Sprintf("\x1b[%dm%v\x1b[0m", c, s)
}

// Print the jsonnet output in the specified format.
// Acceptable formats are: yaml, stream, json
func JsonnetPrint(output string, format string, color bool) {
	switch format {
	case "yaml":
		yaml, err := goyaml.JSONToYAML([]byte(output))
		FatalErrorCheck(err, "Error converting output JSON to YAML")
		fmt.Println(string(yaml))
	case "stream": // output yaml stream
		var o []interface{}
		FatalErrorCheck(json.Unmarshal([]byte(output), &o), "Error unmarshalling output JSON")
		for _, jobj := range o {
			fmt.Println("---")
			buf, err := goyaml.Marshal(jobj)
			FatalErrorCheck(err, "Error marshalling output JSON to YAML")
			fmt.Println(string(buf))
		}
	case "json":
		formatted := Pretty(output, color)
		fmt.Println(formatted)
	default:
		log.Fatal().Msg("Output format must be json, yaml or stream")
	}
}

// Configures the default options for the jsonnet formatter
func GetDefaultFormatOptions() formatter.Options {
	return formatter.Options{
		Indent:           2,
		MaxBlankLines:    2,
		StringStyle:      formatter.StringStyleLeave,
		CommentStyle:     formatter.CommentStyleLeave,
		UseImplicitPlus:  false,
		PrettyFieldNames: true,
		PadArrays:        false,
		PadObjects:       true,
		SortImports:      true,
		StripEverything:  false,
		StripComments:    false,
	}
}

// Formats a jsonnet string using the default options
func FormatJsonnetString(input string) (string, error) {
	return FormatJsonnetStringCustom(input, GetDefaultFormatOptions())
}

// Formats a jsonnet string using custom options
func FormatJsonnetStringCustom(input string, opts formatter.Options) (string, error) {
	return formatter.Format("", input, opts)
}

// Write out a struct to a specified path and file.
// If successful, returns what was written. If not successful, returns an error.
func WriteObjToJsonFile(filename string, path string, objStruct interface{}) (error, string) {
	if err := os.MkdirAll(path, 0755); err != nil {
		return err, "error creating resource directory"
	}

	jsonStr, err := json.MarshalIndent(objStruct, "", "  ")
	if err != nil {
		return err, "error marshalling component resource to json"
	}

	jsonStrFormatted, err := FormatJsonnetString(string(jsonStr))
	if err != nil {
		return err, "error formatting component resource to json"
	}

	return (os.WriteFile(path+"/"+filename, []byte(jsonStrFormatted), 0644)), jsonStrFormatted
}
