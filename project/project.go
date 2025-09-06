package project

import (
	"os"
	"path/filepath"
)

func DetectRoot(startDir string) (string, error) {
	out := ""
	dir := startDir
	for {
		if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
			out = dir
		} else if !os.IsNotExist(err) {
			return "", err
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return out, nil
}
