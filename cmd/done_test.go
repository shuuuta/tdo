package cmd

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/shuuuta/pask/store"
)

func TestDoneTask(t *testing.T) {
	te := setupTestEnv(t)

	t.Run("Done project task", func(t *testing.T) {
		if err := os.Chdir(te.ProjectDir); err != nil {
			t.Fatal(err)
		}
		defer te.Cleanup()

		titles := []string{
			"sample task 1",
			"sample task 2",
			"sample task 3",
			"sample task 4",
			"sample task 5",
		}
		for _, v := range titles {
			if _, err := store.AddTask(te.ProjectDir, v); err != nil {
				t.Fatal(err)
			}
		}

		got1, err := executeCommand("done", "2")
		if err != nil {
			t.Fatal(err)
		}
		exp1 := viewDone(2, "sample task 2")
		if got1 != exp1 {
			t.Fatalf("expect %q, got %q", exp1, got1)
		}

		got2, err := store.LoadProject(te.ProjectDir)
		if err != nil {
			t.Fatal(err)
		}

		var gt2 []string
		for _, v := range got2.Tasks {
			gt2 = append(gt2, v.Title)
		}

		exp2 := []string{
			"sample task 1",
			"sample task 3",
			"sample task 4",
			"sample task 5",
		}
		if !cmp.Equal(gt2, exp2) {
			t.Fatalf("expect\n  %v\ngot\n  %v", exp2, gt2)
		}

		executeCommand("done", "2")

		got3, err := store.LoadProject(te.ProjectDir)
		if err != nil {
			t.Fatal(err)
		}

		var gt3 []string
		for _, v := range got3.Tasks {
			gt3 = append(gt3, v.Title)
		}

		exp3 := []string{
			"sample task 1",
			"sample task 4",
			"sample task 5",
		}
		if !cmp.Equal(gt3, exp3) {
			t.Fatalf("expect\n  %v\ngot\n  %v", exp3, gt3)
		}
	})

	t.Run("Done global task with -g flag", func(t *testing.T) {
		if err := os.Chdir(te.ProjectDir); err != nil {
			t.Fatal(err)
		}
		defer te.Cleanup()

		titles := []string{
			"sample task 1",
			"sample task 2",
			"sample task 3",
			"sample task 4",
			"sample task 5",
		}
		for _, v := range titles {
			if _, err := store.AddGlobalTask(v); err != nil {
				t.Fatal(err)
			}
		}

		executeCommand("done", "-g", "2")

		got1, err := store.LoadGlobalProject()
		if err != nil {
			t.Fatal(err)
		}

		var gt1 []string
		for _, v := range got1.Tasks {
			gt1 = append(gt1, v.Title)
		}

		exp1 := []string{
			"sample task 1",
			"sample task 3",
			"sample task 4",
			"sample task 5",
		}
		if !cmp.Equal(gt1, exp1) {
			t.Fatalf("expect\n  %v\ngot\n  %v", exp1, gt1)
		}

		executeCommand("done", "2")

		got2, err := store.LoadGlobalProject()
		if err != nil {
			t.Fatal(err)
		}

		var gt2 []string
		for _, v := range got2.Tasks {
			gt2 = append(gt2, v.Title)
		}

		exp2 := []string{
			"sample task 1",
			"sample task 4",
			"sample task 5",
		}
		if !cmp.Equal(gt2, exp2) {
			t.Fatalf("expect\n  %v\ngot\n  %v", exp2, gt2)
		}
	})

	t.Run("Done task outside project dir", func(t *testing.T) {
		if err := os.Chdir(te.TmpDir); err != nil {
			t.Fatal(err)
		}
		defer te.Cleanup()

		titles := []string{
			"sample task 1",
			"sample task 2",
			"sample task 3",
			"sample task 4",
			"sample task 5",
		}
		for _, v := range titles {
			if _, err := store.AddGlobalTask(v); err != nil {
				t.Fatal(err)
			}
		}

		executeCommand("done", "2")

		got1, err := store.LoadGlobalProject()
		if err != nil {
			t.Fatal(err)
		}

		var gt1 []string
		for _, v := range got1.Tasks {
			gt1 = append(gt1, v.Title)
		}

		exp1 := []string{
			"sample task 1",
			"sample task 3",
			"sample task 4",
			"sample task 5",
		}
		if !cmp.Equal(gt1, exp1) {
			t.Fatalf("expect\n  %v\ngot\n  %v", exp1, gt1)
		}

		executeCommand("done", "2")

		got2, err := store.LoadGlobalProject()
		if err != nil {
			t.Fatal(err)
		}

		var gt2 []string
		for _, v := range got2.Tasks {
			gt2 = append(gt2, v.Title)
		}

		exp2 := []string{
			"sample task 1",
			"sample task 4",
			"sample task 5",
		}
		if !cmp.Equal(gt2, exp2) {
			t.Fatalf("expect\n  %v\ngot\n  %v", exp2, gt2)
		}
	})

	t.Run("Done multiple tasks by index", func(t *testing.T) {
		if err := os.Chdir(te.ProjectDir); err != nil {
			t.Fatal(err)
		}
		defer te.Cleanup()

		titles := []string{
			"sample task 1",
			"sample task 2",
			"sample task 3",
			"sample task 4",
			"sample task 5",
		}
		for _, v := range titles {
			if _, err := store.AddGlobalTask(v); err != nil {
				t.Fatal(err)
			}
		}

		got1, err := executeCommand("done", "-g", "3", "5", "2")
		if err != nil {
			t.Fatal(err)
		}
		var exp1 string
		for _, id := range []int{2, 3, 5} {
			exp1 = exp1 + viewDone(id, titles[id-1])
		}
		if exp1 != got1 {
			t.Fatalf("expect:\n%s, got:\n%s", exp1, got1)
		}

		got2, err := store.LoadGlobalProject()
		if err != nil {
			t.Fatal(err)
		}

		var gt2 []string
		for _, v := range got2.Tasks {
			gt2 = append(gt2, v.Title)
		}

		exp2 := []string{
			"sample task 1",
			"sample task 4",
		}
		if !cmp.Equal(gt2, exp2) {
			t.Fatalf("expect\n  %v\ngot\n  %v", exp2, gt2)
		}

		executeCommand("done", "-g", "2", "1")

		got3, err := store.LoadGlobalProject()
		if err != nil {
			t.Fatal(err)
		}

		var gt3 []string
		for _, v := range got3.Tasks {
			gt3 = append(gt3, v.Title)
		}

		if len(got3.Tasks) != 0 {
			t.Fatalf("expect empty\ngot\n  %v", gt3)
		}
	})

	t.Run("Handle duplicate indices", func(t *testing.T) {
		if err := os.Chdir(te.ProjectDir); err != nil {
			t.Fatal(err)
		}
		defer te.Cleanup()

		titles := []string{
			"sample task 1",
			"sample task 2",
			"sample task 3",
			"sample task 4",
			"sample task 5",
		}
		for _, v := range titles {
			if _, err := store.AddGlobalTask(v); err != nil {
				t.Fatal(err)
			}
		}

		executeCommand("done", "-g", "2", "2", "2")

		got1, err := store.LoadGlobalProject()
		if err != nil {
			t.Fatal(err)
		}

		var gt1 []string
		for _, v := range got1.Tasks {
			gt1 = append(gt1, v.Title)
		}

		exp1 := []string{
			"sample task 1",
			"sample task 3",
			"sample task 4",
			"sample task 5",
		}
		if !cmp.Equal(gt1, exp1) {
			t.Fatalf("expect\n  %v\ngot\n  %v", exp1, gt1)
		}
	})

	t.Run("Task file does not exist", func(t *testing.T) {
		if err := os.Chdir(te.ProjectDir); err != nil {
			t.Fatal(err)
		}
		defer te.Cleanup()

		got1, err := executeCommand("done", "1")
		if err != nil {
			t.Fatal(err)
		}

		exp1 := "No project tasks found\n"
		if got1 != exp1 {
			t.Fatalf("expect %q, got %q", exp1, got1)
		}

		got2, err := executeCommand("done", "-g", "1")
		if err != nil {
			t.Fatal(err)
		}

		exp2 := "No global tasks found\n"
		if got2 != exp2 {
			t.Fatalf("expect %q, got %q", exp2, got2)
		}
	})

	t.Run("Handle invalid index", func(t *testing.T) {
		if err := os.Chdir(te.TmpDir); err != nil {
			t.Fatal(err)
		}
		defer te.Cleanup()

		titles := []string{
			"sample task 1",
			"sample task 2",
			"sample task 3",
			"sample task 4",
			"sample task 5",
		}
		for _, v := range titles {
			if _, err := store.AddGlobalTask(v); err != nil {
				t.Fatal(err)
			}
		}

		if _, err := executeCommand("done", "6", "2", "2"); err == nil {
			t.Fatal("expect error")
		}

		if _, err := executeCommand("done", "-6", "1", "2"); err == nil {
			t.Fatal("expect error with argument is negative")
		}
		if _, err := executeCommand("done", "0", "2"); err == nil {
			t.Fatal("expect error with argument is zero")
		}
		if _, err := executeCommand("done", "abc", "2"); err == nil {
			t.Fatal("expect error with argument is string")
		}
		if _, err := executeCommand("done", "1.5", "2"); err == nil {
			t.Fatal("expect error with argument is float")
		}

		got1, err := store.LoadGlobalProject()
		if err != nil {
			t.Fatal(err)
		}

		var gt1 []string
		for _, v := range got1.Tasks {
			gt1 = append(gt1, v.Title)
		}

		if !cmp.Equal(gt1, titles) {
			t.Fatalf("expect\n  %v\ngot\n  %v", titles, gt1)
		}
	})
}
