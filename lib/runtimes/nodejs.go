package runtimes

import (
	"bytes"
	"errors"
	"strings"

	"github.com/blang/semver"
	"github.com/buger/jsonparser"
	"github.com/parnurzeal/gorequest"

	"github.com/jysperm/deploybeta/config"
	"github.com/jysperm/deploybeta/lib/utils"
)

type NodejsRuntime struct {
	BuildContext

	packageJson []byte
	NodeVersion string
	UseYarn     bool
	NpmMirror   string
	NvmMirror   string
}

func DecideNodejsRuntime(context *BuildContext) (bool, error) {
	return context.fileExists("package.json")
}

func (runtime *NodejsRuntime) Dockerfile() (*bytes.Buffer, error) {
	packageJson, err := runtime.readFile("package.json")

	if err != nil {
		return nil, err
	}

	nodeVersion, err := decideNodejsVersion(packageJson)

	if err != nil {
		return nil, err
	}

	useYarn, err := runtime.fileExists("yarn.lock")

	if err != nil {
		return nil, err
	}

	runtime.NodeVersion = nodeVersion
	runtime.UseYarn = useYarn
	runtime.NpmMirror = config.NpmMirror
	runtime.NvmMirror = config.NvmMirror

	return utils.ExecuteTemplateToBuffer(utils.GetAssetFilePath("runtime-nodejs/Dockerfile"), runtime)
}

func decideNodejsVersion(packageJson []byte) (string, error) {
	constraint, err := jsonparser.GetString(packageJson, "engines", "node")

	if err != nil {
		return "", err
	}

	versions, err := fetchNodejsVerisons()

	if err != nil {
		return "", err
	}

	constraintRange, err := semver.ParseRange(constraint)

	if err != nil {
		return "", err
	}

	for _, version := range versions {
		if constraintRange(semver.MustParse(version)) {
			return version, nil
		}
	}

	return "", errors.New("No satisfied Node version found: " + constraint)
}

func fetchNodejsVerisons() ([]string, error) {
	res, body, errs := gorequest.New().Proxy(config.HttpProxy).Get("http://nodejs.org/dist/index.json").EndBytes()

	if len(errs) != 0 || res.StatusCode != 200 {
		return nil, errors.New("Bad request for fetching list of versions: " + string(res.StatusCode))
	}

	versions := make([]string, 0)

	_, err := jsonparser.ArrayEach(body, func(value []byte, _ jsonparser.ValueType, _ int, _ error) {
		version, err := jsonparser.GetString(value, "version")

		if err == nil {
			versions = append(versions, strings.TrimLeft(string(version), "v"))
		}
	})

	if err != nil {
		return nil, err
	}

	return versions, nil
}
