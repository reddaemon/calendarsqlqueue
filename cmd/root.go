package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "calendar",
	Short: "Calendar service",
	Long:  `Calendar service`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("Use calendar [command]\nRun 'calendar --help' for usage.\n")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("root execute error: %v", err)
	}
}
