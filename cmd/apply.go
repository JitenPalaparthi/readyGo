package cmd

import (
	"log"
	"readyGo/generate"
	"readyGo/generate/configure"
	templates "readyGo/templates"

	"github.com/spf13/cobra"
)

var configDefault = string(`{
		"version":"0.1",
		"project": "example",
		"type": "http",
		"port": "50054",
		"db": "mongo",
		"models": [
		  {
			"name": "person",
			"fields": [
			   {
				"name": "name",
				"type": "string",
				"isKey": true
			  },
			  {
				"name": "email",
				"type": "string",
				"validateExp": "[a-zA-Z0-9]",
				"isKey": true
			  },
			  {
				"name": "mobile",
				"type": "string"
			  },
			  {
				"name": "status",
				"type": "string"
			  },
			  {
				"name": "last_Modified",
				"type": "string"
			  }
			]
		  }
		]
	  }`)

var applyFile string

func init() {

	applyCmd.Flags().StringVarP(&applyFile, "filename", "f", "default", "config file to generate the project")
	rootCmd.AddCommand(applyCmd)
}

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "applies a configuration file",
	Long:  `applyt applies a configuration file for the project generation`,
	Run: func(cmd *cobra.Command, args []string) {

		/*tm, err := template.New("templates")
		if err != nil {
			log.Fatal("error occured loading templates.---->", err)
		}*/

		// This is another way of reading templates. The above one is to read from template files but due to static file bindings with go single binary ,
		// this has been made directly from map. All templates are copied from tempalte files and loaded to map.
		var tg *generate.Generate

		tm := templates.New()
		if tm == nil {
			log.Fatal("no templates are available")
		}

		templateConfig := "configs/template_config.json"

		tc, err := configure.New(&templateConfig)
		if applyFile == "default" {
			tg, err = generate.NewFromStr(configDefault, tm, tc)
			if err != nil {
				tg.RmDir()
				log.Fatal("seems , things went wrong.. -->", err)
			}
		} else {
			tg, err = generate.New(&applyFile, tm, tc)

			if err != nil {
				tg.RmDir()
				log.Fatal("seems , things went wrong.. -->", err)
			}
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
