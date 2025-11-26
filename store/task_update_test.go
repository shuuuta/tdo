package store

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/shuuuta/pask/model"
)

func TestUpdateTask(t *testing.T) {
	configDir = t.TempDir()

	t.Run("Update existing task", func(t *testing.T) {
		projectPath := "sample/project"
		targetID := 0
		p := &model.Project{
			ProjectPath: projectPath,
			Tasks: []model.Task{
				{
					ID:        targetID,
					Title:     "sample title",
					CreatedAt: time.Now().Add(time.Hour * -24).UTC(),
				},
				{
					ID:    1,
					Title: "remain title",
				},
			},
		}

		if err := SaveProject(p); err != nil {
			t.Fatal(err)
		}

		expTitle := "updated title"
		got1, err := UpdateTask(projectPath, expTitle, targetID)
		if err != nil {
			t.Fatal(err)
		}
		if got1.ID != targetID || got1.Title != expTitle {
			t.Fatalf("expect '%d: %s', got '%d: %s'", targetID, expTitle, got1.ID, got1.Title)
		}

		gotP, err := LoadProject(projectPath)
		if err != nil {
			t.Fatal(err)
		}

		validated := 0
		for _, v := range gotP.Tasks {
			if v.ID == targetID {
				if v.Title != expTitle {
					t.Fatalf("expect %q, got %q", expTitle, v.Title)
				}
				if v.CreatedAt != p.Tasks[targetID].CreatedAt {
					t.Fatalf("unexpected values were updated, exp\n'%+v',\n'%+v'", p.Tasks[1].ID, v.ID)
				}
				validated++
			}
		}

		if validated != 1 {
			t.Fatal("target task did not find")
		}
	})
}

func TestUpdateGlobalTask(t *testing.T) {
	configDir = t.TempDir()

	t.Run("Update global task", func(t *testing.T) {
		targetID := 0
		p := &model.Project{
			IsGlobal: true,
			Tasks: []model.Task{
				{
					ID:        targetID,
					Title:     "sample title",
					CreatedAt: time.Now().Add(time.Hour * -24).UTC(),
				},
				{
					ID:    1,
					Title: "remain title",
				},
			},
		}

		if err := SaveProject(p); err != nil {
			t.Fatal(err)
		}

		expTitle := "updated title"
		got1, err := UpdateGlobalTask(expTitle, targetID)
		if err != nil {
			t.Fatal(err)
		}
		if got1.ID != targetID || got1.Title != expTitle {
			t.Fatalf("expect '%d: %s', got '%d: %s'", targetID, expTitle, got1.ID, got1.Title)
		}

		gotP, err := LoadGlobalProject()
		if err != nil {
			t.Fatal(err)
		}

		validated := 0
		for _, v := range gotP.Tasks {
			if v.ID == targetID {
				if v.Title != expTitle {
					t.Fatalf("expect %q, got %q", expTitle, v.Title)
				}
				if v.CreatedAt != p.Tasks[targetID].CreatedAt {
					t.Fatalf("unrelated values were updated, exp\n'%+v',\n'%+v'", p.Tasks[1].ID, v.ID)
				}
				validated++
			}
		}

		if validated != 1 {
			t.Fatal("target task did not find")
		}
	})

	t.Run("Update does not affect other tasks", func(t *testing.T) {
		targetID := 1
		p := &model.Project{
			IsGlobal: true,
			Tasks: []model.Task{
				{
					ID:        0,
					Title:     "remain title 1",
					CreatedAt: time.Now().Add(time.Hour * -36).UTC(),
				},
				{
					ID:        targetID,
					Title:     "sample title",
					CreatedAt: time.Now().Add(time.Hour * -24).UTC(),
				},
				{
					ID:        2,
					Title:     "remain title 2",
					CreatedAt: time.Now().Add(time.Hour * -12).UTC(),
				},
			},
		}

		if err := SaveProject(p); err != nil {
			t.Fatal(err)
		}

		if _, err := UpdateGlobalTask("update title", targetID); err != nil {
			t.Fatal(err)
		}

		got1, err := LoadGlobalProject()
		if err != nil {
			t.Fatal(err)
		}
		if !cmp.Equal(p.Tasks[0], got1.Tasks[0]) {
			t.Fatalf("unexpected task was updated, exp\n'%+v',\n'%+v'", p.Tasks[0], got1.Tasks[0])
		}
		if !cmp.Equal(p.Tasks[2], got1.Tasks[2]) {
			t.Fatalf("unexpected task was updated, exp\n'%+v',\n'%+v'", p.Tasks[2], got1.Tasks[2])
		}
	})

	t.Run("Handle non-existent task ID", func(t *testing.T) {
		p := &model.Project{
			IsGlobal: true,
			Tasks: []model.Task{
				{
					ID:    0,
					Title: "sample title 1",
				},
				{
					ID:    1,
					Title: "sample title 2",
				},
				{
					ID:    2,
					Title: "sample title 3",
				},
			},
		}
		if err := SaveProject(p); err != nil {
			t.Fatal(err)
		}
		exp1 := "task 3 did not found"
		if _, err := UpdateGlobalTask("update task", 3); err == nil {
			t.Fatalf("expected return error")
		} else if err.Error() == exp1 {
			t.Fatalf("expected %q, got %q", exp1, err.Error())
		}
	})

	t.Run("Handle empty title", func(t *testing.T) {
		p := &model.Project{
			IsGlobal: true,
			Tasks: []model.Task{
				{
					ID:    0,
					Title: "sample title 1",
				},
			},
		}
		if err := SaveProject(p); err != nil {
			t.Fatal(err)
		}

		exp1 := "title string is required"
		if _, err := UpdateGlobalTask(" ", 0); err == nil {
			t.Fatalf("expected return error")
		} else if err.Error() != exp1 {
			t.Fatalf("expected %q, got %q", exp1, err.Error())
		}
	})

}
