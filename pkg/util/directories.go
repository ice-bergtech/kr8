package util

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/google/go-jsonnet"
	"github.com/rs/zerolog/log"

	types "github.com/ice-bergtech/kr8/pkg/types"
)

type ClusterTreeNode struct {
	DirName string
	VM      *jsonnet.VM
	Parent  *ClusterTreeNode
}

// Get a list of cluster from within a directory.
// Walks the directory tree, creating a types.Kr8Cluster for each cluster.jsonnet file found.
func GetClusterFilenames(searchDir string) ([]types.Kr8Cluster, error) {
	fileList := make([]string, 0)
	ClusterData := []types.Kr8Cluster{}

	err := filepath.Walk(searchDir, func(path string, f os.FileInfo, err error) error {
		fileList = append(fileList, path)
		// Pass error through
		return err
	})
	if err != nil {
		return ClusterData, ErrorIfCheck("error building cluster list", err)
	}

	for _, file := range fileList {
		// get the filename
		splitFile := strings.Split(file, "/")
		fileName := splitFile[len(splitFile)-1]
		// check if the filename is cluster.jsonnet
		if fileName == "cluster.jsonnet" {
			entry := types.Kr8Cluster{Name: splitFile[len(splitFile)-2], Path: strings.Join(splitFile[:len(splitFile)-1], "/")}
			ClusterData = append(ClusterData, entry)
		}
	}

	return ClusterData, nil
}

// Get a specific cluster within a directory by name.
// Walks the cluster directory searching for the given clusterName.
// Returns the path to the cluster.
func GetClusterPaths(searchDir string, clusterName string) (string, error) {
	clusterPath := ""

	err := filepath.Walk(searchDir, func(path string, f os.FileInfo, err error) error {
		dir, file := filepath.Split(path)
		if filepath.Base(dir) == clusterName && file == "cluster.jsonnet" {
			clusterPath = path
			// No error
			return nil
		} else {
			// Pass back the error
			return err
		}
	})
	if err != nil {
		return "", ErrorIfCheck("error building cluster list", err)
	}
	if clusterPath == "" {
		return "", types.Kr8Error{Message: "error: could not find cluster: " + clusterName, Value: ""}
	}

	return clusterPath, nil
}

// Get all cluster parameters within a directory.
// Walks through the directory hierarchy and returns all paths to `params.jsonnet` files.
func GetClusterParamsFilenames(basePath string, targetPath string) []string {
	// a slice to store results
	var results []string
	results = append(results, targetPath)

	// remove the cluster.jsonnet
	splitFile := strings.Split(targetPath, "/")

	// gets the targetDir without the cluster.jsonnet
	targetDir := strings.Join(splitFile[:len(splitFile)-1], "/")

	// walk through the directory hierarchy
	for {
		rel, _ := filepath.Rel(basePath, targetDir)

		// check if there's a params.json in the folder
		if _, err := os.Stat(filepath.Join(targetDir, "params.jsonnet")); err == nil {
			results = append(results, filepath.Join(targetDir, "params.jsonnet"))
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
	for i := range len(results) / 2 {
		results[i], results[last-i] = results[last-i], results[i]
	}

	return results
}

// Given a map of filenames, prunes all *.yaml files that are NOT in the map from the directory.
func CleanOutputDir(outputFileMap map[string]bool, componentOutputDir string) error {
	// clean component dir
	dir, err := os.Open(filepath.Clean(componentOutputDir))
	if err != nil {
		return err
	}
	// Lifetime of function
	defer dir.Close()
	names, err := dir.Readdirnames(-1)
	if err != nil {
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
			if err != nil {
				return err
			}
			log.Debug().Msg("Deleted: " + delFile)
		}
	}

	return nil
}
