package store

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/shuuuta/tdo/model"
)

func TestSaveAndLoadProject(t *testing.T) {
	projectPath := "/foo/bar"
	now := time.Now().UTC()
	task1 := model.Task{
		ID:        1,
		Title:     "task 1",
		CreatedAt: now,
	}
	task2 := model.Task{
		ID:        2,
		Title:     "task 2",
		CreatedAt: now,
	}
	project := model.Project{
		ProjectPath: projectPath,
		Tasks: []model.Task{
			task1,
			task2,
		},
	}

	tmpDir := t.TempDir()

	if err := SaveProject(&project, tmpDir); err != nil {
		t.Fatal(err)
	}

	got, err := LoadProject(projectPath, tmpDir)
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(got, &project) {
		t.Fatalf("expect %+v, got %+v", &project, got)
	}
}

func TestSaveAndLoadGlobal(t *testing.T) {
	now := time.Now().UTC()
	task1 := model.Task{
		ID:        1,
		Title:     "task 1",
		CreatedAt: now,
	}
	task2 := model.Task{
		ID:        2,
		Title:     "task 2",
		CreatedAt: now,
	}
	project := model.Project{
		Tasks: []model.Task{
			task1,
			task2,
		},
	}

	tmpDir := t.TempDir()

	if err := SaveGlobal(&project, tmpDir); err != nil {
		t.Fatal(err)
	}

	got, err := LoadGlobal(tmpDir)
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(got, &project) {
		t.Fatalf("expect %+v, got %+v", &project, got)
	}
}
