package helpers

import (
	"errors"
	"os"
	"path/filepath"
)

var ErrUnknowType = errors.New("unknown type of project")

func CheckGo(root string) error {
	if CheckDep(root) || CheckGlide(root) {
		return nil
	}
	return ErrUnknowType
}

func CheckDep(root string) bool {
	return existsInRoot("Gopkg.lock", root) && existsInRoot("Gopkg.toml", root)
}

func CheckGlide(root string) bool {
	return existsInRoot("glide.yaml", root) || existsInRoot("glide.lock", root)
}

func existsInRoot(file string, root string) bool {
	path := filepath.Join(root, file)
	if _, err := os.Stat(path); err != nil {
		return false
	}
	return true
}
