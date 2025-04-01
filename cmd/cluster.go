package cmd

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
)

// clusterCmd represents the cluster command
var clusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "Operate on kr8 clusters",
	Long:  `Manage, list and generate kr8 cluster configurations at the cluster scope`,
	//Run: func(cmd *cobra.Command, args []string) { },
}

var (
	printRaw bool
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List Clusters",
	Long:  "List Clusters in kr8 config hierarchy",
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

var paramsCmd = &cobra.Command{
	Use:   "params",
	Short: "Show Cluster Params",
	Long:  "Show cluster params in kr8 config hierarchy",
	Run: func(cmd *cobra.Command, args []string) {
		viper.BindPFlag("cluster", clusterCmd.PersistentFlags().Lookup("cluster"))
		clusterName := viper.GetString("cluster")

		if clusterName == "" && clusterParams == "" {
			log.Fatal().Msg("Please specify a --cluster name and/or --clusterparams")
		}

		var cList []string
		if componentName != "" {
			cList = append(cList, componentName)
		}
		j := renderClusterParams(cmd, clusterName, cList, clusterParams, false)

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

var componentsCmd = &cobra.Command{
	Use:   "components",
	Short: "Show Cluster Components",
	Long:  "Show the components to be installed in the cluster in the kr8 hierarchy",
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

		j := renderJsonnet(cmd, params, "._components", true, "", "clustercomponents")
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

func init() {
	RootCmd.AddCommand(clusterCmd)
	clusterCmd.AddCommand(listCmd)
	clusterCmd.AddCommand(paramsCmd)
	clusterCmd.AddCommand(componentsCmd)

	clusterCmd.PersistentFlags().StringP("cluster", "C", "", "cluster to operate on")
	clusterCmd.PersistentFlags().StringVarP(&clusterParams, "clusterparams", "", "", "provide cluster params as single file - can be combined with --cluster to override cluster")

	listCmd.PersistentFlags().BoolVarP(&printRaw, "raw", "r", false, "If true, just prints result instead of placing in table.")

	paramsCmd.PersistentFlags().StringVarP(&componentName, "component", "c", "", "component to render params for")
	paramsCmd.Flags().StringVarP(&paramPath, "param", "P", "", "return value of json param from supplied path")
	paramsCmd.Flags().BoolP("notunset", "", false, "Fail if specified param is not set. Otherwise returns blank value if param is not set")
}
