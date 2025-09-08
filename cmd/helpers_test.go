package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
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

func setupTest(t *testing.T, configPath string) func() {
	oListGlobal := listGlobal
	oAddGlobal := addGlobal

	return func() {
		listGlobal = oListGlobal
		addGlobal = oAddGlobal

		files, err := os.ReadDir(configPath)
		if err != nil {
			t.Fatal(err)
		}
		for _, v := range files {
			if err := os.Remove(filepath.Join(configPath, v.Name())); err != nil {
				t.Fatal(err)
			}
		}
	}
}
