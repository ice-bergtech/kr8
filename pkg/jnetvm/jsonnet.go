/*
This code was originally copied almost verbatim from the kubecfg project:
	https://github.com/ksonnet/kubecfg -> https://github.com/kubecfg/kubecfg

Copyright 2018 ksonnet

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package jvm contains the jsonnet rendering logic.
package jnetvm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	//nolint:exptostd
	"golang.org/x/exp/maps"

	jsonnet "github.com/google/go-jsonnet"
	"github.com/google/go-jsonnet/linter"
	"github.com/rs/zerolog"
	"github.com/tidwall/gjson"

	"github.com/ice-bergtech/kr8/pkg/kr8_native_funcs"
	"github.com/ice-bergtech/kr8/pkg/kr8_types"
	types "github.com/ice-bergtech/kr8/pkg/types"
	util "github.com/ice-bergtech/kr8/pkg/util"
)

// Create a Jsonnet VM to run commands in.
// It:
//   - creates a jsonnet VM
//   - registers kr8+ native functions
//   - Add jsonnet library directories
//   - loads external files into extVars
func JsonnetVM(vmConfig types.VMConfig) (*jsonnet.VM, error) {
	jvm := jsonnet.MakeVM()
	kr8_native_funcs.RegisterNativeFuncs(jvm)

	// always add lib directory in base directory to path
	jpath := []string{filepath.Join(vmConfig.BaseDir, "lib")}

	jpath = append(jpath, filepath.SplitList(os.Getenv("KR8_JPATH"))...)
	jpathArgs := vmConfig.JPaths
	jpath = append(jpath, jpathArgs...)

	jvm.Importer(&jsonnet.FileImporter{
		JPaths: jpath,
	})

	for _, extVar := range vmConfig.ExtVars {
		args := strings.SplitN(extVar, "=", 2)
		if len(args) != 2 {
			return nil, types.Kr8Error{Message: "Failed to parse. Missing '=' in parameter`", Value: extVar}
		}
		v, err := os.ReadFile(args[1])
		if err != nil {
			return nil, err
		}
		jvm.ExtVar(args[0], string(v))
	}

	return jvm, nil
}

// Takes a list of jsonnet files and imports each one.
// Formats the string for jsonnet using "+".
// source is only used for error messages.
func JsonnetRenderFiles(
	vmConfig types.VMConfig,
	files []string,
	param string,
	prune bool,
	prepend string,
	source string,
	lint bool,
) (string, error) {
	// copy the slice so that we don't unintentionally modify the original
	jsonnetPaths := make([]string, len(files))

	// range through the files
	for idx, s := range files {
		jsonnetPaths[idx] = fmt.Sprintf("(import '%s')", s)
	}

	// Create a JSonnet VM
	jvm, err := JsonnetVM(vmConfig)
	if err := util.ErrorIfCheck("Error creating jsonnet VM", err); err != nil {
		return "", err
	}

	// Join the slices into a jsonnet compat string. Prepend code from "prepend" variable, if set.
	var jsonnetImport string
	if prepend != "" {
		jsonnetImport = prepend + "+" + strings.Join(jsonnetPaths, "+")
	} else {
		jsonnetImport = strings.Join(jsonnetPaths, "+")
	}

	if param != "" {
		jsonnetImport = "(" + jsonnetImport + ")" + param
	}

	if prune {
		// wrap in std.prune, to remove nulls, empty arrays and hashes
		jsonnetImport = "std.prune(" + jsonnetImport + ")"
	}

	if lint {
		snippets := []linter.Snippet{{FileName: source, Code: jsonnetImport}}
		var buffer bytes.Buffer
		if linter.LintSnippet(jvm, &buffer, snippets) {
			return "", types.Kr8Error{Message: "linting issue", Value: buffer.String()}
		}
	}

	// render the jsonnet
	out, err := jvm.EvaluateAnonymousSnippet(source, jsonnetImport)
	if err := util.ErrorIfCheck("Error evaluating jsonnet snippet", err); err != nil {
		return "", err
	}

	return out, nil
}

// Renders a jsonnet file with the specified options.
func JsonnetRender(
	cmdFlagsJsonnet types.CmdJsonnetOptions,
	filename string,
	vmConfig types.VMConfig,
	logger zerolog.Logger,
) error {
	// Check if cluster and/or clusterparams are specified
	if cmdFlagsJsonnet.Cluster == "" && cmdFlagsJsonnet.ClusterParams == "" {
		return types.Kr8Error{Message: "Please specify a --cluster name and/or --clusterparams", Value: ""}
	}

	// Render the cluster parameters
	config, err := JsonnetRenderClusterParams(
		vmConfig,
		cmdFlagsJsonnet.Cluster,
		[]string{cmdFlagsJsonnet.Component},
		cmdFlagsJsonnet.ClusterParams,
		false,
		cmdFlagsJsonnet.Lint,
	)
	if err := util.ErrorIfCheck("error rendering cluster params", err); err != nil {
		return err
	}

	// Create a new VM instance
	jvm, _ := JsonnetVM(vmConfig)
	// Setup kr8+ config as external vars
	jvm.ExtCode("kr8_cluster", "std.prune("+config+"._cluster)")
	jvm.ExtCode("kr8_components", "std.prune("+config+"._components)")
	jvm.ExtCode("kr8", "std.prune("+config+"."+cmdFlagsJsonnet.Component+")")
	jvm.ExtCode("kr8_unpruned", config+"."+cmdFlagsJsonnet.Component)

	var input string
	// If pruning is enabled, prune the input before rendering
	// This removes all null and empty fields from the imported file
	if cmdFlagsJsonnet.Prune {
		input = "std.prune(import '" + filename + "')"
	} else {
		input = "( import '" + filename + "')"
	}

	logger.Debug().Msg("Processing file through jsonnet vm: " + input)

	//
	// Evaluate the jsonnet snippet and print the result
	// This is where the magic happens! The jsonnet code is evaluated and the result is stored
	//
	j, err := jvm.EvaluateAnonymousSnippet("file", input)
	if err := util.ErrorIfCheck("Error evaluating jsonnet snippet", err); err != nil {
		return err
	}

	return util.JsonnetPrint(j, cmdFlagsJsonnet.Format, cmdFlagsJsonnet.Color)
}

// Only render cluster params (_cluster), without components.
func JsonnetRenderClusterParamsOnly(
	vmConfig types.VMConfig,
	clusterName string,
	clusterParams string,
	prune bool,
	lint bool,
) (string, error) {
	var params []string
	if clusterName != "" {
		clusterPath, err := util.GetClusterPath(vmConfig.BaseDir, clusterName)
		if err != nil {
			return "", err
		}
		params = util.GetClusterParamsFilenames(vmConfig.BaseDir, clusterPath)
	}
	if clusterParams != "" {
		params = append(params, clusterParams)
	}

	return JsonnetRenderFiles(vmConfig, params, "._cluster", prune, "", "clusterparams", lint)
}

// Render cluster params, merged with one or more component's parameters.
// Empty componentName list renders all component parameters.
func JsonnetRenderClusterParams(
	vmConfig types.VMConfig,
	clusterName string,
	componentNames []string,
	clusterParams string,
	prune bool,
	lint bool,
) (string, error) {
	if clusterName == "" && clusterParams == "" {
		return "", types.Kr8Error{Message: "Please specify a --cluster name and/or --clusterparams", Value: ""}
	}

	var params []string
	var componentMap map[string]kr8_types.Kr8ClusterComponentRef

	if clusterName != "" {
		clusterPath, err := util.GetClusterPath(vmConfig.BaseDir, clusterName)
		if err != nil {
			return "", err
		}
		params = util.GetClusterParamsFilenames(vmConfig.BaseDir, clusterPath)
	}
	if clusterParams != "" {
		params = append(params, clusterParams)
	}

	compParams, err := JsonnetRenderFiles(vmConfig, params, "", true, "", "clusterparams", lint)
	if err := util.ErrorIfCheck("failed to render cluster params", err); err != nil {
		return "", err
	}

	compString := gjson.Get(compParams, "_components")
	err = json.Unmarshal([]byte(compString.String()), &componentMap)
	if err := util.ErrorIfCheck("failed to parse component map", err); err != nil {
		return "", err
	}

	// all components
	componentDefaultsMerged, err := MergeComponentDefaults(componentMap, componentNames, vmConfig)
	if err != nil {
		return "", util.ErrorIfCheck("failed to merge component defaults", err)
	}

	return JsonnetRenderFiles(vmConfig, params, "", prune, componentDefaultsMerged, "component params", lint)
}

func MergeComponentDefaults(
	componentMap map[string]kr8_types.Kr8ClusterComponentRef,
	componentNames []string,
	vmConfig types.VMConfig,
) (string, error) {
	componentDefaultsMerged := "{"

	//nolint:exptostd
	listComponentKeys := maps.Keys(componentMap)
	if len(componentNames) > 0 {
		listComponentKeys = componentNames
	}

	for _, key := range listComponentKeys {
		if value, ok := componentMap[key]; ok {
			path := filepath.Join(vmConfig.BaseDir, value.Path, "params.jsonnet")
			fileC, err := os.ReadFile(filepath.Clean(path))
			if err := util.ErrorIfCheck("Error reading file "+path, err); err != nil {
				return "", err
			}
			componentDefaultsMerged += fmt.Sprintf("'%s': %s,", key, string(fileC))
		}
	}
	componentDefaultsMerged += "}"

	return componentDefaultsMerged, nil
}
