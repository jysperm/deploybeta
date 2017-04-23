package golang

import (
	"errors"
	"os"
	"path/filepath"
	"text/template"
)

var IsTesting = false

type Dockerfile struct {
	PackagePath string
	DepManager  string
	PackageName string
}

func GenerateDockerfile(root string, path string, name string) error {
	config := Dockerfile{
		PackagePath: path,
		PackageName: name,
		DepManager:  "",
	}

	execPath, err := os.Executable()
	if err != nil {
		return err
	}
	currentPath, _ := filepath.Split(execPath)
	if IsTesting {
		GOPATH := os.Getenv("GOPATH")
		currentPath = filepath.Join(GOPATH, "src", "github.com", "jysperm", "deploying")
	}
	templatePath := filepath.Join(currentPath, "lib", "builder", "runtimes", "golang", "Dockerfile.template")
	dockerfileTemplate, err := template.ParseFiles(templatePath)
	if err != nil && dockerfileTemplate == nil {
		return err
	}

	if CheckDep(root) {
		config.DepManager = "dep ensure"
	}

	if CheckGlide(root) {
		config.DepManager = "glide install"
	}

	if config.DepManager == "" {
		return errors.New("Not found a avaliable package manager")
	}

	dockerfilePath := filepath.Join(root, "Dockerfile")
	dockerfile, err := os.OpenFile(dockerfilePath, os.O_WRONLY|os.O_CREATE|os.O_EXCL|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer dockerfile.Close()

	if err := dockerfileTemplate.Execute(dockerfile, config); err != nil {
		return err
	}

	return nil
}
