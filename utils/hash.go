package utils

import (
	"crypto/sha1"
	"encoding/hex"
	"os"
	"path/filepath"
)

func HashPath(path string) (string, error) {
	out := path
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		out, err = filepath.EvalSymlinks(path)
		if err != nil {
			return "", err
		}
	}

	apath, err := filepath.Abs(out)
	if err != nil {
		return "", err
	}

	h := sha1.New()
	h.Write([]byte(apath))

	r := hex.EncodeToString(h.Sum(nil))

	return r[:12], nil
}
