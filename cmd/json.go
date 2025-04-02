package cmd

import (
	"github.com/fatih/color"
	"github.com/hokaccha/go-prettyjson"
)

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
	fatalErrorCheck(err, "Error formatting JSON")

	return string(formatted)
}
