package util

import (
	"os"
	"path/filepath"

	"github.com/hashicorp/go-getter"
	"github.com/rs/zerolog/log"
)

// Fetch a git repo from a url and clone it to a destination directory.
// If the performFetch flag is false, it will log the command that would be run and return without doing anything.
func FetchRepoUrl(url string, destination string, performFetch bool) error {
	if !performFetch {
		gitCommand := "git clone -- " + url + " " + destination
		cleanupCmd := "rm -rf \"" + filepath.Join(destination, ".git") + "\""
		log.Info().Msg("Fetch disabled. Would have ran: ")
		log.Info().Msg("`" + gitCommand + "`")
		log.Info().Msg("`" + cleanupCmd + "`")

		return nil
	}

	// Get the current working directory
	pwd, err := os.Getwd()
	FatalErrorCheck("Error getting working directory", err)

	// Download the skeletion directory
	log.Debug().Msg("Downloading skeleton repo from git::" + url)
	client := &getter.Client{
		Src:  "git::" + url,
		Dst:  destination,
		Pwd:  pwd,
		Mode: getter.ClientModeAny,
	}

	FatalErrorCheck("Error getting repo", client.Get())

	// Check for .git folder
	if _, err := os.Stat(filepath.Join(destination, ".git")); !os.IsNotExist(err) {
		log.Debug().Msg("Removing .git directory")
		FatalErrorCheck("Error removing .git directory", os.RemoveAll(destination+"/.git"))
	}
	return nil
}
