package cmd

import (
	"os"
	"testing"

	"github.com/shuuuta/tdo/model"
	"github.com/shuuuta/tdo/store"
)

func TestListCmd(t *testing.T) {
	te := setupTestEnv(t)

	t.Run("Global task does not exist", func(t *testing.T) {
		if err := os.Chdir(te.ProjectDir); err != nil {
			t.Fatal(err)
		}
		defer te.Cleanup()

		got, err := executeCommand("list", "-g")
		if err != nil {
			t.Fatal(err)
		}

		exp := "No global tasks found\n"
		if got != exp {
			t.Fatalf("expect %q, got %q", exp, got)
		}
	})

	t.Run("Project task does not exist", func(t *testing.T) {
		if err := os.Chdir(te.ProjectDir); err != nil {
			t.Fatal(err)
		}
		defer te.Cleanup()

		got, err := executeCommand("list")
		if err != nil {
			t.Fatal(err)
		}

		exp := "No project tasks found\n"
		if got != exp {
			t.Fatalf("expect %q, got %q", exp, got)
		}
	})

	t.Run("Show tasks in project path", func(t *testing.T) {
		if err := os.Chdir(te.ProjectDir); err != nil {
			t.Fatal(err)
		}
		defer te.Cleanup()

		if _, err := store.AddTask(te.ProjectDir, "sample task 1"); err != nil {
			t.Fatal(err)
		}
		if _, err := store.AddTask(te.ProjectDir, "sample task 2"); err != nil {
			t.Fatal(err)
		}

		got, err := executeCommand("list")
		if err != nil {
			t.Fatal(err)
		}
		exp := viewList(&model.Project{
			ProjectPath: te.ProjectDir,
			Tasks: []model.Task{
				{Title: "sample task 1"},
				{Title: "sample task 2"},
			},
		})
		if got != exp {
			t.Fatalf("expect %q, got %q", exp, got)
		}
	})

	t.Run("Show global tasks in project path", func(t *testing.T) {
		if err := os.Chdir(te.ProjectDir); err != nil {
			t.Fatal(err)
		}
		defer te.Cleanup()

		if _, err := store.AddGlobalTask("sample task 1"); err != nil {
			t.Fatal(err)
		}
		if _, err := store.AddGlobalTask("sample task 2"); err != nil {
			t.Fatal(err)
		}

		got, err := executeCommand("list", "-g")
		if err != nil {
			t.Fatal(err)
		}
		exp := viewList(&model.Project{
			IsGlobal: true,
			Tasks: []model.Task{
				{Title: "sample task 1"},
				{Title: "sample task 2"},
			},
		})
		if got != exp {
			t.Fatalf("expect %q, got %q", exp, got)
		}
	})

	t.Run("Show global tasks out of project path", func(t *testing.T) {
		if err := os.Chdir(te.TmpDir); err != nil {
			t.Fatal(err)
		}
		defer te.Cleanup()

		if _, err := store.AddGlobalTask("sample task 1"); err != nil {
			t.Fatal(err)
		}
		if _, err := store.AddGlobalTask("sample task 2"); err != nil {
			t.Fatal(err)
		}

		got, err := executeCommand("list")
		if err != nil {
			t.Fatal(err)
		}
		exp := viewList(&model.Project{
			IsGlobal: true,
			Tasks: []model.Task{
				{Title: "sample task 1"},
				{Title: "sample task 2"},
			},
		})
		if got != exp {
			t.Fatalf("expect %q, got %q", exp, got)
		}
	})
}
