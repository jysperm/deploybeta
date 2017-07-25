package utils

import (
	"os"
	"path/filepath"
)

func CheckYarn(root string) bool {
	return ExistsInRoot("yarn.lock", root)
}

func CheckDep(root string) bool {
	return ExistsInRoot("Gopkg.lock", root) && ExistsInRoot("Gopkg.toml", root)
}

func CheckGlide(root string) bool {
	return ExistsInRoot("glide.yaml", root) || ExistsInRoot("glide.lock", root)
}

func ExistsInRoot(file string, root string) bool {
	path := filepath.Join(root, file)
	if _, err := os.Stat(path); err != nil {
		return false
	}
	return true
}
