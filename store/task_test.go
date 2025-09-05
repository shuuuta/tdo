package store

import (
	"sort"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestAddTask(t *testing.T) {
	configDir = t.TempDir()

	projectPath := "/foo/bar"

	title1 := "creata project and add task"
	if _, err := AddTask(projectPath, title1); err != nil {
		t.Fatal(err)
	}

	lp1, err := LoadProject(projectPath)
	if err != nil {
		t.Fatal(err)
	}
	if len(lp1.Tasks) != 1 || lp1.Tasks[0].Title != title1 {
		t.Fatalf("expect 1 task, got %+v", lp1.Tasks)
	}

	title2 := "add extra task"
	AddTask(projectPath, title2)

	lp2, err := LoadProject(projectPath)
	if err != nil {
		t.Fatal(err)
	}

	var got []string
	for _, v := range lp2.Tasks {
		got = append(got, v.Title)
	}
	exp := []string{title2, title1}

	sort.Slice(got, func(i, j int) bool {
		return got[i] < got[j]
	})
	sort.Slice(exp, func(i, j int) bool {
		return exp[i] < exp[j]
	})

	if !cmp.Equal(got, exp) {
		t.Fatalf("expect %+v task, got %+v", exp, got)
	}
}
