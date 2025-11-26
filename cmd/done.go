package cmd

import (
	"fmt"
	"os"
	"slices"
	"strconv"

	"github.com/shuuuta/pask/model"
	"github.com/shuuuta/pask/project"
	"github.com/shuuuta/pask/store"
	"github.com/spf13/cobra"
)

var doneGlobal bool

func init() {
	rootCmd.AddCommand(doneCmd)

	doneCmd.Flags().BoolVarP(&doneGlobal, "global", "g", false, "mark global tasks as done")
}

var doneCmd = &cobra.Command{
	Use:   "done <task id...>",
	Short: "Mark tasks as done (delete them)",
	Long: `Mark one or more tasks as done by their index numbers, removing them from the list.
  Multiple task indices can be provided as separate arguments.`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runDone(cmd, args)
	},
}

func runDone(cmd *cobra.Command, args []string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("unable to get working directory: %w", err)
	}

	pRoot, err := project.DetectRoot(cwd)
	if err != nil {
		return fmt.Errorf("unable to detect project root: %w", err)
	}

	var p *model.Project
	if pRoot == "" || doneGlobal {
		if p, err = store.LoadGlobalProject(); err != nil {
			if os.IsNotExist(err) {
				cmd.Println("No global tasks found")
				return nil
			}
			return fmt.Errorf("unable to retrieve global tasks: %w\n", err)
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

	var ids []int
	for _, v := range args {
		id, err := strconv.Atoi(v)
		if err != nil {
			return err
		}
		if id > len(p.Tasks) || id < 1 {
			return fmt.Errorf("unable to mark as done task: ID %d\n", id)
		}
		ids = append(ids, id)
	}

	slices.Sort(ids)
	ids = slices.Compact(ids)

	var out string
	offset := 0
	for _, id := range ids {
		out = out + viewDone(id, p.Tasks[id-1-offset].Title)
		p.Tasks = slices.Delete(p.Tasks, id-1-offset, id-offset)
		offset++
	}

	if err := store.SaveProject(p); err != nil {
		return err
	}

	cmd.Print(out)

	return nil
}

func viewDone(id int, text string) string {
	return fmt.Sprintf("[✔ ] %d: %s\n", id, text)
}
