package runtimes

import (
	"os"
	"path/filepath"
	"text/template"
)

func GenerateDockerfile(templatePath string, path string, data interface{}) error {
	dockerfileTemplate, err := template.ParseFiles(templatePath)
	if err != nil {
		return err
	}
	dockerfilePath := filepath.Join(path, "Dockerfile")
	dockerfile, err := os.OpenFile(dockerfilePath, os.O_WRONLY|os.O_CREATE|os.O_EXCL|os.O_TRUNC, 0666)
	defer dockerfile.Close()

	if err := dockerfileTemplate.Execute(dockerfile, data); err != nil {
		return err
	}

	return nil
}
