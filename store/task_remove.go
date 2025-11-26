package store

import (
	"fmt"
	"slices"

	"github.com/shuuuta/pask/model"
)

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

	return slices.Delete(tasks, n, n+1), nil
}
