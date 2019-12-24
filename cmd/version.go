package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// Cobra version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "version",
	Long:  `version`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Sasspiler v0.1.0")
	},
}
