package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestHashPath(t *testing.T) {
	t.Run("output unaffected by trailing slash", func(t *testing.T) {
		p1 := "/foo/bar"
		h1, _ := HashPath(p1)
		if p1 == h1 {
			t.Fatalf("expect not equal path and hash %q", h1)
		}

		p2 := "/foo/bar/"
		h2, _ := HashPath(p2)

		if h1 != h2 {
			t.Fatalf("expect equal %q and %q", h1, h2)
		}

		if len(h1) != 12 {
			t.Fatalf("expect hash length 12, got %d: %v", len(h1), h1)
		}
	})

	t.Run("empty argment return empty string", func(t *testing.T) {
		h1, err := HashPath("")
		if err != nil {
			t.Fatal(err)
		}
		if h1 == "" {
			t.Fatalf("return empty text")
		}
	})

	t.Run("parse symlink before hash", func(t *testing.T) {
		tmpDir := t.TempDir()
		origin := filepath.Join(tmpDir, "origin")
		os.Mkdir(origin, 0755)
		alias := filepath.Join(tmpDir, "alias")
		os.Symlink(origin, alias)

		got1, err := HashPath(origin)
		if err != nil {
			t.Fatal(err)
		}
		got2, err := HashPath(alias)
		if err != nil {
			t.Fatal(err)
		}
		if got1 != got2 {
			t.Fatalf("expect same output %q and %q", got1, got2)
		}
	})
}
