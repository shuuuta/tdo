package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/shuuuta/tdo/store"
	"github.com/spf13/cobra"
)

func TestAddTask(t *testing.T) {
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

	t.Run("Add project task", func(t *testing.T) {
		if err := os.Chdir(projectPath); err != nil {
			t.Fatal(err)
		}
		cleanup := setupTest(t, configPath)
		defer cleanup()

		got, err := executeCommand("add", "sample test")
		if err != nil {
			t.Fatal(err)
		}

		exp := "add project task:\n  - sample test\n"
		if got != exp {
			t.Fatalf("\nexpect %q,\ngot    %q", exp, got)
		}
	})

	t.Run("Add global task with -g flag", func(t *testing.T) {
		if err := os.Chdir(projectPath); err != nil {
			t.Fatal(err)
		}
		cleanup := setupTest(t, configPath)
		defer cleanup()

		got, err := executeCommand("add", "-g", "sample test")
		if err != nil {
			t.Fatal(err)
		}

		exp := "add global task:\n  - sample test\n"
		if got != exp {
			t.Fatalf("\nexpect %q,\ngot    %q", exp, got)
		}
	})

	t.Run("Add task outside git repo", func(t *testing.T) {
		if err := os.Chdir(tmpDir); err != nil {
			t.Fatal(err)
		}
		cleanup := setupTest(t, configPath)
		defer cleanup()

		got, err := executeCommand("add", "-g", "sample test")
		if err != nil {
			t.Fatal(err)
		}

		exp := "add global task:\n  - sample test\n"
		if got != exp {
			t.Fatalf("\nexpect %q,\ngot    %q", exp, got)
		}
	})

	t.Run("Add task with multiple words", func(t *testing.T) {
		if err := os.Chdir(tmpDir); err != nil {
			t.Fatal(err)
		}
		cleanup := setupTest(t, configPath)
		defer cleanup()

		got, err := executeCommand("add", "-g", "sample test", "second test")
		if err != nil {
			t.Fatal(err)
		}

		exp := "add global task:\n  - sample test\n  - second test\n"
		if got != exp {
			t.Fatalf("\nexpect %q,\ngot    %q", exp, got)
		}
	})

	t.Run("Reject empty task title", func(t *testing.T) {
		if err := os.Chdir(tmpDir); err != nil {
			t.Fatal(err)
		}
		cleanup := setupTest(t, configPath)
		defer cleanup()

		if _, err := executeCommand("add", "-g"); err == nil {
			t.Fatal("expect error when no args are provided")
		}

		buf := new(bytes.Buffer)
		testCmd := &cobra.Command{}
		testCmd.SetOut(buf)
		testCmd.SetErr(buf)

		err := runAdd(testCmd, []string{"", ""})
		if err == nil {
			t.Fatal("expect error when empty args are provided")
		}
		exp := "argument cannot be empty"
		if err.Error() != exp {
			t.Fatalf("expect %q, got %q", exp, err.Error())
		}
	})
}
