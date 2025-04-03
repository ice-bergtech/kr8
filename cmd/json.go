package cmd

import (
	"github.com/fatih/color"
	"github.com/hokaccha/go-prettyjson"
	util "github.com/ice-bergtech/kr8/pkg/util"
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
	util.FatalErrorCheck(err, "Error formatting JSON")

	return string(formatted)
}
