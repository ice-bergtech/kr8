package cmd

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	formatter "github.com/google/go-jsonnet/formatter"
	"github.com/panjf2000/ants/v2"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

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
func formatJsonnetString(input string) (string, error) {
	return formatJsonnetStringCustom(input, GetDefaultFormatOptions())
}

// Formats a jsonnet string using custom options
func formatJsonnetStringCustom(input string, opts formatter.Options) (string, error) {
	return formatter.Format("", input, opts)
}

// Fill with string to include and exclude, using kr8's special parsing
type PathFilterOptions struct {
	Includes string
	Excludes string
}

func (pf PathFilterOptions) Filter(input []string) ([]string, error) {
	if pf.Includes == "" && pf.Excludes == "" {
		return nil, fmt.Errorf("no filter conditions provided")
	}
	return input, nil
}

// Contains the paths to include and exclude for a format command
var cmdformatFlags PathFilterOptions

func init() {
	RootCmd.AddCommand(formatCmd)
	formatCmd.Flags().StringVarP(&cmdformatFlags.Includes, "clincludes", "i", "", "filter included cluster by including clusters with matching cluster parameters - comma separate list of key/value conditions separated by = or ~ (for regex match)")
	formatCmd.Flags().StringVarP(&cmdformatFlags.Excludes, "clexcludes", "x", "", "filter included cluster by excluding clusters with matching cluster parameters - comma separate list of key/value conditions separated by = or ~ (for regex match)")
}

var formatCmd = &cobra.Command{
	Use:   "format [flags]",
	Short: "Format jsonnet files",
	Long:  `Format jsonnet configuration files`,

	Args: cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		var fileList []string
		filepath.Walk(rootConfig.BaseDir, func(path string, info fs.FileInfo, err error) error {
			if info.IsDir() && info.Name() == ".git" {
				return filepath.SkipDir
			}

			if !info.IsDir() {
				var errf error
				match := true
				excludeMatch := false
				if cmdformatFlags.Includes != "" {
					match, errf = filepath.Match(cmdformatFlags.Includes, path)
					if errf != nil {
						return nil
					}
				}
				if cmdformatFlags.Excludes != "" {
					excludeMatch, errf = filepath.Match(cmdformatFlags.Excludes, path)
					if errf != nil {
						return nil
					}
				}
				if match && !excludeMatch && (strings.HasSuffix(info.Name(), ".jsonnet") || strings.HasSuffix(info.Name(), ".libsonnet")) {
					fileList = append(fileList, path)
				}
			}
			return nil
		})

		var wg sync.WaitGroup
		parallel, err := cmd.Flags().GetInt("parallel")
		fatalErrorCheck(err, "Error getting parallel flag")
		log.Debug().Msg("Parallel set to " + strconv.Itoa(parallel))
		ants_file, _ := ants.NewPool(parallel)
		for _, filename := range fileList {
			wg.Add(1)
			_ = ants_file.Submit(func() {
				defer wg.Done()
				var bytes []byte
				bytes, err = os.ReadFile(filename)
				output, err := formatter.Format(filename, string(bytes), GetDefaultFormatOptions())
				if err != nil {
					fmt.Fprintln(os.Stderr, err.Error())
					return
				}
				err = os.WriteFile(filename, []byte(output), 0755)
				if err != nil {
					fmt.Fprintln(os.Stderr, err.Error())
					return
				}
			})
		}
		wg.Wait()
	},
}
