//nolint:gochecknoinits,gochecknoglobals
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"

	formatter "github.com/google/go-jsonnet/formatter"
	"github.com/panjf2000/ants/v2"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"

	"github.com/ice-bergtech/kr8/pkg/generate"
	"github.com/ice-bergtech/kr8/pkg/kr8_types"
	"github.com/ice-bergtech/kr8/pkg/types"
	util "github.com/ice-bergtech/kr8/pkg/util"
)

// Contains the paths to include and exclude for a format command.
var cmdFormatFlags util.PathFilterOptions

func init() {
	RootCmd.AddCommand(FormatCmd)
	FormatCmd.Flags().StringVarP(&cmdFormatFlags.Includes,
		"clincludes", "i", "",
		"filter included cluster by including clusters with matching cluster parameters -"+
			" comma separate list of key/value conditions separated by = or ~ (for regex match)",
	)
	FormatCmd.Flags().StringVarP(&cmdFormatFlags.Excludes,
		"clexcludes", "x", "",
		"filter included cluster by excluding clusters with matching cluster parameters -"+
			" comma separate list of key/value conditions separated by = or ~ (for regex match)",
	)
}

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

var FormatCmd = &cobra.Command{
	Use:   "format [flags]",
	Short: "Format jsonnet files",
	Long:  `Format jsonnet configuration files`,

	Args: cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		// First get a list of all files in the base directory and subdirectories. Ignore .git directories.
		var fileList []string
		err := filepath.Walk(RootConfig.BaseDir, func(path string, info fs.FileInfo, err error) error {
			if info.IsDir() {
				if info.Name() == ".git" {
					return filepath.SkipDir
				}

				return nil
			}
			fileList = append(fileList, path)

			return nil
		})
		util.FatalErrorCheck("Error walking the path "+RootConfig.BaseDir, err, log.Logger)

		fileList = util.Filter(fileList, func(s string) bool {
			var result bool
			for _, f := range strings.Split(cmdFormatFlags.Includes, ",") {
				t, _ := filepath.Match(f, s)
				if t {
					return t
				}
				result = result || t
			}

			return result
		})

		fileList = util.Filter(fileList, func(s string) bool {
			var result bool
			for _, f := range strings.Split(cmdFormatFlags.Excludes, ",") {
				t, _ := filepath.Match(f, s)
				if t {
					return !t
				}
				result = result || t
			}

			return !result
		})
		log.Debug().Msg("Filtered file list: " + fmt.Sprintf("%v", fileList))
		log.Debug().Msg("Formatting files...")

		var waitGroup sync.WaitGroup
		parallel, err := cmd.Flags().GetInt("parallel")
		util.FatalErrorCheck("Error getting parallel flag", err, log.Logger)
		log.Debug().Msg("Parallel set to " + strconv.Itoa(parallel))
		ants_file, _ := ants.NewPool(parallel)
		for _, filename := range fileList {
			waitGroup.Add(1)
			_ = ants_file.Submit(func() {
				defer waitGroup.Done()
				var bytes []byte
				bytes, err = os.ReadFile(filepath.Clean(filename))
				output, err := formatter.Format(filename, string(bytes), util.GetDefaultFormatOptions())
				if err != nil {
					fmt.Fprintln(os.Stderr, err.Error())

					return
				}
				err = os.WriteFile(filepath.Clean(filename), []byte(output), 0600)
				if err != nil {
					fmt.Fprintln(os.Stderr, err.Error())

					return
				}
			})
		}
		waitGroup.Wait()
	},
}
