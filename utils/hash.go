package utils

import (
	"crypto/sha1"
	"encoding/hex"
	"path/filepath"
)

func HashPath(path string) (string, error) {
	apath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	h := sha1.New()
	h.Write([]byte(apath))

	r := hex.EncodeToString(h.Sum(nil))

	return r[:12], nil
}
