package utils

import (
	"os"
	"path/filepath"
)

func GetAssetFilePath(relatedPath string) string {
	var err error

	workDir, useDirFromEnv := os.LookupEnv("WORKDIR")

	if !useDirFromEnv {
		workDir, err = os.Getwd()

		if err != nil {
			panic("os.Getwd failed: " + err.Error())
		}
	}

	return filepath.Join(workDir, "assets", relatedPath)
}
