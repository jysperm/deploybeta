package runtimes

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/jysperm/deploybeta/lib/utils"
)

var ErrDepManagerNotFound = errors.New("No satisfied dep manager")

type GolangRuntime struct {
	BuildContext
	DepManager

	PackagePath string
	PackageName string
}

type DepManager struct {
	SpecFile       string
	LockFile       string
	InstallCommand string
}

var depManagers = map[string]DepManager{
	"dep": DepManager{
		SpecFile:       "Gopkg.toml",
		LockFile:       "Gopkg.lock",
		InstallCommand: "dep ensure",
	},
	"glide": DepManager{
		SpecFile:       "glide.yaml",
		LockFile:       "glide.lock",
		InstallCommand: "glide install",
	},
}

func DecideGolangRuntime(context *BuildContext) (bool, error) {
	_, err := decideDepManager(context)

	if err == ErrDepManagerNotFound {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}

func (runtime *GolangRuntime) Dockerfile() (*bytes.Buffer, error) {
	depManager, err := decideDepManager(&runtime.BuildContext)

	if err != nil {
		return nil, err
	}

	runtime.DepManager = depManagers[depManager]
	runtime.PackageName, runtime.PackagePath = extractInfo(runtime.gitUrl)

	return utils.ExecuteTemplateToBuffer(utils.GetAssetFilePath("runtime-golang/Dockerfile"), runtime)
}

func decideDepManager(context *BuildContext) (string, error) {
	for name, depManager := range depManagers {
		hasSpecFile, err := context.fileExists(depManager.SpecFile)

		if err != nil {
			return "", err
		} else if !hasSpecFile {
			continue
		}

		hasLockFile, err := context.fileExists(depManager.LockFile)

		if err != nil {
			return "", err
		} else if hasLockFile {
			return name, nil
		}
	}

	return "", ErrDepManagerNotFound
}

func extractInfo(remoteUrl string) (string, string) {
	byteURL := []byte(remoteUrl)
	if detectUrlType(remoteUrl) {
		slashIndex := strings.LastIndex(remoteUrl, "/") + 1
		gitIndex := strings.LastIndex(remoteUrl, ".git") - 1
		semicolonIndex := strings.LastIndex(remoteUrl, "://") + 3
		return string(byteURL[slashIndex:gitIndex]), string(byteURL[semicolonIndex:gitIndex])
	}
	atIndex := strings.LastIndex(remoteUrl, "@") + 1
	semicolonIndex := strings.LastIndex(remoteUrl, ":")
	slashIndex := strings.LastIndex(remoteUrl, "/") + 1
	gitIndex := strings.LastIndex(remoteUrl, ".git") - 1
	packagePath := fmt.Sprintf("%s/%s", string(byteURL[atIndex:semicolonIndex]), string(byteURL[semicolonIndex+1:gitIndex]))
	return string(byteURL[slashIndex:gitIndex]), packagePath
}

func detectUrlType(remoteUrl string) bool {
	validURLScheme := regexp.MustCompile("(https|http)")
	return validURLScheme.MatchString(remoteUrl)
}
