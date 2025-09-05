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

	id := 0
	for _, v := range p.Tasks {
		if id <= v.ID {
			id = v.ID + 1
		}
	}

	t.ID = id
	t.Title = title
	t.CreatedAt = time.Now().UTC()

	p.Tasks = append(p.Tasks, t)

	SaveProject(p)

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

	id := 0
	for _, v := range p.Tasks {
		if id <= v.ID {
			id = v.ID + 1
		}
	}

	t.ID = id
	t.Title = title
	t.CreatedAt = time.Now().UTC()

	p.Tasks = append(p.Tasks, t)

	SaveProject(p)

	return &t, nil
}

func RemoveTask(projectPath string, id int) error {
	p, err := LoadProject(projectPath)
	if err != nil {
		return err
	}

	n := 0
	hasID := false
	for i, v := range p.Tasks {
		if v.ID == id {
			n = i
			hasID = true
		}
	}
	if !hasID {
		return fmt.Errorf("ID %d is not exist in %s", id, projectPath)
	}

	p.Tasks = p.Tasks[:n+copy(p.Tasks[:n], p.Tasks[n+1:])]

	SaveProject(p)

	return nil
}

func RemoveGlobalTask(id int) error {
	p, err := LoadGlobalProject()
	if err != nil {
		return err
	}

	n := 0
	hasID := false
	for i, v := range p.Tasks {
		if v.ID == id {
			n = i
			hasID = true
		}
	}
	if !hasID {
		return fmt.Errorf("ID %d is not exist in Global task", id)
	}

	p.Tasks = p.Tasks[:n+copy(p.Tasks[:n], p.Tasks[n+1:])]

	SaveProject(p)

	return nil
}
