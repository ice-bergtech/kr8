package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Return the current version of kr8",
	Long:  `return the current version of kr8`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(RootCmd.Use + " Plus Version: " + Version)
		fmt.Println("github.com/google/go-jsonnet v0.20.0")
		fmt.Println("github.com/ghodss/yaml v1.0.0")
		fmt.Println("github.com/Masterminds/sprig/v3 v3.2.3")
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
