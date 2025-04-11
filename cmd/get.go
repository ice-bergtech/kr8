// Copyright © 2019 kubecfg Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"github.com/ice-bergtech/kr8/pkg/jnetvm"
	util "github.com/ice-bergtech/kr8/pkg/util"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"

	"fmt"
	"os"

	"github.com/rs/zerolog/log"
)

// GetCmd represents the get command.
var GetCmd = &cobra.Command{
	Use:   "get",
	Short: "Display one or many kr8 resources",
	Long:  `Displays information about kr8 resources such as clusters and components`,
}

// Holds the options for the get command.
type CmdGetOptions struct {
	// ClusterParams provides a way to provide cluster params as a single file.
	// This can be combined with --cluster to override the cluster.
	ClusterParams string
	// If true, just prints result instead of placing in table
	NoTable bool
	// Field to display from the resource
	FieldName string
	// Cluster to get resources from
	Cluster string
	// Component to get resources from
	Component string
	// Param to display from the resource
	ParamField string
}

var cmdGetFlags CmdGetOptions

func init() {
	RootCmd.AddCommand(GetCmd)
	GetCmd.PersistentFlags().StringVarP(&cmdGetFlags.ClusterParams,
		"clusterparams", "p", "",
		"provide cluster params as single file - can be combined with --cluster to override cluster")

	// clusters
	GetCmd.AddCommand(GetClustersCmd)
	GetClustersCmd.PersistentFlags().BoolVarP(&cmdGetFlags.NoTable,
		"raw", "r", false,
		"If true, just prints result instead of placing in table.")
	// components
	GetCmd.AddCommand(GetComponentsCmd)
	GetComponentsCmd.PersistentFlags().StringVarP(&cmdGetFlags.Cluster,
		"cluster", "C", "",
		"get components for cluster")

	// params
	GetCmd.AddCommand(GetParamsCmd)
	GetParamsCmd.PersistentFlags().StringVarP(&cmdGetFlags.Cluster,
		"cluster", "C", "",
		"get components for cluster")
	GetParamsCmd.PersistentFlags().StringVarP(&cmdGetFlags.Component,
		"component", "c", "",
		"component to render params for")
	GetParamsCmd.PersistentFlags().StringVarP(&cmdGetFlags.ParamField,
		"param", "P", "",
		"return value of json param from supplied path")
}

var GetClustersCmd = &cobra.Command{
	Use:   "clusters [flags]",
	Short: "Get all clusters",
	Long:  "Get all clusters defined in kr8 config hierarchy",
	Run: func(cmd *cobra.Command, args []string) {

		clusters, err := util.GetClusterFilenames(RootConfig.ClusterDir)
		util.FatalErrorCheck("Error getting clusters", err)

		if cmdGetFlags.NoTable {
			for _, c := range clusters {
				println(c.Name + ": " + c.Path)
			}

			return
		}

		var entry []string
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Name", "Path"})

		for _, c := range clusters {
			entry = append(entry, c.Name)
			entry = append(entry, c.Path)
			table.Append(entry)
			entry = entry[:0]
		}
		table.Render()

	},
}

var GetComponentsCmd = &cobra.Command{
	Use:   "components [flags]",
	Short: "Get all components",
	Long:  "Get all available components defined in the kr8 config hierarchy",
	Run: func(cmd *cobra.Command, args []string) {

		if cmdGetFlags.Cluster == "" && cmdGetFlags.ClusterParams == "" {
			log.Fatal().Msg("Please specify a --cluster name and/or --clusterparams file")
		}

		var params []string
		if cmdGetFlags.Cluster != "" {
			clusterPath, err := util.GetClusterPaths(RootConfig.ClusterDir, cmdGetFlags.Cluster)
			util.FatalErrorCheck("error getting cluster path for "+cmdGetFlags.Cluster, err)
			params = util.GetClusterParamsFilenames(RootConfig.ClusterDir, clusterPath)
		}
		if cmdGetFlags.ClusterParams != "" {
			params = append(params, cmdGetFlags.ClusterParams)
		}

		jvm, err := jnetvm.JsonnetRenderFiles(RootConfig.VMConfig, params, "._components", true, "", "components")
		util.FatalErrorCheck("error rendering jsonnet files", err)
		if cmdGetFlags.ParamField != "" {
			value := gjson.Get(jvm, cmdGetFlags.ParamField)
			if value.String() == "" {
				log.Fatal().Msg("Error getting param: " + cmdGetFlags.ParamField)
			} else {
				formatted, err := util.Pretty(jvm, RootConfig.Color)
				util.FatalErrorCheck("error pretty printing jsonnet", err)
				fmt.Println(formatted)
			}
		} else {
			formatted, err := util.Pretty(jvm, RootConfig.Color)
			util.FatalErrorCheck("error pretty printing jsonnet", err)
			fmt.Println(formatted)
		}
	},
}

var GetParamsCmd = &cobra.Command{
	Use:   "params [flags]",
	Short: "Get parameter for components and clusters",
	Long:  "Get parameters assigned to clusters and components in the kr8 config hierarchy",
	Run: func(cmd *cobra.Command, args []string) {
		if cmdGetFlags.Cluster == "" {
			log.Fatal().Msg("Please specify a --cluster")
		}

		var cList []string
		if cmdGetFlags.Component != "" {
			cList = append(cList, cmdGetFlags.Component)
		}

		params, err := jnetvm.JsonnetRenderClusterParams(
			RootConfig.VMConfig,
			cmdGetFlags.Cluster,
			cList,
			cmdGetFlags.ClusterParams,
			true,
		)
		util.FatalErrorCheck("error rendering cluster params", err)

		// if we're not filtering the output, just pretty print and finish
		if cmdGetFlags.ParamField == "" {
			if cmdGetFlags.Component != "" {
				result := gjson.Get(params, cmdGetFlags.Component).String()
				formatted, err := util.Pretty(result, RootConfig.Color)
				util.FatalErrorCheck("error pretty printing jsonnet", err)
				fmt.Println(formatted)
			} else {
				formatted, err := util.Pretty(params, RootConfig.Color)
				util.FatalErrorCheck("error pretty printing jsonnet", err)
				fmt.Println(formatted)
			}

			return
		}

		// Filter on component name first, then field name
		if cmdGetFlags.ParamField != "" {
			value := gjson.Get(params, cmdGetFlags.ParamField)
			if value.String() == "" {
				log.Fatal().Msg("Error getting param: " + cmdGetFlags.ParamField)
			}
			// no formatting because this isn't always json, this is just the value of a field
			fmt.Println(value)
		}
	},
}
