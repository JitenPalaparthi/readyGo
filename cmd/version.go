package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "readyGo v0.0.3",
	Long:  `readyGo v0.0.3.Currently in Alpha stage`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("readyGo  v0.0.3")
	},
}
