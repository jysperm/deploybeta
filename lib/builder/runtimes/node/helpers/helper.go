package helpers

import (
	"errors"
	"os"
	"path/filepath"
)

var ErrUnknowType = errors.New("unknown type of project")

func CheckNodejs(root string) error {
	if existsInRoot("package.json", root) {
		return nil
	}
	return ErrUnknowType
}

func CheckYarn(root string) bool {
	return existsInRoot("yarn.lock", root)
}

func existsInRoot(file string, root string) bool {
	path := filepath.Join(root, file)
	if _, err := os.Stat(path); err != nil {
		return false
	}
	return true
}
