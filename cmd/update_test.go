package cmd

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/shuuuta/tdo/store"
	"github.com/spf13/cobra"
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

func TestUpdateTaskInteractive(t *testing.T) {
	te := setupTestEnv(t)
	t.Run("Update task with interactive mode", func(t *testing.T) {
		if err := os.Chdir(te.ProjectDir); err != nil {
			t.Fatal(err)
		}
		defer te.Cleanup()

		originalTitle := "original task"
		if _, err := store.AddTask(te.ProjectDir, originalTitle); err != nil {
			t.Fatal(err)
		}

		newTitle := "update via interactive mode"
		mockReader := func(prompt string) (LineReader, error) {
			return &mockLineReader{
				readlineFunc: func() (string, error) {
					return newTitle, nil
				},
				writeStdinFunc: func(b []byte) (int, error) {
					return len(b), nil
				},
				closeFunc: func() error {
					return nil
				},
			}, nil
		}

		cmd := &cobra.Command{}
		cmd.SetOut(new(bytes.Buffer))

		if err := runUpdate(cmd, []string{"1"}, mockReader); err != nil {
			t.Fatal(err)
		}

		p, err := store.LoadProject(te.ProjectDir)
		if err != nil {
			t.Fatal(err)
		}

		if p.Tasks[0].Title != newTitle {
			t.Fatalf("expect %q, got %q", newTitle, p.Tasks[0].Title)
		}
	})

	t.Run("Interactive mode preserves original text", func(t *testing.T) {
		if err := os.Chdir(te.ProjectDir); err != nil {
			t.Fatal(err)
		}
		defer te.Cleanup()

		originalTitle := "original task"
		if _, err := store.AddTask(te.ProjectDir, originalTitle); err != nil {
			t.Fatal(err)
		}

		var capturedWriteStdin []byte
		mockReaderWithCapture := func(prompt string) (LineReader, error) {
			return &mockLineReader{
				readlineFunc: func() (string, error) {
					return "new title", nil
				},
				writeStdinFunc: func(b []byte) (int, error) {
					capturedWriteStdin = b
					return len(b), nil
				},
				closeFunc: func() error {
					return nil
				},
			}, nil
		}

		cmd := &cobra.Command{}
		cmd.SetOut(new(bytes.Buffer))

		err := runUpdate(cmd, []string{"1"}, mockReaderWithCapture)
		if err != nil {
			t.Fatal(err)
		}

		if string(capturedWriteStdin) != originalTitle {
			t.Fatalf("expect WriteStdin to be called with %q, got %q", originalTitle, string(capturedWriteStdin))
		}
	})

	t.Run("Interactive mode handles readline error", func(t *testing.T) {
		if err := os.Chdir(te.ProjectDir); err != nil {
			t.Fatal(err)
		}
		defer te.Cleanup()

		if _, err := store.AddTask(te.ProjectDir, "sample task"); err != nil {
			t.Fatal(err)
		}

		expectedErr := fmt.Errorf("readline error")
		mockReaderWithError := func(prompt string) (LineReader, error) {
			return &mockLineReader{
				readlineFunc: func() (string, error) {
					return "", expectedErr
				},
			}, nil
		}

		cmd := &cobra.Command{}
		cmd.SetOut(new(bytes.Buffer))

		err := runUpdate(cmd, []string{"1"}, mockReaderWithError)
		if err != expectedErr {
			t.Fatalf("expect error %v, got %v", expectedErr, err)
		}
	})
}

type mockLineReader struct {
	readlineFunc   func() (string, error)
	writeStdinFunc func([]byte) (int, error)
	closeFunc      func() error
}

func (m *mockLineReader) Readline() (string, error) {
	if m.readlineFunc != nil {
		return m.readlineFunc()
	}
	return "", nil
}
func (m *mockLineReader) WriteStdin(b []byte) (int, error) {
	if m.writeStdinFunc != nil {
		return m.writeStdinFunc(b)
	}
	return len(b), nil
}
func (m *mockLineReader) Close() error {
	if m.closeFunc != nil {
		return m.closeFunc()
	}
	return nil
}
