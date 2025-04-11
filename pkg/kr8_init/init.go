package kr8init

import (
	"os"
	"strings"

	types "github.com/ice-bergtech/kr8p/pkg/types"
	util "github.com/ice-bergtech/kr8p/pkg/util"
)

// Kr8InitOptions defines the options used by the init subcommands.
type Kr8InitOptions struct {
	// URL to fetch the skeleton directory from
	InitUrl string
	// Name of the cluster to initialize
	ClusterName string
	// Name of the component to initialize
	ComponentName string
	// Type of component to initialize (e.g. jsonnet, yml, chart, compose)
	ComponentType string
	// Whether to run in interactive mode or not
	Interactive bool
	// Whether to fetch remote resources or not
	Fetch bool
}

// Generate a cluster.jsonnet file based on the provided Kr8ClusterSpec and store it in the specified directory.
func GenerateClusterJsonnet(cSpec types.Kr8ClusterSpec, dstDir string) error {
	filename := "cluster.jsonnet"
	clusterJson := types.Kr8ClusterJsonnet{
		ClusterSpec: cSpec,
		// Bug() Unsure if Path is correct
		Cluster:    types.Kr8Cluster{Name: cSpec.Name, Path: cSpec.ClusterDir},
		Components: map[string]types.Kr8ClusterComponentRef{},
	}
	_, err := util.WriteObjToJsonFile(filename, dstDir+"/"+cSpec.Name, clusterJson)

	return err
}

// Generate default component kr8_spec values and store in params.jsonnet.
// Based on the type:
//
// jsonnet: create a component.jsonnet file and reference it from the params.jsonnet file
//
// yml: leave a note in the params.jsonnet file about where and how the yml files can be referenced
//
// chart: generate a simple taskfile that handles vendoring the chart data
func GenerateComponentJsonnet(componentOptions Kr8InitOptions, dstDir string) error {
	compJson := types.Kr8ComponentJsonnet{
		Kr8Spec: types.Kr8ComponentSpec{
			Kr8_allparams:         false,
			Kr8_allclusters:       false,
			DisableOutputDirClean: false,
			Includes:              []types.Kr8ComponentSpecIncludeObject{},
			ExtFiles:              map[string]string{},
			JPaths:                []string{},
		},
		ReleaseName: strings.ReplaceAll(componentOptions.ComponentName, "_", "-"),
		Namespace:   "Default",
		Version:     "1.0.0",
		CalledFrom:  "",
	}
	switch componentOptions.ComponentType {
	case "jsonnet":
		compJson.Kr8Spec.Includes = append(
			compJson.Kr8Spec.Includes,
			types.Kr8ComponentSpecIncludeObject{File: "component.jsonnet", DestName: "component", DestExt: "yaml"},
		)
	case "yml":
		compJson.Kr8Spec.Includes = append(compJson.Kr8Spec.Includes,
			types.Kr8ComponentSpecIncludeObject{
				File:     "input.yml",
				DestDir:  "",
				DestName: "glhf",
				DestExt:  "yml",
			},
		)
	case "tpl":
		compJson.Kr8Spec.Includes = append(compJson.Kr8Spec.Includes,
			types.Kr8ComponentSpecIncludeObject{
				File:     "README.tpl",
				DestDir:  "docs",
				DestName: "ReadMe",
				DestExt:  "md",
			},
		)
	case "chart":
		break
	default:
		break
	}

	_, err := util.WriteObjToJsonFile("params.jsonnet", dstDir+"/"+componentOptions.ComponentName, compJson)

	return err
}

// Downloads a starter kr8p jsonnet lib from github.
// If fetch is true, it will download the repo in the /lib directory.
// If false, it will print the git commands to run.
// Repo: https://github.com/ice-bergtech/kr8-libsonnet .
// return util.FetchRepoUrl("https://github.com/ice-bergtech/kr8-libsonnet", dstDir+"/kr8-lib", !fetch).
func GenerateLib(fetch bool, dstDir string) error {
	if err := util.GenErrorIfCheck("error creating lib directory", os.MkdirAll(dstDir, 0750)); err != nil {
		return err
	}

	return util.FetchRepoUrl("https://github.com/kube-libsonnet/kube-libsonnet.git", dstDir+"/kube-lib", !fetch)
}

// Generates a starter readme for the repo, and writes it to the destination directory.
func GenerateReadme(dstDir string, cmdOptions Kr8InitOptions, clusterSpec types.Kr8ClusterSpec) error {
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
		"2. Define tiered cluster configuration in the `" + clusterSpec.ClusterDir + "` directory.",
		"3. Run `kr8p generate` to generate component configuration files.",
		"",
		"## Info ",
		"",
		"This project is initialized with the following parameters:",
		"",
		"	* ClusterName: `" + cmdOptions.ClusterName + "`",
		"	* Fetch External Libs: " + fetch,
		"   * Cluster config root directory: `" + clusterSpec.ClusterDir + "`",
		"   * Component root directory: `components`",
		"   * Cluster config root directory: `" + clusterSpec.ClusterDir + "`",
		"   * Generated config outpu directory: `" + clusterSpec.GenerateDir + "`",
		"",
		"Generated using [kr8+](https://github.com/ice-bergtech/kr8)",
	}, "\n")

	return os.WriteFile(dstDir+"/Readme.md", []byte(readmeTemplate), 0600)
}
