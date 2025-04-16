package kr8_init

import (
	types "github.com/ice-bergtech/kr8/pkg/types"
	util "github.com/ice-bergtech/kr8/pkg/util"
)

// Generate a cluster.jsonnet file based on the provided Kr8ClusterSpec and store it in the specified directory.
func GenerateClusterJsonnet(cSpec types.Kr8ClusterSpec, dstDir string) error {
	filename := "cluster.jsonnet"
	clusterJson := types.Kr8ClusterJsonnet{
		ClusterSpec: cSpec,
		// Bug() Unsure if Path is correct
		Cluster:    types.Kr8Cluster{Name: cSpec.Name, Path: cSpec.ClusterOutputDir},
		Components: map[string]types.Kr8ClusterComponentRef{},
	}
	_, err := util.WriteObjToJsonFile(filename, dstDir+"/"+cSpec.Name, clusterJson)

	return err
}
