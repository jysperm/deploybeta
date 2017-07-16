package golang

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/jysperm/deploying/lib/builder/runtimes"
	"github.com/jysperm/deploying/lib/utils"
)

type Dockerfile struct {
	PackagePath string
	DepManager  string
	PackageName string
}

func GenerateDockerfile(root string, remoteURL string) error {
	name, path := extractInfo(remoteURL)
	config := Dockerfile{
		PackagePath: path,
		PackageName: name,
		DepManager:  "",
	}

	templatePath, err := utils.GetAssetFilePath("lib/builder/runtimes/golang/Dockerfile.template")
	if err != nil {
		return err
	}

	if runtimes.CheckDep(root) {
		config.DepManager = "dep ensure"
	}

	if runtimes.CheckGlide(root) {
		config.DepManager = "glide install"
	}

	if config.DepManager == "" {
		return errors.New("Not found a avaliable package manager")
	}

	if err := runtimes.GenerateDockerfile(templatePath, root, config); err != nil {
		return err
	}

	return nil
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
