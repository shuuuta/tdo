package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/shuuuta/tdo/model"
	"github.com/shuuuta/tdo/store"
)

func setupTest(t *testing.T, configPath string) func() {
	oListGlobal := listGlobal

	return func() {
		listGlobal = oListGlobal

		files, err := os.ReadDir(configPath)
		if err != nil {
			t.Fatal(err)
		}
		for _, v := range files {
			if err := os.Remove(filepath.Join(configPath, v.Name())); err != nil {
				t.Fatal(err)
			}
		}
	}
}

func TestListCmd(t *testing.T) {
	tmpDir := t.TempDir()

	configPath := filepath.Join(tmpDir, "conf")
	os.MkdirAll(configPath, 0755)
	store.SetConfigDir(configPath)

	projectPath := filepath.Join(tmpDir, "project")
	os.MkdirAll(filepath.Join(projectPath, ".git"), 0755)

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(cwd)

	t.Run("Global task does not exist", func(t *testing.T) {
		if err := os.Chdir(projectPath); err != nil {
			t.Fatal(err)
		}
		cleanup := setupTest(t, configPath)
		defer cleanup()

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
		if err := os.Chdir(projectPath); err != nil {
			t.Fatal(err)
		}
		cleanup := setupTest(t, configPath)
		defer cleanup()

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
		if err := os.Chdir(projectPath); err != nil {
			t.Fatal(err)
		}
		cleanup := setupTest(t, configPath)
		defer cleanup()

		realPPath, err := filepath.EvalSymlinks(projectPath)
		if err != nil {
			t.Fatal(err)
		}

		_, err = store.AddTask(realPPath, "sample task 1")
		if err != nil {
			t.Fatal(err)
		}
		_, err = store.AddTask(realPPath, "sample task 2")
		if err != nil {
			t.Fatal(err)
		}

		got, err := executeCommand("list")
		if err != nil {
			t.Fatal(err)
		}
		exp := viewList(&model.Project{
			ProjectPath: realPPath,
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
		if err := os.Chdir(projectPath); err != nil {
			t.Fatal(err)
		}
		cleanup := setupTest(t, configPath)
		defer cleanup()

		_, err = store.AddGlobalTask("sample task 1")
		if err != nil {
			t.Fatal(err)
		}
		_, err = store.AddGlobalTask("sample task 2")
		if err != nil {
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
		if err := os.Chdir(tmpDir); err != nil {
			t.Fatal(err)
		}
		cleanup := setupTest(t, configPath)
		defer cleanup()

		_, err = store.AddGlobalTask("sample task 1")
		if err != nil {
			t.Fatal(err)
		}
		_, err = store.AddGlobalTask("sample task 2")
		if err != nil {
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
