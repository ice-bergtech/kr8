package generate

import (
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"

	"github.com/ice-bergtech/kr8/pkg/kr8_types"
	util "github.com/ice-bergtech/kr8/pkg/util"
)

func CleanOutputDir(outputFileMap map[string]bool, componentOutputDir string) error {
	// clean component dir
	dir, err := os.Open(filepath.Clean(componentOutputDir))
	if err := util.ErrorIfCheck("", err); err != nil {
		return err
	}
	// Lifetime of function
	defer dir.Close()
	names, err := dir.Readdirnames(-1)
	if err := util.ErrorIfCheck("", err); err != nil {
		return err
	}
	for _, name := range names {
		if _, ok := outputFileMap[name]; ok {
			// file is managed
			continue
		}
		if filepath.Ext(name) == ".yaml" {
			delFile := filepath.Join(componentOutputDir, name)
			err = os.RemoveAll(delFile)
			if err := util.ErrorIfCheck("", err); err != nil {
				return err
			}
			log.Debug().Msg("Deleted unmanaged file: " + delFile)
		}
	}

	return nil
}

func setupClusterGenerateDirs(kr8Spec kr8_types.Kr8ClusterSpec) ([]string, error) {
	// create cluster dir
	if _, err := os.Stat(kr8Spec.ClusterOutputDir); os.IsNotExist(err) {
		err = os.MkdirAll(kr8Spec.ClusterOutputDir, 0750)
		if err := util.ErrorIfCheck("Error creating cluster generateDir", err); err != nil {
			return []string{}, err
		}
	}

	// get list of current generated components directories
	dir, err := os.Open(kr8Spec.ClusterOutputDir)
	if err := util.ErrorIfCheck("Error opening clusterDir", err); err != nil {
		return []string{}, err
	}
	defer dir.Close()

	read_all_dirs := -1
	generatedCompList, err := dir.Readdirnames(read_all_dirs)
	if err := util.ErrorIfCheck("Error reading directories", err); err != nil {
		return []string{}, err
	}

	return generatedCompList, nil
}

// Check if a file needs updating based on its current contents and the new contents.
func CheckIfUpdateNeeded(outFile string, outStr string) (bool, error) {
	outFile = filepath.Clean(outFile)
	if _, err := os.Stat(outFile); os.IsNotExist(err) {
		log.Debug().Msg("File needs to be created: " + outFile)

		return true, nil
	}

	currentContents, err := os.ReadFile(outFile)
	if err := util.ErrorIfCheck("Error reading file "+outFile, err); err != nil {
		return false, err
	}
	if string(currentContents) != outStr {
		log.Debug().Msg("File needs to be updated: " + outFile)

		return true, nil
	}

	return false, nil
}
