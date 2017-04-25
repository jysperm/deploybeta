package golang

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
)

var IsTesting = false

type Dockerfile struct {
	PackagePath string
	DepManager  string
	PackageName string
}

func detectType(remoteURL string) bool {
	validURLScheme := regexp.MustCompile("(https|http)")
	return validURLScheme.MatchString(remoteURL)
}

func extractInfo(remoteURL string) (string, string) {
	byteURL := []byte(remoteURL)
	if detectType(remoteURL) {
		slashIndex := strings.LastIndex(remoteURL, "/") + 1
		gitIndex := strings.LastIndex(remoteURL, ".git") - 1
		semicolonIndex := strings.LastIndex(remoteURL, "://") + 3
		return string(byteURL[slashIndex:gitIndex]), string(byteURL[semicolonIndex:gitIndex])
	}
	atIndex := strings.LastIndex(remoteURL, "@") + 1
	semicolonIndex := strings.LastIndex(remoteURL, ":")
	slashIndex := strings.LastIndex(remoteURL, "/") + 1
	gitIndex := strings.LastIndex(remoteURL, ".git") - 1
	packagePath := fmt.Sprintf("%s/%s", string(byteURL[atIndex:semicolonIndex]), string(byteURL[semicolonIndex+1:gitIndex]))
	return string(byteURL[slashIndex:gitIndex]), packagePath
}

func GenerateDockerfile(root string, remoteURL string) error {
	name, path := extractInfo(remoteURL)
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
