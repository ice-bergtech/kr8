package util

import (
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	types "github.com/ice-bergtech/kr8/pkg/types"
)

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
// [filepath.Walk]s the cluster directory tree searching for the given clusterName.
// Returns the path to the cluster.jsonnet file.
func GetClusterPath(searchDir string, clusterName string) (string, error) {
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

// Walk a directory to build a list of all files in the tree.
func BuildDirFileList(directory string) ([]string, error) {
	clusterPaths := []string{}

	err := filepath.Walk(directory, func(path string, f os.FileInfo, err error) error {
		if f != nil && !f.IsDir() {
			clusterPaths = append(clusterPaths, path)
		}

		return nil
	})
	if err != nil {
		return []string{}, ErrorIfCheck("error building dir file list", err)
	}

	sort.Strings(clusterPaths)

	return clusterPaths, nil
}

// Write bytes to file (path included).
func WriteFile(input []byte, file string) error {
	f, err := os.Create(filepath.Clean(file))
	if err != nil {
		return err
	}
	defer f.Close()

	return os.WriteFile(file, input, 0600)
}

// Read bytes from file (path included).
func ReadFile(file string) ([]byte, error) {
	fCache, err := os.Open(filepath.Clean(file))
	if err != nil {
		return nil, err
	}
	defer fCache.Close()

	fileInfo, err := fCache.Stat()
	if err != nil {
		return nil, err
	}
	text := make([]byte, fileInfo.Size())
	_, err = fCache.Read(text)
	if err != nil {
		return nil, err
	}

	return text, nil
}

// Write bytes to gzip file (path included).
func WriteGzip(input []byte, file string) error {
	outFile, err := os.Create(filepath.Clean(file))
	if err != nil {
		return err
	}
	defer outFile.Close()
	gzipWriter := gzip.NewWriter(outFile)
	_, err = gzipWriter.Write(input)
	_ = gzipWriter.Close()

	return err
}

// Read bytes from a gzip file (path included).
func ReadGzip(filename string) ([]byte, error) {
	inFile, err := os.Open(filepath.Clean(filename))
	if err != nil {
		return []byte{}, err
	}
	reader, err := gzip.NewReader(inFile)
	if err != nil {
		return []byte{}, err
	}
	data, err := io.ReadAll(reader)
	if err != nil {
		return []byte{}, err
	}

	return data, nil
}

// Lists each file in a directory and formats all .jsonnet and .libsonnet files.
// If recursive flag is enabled, will explore the directory tree.
func FileFuncInDir(
	inputPath string,
	recursive bool,
	fileFunc func(string, zerolog.Logger) error,
	logger zerolog.Logger,
) error {
	fileInfo, err := os.Stat(inputPath)
	if err != nil {
		return err
	}

	filePaths := []string{inputPath}

	if fileInfo.IsDir() {
		filePaths, err = dirFilesApplyFunc(inputPath, recursive, fileFunc, logger)
		if err != nil {
			return err
		}
	}

	log.Debug().Any("files", filePaths).Msg("collected files")

	for _, file := range filePaths {
		err := fileFunc(file, logger.With().Str("file", file).Logger())
		if err != nil {
			logger.Error().Err(err).Msg("issue processing file")
		} else {
			logger.Info().Msg(file)
		}
	}

	return nil
}

// Returns a list of files in the inputPath.
// If a directory is encountered and recursive is true, will format the directory.
func dirFilesApplyFunc(
	inputPath string,
	recursive bool,
	fileFunc func(string, zerolog.Logger) error,
	logger zerolog.Logger,
) ([]string, error) {
	filePaths := []string{}
	dirEntries, err := os.ReadDir(inputPath)
	if err != nil {
		return nil, err
	}
	for _, entry := range dirEntries {
		if entry.IsDir() {
			if recursive {
				err := FileFuncInDir(entry.Name(), recursive, fileFunc, logger)
				if err != nil {
					return nil, err
				}
			}

			continue
		}
		filePaths = append(filePaths, entry.Name())
	}

	return filePaths, nil
}
