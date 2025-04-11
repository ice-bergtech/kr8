package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	types "github.com/ice-bergtech/kr8p/pkg/types"
	util "github.com/ice-bergtech/kr8p/pkg/util"
)

// exported version variable.
var version string

// RootCmd represents the base command when called without any subcommands.
var RootCmd = &cobra.Command{
	Use:   "kr8p",
	Short: "Kubernetes config parameter framework",
	Long: `A tool to generate Kubernetes configuration from a hierarchy
	of jsonnet files`,
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
	// kr8p config base directory
	BaseDir string
	// kr8p cluster directory
	ClusterDir string
	// kr8p component directory
	ComponentDir string
	// A config file with kr8p configuration
	ConfigFile string
	// parallelism - defaults to runtime.GOMAXPROCS(0)
	Parallel int
	// log more information about what kr8p is doing. Overrides --loglevel
	Debug bool
	// set log level
	LogLevel string
	// enable colorized output (default true). Set to false to disable")
	Color bool
	// contains ingormation to configure jsonnet vm
	VMConfig types.VMConfig
}

var RootConfig CmdRootOptions

func init() {
	cobra.OnInitialize(InitConfig)

	RootCmd.PersistentFlags().BoolVar(&RootConfig.Debug,
		"debug", false,
		"log more information about what kr8p is doing. Overrides --loglevel")
	RootCmd.PersistentFlags().StringVarP(&RootConfig.LogLevel,
		"loglevel", "L", "info",
		"set log level")
	RootCmd.PersistentFlags().StringVarP(&RootConfig.BaseDir, "base", "B", "./", "kr8p config base directory")
	RootCmd.PersistentFlags().StringVarP(&RootConfig.ClusterDir,
		"clusterdir", "D", "",
		"kr8p cluster directory")
	RootCmd.PersistentFlags().StringVarP(&RootConfig.ComponentDir,
		"componentdir", "d", "",
		"kr8p component directory")
	RootCmd.PersistentFlags().BoolVar(&RootConfig.Color,
		"color", true,
		"enable colorized output. Set to false to disable")
	RootCmd.PersistentFlags().StringArrayVarP(&RootConfig.VMConfig.Jpaths,
		"jpath", "J", nil,
		"Directories to add to jsonnet include path. Repeat arg for multiple directories")
	RootCmd.PersistentFlags().StringSliceVar(&RootConfig.VMConfig.ExtVars,
		"ext-str-file", nil,
		"Set jsonnet extvar from file contents")
	RootCmd.PersistentFlags().IntVarP(&RootConfig.Parallel,
		"parallel", "", -1,
		"parallelism - defaults to runtime.GOMAXPROCS(0)")
	RootCmd.PersistentFlags().StringVarP(&RootConfig.ConfigFile,
		"config", "", "",
		"A config file with kr8p configuration")
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

	log.Logger = log.Output(
		zerolog.ConsoleWriter{
			Out:     os.Stderr,
			NoColor: !RootConfig.Color,
			FormatErrFieldValue: func(err interface{}) string {
				// https://github.com/rs/zerolog/blob/a21d6107dcda23e36bc5cfd00ce8fdbe8f3ddc23/console.go#L21
				colorRed := 31
				colorBold := 1
				s := strings.ReplaceAll(strings.ReplaceAll(strings.TrimRight(err.(string), "\\n"), "\\t", " "), "\\n", " |")

				return util.Colorize(util.Colorize(s, colorBold, !RootConfig.Color), colorRed, !RootConfig.Color)
			},
		},
	)
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
