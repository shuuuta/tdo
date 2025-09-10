package cmd

import (
	"fmt"
	"os"

	"github.com/shuuuta/tdo/model"
	"github.com/shuuuta/tdo/project"
	"github.com/shuuuta/tdo/store"
	"github.com/spf13/cobra"
)

var addGlobal bool

func init() {
	rootCmd.AddCommand(addCmd)

	addCmd.Flags().BoolVarP(&addGlobal, "global", "g", false, "add global task")
}

var addCmd = &cobra.Command{
	Use:   "add [task titles...]",
	Short: "Add one or more tasks to current project or global tasks",
	Long: `Add tasks to the current Git project, or use --global to add global tasks.
  Multiple task titles can be provided as separate arguments. If not in a Git
  repository, tasks are automatically added as global tasks.`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runAdd(cmd, args)
	},
}

func runAdd(cmd *cobra.Command, args []string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("unable to get working directory: %w", err)
	}

	pRoot, err := project.DetectRoot(cwd)
	if err != nil {
		return fmt.Errorf("unable to detect project root: %w", err)
	}

	var title string
	var tasks []*model.Task
	var errors []string
	if pRoot == "" || addGlobal {
		title = "add global task:"
		for _, v := range args {
			if v == "" {
				continue
			}
			t, err := store.AddGlobalTask(v)
			if err != nil {
				errors = append(errors, fmt.Sprintf("Failed to add %s: %v", v, err))
				continue
			}
			tasks = append(tasks, t)
		}
	} else {
		title = "add project task:"
		for _, v := range args {
			if v == "" {
				continue
			}
			t, err := store.AddTask(pRoot, v)
			if err != nil {
				errors = append(errors, fmt.Sprintf("Failed to add %s: %v", v, err))
				continue
			}
			tasks = append(tasks, t)
		}
	}

	if len(tasks) > 0 {
		cmd.Println(title)
		for _, v := range tasks {
			cmd.Printf("  - %s\n", v.Title)
		}
	}

	if len(errors) > 0 {
		for _, m := range errors {
			cmd.PrintErrln(m)
		}

		s := ""
		if len(errors) > 1 {
			s = "s"
		}
		return fmt.Errorf("%d task%s failed to add", len(errors), s)
	}

	if len(tasks) == 0 && len(errors) == 0 {
		return fmt.Errorf("argument cannot be empty")
	}

	return nil
}
