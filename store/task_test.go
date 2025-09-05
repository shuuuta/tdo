package store

import (
	"sort"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/shuuuta/tdo/model"
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

func TestAddGlobalTask(t *testing.T) {
	configDir = t.TempDir()

	title1 := "creata project and add task"
	if _, err := AddGlobalTask(title1); err != nil {
		t.Fatal(err)
	}

	lp1, err := LoadGlobalProject()
	if err != nil {
		t.Fatal(err)
	}
	if len(lp1.Tasks) != 1 || lp1.Tasks[0].Title != title1 {
		t.Fatalf("expect 1 task, got %+v", lp1.Tasks)
	}

	title2 := "add extra task"
	AddGlobalTask(title2)

	lp2, err := LoadGlobalProject()
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

func TestGetNextID(t *testing.T) {
	var tasks []model.Task
	id1 := getNextID(tasks)
	if id1 != 0 {
		t.Fatalf("expect 0, got %d", id1)
	}

	tasks = append(tasks, model.Task{
		ID: 6,
	})
	id2 := getNextID(tasks)
	if id2 != 7 {
		t.Fatalf("expect 7, got %d", id2)
	}

	tasks = append(tasks, model.Task{ID: 1}, model.Task{ID: 3}, model.Task{ID: 8})
	id3 := getNextID(tasks)
	if id3 != 9 {
		t.Fatalf("expect 9, got %d", id3)
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
