package cmd

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/shuuuta/tdo/store"
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

		executeCommand("done", "2")

		got1, err := store.LoadProject(te.ProjectDir)
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
			"sample task 4",
			"sample task 5",
		}
		if !cmp.Equal(gt2, exp2) {
			t.Fatalf("expect\n  %v\ngot\n  %v", exp2, gt2)
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

	//5. 複数インデックス削除
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

		executeCommand("done", "-g", "3", "5", "2")

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
			"sample task 4",
		}
		if !cmp.Equal(gt1, exp1) {
			t.Fatalf("expect\n  %v\ngot\n  %v", exp1, gt1)
		}

		executeCommand("done", "-g", "2", "1")

		got2, err := store.LoadGlobalProject()
		if err != nil {
			t.Fatal(err)
		}

		var gt2 []string
		for _, v := range got2.Tasks {
			gt2 = append(gt2, v.Title)
		}

		if len(got2.Tasks) != 0 {
			t.Fatalf("expect empty\ngot\n  %v", gt2)
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

		if err := runDone([]string{"6", "2", "2"}); err == nil {
			t.Fatal("expect error")
		}

		if err := runDone([]string{"-6", "1", "2"}); err == nil {
			t.Fatal("expect error with argument is negative")
		}
		if err := runDone([]string{"0", "2"}); err == nil {
			t.Fatal("expect error with argument is zero")
		}
		if err := runDone([]string{"abc", "2"}); err == nil {
			t.Fatal("expect error with argument is string")
		}
		if err := runDone([]string{"1.5", "2"}); err == nil {
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
