package utils

import (
	"os"
	"path/filepath"
)

func GetAssetFilePath(relatedPath string) (string, error) {
	if dir := os.Getenv("WORKDIR"); dir != "" {
		return filepath.Join(dir, relatedPath), nil
	}

	workDir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Join(workDir, relatedPath), nil
}
