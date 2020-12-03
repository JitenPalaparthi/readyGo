package cmd

import (
	"fmt"
	"log"
	"readyGo/boxops"
	"readyGo/generate"
	"readyGo/mapping"

	"github.com/spf13/cobra"
)

var applyFile, projectType string

func init() {

	applyCmd.Flags().StringVarP(&applyFile, "filename", "f", "default", "config file to generate the project")
	applyCmd.Flags().StringVarP(&projectType, "type", "t", "http_mongo", "type of the project http_mongo | http_sql_pg | grpc_mogo | grpc_sql | cloudevent")

	rootCmd.AddCommand(applyCmd)
}

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "applies a configuration file",
	Long:  `apply applies a configuration file for the project generation.User must supply a configuration file`,
	Run: func(cmd *cobra.Command, args []string) {

		ops := boxops.New("../box")
		mapping, err := mapping.New(ops, "configs/mappings.json", projectType)
		fmt.Println(*mapping)
		if err != nil {
			log.Fatal(err)
		}

		if applyFile == "default" {
			log.Fatal("apply must supply corrosponding configuration file")
		}

		tg, err := generate.New(&applyFile, mapping)
		if err != nil {
			log.Fatal(err)
		}
		err = tg.CreateAll()
		if err != nil {
			log.Fatal(err)
		}

	},
}
