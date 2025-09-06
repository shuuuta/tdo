package project

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDetectRoot(t *testing.T) {
	tmpDir := t.TempDir()

	got1, err := DetectRoot(tmpDir)
	if err != nil {
		t.Fatal(err)
	}
	if got1 != "" {
		t.Fatalf("expect empty string, got %q", got1)
	}

	exp2 := filepath.Join(tmpDir, "exp2")
	fakeDir := filepath.Join(exp2, "bar", "fakeDir")
	startDir := filepath.Join(fakeDir, "baz", "qux")

	if err := os.MkdirAll(startDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(filepath.Join(exp2, ".git"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(filepath.Join(fakeDir, ".git"), 0755); err != nil {
		t.Fatal(err)
	}

	got2, err := DetectRoot(startDir)
	if err != nil {
		t.Fatal(err)
	}
	if got2 != exp2 {
		t.Fatalf("expect %q, got %q", exp2, got2)
	}
}
