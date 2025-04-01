package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/tidwall/gjson"
)

type componentDef struct {
	Path string `json:"path"`
}

func (c *Clusters) addItem(item Cluster) Clusters {
	c.Cluster = append(c.Cluster, item)
	return *c
}

func getClusters(searchDir string) (Clusters, error) {

	fileList := make([]string, 0)

	e := filepath.Walk(searchDir, func(path string, f os.FileInfo, err error) error {
		fileList = append(fileList, path)
		return err
	})

	if e != nil {
		log.Fatal().Err(e).Msg("Error building cluster list: ")

	}

	ClusterData := []Cluster{}
	c := Clusters{ClusterData}

	for _, file := range fileList {

		splitFile := strings.Split(file, "/")
		// get the filename
		fileName := splitFile[len(splitFile)-1]

		if fileName == "cluster.jsonnet" {
			entry := Cluster{Name: splitFile[len(splitFile)-2], Path: strings.Join(splitFile[:len(splitFile)-1], "/")}
			c.addItem(entry)

		}
	}

	return c, nil

}

func getCluster(searchDir string, clusterName string) string {

	var clusterPath string

	e := filepath.Walk(searchDir, func(path string, f os.FileInfo, err error) error {
		dir, file := filepath.Split(path)
		if filepath.Base(dir) == clusterName && file == "cluster.jsonnet" {
			clusterPath = path
			return nil
		} else {
			return err
		}
	})

	if e != nil {
		log.Fatal().Err(e).Msg("Error building cluster list: ")

	}

	if clusterPath == "" {
		log.Fatal().Msg("Could not find cluster: " + clusterName)
	}

	return clusterPath

}

func getClusterParams(basePath string, targetPath string) []string {

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

// only render cluster params (_cluster), without components
func renderClusterParamsOnly(vmconfig VMConfig, clusterName string, clusterParams string, prune bool) string {
	var params []string
	if clusterName != "" {
		clusterPath := getCluster(flagClusterDir, clusterName)
		params = getClusterParams(flagClusterDir, clusterPath)
	}
	if clusterParams != "" {
		params = append(params, clusterParams)
	}
	renderedParams := renderJsonnet(vmconfig, params, "._cluster", prune, "", "clusterparams")

	return renderedParams
}

// render cluster params, merged with one or more component's parameters. Empty componentName list renders all component parameters
func renderClusterParams(vmconfig VMConfig, clusterName string, componentNames []string, clusterParams string, prune bool) string {
	if clusterName == "" && clusterParams == "" {
		log.Fatal().Msg("Please specify a --cluster name and/or --clusterparams")
	}

	var params []string
	var componentMap map[string]componentDef

	if clusterName != "" {
		clusterPath := getCluster(flagClusterDir, clusterName)
		params = getClusterParams(flagClusterDir, clusterPath)
	}
	if clusterParams != "" {
		params = append(params, clusterParams)
	}

	compParams := renderJsonnet(vmconfig, params, "", true, "", "clusterparams")

	compString := gjson.Get(compParams, "_components")
	err := json.Unmarshal([]byte(compString.String()), &componentMap)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to parse component map")
	}
	componentDefaultsMerged := "{"
	if len(componentNames) > 0 {
		// we are passed a list of components
		for _, key := range componentNames {
			if value, ok := componentMap[key]; ok {
				path := rootFlagBaseDir + "/" + value.Path + "/params.jsonnet"
				fileC, err := os.ReadFile(path)
				if err != nil {
					log.Fatal().Err(err).Msg("Error reading " + path)
				}
				componentDefaultsMerged = componentDefaultsMerged + fmt.Sprintf("'%s': %s,", key, string(fileC))
			}
		}
	} else {
		// all components
		for key, value := range componentMap {
			if flagComponentName != "" && key != flagComponentName {
				continue
			}
			path := rootFlagBaseDir + "/" + value.Path + "/params.jsonnet"
			fileC, err := os.ReadFile(path)
			if err != nil {
				log.Fatal().Err(err).Msg("Error reading " + path)
			}
			componentDefaultsMerged = componentDefaultsMerged + fmt.Sprintf("'%s': %s,", key, string(fileC))
		}
	}
	componentDefaultsMerged = componentDefaultsMerged + "}"

	compParams = renderJsonnet(vmconfig, params, "", prune, componentDefaultsMerged, "componentparams")

	return compParams
}
