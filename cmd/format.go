//nolint:gochecknoinits,gochecknoglobals
package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"

	formatter "github.com/google/go-jsonnet/formatter"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	util "github.com/ice-bergtech/kr8/pkg/util"
)

// Contains the paths to include and exclude for a format command.
type CmdFormatOptions struct {
	Recursive bool
}

var cmdFormatOptions CmdFormatOptions

func init() {
	RootCmd.AddCommand(FormatCmd)
	FormatCmd.Flags().BoolVarP(&cmdFormatOptions.Recursive, "recursive", "r", false,
		"recursively explore the parameter directories",
	)
}

// Read, format, and write back a file.
// github.com/google/go-jsonnet/formatter is used to format files.
func FormatFile(filename string) error {
	bytes, err := os.ReadFile(filepath.Clean(filename))
	if err != nil {
		return err
	}
	output, err := formatter.Format(filename, string(bytes), util.GetDefaultFormatOptions())
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Clean(filename), []byte(output), 0600)
}

// Lists each file in a directory and formats all .jsonnet and .libsonnet files.
// If recursive flag is enabled, will explore the directory tree.
func FormatDir(inputPath string, recursive bool, logger zerolog.Logger) error {
	fileInfo, err := os.Stat(inputPath)
	if err != nil {
		return err
	}

	filePaths := []string{inputPath}

	if fileInfo.IsDir() {
		filePaths, err = dirFileListAndFormat(inputPath, recursive, logger)
		if err != nil {
			return err
		}
	}

	log.Debug().Any("files", filePaths).Msg("collected files")

	for _, file := range filePaths {
		ext := filepath.Ext(file)
		if ext == ".jsonnet" || ext == ".libsonnet" {
			err := FormatFile(file)
			if err != nil {
				logger.Error().Err(err).Msg("issue formatting " + file)
			} else {
				logger.Info().Msg("formatted " + file)
			}
		}
	}

	return nil
}

// Returns a list of files in the inputPath.
// If a directory is encountered and recursive is true, will format the directory.
func dirFileListAndFormat(inputPath string, recursive bool, logger zerolog.Logger) ([]string, error) {
	filePaths := []string{}
	dirEntries, err := os.ReadDir(inputPath)
	if err != nil {
		return nil, err
	}
	for _, entry := range dirEntries {
		if entry.IsDir() {
			if recursive {
				err := FormatDir(entry.Name(), recursive, logger)
				if err != nil {
					return nil, err
				}
			}

			continue
		}
		filePaths = append(filePaths, entry.Name())
	}

	return filePaths, nil
}

// Take the formatting options and format them for help output.
func prettyPrintFormattingOpts() string {
	str, err := json.MarshalIndent(util.GetDefaultFormatOptions(), "", "  ")
	if err != nil {
		return err.Error()
	}
	output, err := util.FormatJsonnetString(string(str))
	if err != nil {
		return err.Error()
	}

	return output
}

var FormatCmd = &cobra.Command{
	Use:     "format [flags] [files or directories]",
	Aliases: []string{"fmt"},
	Short:   "Format jsonnet files in a directory.  Defaults to `./`",
	Long: `Formats jsonnet and libsonnet files.
A list of files and/or directories. Defaults to current directory (./).
If path is a directory, scans directories for files with matching extensions.
Formats files with the following options: ` + prettyPrintFormattingOpts(),
	Args: cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		log.Debug().Any("options", cmdFormatOptions).Msg("Formatting files...")
		paths := args
		if len(paths) == 0 {
			paths = []string{"./"}
		}

		for _, path := range paths {
			logger := log.With().Str("param", path).Logger()
			err := FormatDir(path, cmdFormatOptions.Recursive, logger)
			if err != nil {
				logger.Error().Err(err).Msg("issue formatting path")
			}
		}
	},
}
