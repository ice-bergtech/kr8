// Defines the cli-interface commands available to the user.
//
//nolint:gochecknoinits,gochecknoglobals
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	types "github.com/ice-bergtech/kr8/pkg/types"
	util "github.com/ice-bergtech/kr8/pkg/util"
)

// exported version variable.
var version string

// RootCmd represents the base command when called without any subcommands.
var RootCmd = &cobra.Command{
	Use:   "kr8",
	Short: "A jsonnet-powered config management tool",
	Long:  `An opinionated configuration management tool for Kubernetes Clusters powered by jsonnet`,
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(ver string) {
	version = ver
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

// Default options that are available to all commands.
type CmdRootOptions struct {
	// kr8+ config base directory
	BaseDir string
	// kr8+ cluster directory
	ClusterDir string
	// kr8+ component directory
	ComponentDir string
	// A config file with kr8+ configuration
	ConfigFile string
	// parallelism - defaults to runtime.GOMAXPROCS(0)
	Parallel int
	// log more information about what kr8+ is doing. Overrides --loglevel
	Debug bool
	// set log level
	LogLevel string
	// enable colorized output (default true). Set to false to disable")
	Color bool
	// contains information to configure jsonnet vm
	VMConfig types.VMConfig
	// Profiling output directory.  Only captured if set.
	ProfilingDir string
	// CPU profiling output file handle.
	ProfilingCPUFile *os.File
}

var RootConfig CmdRootOptions

func init() {
	// Ran before each command is ran
	cobra.OnInitialize(InitConfig, ProfilingInitializer)
	cobra.OnFinalize(ProfilingFinalizer)

	RootCmd.PersistentFlags().BoolVar(&RootConfig.Debug,
		"debug", false,
		"log additional information about what kr8+ is doing. Overrides --loglevel")
	RootCmd.PersistentFlags().StringVarP(&RootConfig.LogLevel,
		"loglevel", "L", "info",
		"set zerolog log level")
	RootCmd.PersistentFlags().StringVarP(&RootConfig.BaseDir, "base", "B", "./", "kr8+ root configuration directory")
	RootCmd.PersistentFlags().StringVarP(&RootConfig.ClusterDir,
		"clusterdir", "D", "",
		"kr8+ cluster directory")
	RootCmd.PersistentFlags().StringVarP(&RootConfig.ComponentDir,
		"componentdir", "d", "",
		"kr8+ component directory")
	RootCmd.PersistentFlags().BoolVar(&RootConfig.Color,
		"color", true,
		"enable colorized output")
	RootCmd.PersistentFlags().StringArrayVarP(&RootConfig.VMConfig.JPaths,
		"jpath", "J", nil,
		"additional jsonnet library directories")
	RootCmd.PersistentFlags().StringSliceVar(&RootConfig.VMConfig.ExtVars,
		"ext-str-file", nil,
		"set comma-separated jsonnet extVars from file contents in the format `key=file`")
	RootCmd.PersistentFlags().IntVarP(&RootConfig.Parallel,
		"parallel", "", -1,
		"parallelism - defaults to runtime.GOMAXPROCS(0)")
	RootCmd.PersistentFlags().StringVarP(&RootConfig.ConfigFile,
		"config", "", "",
		"a single config file with kr8+ configuration")
	RootCmd.PersistentFlags().StringVarP(&RootConfig.ProfilingDir,
		"profiledir", "", "",
		"directory to write pprof profile data to")
}

func ConfigureLogger(debug bool) {
	if debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		switch RootConfig.LogLevel {
		case "debug":
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
		case "info":
			zerolog.SetGlobalLevel(zerolog.InfoLevel)
		case "warn":
			zerolog.SetGlobalLevel(zerolog.WarnLevel)
		case "error":
			zerolog.SetGlobalLevel(zerolog.ErrorLevel)
		case "fatal":
			zerolog.SetGlobalLevel(zerolog.FatalLevel)
		case "panic":
			zerolog.SetGlobalLevel(zerolog.PanicLevel)
		default:
			log.Fatal().Msg("invalid log level: " + RootConfig.LogLevel)
		}
	}

	log.Logger = util.SetupLogger(RootConfig.Color)
}

// InitConfig reads in config file and ENV variables if set.
func InitConfig() {
	ConfigureLogger(RootConfig.Debug)
	// enable ability to specify config file via flag
	if RootConfig.ConfigFile != "" {
		viper.SetConfigFile(RootConfig.ConfigFile)
	}

	if RootConfig.Parallel == -1 {
		RootConfig.Parallel = runtime.GOMAXPROCS(0)
	}

	viper.SetConfigName(".kr8") // name of config file (without extension)
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME") // adding home directory as first search path
	viper.SetEnvPrefix("KR8")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Debug().Msg("Using config file:" + viper.ConfigFileUsed())
	}

	// Setup configuration defaults
	log.Debug().Msg("Using base directory: " + RootConfig.BaseDir)

	if RootConfig.ClusterDir == "" {
		RootConfig.ClusterDir = filepath.Join(RootConfig.BaseDir, "clusters")
	}
	log.Debug().Msg("Using cluster directory: " + RootConfig.ClusterDir)

	if RootConfig.ComponentDir == "" {
		RootConfig.ComponentDir = filepath.Join(RootConfig.BaseDir, "components")
	}
	// Set base config for jvm repo as well.
	RootConfig.VMConfig.BaseDir = RootConfig.BaseDir
	log.Debug().Msg("Using component directory: " + RootConfig.ComponentDir)
}

// Stop profiling and write cpu and memory profiling files if configured.
func ProfilingFinalizer() {
	if RootConfig.ProfilingDir != "" {
		pprof.StopCPUProfile()
		if RootConfig.ProfilingCPUFile != nil {
			_ = RootConfig.ProfilingCPUFile.Close()
		}

		runtime.GC() // get up-to-date statistics

		// Various types of profiles that can be collected:
		// https://cs.opensource.google/go/go/+/go1.24.2:src/runtime/pprof/pprof.go;l=178
		var err error
		heapFile, err := os.Create(filepath.Join(RootConfig.ProfilingDir, "profile_heap.pb.gz"))
		if err != nil {
			log.Fatal().Err(err).Msg("could not write memory profile: ")
		}
		if err = pprof.WriteHeapProfile(heapFile); err != nil {
			_ = heapFile.Close()
			log.Fatal().Err(err).Msg("could not write memory profile: ")
		}
		_ = heapFile.Close()
	}
}

// Sets up program profiling.
func ProfilingInitializer() {
	var err error
	if RootConfig.ProfilingDir != "" {
		RootConfig.ProfilingCPUFile, err = os.Create(filepath.Join(RootConfig.ProfilingDir, "profile_cpu.pb.gz"))
		if err != nil {
			log.Fatal().Err(err).Msg("could not create CPU profile: ")
		}
		if err := pprof.StartCPUProfile(RootConfig.ProfilingCPUFile); err != nil {
			log.Fatal().Err(err).Msg("could not create CPU profile: ")
		}
	}
}
