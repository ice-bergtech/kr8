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

type CmdGetOptions struct {
	ClusterParams string
	NoTable       bool
	FieldName     string
	Cluster       string
	Component     string
	ParamField    string
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
	getParamsCmd.Flags().StringVarP(&cmdGetFlags.ParamField, "param", "P", "", "return value of json param from supplied path")

}

var getClustersCmd = &cobra.Command{
	Use:   "clusters",
	Short: "Get all clusters",
	Long:  "Get all clusters defined in kr8 config hierarchy",
	Run: func(cmd *cobra.Command, args []string) {

		clusters, err := getClusters(rootConfig.ClusterDir)
		fatalErrorCheck(err, "Error getting clusters")

		if cmdGetFlags.NoTable {
			for _, c := range clusters.Cluster {
				println(c.Name + ": " + c.Path)
			}
			return
		}

		var entry []string
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Name", "Path"})

		for _, c := range clusters.Cluster {
			entry = append(entry, c.Name)
			entry = append(entry, c.Path)
			table.Append(entry)
			entry = entry[:0]
		}
		table.Render()

	},
}

var getComponentsCmd = &cobra.Command{
	Use:   "components",
	Short: "Get all components",
	Long:  "Get all available components defined in the kr8 config hierarchy",
	Run: func(cmd *cobra.Command, args []string) {

		if cmdGetFlags.Cluster == "" && cmdGetFlags.ClusterParams == "" {
			log.Fatal().Msg("Please specify a --cluster name and/or --clusterparams file")
		}

		var params []string
		if cmdGetFlags.Cluster != "" {
			clusterPath := getCluster(rootConfig.ClusterDir, cmdGetFlags.Cluster)
			params = getClusterParams(rootConfig.ClusterDir, clusterPath)
		}
		if cmdGetFlags.ClusterParams != "" {
			params = append(params, cmdGetFlags.ClusterParams)
		}

		j := renderJsonnet(rootConfig.VMConfig, params, "._components", true, "", "components")
		if cmdGetFlags.ParamField != "" {
			value := gjson.Get(j, cmdGetFlags.ParamField)
			if value.String() == "" {
				log.Fatal().Msg("Error getting param: " + cmdGetFlags.ParamField)
			} else {
				formatted := Pretty(j, rootConfig.Color)
				fmt.Println(formatted)
			}
		} else {
			formatted := Pretty(j, rootConfig.Color)
			fmt.Println(formatted)
		}
	},
}

var getParamsCmd = &cobra.Command{
	Use:   "params",
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

		params := renderClusterParams(rootConfig.VMConfig, cmdGetFlags.Cluster, cList, cmdGetFlags.ClusterParams, true)

		if cmdGetFlags.ParamField != "" {
			value := gjson.Get(params, cmdGetFlags.ParamField)
			notUnset, _ := cmd.Flags().GetBool("notunset")
			if notUnset && value.String() == "" {
				log.Fatal().Msg("Error getting param: " + cmdGetFlags.ParamField)
			} else {
				fmt.Println(value) // no formatting because this isn't always json, this is just the value of a field
			}
		} else {
			formatted := Pretty(params, rootConfig.Color)
			fmt.Println(formatted)
		}

	},
}
