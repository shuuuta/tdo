package store

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/shuuuta/pask/model"
)

func TestRemoveTask(t *testing.T) {
	configDir = t.TempDir()

	ppath := "/foo/bar"
	targetT := "target title"
	remainT := "remain title"

	p := &model.Project{
		ProjectPath: ppath,
		Tasks: []model.Task{
			{
				ID:    0,
				Title: remainT,
			},
			{
				ID:    1,
				Title: targetT,
			},
		},
	}

	if err := SaveProject(p); err != nil {
		t.Fatal(err)
	}

	if err := RemoveTask(ppath, 1); err != nil {
		t.Fatal(err)
	}

	got1, err := LoadProject(ppath)
	if err != nil {
		t.Fatal(err)
	}

	if len(got1.Tasks) != 1 {
		t.Fatalf("expect number of task is 1, got %d", len(got1.Tasks))
	}
	if got1.Tasks[0].Title != remainT {
		t.Fatalf("expect %q, got %q", remainT, got1.Tasks[0].Title)
	}

	if err := RemoveTask(ppath, 9); err == nil {
		t.Fatal("expected error for non-existent ID")
	}
}

func TestRemoveGlobalTask(t *testing.T) {
	configDir = t.TempDir()

	targetT := "target title"
	remainT := "remain title"

	p := &model.Project{
		IsGlobal: true,
		Tasks: []model.Task{
			{
				ID:    0,
				Title: remainT,
			},
			{
				ID:    1,
				Title: targetT,
			},
		},
	}

	if err := SaveProject(p); err != nil {
		t.Fatal(err)
	}

	if err := RemoveGlobalTask(1); err != nil {
		t.Fatal(err)
	}

	got1, err := LoadGlobalProject()
	if err != nil {
		t.Fatal(err)
	}

	if len(got1.Tasks) != 1 {
		t.Fatalf("expect number of task is 1, got %d", len(got1.Tasks))
	}
	if got1.Tasks[0].Title != remainT {
		t.Fatalf("expect %q, got %q", remainT, got1.Tasks[0].Title)
	}

	if err := RemoveGlobalTask(9); err == nil {
		t.Fatal("expected error for non-existent ID")
	}
}

func TestRemoveTaskByID(t *testing.T) {
	var tasks1 []model.Task
	if _, err := removeTaskByID(tasks1, 1); err == nil {
		t.Fatal("expect error for non-existent ID")
	}

	tasks2 := []model.Task{
		{ID: 0},
		{ID: 2},
		{ID: 4},
		{ID: 7},
		{ID: 9},
	}
	tasks2, _ = removeTaskByID(tasks2, 0)
	tasks2, _ = removeTaskByID(tasks2, 4)
	tasks2, _ = removeTaskByID(tasks2, 9)
	exp2 := []model.Task{
		{ID: 2},
		{ID: 7},
	}
	if !cmp.Equal(tasks2, exp2) {
		t.Fatalf("expect %+v, got %+v", exp2, tasks2)
	}

	tasks3 := []model.Task{
		{ID: 1},
		{ID: 2},
		{ID: 2},
		{ID: 7},
	}
	tasks3, _ = removeTaskByID(tasks3, 2)
	exp3 := []model.Task{
		{ID: 1},
		{ID: 2},
		{ID: 7},
	}
	if !cmp.Equal(tasks3, exp3) {
		t.Fatalf("expect %+v, got %+v", exp3, tasks3)
	}
}
