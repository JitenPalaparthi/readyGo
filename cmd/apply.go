package cmd

import (
	"log"
	"readyGo/generate"
	"readyGo/generate/configure"
	"readyGo/generate/template"

	"github.com/spf13/cobra"
)

var applyFile string

func init() {

	applyCmd.Flags().StringVarP(&applyFile, "filename", "f", "configs/config.yaml", "config file to generate the project")
	rootCmd.AddCommand(applyCmd)
}

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "applies a configuration file",
	Long:  `applyt applies a configuration file for the project generation`,
	Run: func(cmd *cobra.Command, args []string) {

		tm, err := template.New("templates")

		if err != nil {
			log.Fatal("error occured loading templates.---->", err)
		}

		templateConfig := "configs/template_config.json"

		tc, err := configure.New(&templateConfig)

		tg, err := generate.New(&applyFile, tm, tc)

		if err != nil {
			tg.RmDir()
			log.Fatal("seems , things went wrong.. -->", err)
		}

		err = tg.GenerateAll("http_mongo")
		if err != nil {
			log.Println("seems , things went wrong.Rolling back all generated files -->", err)
			err = tg.RmDir()
			if err != nil {
				log.Println("Unable to remove files. Please remove created directory manually", err)
			}
		}

	},
}
