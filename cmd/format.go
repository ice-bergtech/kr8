//nolint:gochecknoinits,gochecknoglobals
package cmd

import (
	"io/fs"
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
type CmdFormatOptions struct {
	Paths     []string
	Recursive bool
}

var cmdFormatOptions CmdFormatOptions

func init() {
	RootCmd.AddCommand(FormatCmd)
	FormatCmd.Flags().StringSliceVarP(&cmdFormatOptions.Paths, "paths", "i", []string{"./"},
		"A list of files and/or directories. Defaults to current directory."+
			"If path is a directory, scans directories for files with matching extensions",
	)
	FormatCmd.Flags().BoolVarP(&cmdFormatOptions.Recursive, "recursive", "r", false,
		"If true, will explore directories, formatting files.",
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

// Based on command parameters, builds a list of cluster files that are formatted.
func formatClusterFiles() map[string]string {
	// First get a list of all files in the base directory and subdirectories. Ignore .git directories.
	log.Debug().Msg("Formatting cluster configuration files...")

	allClusterFiles, err := util.GetClusterFilenames(RootConfig.BaseDir)
	util.FatalErrorCheck("issue finding cluster files", err, log.Logger)

	clusterPaths := make(map[string]string, len(allClusterFiles))

	for _, cluster := range allClusterFiles {
		err := FormatFile(filepath.Join(cluster.Path, "cluster.jsonnet"))
		if err != nil {
			log.Error().Str(cluster.Name, cluster.Path).Msg("issue formatting file")
		} else {
			log.Info().Str(cluster.Name, cluster.Path).Msg("formatted")
			clusterPaths[cluster.Name] = cluster.Path
		}
	}

	return clusterPaths
}

func FormatDir(directory string, recursive bool) error {
	err := filepath.WalkDir(RootConfig.BaseDir, func(path string, info fs.DirEntry, err error) error {
		if info.IsDir() {
			if recursive {
				return FormatDir(path, recursive)
			}
		} else {

			fileList = append(fileList, path)
		}

		return nil
	})
}

var FormatCmd = &cobra.Command{
	Use:   "format [flags]",
	Short: "Format jsonnet files",
	Long:  `Format jsonnet and libsonnet configuration files`,

	Args: cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		log.Debug().Msg("Formatting files...")

		// First format cluster paths and get cluster list
		clusterPaths := formatClusterFiles()

		kr8Opts := types.Kr8Opts{
			BaseDir:      RootConfig.BaseDir,
			ComponentDir: RootConfig.ComponentDir,
			ClusterDir:   RootConfig.ClusterDir,
		}

		// Setup pooling options
		parallel, err := cmd.Flags().GetInt("parallel")
		util.FatalErrorCheck("Error getting parallel flag", err, log.Logger)
		if parallel == -1 {
			parallel = runtime.GOMAXPROCS(0)
		}
		log.Debug().Msg("Parallel set to " + strconv.Itoa(parallel))

		var waitGroup sync.WaitGroup
		ants_file, _ := ants.NewPool(parallel)

		// Go through each cluster and format component files.
		for clusterName, clusterDir := range clusterPaths {
			logger := log.With().Str("cluster", clusterName).Logger()

			// Build the cluster-level component parameters for reference.
			_, clusterComponents, err := generate.CompileClusterConfiguration(
				clusterName,
				clusterDir,
				kr8Opts,
				RootConfig.VMConfig,
				"",
				logger,
			)
			util.FatalErrorCheck("issue building cluster config", err, logger)

			// Filter list of components to process
			compList := generate.CalculateClusterComponentList(clusterComponents, cmdFormatFlags)

			// Process each component in a waitgroup function.
			for _, component := range compList {
				subLogger := logger.With().Str("component", component).Logger()
				waitGroup.Add(1)
				_ = ants_file.Submit(func() {
					defer waitGroup.Done()
					path, ok := clusterComponents[component].Map()["Path"]
					if !ok || !path.Exists() {
						subLogger.Error().Msg("issue getting component path")
					}
					err := FormatFile(path.String())
					if util.LogErrorIfCheck("issue formatting file", err, subLogger) != nil {
						return
					}

					// get kr8_spec from cluster's component params
					compSpec, err := kr8_types.CreateComponentSpec(gjson.Get(
						clusterComponents[component].Raw,
						component+".kr8_spec",
					), logger)
					if util.LogErrorIfCheck("Error creating component spec", err, logger) != nil {
						return
					}

					for _, includes := range compSpec.Includes {
						ext := filepath.Ext(includes.File)
						if ext == ".jsonnet" || ext == ".libsonnet" {
							err = FormatFile(includes.File)
							if err != nil {
								logger.Error().Msg("issue formatting file")
							} else {
								logger.Info().Str("file", includes.File).Msg("formatted")
							}
						}
					}
				})
			}
		}

		waitGroup.Wait()
	},
}
