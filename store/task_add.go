package store

import (
	"os"
	"time"

	"github.com/shuuuta/pask/model"
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

func getNextID(tasks []model.Task) int {
	id := 0
	for _, v := range tasks {
		if id <= v.ID {
			id = v.ID + 1
		}
	}

	return id
}
