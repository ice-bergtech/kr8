package util

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"

	types "github.com/ice-bergtech/kr8/pkg/types"
)

func GetClusters(searchDir string) ([]types.Kr8Cluster, error) {

	fileList := make([]string, 0)

	FatalErrorCheck(
		filepath.Walk(searchDir, func(path string, f os.FileInfo, err error) error {
			fileList = append(fileList, path)
			return err
		}),
		"Error building cluster list",
	)

	ClusterData := []types.Kr8Cluster{}

	for _, file := range fileList {

		splitFile := strings.Split(file, "/")
		// get the filename
		fileName := splitFile[len(splitFile)-1]

		if fileName == "cluster.jsonnet" {
			entry := types.Kr8Cluster{Name: splitFile[len(splitFile)-2], Path: strings.Join(splitFile[:len(splitFile)-1], "/")}
			ClusterData = append(ClusterData, entry)
		}
	}

	return ClusterData, nil

}

func GetCluster(searchDir string, clusterName string) string {
	clusterPath := ""

	FatalErrorCheck(
		filepath.Walk(searchDir, func(path string, f os.FileInfo, err error) error {
			dir, file := filepath.Split(path)
			if filepath.Base(dir) == clusterName && file == "cluster.jsonnet" {
				clusterPath = path
				return nil
			} else {
				return err
			}
		}),
		"Error building cluster list",
	)

	if clusterPath == "" {
		log.Fatal().Msg("Could not find cluster: " + clusterName)
	}

	return clusterPath

}

func GetClusterParams(basePath string, targetPath string) []string {

	// a slice to store results
	var results []string
	results = append(results, targetPath)

	// remove the cluster.jsonnet
	splitFile := strings.Split(targetPath, "/")

	// gets the targetDir without the cluster.jsonnet
	targetDir := strings.Join(splitFile[:len(splitFile)-1], "/")

	// walk through the directory hierachy
	for {
		rel, _ := filepath.Rel(basePath, targetDir)

		// check if there's a params.json in the folder
		if _, err := os.Stat(targetDir + "/params.jsonnet"); err == nil {
			results = append(results, targetDir+"/params.jsonnet")
		}

		// stop if we're in the basePath
		if rel == "." {
			break
		}

		// next!
		targetDir += "/.."
	}

	// jsonnet's import order matters, so we need to reverse the slice
	last := len(results) - 1
	for i := 0; i < len(results)/2; i++ {
		results[i], results[last-i] = results[last-i], results[i]
	}

	return results

}
