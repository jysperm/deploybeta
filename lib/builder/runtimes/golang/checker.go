package golang

import (
	"os"
	"path/filepath"
)

func existsInRoot(file string, root string) bool {
	path := filepath.Join(root, file)
	if _, err := os.Stat(path); err != nil {
		return false
	}
	return true
}

func CheckDep(root string) bool {
	return existsInRoot("Gopkg.lock", root) && existsInRoot("Gopkg.toml", root)
}

func CheckGlide(root string) bool {
	return existsInRoot("glide.yaml", root) || existsInRoot("glide.lock", root)
}
