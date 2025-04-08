package generate

import (
	"os"
	"path/filepath"

	types "github.com/ice-bergtech/kr8/pkg/types"
	util "github.com/ice-bergtech/kr8/pkg/util"
	"github.com/rs/zerolog/log"
)

func CleanOutputDir(outputFileMap map[string]bool, componentOutputDir string) {
	// clean component dir
	d, err := os.Open(filepath.Clean(componentOutputDir))
	util.FatalErrorCheck("", err)
	// Lifetime of function
	defer d.Close()
	names, err := d.Readdirnames(-1)
	util.FatalErrorCheck("", err)
	for _, name := range names {
		if _, ok := outputFileMap[name]; ok {
			// file is managed
			continue
		}
		if filepath.Ext(name) == ".yaml" {
			delFile := filepath.Join(componentOutputDir, name)
			err = os.RemoveAll(delFile)
			util.FatalErrorCheck("", err)
			log.Debug().Msg("Deleted: " + delFile)
		}
	}
}

func setupClusterGenerateDirs(kr8Spec types.Kr8ClusterSpec) []string {
	// create cluster dir
	if _, err := os.Stat(kr8Spec.ClusterDir); os.IsNotExist(err) {
		err = os.MkdirAll(kr8Spec.ClusterDir, 0750)
		util.FatalErrorCheck("Error creating cluster generateDir", err)
	}

	// get list of current generated components directories
	d, err := os.Open(kr8Spec.ClusterDir)
	util.FatalErrorCheck("Error opening clusterDir", err)
	defer d.Close()

	read_all_dirs := -1
	generatedCompList, err := d.Readdirnames(read_all_dirs)
	util.FatalErrorCheck("Error reading directories", err)

	return generatedCompList
}

// Check if a file needs updating based on its current contents and the new contents.
func CheckIfUpdateNeeded(outFile string, outStr string) bool {
	var updateNeeded bool
	outFile = filepath.Clean(outFile)
	if _, err := os.Stat(outFile); os.IsNotExist(err) {
		log.Debug().Msg("Creating " + outFile)
		updateNeeded = true
	} else {
		currentContents, err := os.ReadFile(outFile)
		util.FatalErrorCheck("Error reading file", err)
		if string(currentContents) != outStr {
			updateNeeded = true
			log.Debug().Msg("Updating: " + outFile)
		}
	}

	return updateNeeded
}
