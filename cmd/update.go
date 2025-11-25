package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/chzyer/readline"
	"github.com/shuuuta/tdo/log"
	"github.com/shuuuta/tdo/model"
	"github.com/shuuuta/tdo/project"
	"github.com/shuuuta/tdo/store"
	"github.com/spf13/cobra"
)

var updateGlobal bool

func init() {
	rootCmd.AddCommand(updateCmd)

	updateCmd.Flags().BoolVarP(&updateGlobal, "global", "g", false, "update global task")
}

var updateCmd = &cobra.Command{
	Use:   "update <task id> [new title]",
	Short: "Update a task's title",
	Long: `Update the title of an existing task by its index number.
  If the new title is not provided, an interactive editor will open with the
  current task title. Use --global to update global tasks.`,
	Args: cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runUpdate(cmd, args)
	},
}

func runUpdate(cmd *cobra.Command, args []string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("unable to get working directory: %w", err)
	}

	pRoot, err := project.DetectRoot(cwd)
	if err != nil {
		return fmt.Errorf("unable to detect project root: %w", err)
	}

	isGlobal := pRoot == "" || updateGlobal

	var p *model.Project

	if isGlobal {
		if p, err = store.LoadGlobalProject(); err != nil {
			if os.IsNotExist(err) {
				cmd.Println("No global tasks found")
				return nil
			}
			return fmt.Errorf("unable to retrieve project tasks: %w\n", err)
		}
	} else {
		if p, err = store.LoadProject(pRoot); err != nil {
			if os.IsNotExist(err) {
				cmd.Println("No project tasks found")
				return nil
			}
			return fmt.Errorf("unable to retrieve project tasks: %w\n", err)
		}
	}

	targetNum, err := strconv.Atoi(args[0])
	if err != nil {
		log.Logf("%s", err.Error())
		return fmt.Errorf("task ID must be a number: %s", args[0])
	}

	if targetNum > len(p.Tasks) || targetNum < 1 {
		return fmt.Errorf("unable to find task: ID %d\n", targetNum)
	}

	var title string
	if len(args) == 1 {
		rl, err := readline.New("> ")
		if err != nil {
			return err
		}
		defer rl.Close()

		current := p.Tasks[targetNum-1].Title
		if _, err := rl.WriteStdin([]byte(current)); err != nil {
			return nil
		}

		if title, err = rl.Readline(); err != nil {
			return err
		}
	} else {
		title = args[1]
	}

	var gotTask *model.Task

	if isGlobal {
		if gotTask, err = store.UpdateGlobalTask(title, p.Tasks[targetNum-1].ID); err != nil {
			return err
		}
	} else {
		if gotTask, err = store.UpdateTask(pRoot, title, p.Tasks[targetNum-1].ID); err != nil {
			return err
		}
	}

	var head string
	if pRoot == "" || updateGlobal {
		head = "update global task:"
	} else {
		head = "update project task:"
	}

	cmd.Println(head)
	cmd.Print(viewUpdate(targetNum, gotTask.Title))

	return nil
}

func viewUpdate(id int, text string) string {
	return fmt.Sprintf("[ ] %d: %s\n", id, text)
}
