// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
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
	"github.com/spf13/viper"
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

var (
	printRaw bool
)

var getClustersCmd = &cobra.Command{
	Use:   "clusters",
	Short: "Get all clusters",
	Long:  "Get all clusters defined in kr8 config hierarchy",
	Run: func(cmd *cobra.Command, args []string) {

		clusters, err := getClusters(clusterDir)

		if err != nil {
			log.Fatal().Err(err).Msg("Error getting cluster")
		}

		if printRaw {
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

		viper.BindPFlag("cluster", clusterCmd.PersistentFlags().Lookup("cluster"))
		clusterName := viper.GetString("cluster")

		if clusterName == "" && clusterParams == "" {
			log.Fatal().Msg("Please specify a --cluster name and/or --clusterparams")
		}

		var params []string
		if clusterName != "" {
			clusterPath := getCluster(clusterDir, clusterName)
			params = getClusterParams(clusterDir, clusterPath)
		}
		if clusterParams != "" {
			params = append(params, clusterParams)
		}

		j := renderJsonnet(cmd, params, "._components", true, "", "components")
		if paramPath != "" {
			value := gjson.Get(j, paramPath)
			if value.String() == "" {
				log.Fatal().Msg("Error getting param: " + paramPath)
			} else {
				formatted := Pretty(j, colorOutput)
				fmt.Println(formatted)
			}
		} else {
			formatted := Pretty(j, colorOutput)
			fmt.Println(formatted)
		}
	},
}

var getParamsCmd = &cobra.Command{
	Use:   "params",
	Short: "Get parameter for components and clusters",
	Long:  "Get parameters assigned to clusters and components in the kr8 config hierarchy",
	Run: func(cmd *cobra.Command, args []string) {

		clusterName := cluster

		if clusterName == "" {
			log.Fatal().Msg("Please specify a --cluster")
		}

		var cList []string
		if componentName != "" {
			cList = append(cList, componentName)
		}

		j := renderClusterParams(cmd, clusterName, cList, clusterParams, true)

		if paramPath != "" {
			value := gjson.Get(j, paramPath)
			notUnset, _ := cmd.Flags().GetBool("notunset")
			if notUnset && value.String() == "" {
				log.Fatal().Msg("Error getting param: " + paramPath)
			} else {
				fmt.Println(value) // no formatting because this isn't always json, this is just the value of a field
			}
		} else {
			formatted := Pretty(j, colorOutput)
			fmt.Println(formatted)
		}

	},
}

func init() {
	RootCmd.AddCommand(getCmd)

	// clusters
	getCmd.AddCommand(getClustersCmd)
	getClustersCmd.PersistentFlags().BoolVarP(&printRaw, "raw", "r", false, "If true, just prints result instead of placing in table.")
	getClustersCmd.PersistentFlags().StringVarP(&clusterParams, "clusterparams", "", "", "provide cluster params as single file - can be combined with --cluster to override cluster")
	// components
	getCmd.AddCommand(getComponentsCmd)
	getComponentsCmd.PersistentFlags().StringVarP(&cluster, "cluster", "C", "", "get components for cluster")
	// params
	getCmd.AddCommand(getParamsCmd)
	getParamsCmd.PersistentFlags().StringVarP(&cluster, "cluster", "C", "", "get components for cluster")
	getParamsCmd.PersistentFlags().StringVarP(&componentName, "component", "c", "", "component to render params for")
	getParamsCmd.Flags().StringVarP(&paramPath, "param", "P", "", "return value of json param from supplied path")

}
