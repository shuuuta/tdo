package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/shuuuta/tdo/store"
)

func executeCommand(args ...string) (string, error) {
	cmd := rootCmd
	cmd.SetArgs(args)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)

	err := cmd.Execute()

	return buf.String(), err
}

type testEnv struct {
	TmpDir     string
	ConfigDir  string
	ProjectDir string
	Cleanup    func()
}

func setupTestEnv(t *testing.T) *testEnv {
	tmpDir := t.TempDir()

	oListGlobal := listGlobal
	oAddGlobal := addGlobal

	configDir := filepath.Join(tmpDir, "conf")
	os.MkdirAll(configDir, 0755)
	store.SetConfigDir(configDir)

	projectDir := filepath.Join(tmpDir, "project")
	os.MkdirAll(filepath.Join(projectDir, ".git"), 0755)

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	cleanup := func() {
		listGlobal = oListGlobal
		addGlobal = oAddGlobal

		files, err := os.ReadDir(configDir)
		if err != nil {
			t.Fatal(err)
		}
		for _, v := range files {
			if err := os.Remove(filepath.Join(configDir, v.Name())); err != nil {
				t.Fatal(err)
			}
		}
		os.Chdir(cwd)
	}

	return &testEnv{
		TmpDir:     tmpDir,
		ConfigDir:  configDir,
		ProjectDir: projectDir,
		Cleanup:    cleanup,
	}
}
