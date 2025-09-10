package cmd

import (
	"os"
	"testing"
)

func TestAddTask(t *testing.T) {
	te := setupTestEnv(t)

	t.Run("Add project task", func(t *testing.T) {
		if err := os.Chdir(te.ProjectDir); err != nil {
			t.Fatal(err)
		}
		defer te.Cleanup()

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
		if err := os.Chdir(te.ProjectDir); err != nil {
			t.Fatal(err)
		}
		defer te.Cleanup()

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
		if err := os.Chdir(te.TmpDir); err != nil {
			t.Fatal(err)
		}
		defer te.Cleanup()

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
		if err := os.Chdir(te.TmpDir); err != nil {
			t.Fatal(err)
		}
		defer te.Cleanup()

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
		if err := os.Chdir(te.TmpDir); err != nil {
			t.Fatal(err)
		}
		defer te.Cleanup()

		if _, err := executeCommand("add", "-g"); err == nil {
			t.Fatal("expect error when no args are provided")
		}

		_, err2 := executeCommand("add", "", "")
		if err2 == nil {
			t.Fatal("expect error when empty args are provided")
		}
		exp := "argument cannot be empty"
		if err2.Error() != exp {
			t.Fatalf("expect %q, got %q", exp, err2.Error())
		}
	})
}
