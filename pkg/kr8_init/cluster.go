package kr8_init

import (
	"github.com/ice-bergtech/kr8/pkg/kr8_types"
	util "github.com/ice-bergtech/kr8/pkg/util"
)

// Generate a cluster.jsonnet file based on the provided Kr8ClusterSpec and store it in the specified directory.
func GenerateClusterJsonnet(cSpec kr8_types.Kr8ClusterSpec, dstDir string) error {
	filename := "cluster.jsonnet"
	clusterJson := kr8_types.Kr8ClusterJsonnet{
		ClusterSpec: cSpec,
		// Bug() Unsure if Path is correct
		Cluster:    kr8_types.Kr8Cluster{Name: cSpec.Name, Path: cSpec.ClusterOutputDir},
		Components: map[string]kr8_types.Kr8ClusterComponentRef{},
	}
	_, err := util.WriteObjToJsonFile(filename, dstDir+"/"+cSpec.Name, clusterJson)

	return err
}
