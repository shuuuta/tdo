package store

import (
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

// // タスク完了（MVPでは削除）
// func DoneTask(projectPath string, id int) error
func RemoveTask() {
	// load project json
	// if not exist

	// save file
}
