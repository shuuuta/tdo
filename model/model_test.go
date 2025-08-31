package model

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestTaskJSONRoundTrip(t *testing.T) {
	now := time.Now().UTC()
	task := Task{ID: 1, Title: "sample title", CreatedAt: now}

	b, err := json.Marshal(task)
	if err != nil {
		t.Fatal(err)
	}

	var got Task
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatal(err)
	}

	if !cmp.Equal(got, task) {
		t.Errorf("expected: \n%+v, \ngot: \n%+v", task, got)
	}
}

func TestProjectJSONRoundTrip(t *testing.T) {
	now := time.Now().UTC()
	task1 := Task{ID: 1, Title: "sample title 1", CreatedAt: now}
	task2 := Task{ID: 2, Title: "sample title 2", CreatedAt: now}
	project := Project{
		ProjectPath: "./path/to/project",
		Tasks: []Task{
			task1,
			task2,
		},
	}

	b, err := json.Marshal(project)
	if err != nil {
		t.Fatal(err)
	}

	var got Project
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatal(err)
	}

	if !cmp.Equal(got, project) {
		t.Errorf("expected: \n%+v, \ngot: \n%+v", project, got)
	}
}
