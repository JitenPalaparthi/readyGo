package cmd

import (
	"log"
	"readyGo/boxops"
	"readyGo/generate"
	"readyGo/lang/implement"
	"readyGo/scalar"

	"github.com/spf13/cobra"
)

var applyFileValidate string

func init() {
	validateCmd.Flags().StringVarP(&applyFileValidate, "filename", "f", "", "user has to privide the file.There is no default file.")
	rootCmd.AddCommand(validateCmd)
}

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "validates a configuration file",
	Long:  `validate validates a configuration file for the project generation.User must supply a configuration file`,
	Run: func(cmd *cobra.Command, args []string) {

		ops := boxops.New("../box")

		scalar, err := scalar.New(ops, "configs/scalars.json")

		if err != nil {
			log.Fatal(Fata(err))
		}

		if genFile == "default" {
			log.Fatal(Fata("apply must supply corrosponding configuration file"))
		}

		imlementer := implement.New()

		_, err = generate.New(&applyFileValidate, scalar, imlementer)
		if err != nil {
			log.Println(Warn("There are errors.Validation failed"))
			log.Println(Fata(err))
		} else {
			log.Println(Info("Successfully validate.Use apply command to generate required files"))
		}
	},
}
