package cmd

import (
	"fmt"
	"os"

	"github.com/shuuuta/tdo/log"
	"github.com/shuuuta/tdo/model"
	"github.com/shuuuta/tdo/project"
	"github.com/shuuuta/tdo/store"
	"github.com/spf13/cobra"
)

var listGlobal bool

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().BoolVarP(&listGlobal, "global", "g", false, "show global tasks")
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List tasks for current project or global tasks",
	Long: `List all tasks for the current Git project. If no Git repository
is detected or the --global flag is used, global tasks will be shown instead.

Tasks are numbered for easy reference when using other commands.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := runList(cmd); err != nil {
			cmd.PrintErrln(err)
			os.Exit(1)
		}
	},
}

func runList(cmd *cobra.Command) error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("unable to get working directory: %w", err)
	}

	pRoot, err := project.DetectRoot(cwd)
	if err != nil {
		return fmt.Errorf("unable to detect project root: %w", err)
	}

	var p *model.Project

	if pRoot == "" || listGlobal {
		p, err = store.LoadGlobalProject()
		if err != nil {
			if os.IsNotExist(err) {
				cmd.Println("No global tasks found")
				return nil
			}
			log.Logf("LoadGlobalProject error: %v", err)
			return fmt.Errorf("unable to retrieve global tasks: %w\n", err)
		}
	} else {
		p, err = store.LoadProject(pRoot)
		if err != nil {
			if os.IsNotExist(err) {
				cmd.Println("No project tasks found")
				return nil
			}
			log.Logf("LoadProject error: %v", err)
			return fmt.Errorf("unable to retrieve project tasks: %w\n", err)
		}
	}

	cmd.Printf("%s", viewList(p))

	return nil
}

func viewList(project *model.Project) string {
	var out string
	if project.IsGlobal {
		out = fmt.Sprintln("Global tasks")
	} else {
		out = fmt.Sprintln(project.ProjectPath)
	}
	out = out + fmt.Sprintln("========================")
	for i, v := range project.Tasks {
		out = out + fmt.Sprintf("[ ] %d: %s\n", i+1, v.Title)
	}

	return out
}
