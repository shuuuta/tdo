package store

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/shuuuta/tdo/model"
	"github.com/shuuuta/tdo/utils"
)

//// タスク追加
//func AddTask(projectPath, title string) (*model.Task, error)
//
//// タスク完了（MVPでは削除）
//func DoneTask(projectPath string, id int) error

func SaveProject(project *model.Project, saveDir string) error {
	if project.IsGlobal && project.ProjectPath != "" {
		return fmt.Errorf("global project must not have project path.")
	}
	if !project.IsGlobal && project.ProjectPath == "" {
		return fmt.Errorf("project must have project path.")
	}

	j, err := json.MarshalIndent(project, "", "  ")
	if err != nil {
		return err
	}

	fname := "global.json"

	if !project.IsGlobal {
		h, err := utils.HashPath(project.ProjectPath)
		if err != nil {
			return err
		}
		fname = h + ".json"
	}

	fpath := filepath.Join(saveDir, fname)
	if err := os.WriteFile(fpath, j, 0644); err != nil {
		return err
	}
	return nil
}

func LoadProject(projectPath, saveDir string) (*model.Project, error) {
	var p model.Project

	h, err := utils.HashPath(projectPath)
	if err != nil {
		return &p, err
	}

	d, err := os.ReadFile(filepath.Join(saveDir, h+".json"))
	if err != nil {
		return &p, err
	}

	if err := json.Unmarshal(d, &p); err != nil {
		return &p, err
	}

	return &p, nil
}

func LoadGlobal(saveDir string) (*model.Project, error) {
	var p model.Project

	d, err := os.ReadFile(filepath.Join(saveDir, "global.json"))
	if err != nil {
		return &p, err
	}

	if err := json.Unmarshal(d, &p); err != nil {
		return &p, err
	}

	return &p, nil
}

func LoadAllProjects(saveDir string) ([]*model.Project, error) {
	var ps []*model.Project

	e, err := os.ReadDir(saveDir)
	if err != nil {
		return ps, err
	}

	var wg sync.WaitGroup
	c := make(chan *model.Project)
	for _, v := range e {
		wg.Add(1)
		go func(v os.DirEntry) {
			defer wg.Done()
			if v.IsDir() {
				return
			}
			fname := v.Name()
			ext := filepath.Ext(fname)
			if ext != ".json" {
				return
			}

			fpath := filepath.Join(saveDir, fname)
			fmt.Println(fpath)
			d, err := os.ReadFile(fpath)
			if err != nil {
				return
			}

			var p model.Project
			if err := json.Unmarshal(d, &p); err != nil {
				return
			}
			c <- &p
			return
		}(v)
	}

	go func() {
		wg.Wait()
		close(c)
	}()

	for p := range c {
		ps = append(ps, p)
	}

	return ps, nil
}
