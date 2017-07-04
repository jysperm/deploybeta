package utils

import (
	"os"
	"path/filepath"
)

func GetAssetFilePath(relatedPath string) (string, error) {
	if os.Getenv("WORKDIR") == "" {
		execPath, err := os.Executable()
		if err != nil {
			return "", err
		}
		currentPath, _ := filepath.Split(execPath)
		return filepath.Join(currentPath, relatedPath), nil
	}

	return filepath.Join(os.Getenv("WORKDIR"), relatedPath), nil
}
