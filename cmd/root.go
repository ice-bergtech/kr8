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
)

var (
	cfgFile           string
	flagBaseDir       string
	flagClusterDir    string
	flagComponentDir  string
	flagClusterParams string
	flagCluster       string
	flagLogLevel      string

	flagVMConfig VMConfig

	flagDebug       bool
	flagColorOutput bool
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

	RootCmd.PersistentFlags().BoolVar(&flagDebug, "debug", false, "log more information about what kr8 is doing. Overrides --loglevel")
	RootCmd.PersistentFlags().StringVarP(&flagLogLevel, "loglevel", "L", "info", "set log level")
	RootCmd.PersistentFlags().StringVarP(&flagBaseDir, "base", "B", ".", "kr8 config base directory")
	RootCmd.PersistentFlags().StringVarP(&flagClusterDir, "clusterdir", "D", "", "kr8 cluster directory")
	RootCmd.PersistentFlags().StringVarP(&flagComponentDir, "componentdir", "d", "", "kr8 component directory")
	RootCmd.PersistentFlags().BoolVar(&flagColorOutput, "color", true, "enable colorized output (default). Set to false to disable")
	RootCmd.PersistentFlags().StringArrayVarP(&flagVMConfig.Jpaths, "jpath", "J", nil, "Directories to add to jsonnet include path. Repeat arg for multiple directories")
	RootCmd.PersistentFlags().StringSliceVar(&flagVMConfig.ExtVars, "ext-str-file", nil, "Set jsonnet extvar from file contents")
	RootCmd.PersistentFlags().IntVarP(&flagParallel, "parallel", "", runtime.GOMAXPROCS(0), "parallelism - defaults to GOMAXPROCS")

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

	if flagDebug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		switch flagLogLevel {
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
			log.Fatal().Msg("invalid log level: " + flagLogLevel)
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
	flagColorOutput = viper.GetBool("color")
	log.Logger = log.Output(
		zerolog.ConsoleWriter{
			Out:     os.Stderr,
			NoColor: !flagColorOutput,
			FormatErrFieldValue: func(err interface{}) string {
				return strings.ReplaceAll(strings.Join(strings.Split(fmt.Sprintf("%v", err), "\\n"), " | "), "\\t", "")
			}})

	flagBaseDir = viper.GetString("base")
	log.Debug().Msg("Using base directory: " + flagBaseDir)
	flagClusterDir = viper.GetString("clusterdir")
	if flagClusterDir == "" {
		flagClusterDir = flagBaseDir + "/clusters"
	}
	log.Debug().Msg("Using cluster directory: " + flagClusterDir)
	if flagComponentDir == "" {
		flagComponentDir = flagBaseDir + "/components"
	}
	log.Debug().Msg("Using component directory: " + flagComponentDir)

}
