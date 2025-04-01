// Copyright © 2018 Lee Briggs <lee@leebriggs.co.uk>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"os"

	"github.com/hashicorp/go-getter"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	dl_url   string
	real_url string
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize kr8 config repos, components and clusters",
	Long: `kr8 requires specific directories and exists for its config to work.
This init command helps in creating directory structure for repos, clusters and 
components`,
	//Run: func(cmd *cobra.Command, args []string) {},
}

var repoCmd = &cobra.Command{
	Use:   "repo dir",
	Args:  cobra.MinimumNArgs(1),
	Short: "Initialize a new kr8 config repo",
	Long: `Initialize a new kr8 config repo by downloading the kr8 config skeleton repo
and initialize a git repo so you can get started`,
	Run: func(cmd *cobra.Command, args []string) {
		if dl_url != "" {
			real_url = dl_url
		} else {
			log.Fatal().Msg("Must specify a URL arg")
		}
		// Get the current working directory
		pwd, err := os.Getwd()
		fatalErrorCheck(err, "Error getting working directory")

		// Download the skeletion directory
		log.Debug().Msg("Downloading skeleton repo from " + real_url)
		client := &getter.Client{
			Src:  real_url,
			Dst:  args[0],
			Pwd:  pwd,
			Mode: getter.ClientModeAny,
		}

		fatalErrorCheck(client.Get(), "Error getting repo")

		// Check for .git folder
		if _, err := os.Stat(args[0] + "/.git"); !os.IsNotExist(err) {
			log.Debug().Msg("Removing .git directory")
			os.RemoveAll(args[0] + "/.git")
		}
	},
}

func init() {
	RootCmd.AddCommand(initCmd)
	initCmd.AddCommand(repoCmd)

	repoCmd.PersistentFlags().StringVar(&dl_url, "url", "", "Source of skeleton directory to create repo from")

}
