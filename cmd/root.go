package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	cfgFile       string
	baseDir       string
	clusterDir    string
	componentDir  string
	clusterParams string
	cluster       string
	logLevel      string

	rootVMConfig VMConfig

	debug       bool
	colorOutput bool
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

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "log more information about what kr8 is doing. Overrides --loglevel")
	RootCmd.PersistentFlags().StringVarP(&logLevel, "loglevel", "L", "info", "set log level")
	RootCmd.PersistentFlags().StringVarP(&baseDir, "base", "B", ".", "kr8 config base directory")
	RootCmd.PersistentFlags().StringVarP(&clusterDir, "clusterdir", "D", "", "kr8 cluster directory")
	RootCmd.PersistentFlags().StringVarP(&componentDir, "componentdir", "d", "", "kr8 component directory")
	RootCmd.PersistentFlags().BoolVar(&colorOutput, "color", true, "enable colorized output (default). Set to false to disable")
	RootCmd.PersistentFlags().StringArrayVarP(&rootVMConfig.Jpaths, "jpath", "J", nil, "Directories to add to jsonnet include path. Repeat arg for multiple directories")
	RootCmd.PersistentFlags().StringSliceVar(&rootVMConfig.ExtVars, "ext-str-file", nil, "Set jsonnet extvar from file contents")

	viper.BindPFlag("base", RootCmd.PersistentFlags().Lookup("base"))
	viper.BindPFlag("clusterdir", RootCmd.PersistentFlags().Lookup("clusterdir"))
	viper.BindPFlag("componentdir", RootCmd.PersistentFlags().Lookup("componentdir"))
	viper.BindPFlag("color", RootCmd.PersistentFlags().Lookup("color"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	if debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		switch logLevel {
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
			log.Fatal().Msg("invalid log level: " + logLevel)
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
	colorOutput = viper.GetBool("color")
	log.Logger = log.Output(
		zerolog.ConsoleWriter{
			Out:     os.Stderr,
			NoColor: !colorOutput,
			FormatErrFieldValue: func(err interface{}) string {
				return strings.ReplaceAll(strings.Join(strings.Split(fmt.Sprintf("%v", err), "\\n"), " | "), "\\t", "")
			}})

	baseDir = viper.GetString("base")
	log.Debug().Msg("Using base directory: " + baseDir)
	clusterDir = viper.GetString("clusterdir")
	if clusterDir == "" {
		clusterDir = baseDir + "/clusters"
	}
	log.Debug().Msg("Using cluster directory: " + clusterDir)
	if componentDir == "" {
		componentDir = baseDir + "/components"
	}
	log.Debug().Msg("Using component directory: " + componentDir)

}
