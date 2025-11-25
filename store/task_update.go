package store

import (
	"fmt"
	"strings"

	"github.com/shuuuta/tdo/model"
)

func UpdateTask(projectPath, title string, id int) (*model.Task, error) {
	t := model.Task{}

	p, err := LoadProject(projectPath)
	if err != nil {
		return &t, err
	}

	return updateTaskFromProject(p, title, id)
}

func UpdateGlobalTask(title string, id int) (*model.Task, error) {
	t := model.Task{}

	p, err := LoadGlobalProject()
	if err != nil {
		return &t, err
	}

	return updateTaskFromProject(p, title, id)
}

func updateTaskFromProject(project *model.Project, title string, id int) (*model.Task, error) {
	t := model.Task{}

	ttl := strings.TrimSpace(title)
	if ttl == "" {
		return &t, fmt.Errorf("title string is required")
	}

	found := false
	for i, v := range project.Tasks {
		if v.ID == id {
			project.Tasks[i].Title = ttl
			t = project.Tasks[i]
			found = true
			break
		}
	}

	if !found {
		return &t, fmt.Errorf("task %d was not find", id)
	}

	if err := SaveProject(project); err != nil {
		return &t, err
	}
	return &t, nil
}
