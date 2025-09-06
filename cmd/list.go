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
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		pRoot, err := project.DetectRoot(cwd)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		var p *model.Project

		if pRoot == "" || listGlobal {
			p, err = store.LoadGlobalProject()
			if err != nil {
				if os.IsNotExist(err) {
					fmt.Println("Global task is Not exist")
					return
				}
				log.Logf("LoadGlobalProject error: %v", err)
				fmt.Fprintln(os.Stderr, err)
				fmt.Fprintln(os.Stderr, "[error] Unable to retrieve global tasks")
				os.Exit(1)
			}
		} else {
			p, err = store.LoadProject(pRoot)
			if err != nil {
				if os.IsNotExist(err) {
					fmt.Println("Project task is not exist")
					return
				}
				log.Logf("LoadProject error: %v", err)
				fmt.Fprintln(os.Stderr, "[error] Unable to retrieve project tasks")
				os.Exit(1)
			}
		}

		if p.IsGlobal {
			fmt.Println("Global tasks")
		} else {
			fmt.Println(p.ProjectPath)
		}
		fmt.Println("========================")
		for i, v := range p.Tasks {
			fmt.Printf("[ ] %d: %s\n", i+1, v.Title)
		}
	},
}
