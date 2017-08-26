package utils

import (
	"os"
	"path/filepath"
)

func GetAssetFilePath(relatedPath string) (string, error) {
	workDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return filepath.Join(workDir, relatedPath), nil
}
