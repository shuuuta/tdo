package store

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sync"

	"github.com/shuuuta/tdo/log"
	"github.com/shuuuta/tdo/model"
	"github.com/shuuuta/tdo/utils"
)

var configDir string

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(fmt.Sprintf("cannot determine home directory: %v", err))
	}
	configDir = path.Join(home, ".config", "tdo")
}

func SaveProject(project *model.Project) error {
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

	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	fpath := filepath.Join(configDir, fname)
	if err := os.WriteFile(fpath, j, 0644); err != nil {
		return err
	}
	return nil
}

func LoadProject(projectPath string) (*model.Project, error) {
	var p model.Project

	h, err := utils.HashPath(projectPath)
	if err != nil {
		return &p, err
	}

	d, err := os.ReadFile(filepath.Join(configDir, h+".json"))
	if err != nil {
		return &p, err
	}

	if err := json.Unmarshal(d, &p); err != nil {
		return &p, err
	}

	return &p, nil
}

func LoadGlobalProject() (*model.Project, error) {
	var p model.Project

	d, err := os.ReadFile(filepath.Join(configDir, "global.json"))
	if err != nil {
		return &p, err
	}

	if err := json.Unmarshal(d, &p); err != nil {
		return &p, err
	}

	return &p, nil
}

func LoadAllProjects() ([]*model.Project, error) {
	var ps []*model.Project

	e, err := os.ReadDir(configDir)
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

			fpath := filepath.Join(configDir, fname)
			d, err := os.ReadFile(fpath)
			if err != nil {
				log.Logf("%q could not open: %s", fpath, err.Error())
				return
			}

			var p model.Project
			if err := json.Unmarshal(d, &p); err != nil {
				log.Logf("[INFO] %q is not project file: %s", fpath, err.Error())
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
