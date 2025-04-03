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
	jvm "github.com/ice-bergtech/kr8/pkg/jvm"
	util "github.com/ice-bergtech/kr8/pkg/util"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"

	"fmt"
	"os"

	"github.com/rs/zerolog/log"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Display one or many kr8 resources",
	Long:  `Displays information about kr8 resources such as clusters and components`,
}

// Holds the options for the get command.
type CmdGetOptions struct {
	// ClusterParams provides a way to provide cluster params as a single file. This can be combined with --cluster to override the cluster.
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
	RootCmd.AddCommand(getCmd)
	getCmd.PersistentFlags().StringVarP(&cmdGetFlags.ClusterParams, "clusterparams", "p", "", "provide cluster params as single file - can be combined with --cluster to override cluster")

	// clusters
	getCmd.AddCommand(getClustersCmd)
	getClustersCmd.PersistentFlags().BoolVarP(&cmdGetFlags.NoTable, "raw", "r", false, "If true, just prints result instead of placing in table.")
	// components
	getCmd.AddCommand(getComponentsCmd)
	getComponentsCmd.PersistentFlags().StringVarP(&cmdGetFlags.Cluster, "cluster", "C", "", "get components for cluster")

	// params
	getCmd.AddCommand(getParamsCmd)
	getParamsCmd.PersistentFlags().StringVarP(&cmdGetFlags.Cluster, "cluster", "C", "", "get components for cluster")
	getParamsCmd.PersistentFlags().StringVarP(&cmdGetFlags.Component, "component", "c", "", "component to render params for")
	getParamsCmd.PersistentFlags().StringVarP(&cmdGetFlags.ParamField, "param", "P", "", "return value of json param from supplied path")
}

var getClustersCmd = &cobra.Command{
	Use:   "clusters [flags]",
	Short: "Get all clusters",
	Long:  "Get all clusters defined in kr8 config hierarchy",
	Run: func(cmd *cobra.Command, args []string) {

		clusters, err := util.GetClusters(rootConfig.ClusterDir)
		util.FatalErrorCheck(err, "Error getting clusters")

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

var getComponentsCmd = &cobra.Command{
	Use:   "components [flags]",
	Short: "Get all components",
	Long:  "Get all available components defined in the kr8 config hierarchy",
	Run: func(cmd *cobra.Command, args []string) {

		if cmdGetFlags.Cluster == "" && cmdGetFlags.ClusterParams == "" {
			log.Fatal().Msg("Please specify a --cluster name and/or --clusterparams file")
		}

		var params []string
		if cmdGetFlags.Cluster != "" {
			clusterPath := util.GetCluster(rootConfig.ClusterDir, cmdGetFlags.Cluster)
			params = util.GetClusterParams(rootConfig.ClusterDir, clusterPath)
		}
		if cmdGetFlags.ClusterParams != "" {
			params = append(params, cmdGetFlags.ClusterParams)
		}

		j := jvm.RenderJsonnet(rootConfig.VMConfig, params, "._components", true, "", "components")
		if cmdGetFlags.ParamField != "" {
			value := gjson.Get(j, cmdGetFlags.ParamField)
			if value.String() == "" {
				log.Fatal().Msg("Error getting param: " + cmdGetFlags.ParamField)
			} else {
				formatted := util.Pretty(j, rootConfig.Color)
				fmt.Println(formatted)
			}
		} else {
			formatted := util.Pretty(j, rootConfig.Color)
			fmt.Println(formatted)
		}
	},
}

var getParamsCmd = &cobra.Command{
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

		params := jvm.RenderClusterParams(rootConfig.VMConfig, cmdGetFlags.Cluster, cList, cmdGetFlags.ClusterParams, true)

		// if we're not filtering the output, just pretty print and finish
		if cmdGetFlags.ParamField == "" {
			if cmdGetFlags.Component != "" {
				result := gjson.Get(params, cmdGetFlags.Component).String()
				fmt.Println(util.Pretty(result, rootConfig.Color))
			} else {
				fmt.Println(util.Pretty(params, rootConfig.Color))
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
