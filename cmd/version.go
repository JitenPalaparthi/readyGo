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
	Short: "readyGo v0.0.10",
	Long: `readyGo v0.0.10.Currently in Active Development\n grpc Supported versions:
	 	\nprotoc-gen-go v1.25.0
	 	\nprotoc        v3.14.0`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("readyGo  v0.0.10\n\nCompiled and tested for grpc versions:\t\nprotoc-gen-go v1.25.0\t\nprotoc        v3.14.0")
	},
}
