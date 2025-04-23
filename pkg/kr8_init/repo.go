package kr8_init

import (
	"os"
	"strings"

	"github.com/ice-bergtech/kr8/pkg/kr8_types"
)

// Generates a starter readme for the repo, and writes it to the destination directory.
func GenerateReadme(dstDir string, cmdOptions Kr8InitOptions, clusterSpec kr8_types.Kr8ClusterSpec) error {
	var fetch string
	if cmdOptions.Fetch {
		fetch = "true"
	} else {
		fetch = "false"
	}

	readmeTemplate := strings.Join([]string{
		"# Stack " + cmdOptions.ClusterName + " Readme",
		"",
		"## Project Overview",
		"",
		"This project is a cluster stack initialized by kr8+",
		"",
		"* Generate and customize component configuration for Kubernetes clusters across environments, regions and platforms",
		"* Opinionated config, flexible deployment. kr8+ simply generates manifests for you, you decide how to deploy them",
		"* Render and override component config from multiple sources",
		"  * Helm, Kustomize, Static manifests, raw configuration",
		"* Generate static configuration across clusters that is CI/CD friendly",
		"  * Kubernetes manifests, Helm charts, Kustomize overlays, Documentation, text files",
		"",
		"## Usage",
		"",
		"1. Define components in the `components` directory.",
		"2. Define tiered cluster configuration in the `" + clusterSpec.ClusterOutputDir + "` directory.",
		"3. Run `kr8 generate` to generate component configuration files.",
		"",
		"## Info ",
		"",
		"This project is initialized with the following parameters:",
		"",
		"	* ClusterName: `" + cmdOptions.ClusterName + "`",
		"	* Fetch External Libs: " + fetch,
		"   * Cluster config root directory: `" + clusterSpec.ClusterOutputDir + "`",
		"   * Component root directory: `components`",
		"   * Cluster config root directory: `" + clusterSpec.ClusterOutputDir + "`",
		"   * Generated config output directory: `" + clusterSpec.GenerateDir + "`",
		"",
		"Generated using [kr8+](https://github.com/ice-bergtech/kr8)",
	}, "\n")

	return os.WriteFile(dstDir+"/Readme.md", []byte(readmeTemplate), 0600)
}
