package store

import (
	"fmt"
	"os"
	"time"

	"github.com/shuuuta/tdo/model"
)

func AddTask(projectPath, title string) (*model.Task, error) {
	t := model.Task{}

	p, err := LoadProject(projectPath)
	if err != nil {
		if os.IsNotExist(err) {
			t.ID = 0
			t.Title = title
			t.CreatedAt = time.Now().UTC()

			p := model.Project{
				ProjectPath: projectPath,
				IsGlobal:    false,
				Tasks:       []model.Task{t},
			}

			if err := SaveProject(&p); err != nil {
				return &t, err
			}
			return &t, nil
		} else {
			return &t, err
		}
	}

	id := getNextID(p.Tasks)

	t.ID = id
	t.Title = title
	t.CreatedAt = time.Now().UTC()

	p.Tasks = append(p.Tasks, t)

	if err := SaveProject(p); err != nil {
		return &t, err
	}

	return &t, nil
}

func AddGlobalTask(title string) (*model.Task, error) {
	t := model.Task{}

	p, err := LoadGlobalProject()
	if err != nil {
		if os.IsNotExist(err) {
			t.ID = 0
			t.Title = title
			t.CreatedAt = time.Now().UTC()

			p := model.Project{
				IsGlobal: true,
				Tasks:    []model.Task{t},
			}

			if err := SaveProject(&p); err != nil {
				return &t, err
			}
			return &t, nil
		} else {
			return &t, err
		}
	}

	id := getNextID(p.Tasks)

	t.ID = id
	t.Title = title
	t.CreatedAt = time.Now().UTC()

	p.Tasks = append(p.Tasks, t)

	if err := SaveProject(p); err != nil {
		return &t, err
	}

	return &t, nil
}

func RemoveTask(projectPath string, id int) error {
	p, err := LoadProject(projectPath)
	if err != nil {
		return err
	}

	t, err := removeTaskByID(p.Tasks, id)
	if err != nil {
		return err
	}

	p.Tasks = t

	if err := SaveProject(p); err != nil {
		return err
	}

	return nil
}

func RemoveGlobalTask(id int) error {
	p, err := LoadGlobalProject()
	if err != nil {
		return err
	}

	t, err := removeTaskByID(p.Tasks, id)
	if err != nil {
		return err
	}

	p.Tasks = t

	if err := SaveProject(p); err != nil {
		return err
	}

	return nil
}

func getNextID(tasks []model.Task) int {
	id := 0
	for _, v := range tasks {
		if id <= v.ID {
			id = v.ID + 1
		}
	}

	return id
}

func removeTaskByID(tasks []model.Task, id int) ([]model.Task, error) {
	n := 0
	hasID := false
	for i, v := range tasks {
		if v.ID == id {
			n = i
			hasID = true
		}
	}
	if !hasID {
		return tasks, fmt.Errorf("ID %d is not exist", id)
	}

	return tasks[:n+copy(tasks[:n], tasks[n+1:])], nil
}
