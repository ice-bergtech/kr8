package cmd

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"

	formatter "github.com/google/go-jsonnet/formatter"
	"github.com/panjf2000/ants/v2"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	formatDir string
	pIncludes string
	pExcludes string
)

var formatCmd = &cobra.Command{
	Use:   "format",
	Short: "Format jsonnet files",
	Long:  `Format jsonnet configuration files`,

	Args: cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		var fileList []string
		filepath.Walk(formatDir, func(path string, info fs.FileInfo, err error) error {
			if info.IsDir() && info.Name() == ".git" {
				return filepath.SkipDir
			}

			if !info.IsDir() {
				var errf error
				match := true
				excludeMatch := false
				if pIncludes != "" {
					match, errf = filepath.Match(pIncludes, path)
					if errf != nil {
						return nil
					}
				}
				if pExcludes != "" {
					excludeMatch, errf = filepath.Match(pExcludes, path)
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

		options := formatter.DefaultOptions()
		options.StringStyle = formatter.StringStyleLeave
		for _, filename := range fileList {
			wg.Add(1)
			_ = ants_file.Submit(func() {
				defer wg.Done()
				var bytes []byte
				bytes, err = os.ReadFile(filename)
				output, err := formatter.Format(filename, string(bytes), options)
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

func init() {
	RootCmd.AddCommand(formatCmd)
	formatCmd.Flags().StringVarP(&formatDir, "dir", "", "./", "Root directory to walk and format")

	formatCmd.Flags().StringVarP(&pIncludes, "pincludes", "", "", "filter included paths by including paths - filepath.Match format - https://pkg.go.dev/path/filepath#Match")
	formatCmd.Flags().StringVarP(&pExcludes, "pexcludes", "", "", "filter included paths by excluding paths - filepath.Match format - https://pkg.go.dev/path/filepath#Match")
	formatCmd.Flags().IntP("parallel", "", runtime.GOMAXPROCS(0), "parallelism - defaults to GOMAXPROCS")
	viper.BindPFlag("pIncludes", formatCmd.PersistentFlags().Lookup("pIncludes"))
	viper.BindPFlag("pExcludes", formatCmd.PersistentFlags().Lookup("pExcludes"))
}
