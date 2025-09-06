package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "tdo",
	Short: "short",
	Long:  `Long`,
	Run: func(cmd *cobra.Command, arts []string) {
		fmt.Println("bar")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
