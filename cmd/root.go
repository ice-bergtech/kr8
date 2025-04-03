package cmd

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	types "github.com/ice-bergtech/kr8/pkg/types"
	util "github.com/ice-bergtech/kr8/pkg/util"
)

// exported Version variable
var Version string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "kr8",
	Short: "Kubernetes config parameter framework",
	Long: `A tool to generate Kubernetes configuration from a hierarchy
	of jsonnet files`,
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(version string) {
	Version = version
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

type cmdRootOptions struct {
	BaseDir      string
	ClusterDir   string
	ComponentDir string
	ConfigFile   string
	Parallel     int
	Debug        bool
	LogLevel     string
	Color        bool
	VMConfig     types.VMConfig
}

var rootConfig cmdRootOptions

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().BoolVar(&rootConfig.Debug, "debug", false, "log more information about what kr8 is doing. Overrides --loglevel")
	RootCmd.PersistentFlags().StringVarP(&rootConfig.LogLevel, "loglevel", "L", "info", "set log level")
	RootCmd.PersistentFlags().StringVarP(&rootConfig.BaseDir, "base", "B", ".", "kr8 config base directory")
	RootCmd.PersistentFlags().StringVarP(&rootConfig.ClusterDir, "clusterdir", "D", "", "kr8 cluster directory")
	RootCmd.PersistentFlags().StringVarP(&rootConfig.ComponentDir, "componentdir", "d", "", "kr8 component directory")
	RootCmd.PersistentFlags().BoolVar(&rootConfig.Color, "color", true, "enable colorized output (default). Set to false to disable")
	RootCmd.PersistentFlags().StringArrayVarP(&rootConfig.VMConfig.Jpaths, "jpath", "J", nil, "Directories to add to jsonnet include path. Repeat arg for multiple directories")
	RootCmd.PersistentFlags().StringSliceVar(&rootConfig.VMConfig.ExtVars, "ext-str-file", nil, "Set jsonnet extvar from file contents")
	RootCmd.PersistentFlags().IntVarP(&rootConfig.Parallel, "parallel", "", runtime.GOMAXPROCS(0), "parallelism - defaults to GOMAXPROCS")
	RootCmd.PersistentFlags().StringVarP(&rootConfig.ConfigFile, "config", "", "", "A config file with kr8 configuration")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if rootConfig.ConfigFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(rootConfig.ConfigFile)
	}

	if rootConfig.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		switch rootConfig.LogLevel {
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
			log.Fatal().Msg("invalid log level: " + rootConfig.LogLevel)
		}
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

	log.Logger = log.Output(
		zerolog.ConsoleWriter{
			Out:     os.Stderr,
			NoColor: !rootConfig.Color,
			FormatErrFieldValue: func(err interface{}) string {
				// https://github.com/rs/zerolog/blob/a21d6107dcda23e36bc5cfd00ce8fdbe8f3ddc23/console.go#L21
				colorRed := 31
				colorBold := 1
				s := strings.ReplaceAll(strings.ReplaceAll(strings.TrimRight(err.(string), "\\n"), "\\t", " "), "\\n", " |")
				return util.Colorize(util.Colorize(fmt.Sprintf("%s", s), colorBold, !rootConfig.Color), colorRed, !rootConfig.Color)
			},
		},
	)

	// Setup configuration defaults
	//s.BaseDir = viper.GetString("base")
	log.Debug().Msg("Using base directory: " + rootConfig.BaseDir)

	//s.ClusterDir = viper.GetString("clusterdir")
	if rootConfig.ClusterDir == "" {
		rootConfig.ClusterDir = rootConfig.BaseDir + "/clusters"
	}
	log.Debug().Msg("Using cluster directory: " + rootConfig.ClusterDir)

	if rootConfig.ComponentDir == "" {
		rootConfig.ComponentDir = rootConfig.BaseDir + "/components"
	}
	// Set base config for jvm repo as well.
	rootConfig.VMConfig.BaseDir = rootConfig.BaseDir
	log.Debug().Msg("Using component directory: " + rootConfig.ComponentDir)
}
