package cmd

import (
	"fmt"
	"os"

	"github.com/shuuuta/tdo/log"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "tdo",
	Short: "A simple todo manager for projects and global tasks",
	Long: `Tdo is a command-line todo manager that helps you organize tasks
both globally and per Git project. Tasks are automatically associated
with the current Git repository, or you can manage global tasks that
are available everywhere.`,
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
