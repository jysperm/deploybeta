package golang

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/jysperm/deploybeta/config"
	"github.com/jysperm/deploybeta/lib/utils"
)

type Dockerfile struct {
	PackagePath  string
	DepManager   string
	PackageName  string
	ProxyCommand string
	AptMirror    string
}

var ErrUnknowType = errors.New("unknown type of project")

func Check(root string) error {
	if checkDep(root) || checkGlide(root) {
		return nil
	}
	return ErrUnknowType
}

func GenerateDockerfile(root string, remoteURL string) (*bytes.Buffer, error) {
	name, path := extractInfo(remoteURL)
	cfg := Dockerfile{
		PackagePath:  path,
		PackageName:  name,
		DepManager:   "",
		ProxyCommand: "http_proxy=" + config.HttpProxy + " https_proxy=" + config.HttpProxy,
		AptMirror:    config.AptMirror,
	}

	templatePath := utils.GetAssetFilePath("runtime-go/Dockerfile.template")

	if checkDep(root) {
		cfg.DepManager = "dep ensure"
	}

	if checkGlide(root) {
		cfg.DepManager = "glide install"
	}

	dockerfileTemplate, err := template.ParseFiles(templatePath)
	if err != nil {
		return nil, err
	}

	fileBuffer := new(bytes.Buffer)
	fileWriter := bufio.NewWriter(fileBuffer)

	if err := dockerfileTemplate.Execute(fileWriter, cfg); err != nil {
		return nil, err
	}

	fileWriter.Flush()

	return fileBuffer, nil
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

func checkDep(root string) bool {
	return existsInRoot("Gopkg.lock", root) && existsInRoot("Gopkg.toml", root)
}

func checkGlide(root string) bool {
	return existsInRoot("glide.yaml", root) || existsInRoot("glide.lock", root)
}

func existsInRoot(file string, root string) bool {
	path := filepath.Join(root, file)
	if _, err := os.Stat(path); err != nil {
		return false
	}
	return true
}
