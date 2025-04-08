package cmd

import (
	"fmt"
	"runtime/debug"
	"sort"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

type Stamp struct {
	InfoGoVersion  string
	InfoGoCompiler string
	InfoGOARCH     string
	InfoGOOS       string
	InfoBuildTime  string
	VCSRevision    string
}

func retrieveStamp(info *debug.BuildInfo) *Stamp {
	stamp := Stamp{}
	for _, setting := range info.Settings {
		switch setting.Key {
		case "info.goVersion":
			stamp.InfoGoVersion = setting.Value
		case "-compiler":
			stamp.InfoGoCompiler = setting.Value
		case "GOARCH":
			stamp.InfoGOARCH = setting.Value
		case "GOOS":
			stamp.InfoGOOS = setting.Value
		case "vcs.time":
			stamp.InfoBuildTime = setting.Value
		case "vcs.revision":
			stamp.VCSRevision = setting.Value
		}
	}

	return &stamp
}

func retrieveDepends(info *debug.BuildInfo) []string {
	var ver string

	Depends := []string{}

	filter := map[string]string{
		"github.com/google/go-jsonnet":    "jsonnet ",
		"github.com/ghodss/yaml":          "yaml    ",
		"github.com/grafana/tanka":        "helm    ",
		"github.com/kubernetes/kompose":   "kompose ",
		"github.com/Masterminds/sprig/v3": "template",
	}

	for _, module := range info.Deps {
		log.Debug().Msg(strings.Join([]string{module.Path, module.Version, module.Sum}, " "))
		if _, ok := filter[module.Path]; !ok {
			continue
		}

		if len(module.Version) == 0 {
			ver = module.Sum
		} else {
			ver = module.Version
		}
		Depends = append(Depends, fmt.Sprintf("%s: %s %s", filter[module.Path], module.Path, ver))
	}

	sort.Strings(Depends)
	return Depends
}

// Print out versions of packages in use.
// Chore() - Updated manually.
var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Return the current version of kr8+",
	Long:  `return the current version of kr8+`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(RootCmd.Use + "+ Version: " + version)
		info, ok := debug.ReadBuildInfo()
		if !ok {
			panic("Could not read build info")
		}
		stamp := retrieveStamp(info)
		fmt.Printf("  Built with %s on %s\n", stamp.InfoGoCompiler, stamp.InfoBuildTime)
		fmt.Printf("  VCS revision: %s\n", stamp.VCSRevision)
		fmt.Printf("  Go version %s, GOOS %s, GOARCH %s\n", info.GoVersion, stamp.InfoGOOS, stamp.InfoGOARCH)
		fmt.Print("  Dependencies:\n")
		for _, mod := range retrieveDepends(info) {
			fmt.Printf("    %s\n", mod)
		}

	},
}

func init() {
	RootCmd.AddCommand(VersionCmd)
}
