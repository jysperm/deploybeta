package generator

import (
	"os"
	"path/filepath"
	"text/template"
)

const (
	Go175 = "go:1.7.5"
	Go181 = "go:1.8.1"
)

type Dockerfile struct {
	GoVersion   string
	PackagePath string
	DepManager  string
	PackageName string
}

func GenerateDockerfile(root string, version string, path string, name string) error {
	config := Dockerfile{
		GoVersion:   version,
		PackagePath: path,
		PackageName: name,
	}
	currentPath, err := os.Executable()
	if err != nil {
		return err
	}
	templatePath := filepath.Join(currentPath, "builder", "runtimes", "golang", "Dockerfile.template")
	dockerfileTemplate, err := template.ParseFiles(templatePath)
	if err != nil && dockerfileTemplate == nil {
		return err
	}

	if err := CheckDep(root); err != nil {
		return err
	} else {
		config.DepManager = "dep ensure"
	}

	if err := CheckGlide(root); err != nil {
		return err
	} else {
		config.DepManager = "glide install"
	}

	dockerfilePath := filepath.Join(root, "Dockerfile")
	dockerfile, err := os.OpenFile(dockerfilePath, os.O_WRONLY|os.O_CREATE|os.O_EXCL|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}

	if err := dockerfileTemplate.Execute(dockerfile, config); err != nil {
		return err
	}

	return nil
}
