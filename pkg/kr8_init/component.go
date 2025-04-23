package kr8_init

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/ice-bergtech/kr8/pkg/kr8_types"
	util "github.com/ice-bergtech/kr8/pkg/util"
	"github.com/rs/zerolog/log"
)

// Generate default component kr8_spec values and store in params.jsonnet.
// Based on the type:
//
// jsonnet: create a component.jsonnet file and reference it from the params.jsonnet file
//
// yml: leave a note in the params.jsonnet file about where and how the yml files can be referenced
//
// chart: generate a simple taskfile that handles vendoring the chart data
func GenerateComponentJsonnet(componentOptions Kr8InitOptions, dstDir string) error {
	compJson := kr8_types.Kr8ComponentJsonnet{
		Kr8Spec: kr8_types.Kr8ComponentSpec{
			Kr8_allparams:         false,
			Kr8_allclusters:       false,
			DisableOutputDirClean: false,
			Includes:              []kr8_types.Kr8ComponentSpecIncludeObject{},
			ExtFiles:              map[string]string{},
			JPaths:                []string{},
		},
		ReleaseName: strings.ReplaceAll(componentOptions.ComponentName, "_", "-"),
		Namespace:   "default",
		Version:     "1.0.0",
		CalledFrom:  "std.thisFile",
	}
	switch componentOptions.ComponentType {
	case "jsonnet":
		return InitComponentJsonnet(compJson, dstDir, componentOptions)
	case "yml":
		return InitComponentYaml(compJson, dstDir, componentOptions)
	case "tpl":
		return InitComponentTemplate(compJson, dstDir, componentOptions)
	case "chart":
		return InitComponentChart(dstDir, componentOptions, compJson)
	default:
		break
	}

	return nil
}

// Initializes the basic parts of a helm chart component.
func InitComponentChart(dstDir string, componentOptions Kr8InitOptions, compJson kr8_types.Kr8ComponentJsonnet) error {
	folderDir := filepath.Join(dstDir, componentOptions.ComponentName)
	if err := os.MkdirAll(folderDir, 0750); err != nil {
		log.Error().Err(err).Msg("component directory not created")
	}
	if err := GenerateChartJsonnet(compJson, componentOptions, folderDir); err != nil {
		log.Error().Err(err).Msg("component jsonnet not created")
	}
	if err := GenerateChartTaskfile(compJson, componentOptions, folderDir); err != nil {
		log.Error().Err(err).Msg("component taskfile not created")
	}
	compJson.Kr8Spec.Includes = append(compJson.Kr8Spec.Includes,
		kr8_types.Kr8ComponentSpecIncludeObject{
			File:     componentOptions.ComponentName + "-chart.jsonnet",
			DestDir:  "",
			DestName: componentOptions.ComponentName,
			DestExt:  "yml",
		},
	)
	_, err := util.WriteObjToJsonFile("params.jsonnet", dstDir+"/"+componentOptions.ComponentName, compJson)

	return err
}

// Initializes the based parts of a template-based component.
func InitComponentTemplate(
	compJson kr8_types.Kr8ComponentJsonnet,
	dstDir string,
	componentOptions Kr8InitOptions,
) error {
	compJson.Kr8Spec.Includes = append(compJson.Kr8Spec.Includes,
		kr8_types.Kr8ComponentSpecIncludeObject{
			File:     "README.tpl",
			DestDir:  "docs",
			DestName: "ReadMe",
			DestExt:  "md",
		},
	)
	_, err := util.WriteObjToJsonFile("params.jsonnet", dstDir+"/"+componentOptions.ComponentName, compJson)

	return err
}

// Initializes the basic parts of a yaml-based component.
func InitComponentYaml(compJson kr8_types.Kr8ComponentJsonnet, dstDir string, componentOptions Kr8InitOptions) error {
	compJson.Kr8Spec.Includes = append(compJson.Kr8Spec.Includes,
		kr8_types.Kr8ComponentSpecIncludeObject{
			File:     "input.yml",
			DestDir:  "",
			DestName: "glhf",
			DestExt:  "yml",
		},
	)
	_, err := util.WriteObjToJsonFile("params.jsonnet", dstDir+"/"+componentOptions.ComponentName, compJson)

	return err
}

// Initializes the basic parts of a jsonnet-based component.
func InitComponentJsonnet(
	compJson kr8_types.Kr8ComponentJsonnet,
	dstDir string,
	componentOptions Kr8InitOptions,
) error {
	compJson.Kr8Spec.Includes = append(
		compJson.Kr8Spec.Includes,
		kr8_types.Kr8ComponentSpecIncludeObject{File: "component.jsonnet", DestName: "component", DestExt: "yaml"},
	)
	_, err := util.WriteObjToJsonFile("params.jsonnet", dstDir+"/"+componentOptions.ComponentName, compJson)

	return err
}

// Generates a jsonnet files that references a local helm chart.
func GenerateChartJsonnet(
	compJson kr8_types.Kr8ComponentJsonnet,
	componentOptions Kr8InitOptions,
	folderDir string,
) error {
	chartJsonnetText := strings.Join([]string{
		"# This loads the component configuration into the `config` var",
		"local config = std.extVar('kr8');",
		"",
		"local helm_template = std.native('helmTemplate') " +
			"(config.release_name, './vendor/'+'" + componentOptions.ComponentName + "-'+config.Version, {",
		"		calledFrom: std.thisFile,",
		"		namespace: config.namespace,",
		"		values:  if 'helm_values' in config then config.helm_values else {},",
		"});",
		"",
		"[",
		"		object",
		"		for object in std.objectValues(helm_template)",
		"		if 'kind' in object && object.kind != 'Secret'",
		"]",
	}, "\n")

	return os.WriteFile(
		filepath.Join(folderDir, componentOptions.ComponentName+"-chart.jsonnet"),
		[]byte(chartJsonnetText), 0600,
	)
}

// Generates a go-task taskfile that's setup to download a helm chart into a local `vendor` directory.
func GenerateChartTaskfile(
	comp kr8_types.Kr8ComponentJsonnet,
	componentOptions Kr8InitOptions,
	folderDir string,
) error {
	taskfileText := strings.Join([]string{
		"# https://taskfile.dev/usage",
		"version: '3'",
		"",
		"vars:",
		"  CHART_NAME: '" + componentOptions.ComponentName + "'",
		"  CHART_REPO: ''",
		"",
		"tasks:",
		"  default:",
		"    cmds:",
		"      - task: fetch-" + comp.Version + "",
		"",
		"  # Copy/paste this for each new version",
		"  fetch-" + comp.Version + ":",
		"    desc: 'fetch chart version " + comp.Version + "'",
		"    vars:",
		"      CHART_VER: '" + comp.Version + "'",
		"    cmds:",
		"      - task: fetch-chart",
		"        vars: { VER: '{{.CHART_VER}}'}",
		"",
		"  fetch-chart:",
		"    desc: 'fetch a helm chart'",
		"    vars:",
		"      VER: '{{default 'unset' .VER}}'",
		"    cmds:",
		"      - mkdir -p ./vendor/{{.CHART_NAME}}-{{.VER}} && rm -rf ./vendor/{{.CHART_NAME}}-{{.VER}}/*",
		"      - mkdir -p ./vendor/tmp && rm -rf ./vendor/tmp/*",
		"      # add the helm repo and fetch it locally into vendor directory",
		"      - helm fetch --repo {{.CHART_REPO}} --untar --untardir ./vendor/tmp --version '{{.VER}}' '{{.CHART_NAME}}'",
		"      - mv ./vendor/tmp/{{.CHART_NAME}}/* ./vendor/{{.CHART_NAME}}-{{.VER}}/ && rm -rf ./vendor/tmp",
	}, "\n")

	return os.WriteFile(filepath.Join(folderDir, "Taskfile.yml"), []byte(taskfileText), 0600)
}
