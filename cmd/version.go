package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Print out versions of packages in use.
// Chore() - Updated manually.
var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Return the current version of kr8",
	Long:  `return the current version of kr8`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(RootCmd.Use + " Plus Version: " + version)
		fmt.Println("jsonnet: github.com/google/go-jsonnet v0.20.0")
		fmt.Println("yml: github.com/ghodss/yaml v1.0.0")
		fmt.Println("helm: github.com/grafana/tanka v0.27.1")
		fmt.Println("kompose: github.com/kubernetes/kompose v1.35.0")
		fmt.Println("template: github.com/Masterminds/sprig/v3 v3.2.3")
	},
}

func init() {
	RootCmd.AddCommand(VersionCmd)
}
