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
	Short: "readyGo v0.0.5",
	Long:  `readyGo v0.0.5.Currently in Active Development`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("readyGo  v0.0.5")
	},
}
