package store

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/shuuuta/tdo/model"
)

func TestSaveAndLoadProject(t *testing.T) {
	projectPath := "/foo/bar"
	now := time.Now().UTC()
	project := model.Project{
		ProjectPath: projectPath,
		Tasks: []model.Task{
			{
				ID:        1,
				Title:     "task 1",
				CreatedAt: now,
			},
			{
				ID:        2,
				Title:     "task 2",
				CreatedAt: now,
			},
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
	project := model.Project{
		IsGlobal: true,
		Tasks: []model.Task{
			{
				ID:        1,
				Title:     "task 1",
				CreatedAt: now,
			},
			{
				ID:        2,
				Title:     "task 2",
				CreatedAt: now,
			},
		},
	}

	tmpDir := t.TempDir()

	if err := SaveProject(&project, tmpDir); err != nil {
		t.Fatal(err)
	}

	got, err := LoadGlobal(tmpDir)
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(got.Tasks, project.Tasks) {
		t.Fatalf("expect %+v, got %+v", project.Tasks, got.Tasks)
	}
	if !got.IsGlobal {
		t.Fatalf("expect IsGlobal true, got %v", got.IsGlobal)
	}
}

func TestSaveProjectWithoutPath(t *testing.T) {
	project := model.Project{
		ProjectPath: "",
		Tasks: []model.Task{
			{
				ID:        1,
				Title:     "task 1",
				CreatedAt: time.Now().UTC(),
			},
		},
	}

	if err := SaveProject(&project, t.TempDir()); err == nil {
		t.Fatal("expect return error")
	}
}

func TestSaveGlobalWithPath(t *testing.T) {
	project := model.Project{
		IsGlobal:    true,
		ProjectPath: "/path/to/project",
		Tasks: []model.Task{
			{
				ID:        1,
				Title:     "task 1",
				CreatedAt: time.Now().UTC(),
			},
		},
	}

	if err := SaveProject(&project, t.TempDir()); err == nil {
		t.Fatal("expect return error")
	}
}

func TestLoadAll(t *testing.T) {
	now := time.Now().UTC()
	project1 := model.Project{
		ProjectPath: "/path/to/project1",
		Tasks: []model.Task{
			{
				ID:        1,
				Title:     "task 1",
				CreatedAt: now,
			},
			{
				ID:        2,
				Title:     "task 2",
				CreatedAt: now,
			},
		},
	}
	project2 := model.Project{
		ProjectPath: "/path/to/project2",
		Tasks: []model.Task{
			{
				ID:        3,
				Title:     "task 3",
				CreatedAt: now,
			},
			{
				ID:        4,
				Title:     "task 4",
				CreatedAt: now,
			},
		},
	}
	projectGlobal := model.Project{
		IsGlobal: true,
		Tasks: []model.Task{
			{
				ID:        5,
				Title:     "task 3",
				CreatedAt: now,
			},
			{
				ID:        6,
				Title:     "task 6",
				CreatedAt: now,
			},
		},
	}

	tmpDir := t.TempDir()

	if err := SaveProject(&project1, tmpDir); err != nil {
		t.Fatal(err)
	}
	if err := SaveProject(&project2, tmpDir); err != nil {
		t.Fatal(err)
	}
	if err := SaveProject(&projectGlobal, tmpDir); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "dummy.json"), []byte("{dummy: 1}"), 0644); err != nil {
		t.Fatal(err)
	}

	allP := []*model.Project{&project1, &project2, &projectGlobal}

	got, err := LoadAllProjects(tmpDir)
	if err != nil {
		t.Fatal(err)
	}

	trans := cmp.Transformer("Sort", func(in []*model.Project) []*model.Project {
		sort.Slice(in, func(i, j int) bool {
			return in[i].ProjectPath < in[j].ProjectPath
		})
		return in
	})

	if !cmp.Equal(got, allP, trans) {
		var expName []string
		var gotName []string
		for _, v := range allP {
			expName = append(expName, v.ProjectPath)
		}
		for _, v := range got {
			gotName = append(gotName, v.ProjectPath)
		}
		t.Fatalf("expect %+v, got %+v", expName, gotName)
	}
	var expName []string
	var gotName []string
	for _, v := range allP {
		expName = append(expName, v.ProjectPath)
	}
	for _, v := range got {
		gotName = append(gotName, v.ProjectPath)
	}
	fmt.Printf("expect %+v, got %+v\n", expName, gotName)
}
