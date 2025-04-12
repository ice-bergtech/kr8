package kr8_init

import (
	"os"
	"path/filepath"
	"strings"

	types "github.com/ice-bergtech/kr8/pkg/types"
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
		Namespace:   "default",
		Version:     "1.0.0",
		CalledFrom:  "std.thisFile",
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
		folderDir := filepath.Join(dstDir, componentOptions.ComponentName)
		if err := os.MkdirAll(folderDir, 0750); err != nil {
			log.Error().Err(err).Msg("component directory not created")
		}
		if err := GenerateChartJsonnet(compJson, componentOptions, folderDir); err != nil {
			log.Error().Err(err).Msg("component directory not created")
		}
		if err := GenerateChartTaskfile(compJson, componentOptions, folderDir); err != nil {
			log.Error().Err(err).Msg("component directory not created")
		}
		compJson.Kr8Spec.Includes = append(compJson.Kr8Spec.Includes,
			types.Kr8ComponentSpecIncludeObject{
				File:     componentOptions.ComponentName + ".jsonnet",
				DestDir:  "",
				DestName: componentOptions.ComponentName,
				DestExt:  "yml",
			},
		)
	default:
		break
	}

	_, err := util.WriteObjToJsonFile("params.jsonnet", dstDir+"/"+componentOptions.ComponentName, compJson)

	return err
}

func GenerateChartJsonnet(compJson types.Kr8ComponentJsonnet, componentOptions Kr8InitOptions, folderDir string) error {
	chartJsonnetText := strings.Join([]string{
		"# This loads the component configuration into the `config` var",
		"local config = std.extVar('kr8');",
		"",
		"local helm_template = std.native('helmTemplate')(config.release_name, './vendor/'+'" + componentOptions.ComponentName + "-'+config.Version, {",
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

	return os.WriteFile(filepath.Join(folderDir, componentOptions.ComponentName+".jsonnet"), []byte(chartJsonnetText), 0600)
}

func GenerateChartTaskfile(comp types.Kr8ComponentJsonnet, componentOptions Kr8InitOptions, folderDir string) error {
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
