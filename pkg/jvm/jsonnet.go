/*
This code was originally copied almost verbatim from the kubecfg project: https://github.com/ksonnet/kubecfg

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
package jvm

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/exp/maps"

	jsonnet "github.com/google/go-jsonnet"
	"github.com/rs/zerolog/log"
	"github.com/tidwall/gjson"

	types "github.com/ice-bergtech/kr8/pkg/types"
	util "github.com/ice-bergtech/kr8/pkg/util"
)

// Create a Jsonnet VM to run commands in
func JsonnetVM(vmconfig types.VMConfig) (*jsonnet.VM, error) {
	vm := jsonnet.MakeVM()
	RegisterNativeFuncs(vm)

	// always add lib directory in base directory to path
	jpath := []string{vmconfig.BaseDir + "/lib"}

	jpath = append(jpath, filepath.SplitList(os.Getenv("KR8_JPATH"))...)
	jpathArgs := vmconfig.Jpaths
	jpath = append(jpath, jpathArgs...)

	vm.Importer(&jsonnet.FileImporter{
		JPaths: jpath,
	})

	for _, extvar := range vmconfig.ExtVars {
		kv := strings.SplitN(extvar, "=", 2)
		if len(kv) != 2 {
			log.Fatal().Str("ext-str-file", extvar).Msg("Failed to parse. Missing '=' in parameter`")
		}
		v, err := os.ReadFile(kv[1])
		if err != nil {
			panic(err)
		}
		vm.ExtVar(kv[0], string(v))
	}
	return vm, nil
}

// Takes a list of jsonnet files and imports each one and mixes them with "+"
func JsonnetRenderFiles(vmConfig types.VMConfig, files []string, param string, prune bool, prepend string, source string) string {

	// copy the slice so that we don't unitentionally modify the original
	jsonnetPaths := make([]string, len(files[:0]))
	copy(jsonnetPaths, files[:0])

	// range through the files
	for _, s := range files {
		jsonnetPaths = append(jsonnetPaths, fmt.Sprintf("(import '%s')", s))
	}

	// Create a JSonnet VM
	vm, err := JsonnetVM(vmConfig)
	util.FatalErrorCheck(err, "Error creating jsonnet VM")

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

	// render the jsonnet
	out, err := vm.EvaluateAnonymousSnippet(source, jsonnetImport)
	util.FatalErrorCheck(err, "Error evaluating jsonnet snippet")

	return out

}

// Renders a jsonnet file with the specified options.
func JsonnetRender(cmdFlagsJsonnet types.CmdJsonnetOptions, filename string, vmConfig types.VMConfig) {
	// Check if cluster and/or clusterparams are specified
	if cmdFlagsJsonnet.Cluster == "" && cmdFlagsJsonnet.ClusterParams == "" {
		log.Fatal().Msg("Please specify a --cluster name and/or --clusterparams")
	}

	// Render the cluster parameters
	config := JsonnetRenderClusterParams(vmConfig, cmdFlagsJsonnet.Cluster, []string{cmdFlagsJsonnet.Component}, cmdFlagsJsonnet.ClusterParams, false)

	// Create a new VM instance
	vm, _ := JsonnetVM(vmConfig)
	// Setup kr8 config as external vars
	vm.ExtCode("kr8_cluster", "std.prune("+config+"._cluster)")
	vm.ExtCode("kr8_components", "std.prune("+config+"._components)")
	vm.ExtCode("kr8", "std.prune("+config+"."+cmdFlagsJsonnet.Component+")")
	vm.ExtCode("kr8_unpruned", config+"."+cmdFlagsJsonnet.Component)

	var input string
	// If pruning is enabled, prune the input before rendering
	// This removes all null and empty fields from the imported file
	if cmdFlagsJsonnet.Prune {
		input = "std.prune(import '" + filename + "')"
	} else {
		input = "( import '" + filename + "')"
	}

	//
	// Evaluate the jsonnet snippet and print the result
	// This is where the magic happens! The jsonnet code is evaluated and the result is stored
	//
	j, err := vm.EvaluateAnonymousSnippet("file", input)
	util.FatalErrorCheck(err, "Error evaluating jsonnet snippet")

	util.JsonnetPrint(j, cmdFlagsJsonnet.Format, cmdFlagsJsonnet.Color)
}

// Only render cluster params (_cluster), without components
func JsonnetRenderClusterParamsOnly(vmconfig types.VMConfig, clusterName string, clusterParams string, prune bool) string {
	var params []string
	if clusterName != "" {
		clusterPath := util.GetClusterPaths(vmconfig.BaseDir, clusterName)
		params = util.GetClusterParamsFilenames(vmconfig.BaseDir, clusterPath)
	}
	if clusterParams != "" {
		params = append(params, clusterParams)
	}
	renderedParams := JsonnetRenderFiles(vmconfig, params, "._cluster", prune, "", "clusterparams")

	return renderedParams
}

// Render cluster params, merged with one or more component's parameters. Empty componentName list renders all component parameters
func JsonnetRenderClusterParams(vmconfig types.VMConfig, clusterName string, componentNames []string, clusterParams string, prune bool) string {
	if clusterName == "" && clusterParams == "" {
		log.Fatal().Msg("Please specify a --cluster name and/or --clusterparams")
	}

	var params []string
	var componentMap map[string]types.Kr8ClusterComponentRef

	if clusterName != "" {
		clusterPath := util.GetClusterPaths(vmconfig.BaseDir, clusterName)
		params = util.GetClusterParamsFilenames(vmconfig.BaseDir, clusterPath)
	}
	if clusterParams != "" {
		params = append(params, clusterParams)
	}

	compParams := JsonnetRenderFiles(vmconfig, params, "", true, "", "clusterparams")

	compString := gjson.Get(compParams, "_components")
	err := json.Unmarshal([]byte(compString.String()), &componentMap)
	util.FatalErrorCheck(err, "failed to parse component map")
	componentDefaultsMerged := "{"

	listComponentKeys := maps.Keys(componentMap)
	if len(componentNames) > 0 {
		listComponentKeys = componentNames
	}

	// all components
	for _, key := range listComponentKeys {
		if value, ok := componentMap[key]; ok {
			path := vmconfig.BaseDir + "/" + value.Path + "/params.jsonnet"
			fileC, err := os.ReadFile(path)
			util.FatalErrorCheck(err, "Error reading file "+path)
			componentDefaultsMerged = componentDefaultsMerged + fmt.Sprintf("'%s': %s,", key, string(fileC))
		}
	}
	componentDefaultsMerged = componentDefaultsMerged + "}"

	compParams = JsonnetRenderFiles(vmconfig, params, "", prune, componentDefaultsMerged, "componentparams")

	return compParams
}
