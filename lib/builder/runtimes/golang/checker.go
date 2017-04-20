package generator

import (
	"os"
	"path/filepath"
)

func CheckDep(root string) error {
	manifest := filepath.Join(root, "manifest.json")
	lock := filepath.Join(root, "lock.json")
	if _, err := os.Stat(manifest); err != nil {
		return err
	}
	if _, err := os.Stat(lock); err != nil {
		return err
	}
	return nil
}

func CheckGlide(root string) error {
	glide := filepath.Join(root, "glide.yaml")
	lock := filepath.Join(root, "glide.lock")
	if _, err := os.Stat(glide); err != nil {
		return err
	}
	if _, err := os.Stat(lock); err != nil {
		return err
	}
	return nil
}
