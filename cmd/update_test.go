package cmd

import (
	"fmt"
	"os"
	"testing"

	"github.com/shuuuta/tdo/store"
)

func TestUpdateTask(t *testing.T) {
	te := setupTestEnv(t)

	t.Run("Update project task", func(t *testing.T) {
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
		defer te.Cleanup()

		exp1 := "update test"

		got1, err := executeCommand("update", "1", exp1)
		if err != nil {
			t.Fatal(err)
		}

		exp2 := fmt.Sprintf("update project task:\n%s", viewUpdate(1, exp1))
		if got1 != exp2 {
			t.Fatalf("\nexpect %q,\ngot    %q", exp2, got1)
		}

		got2, err := store.LoadProject(te.ProjectDir)
		if err != nil {
			t.Fatal(err)
		}

		if got2.Tasks[0].Title != exp1 {
			t.Fatalf("expect %q, got %q", exp1, got2.Tasks[0].Title)
		}
	})

	t.Run("Update global task with -g flag", func(t *testing.T) {
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
		defer te.Cleanup()

		exp1 := "update test"

		got1, err := executeCommand("update", "-g", "1", exp1)
		if err != nil {
			t.Fatal(err)
		}

		exp2 := fmt.Sprintf("update global task:\n%s", viewUpdate(1, exp1))
		if got1 != exp2 {
			t.Fatalf("\nexpect %q,\ngot    %q", exp2, got1)
		}

		got2, err := store.LoadGlobalProject()
		if err != nil {
			t.Fatal(err)
		}

		if got2.Tasks[0].Title != exp1 {
			t.Fatalf("expect %q, got %q", exp1, got2.Tasks[0].Title)
		}
	})

	t.Run("Update task outside git repo (fallback to global)", func(t *testing.T) {
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
		defer te.Cleanup()

		exp1 := "update test"

		got1, err := executeCommand("update", "1", exp1)
		if err != nil {
			t.Fatal(err)
		}

		exp2 := fmt.Sprintf("update global task:\n%s", viewUpdate(1, exp1))
		if got1 != exp2 {
			t.Fatalf("\nexpect %q,\ngot    %q", exp2, got1)
		}

		got2, err := store.LoadGlobalProject()
		if err != nil {
			t.Fatal(err)
		}

		if got2.Tasks[0].Title != exp1 {
			t.Fatalf("expect %q, got %q", exp1, got2.Tasks[0].Title)
		}
	})

	t.Run("Reject empty task title", func(t *testing.T) {
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
		defer te.Cleanup()

		_, err1 := executeCommand("update", "1", "  ")
		if err1 == nil {
			t.Fatal("expect error when empty args are provided")
		}
		exp := "title string is required"
		if err1.Error() != exp {
			t.Fatalf("expect %q, got %q", exp, err1.Error())
		}
	})

	t.Run("Handle invalid index", func(t *testing.T) {
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
		defer te.Cleanup()

		_, err1 := executeCommand("update", "3", "sample task")
		if err1 == nil {
			t.Fatal("expect error when empty args are provided")
		}
		exp1 := "unable to find task: ID 3\n"
		if err1.Error() != exp1 {
			t.Fatalf("expect %q, got %q", exp1, err1.Error())
		}

		_, err2 := executeCommand("update", "'-1'", "sample task")
		if err2 == nil {
			t.Fatal("expect error when empty args are provided")
		}
		exp2 := "task ID must be a number: '-1'"
		if err2.Error() != exp2 {
			t.Fatalf("expect %q, got %q", exp2, err2.Error())
		}
	})
}
