package store

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/shuuuta/tdo/model"
	"github.com/shuuuta/tdo/utils"
)

//// プロジェクト取得
//func LoadProject(path string) (*model.Project, error)
//
//// プロジェクト保存
//func SaveProject(project *model.Project) error
//
//// タスク追加
//func AddTask(projectPath, title string) (*model.Task, error)
//
//// タスク完了（MVPでは削除）
//func DoneTask(projectPath string, id int) error
//
//// プロジェクト横断タスク取得（-a オプション用）
//func LoadAllProjects() ([]*model.Project, error)
//
//// グローバルタスク読み書き
//func LoadGlobalTasks() ([]model.Task, error)
//func AddGlobalTask(title string) (*model.Task, error)

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

func SaveProject(project *model.Project, path string) error {
	j, err := json.MarshalIndent(project, "", "  ")
	if err != nil {
		return err
	}

	h, err := utils.HashPath(project.ProjectPath)
	if err != nil {
		return err
	}

	fpath := filepath.Join(path, h+".json")
	if err := os.WriteFile(fpath, j, 0644); err != nil {
		return err
	}
	return nil
}
