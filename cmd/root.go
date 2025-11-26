package cmd

import (
	"fmt"
	"os"

	"github.com/shuuuta/pask/log"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "pask",
	Short: "A simple task manager for projects and global tasks",
	Long: `Pask (short for "project task") is a command-line task manager that helps
you organize tasks both globally and per Git project. Tasks are automatically
associated with the current Git repository, or you can manage global tasks
that are available everywhere.`,
	Run: func(cmd *cobra.Command, arts []string) {
		cmd.Usage()
	},
}

func Execute() {
	log.Init()
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
