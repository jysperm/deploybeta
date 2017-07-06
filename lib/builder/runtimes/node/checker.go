package node

import (
	"os"
	"path/filepath"
)

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
