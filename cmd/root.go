package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "readyGo",
	Short: "readyGo generates http,grpc,cloudEvent based projects.It's not a template engine.It generates working project.",
	Long:  `readyGo is a command line interface( probably the name of readyGo CLI would be rgo) application, it is designed to scaffold creation of different types of go based projects.readyGo is designed for developers in mind. Ideally readyGo should provide ready to use application code. The code is generated based on configurations provided by the end user`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here

	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
