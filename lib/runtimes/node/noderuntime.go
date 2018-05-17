package node

import (
	"bufio"
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"

	"github.com/blang/semver"
	"github.com/buger/jsonparser"
	"github.com/parnurzeal/gorequest"

	"github.com/jysperm/deploybeta/config"
	"github.com/jysperm/deploybeta/lib/utils"
)

type Dockerfile struct {
	NodeVersion string
	HasYarn     bool
	HTTPProxy   string
	HTTPSProxy  string
	AptCnMirror string
	NpmCnMirror string
	NvmCnMirror string
}

var ErrUnknowType = errors.New("unknown type of project")

func Check(root string) error {
	if existsInRoot("package.json", root) {
		return nil
	}
	return ErrUnknowType
}

func GenerateDockerfile(root string) (*bytes.Buffer, error) {
	cfg := Dockerfile{
		HasYarn:     false,
		HTTPProxy:   config.HttpProxy,
		HTTPSProxy:  config.HttpsProxy,
		AptCnMirror: config.AptCnMirror,
		NpmCnMirror: config.NpmCnMirror,
		NvmCnMirror: config.NvmCnMirror,
	}

	node, err := extraVersion(root)
	if err != nil {
		return nil, err
	}
	cfg.NodeVersion = node

	templatePath := utils.GetAssetFilePath("runtime-node/Dockerfile.template")

	if checkYarn(root) {
		cfg.HasYarn = true
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

	if err := fileWriter.Flush(); err != nil {
		return nil, err
	}

	return fileBuffer, nil
}

func fetchVerisonList() ([]byte, error) {
	res, body, errs := gorequest.New().Get("http://nodejs.org/dist/index.json").EndBytes()
	if len(errs) != 0 || res.StatusCode != 200 {
		return nil, errors.New("Bad request for fetching list of versions")
	}

	return body, nil
}

func parseVerion(path string) (string, error) {
	packagePath := filepath.Join(path, "package.json")
	packageInfo, err := ioutil.ReadFile(packagePath)
	if err != nil {
		return "", err
	}

	nodeVersion, err := jsonparser.GetString(packageInfo, "engines", "node")
	if err != nil {
		return "", err
	}

	return nodeVersion, nil
}

func extraVersion(path string) (string, error) {
	list, err := fetchVerisonList()
	if err != nil {
		return "", err
	}

	nodeVersion, err := parseVerion(path)
	if err != nil {
		return "", err
	}

	nodeRange, err := semver.ParseRange(nodeVersion)
	if err != nil {
		return "", err
	}

	exactNode := ""
	_, err = jsonparser.ArrayEach(list, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		node, _, _, _ := jsonparser.Get(value, "version")
		node = node[1:]

		if nodeRange(semver.MustParse(string(node))) && exactNode == "" {
			exactNode = string(node)
		}
	})

	if err != nil {
		return "", err
	}

	return exactNode, nil
}

func checkYarn(root string) bool {
	return existsInRoot("yarn.lock", root)
}

func existsInRoot(file string, root string) bool {
	path := filepath.Join(root, file)
	if _, err := os.Stat(path); err != nil {
		return false
	}
	return true
}
