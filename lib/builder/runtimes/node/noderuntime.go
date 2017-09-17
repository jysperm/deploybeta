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

	"github.com/jysperm/deploying/lib/utils"
)

type Dockerfile struct {
	NodeVersion string
	HasYarn     bool
}

var ErrUnknowType = errors.New("unknown type of project")

func Check(root string) error {
	if existsInRoot("package.json", root) {
		return nil
	}
	return ErrUnknowType
}

func GenerateDockerfile(root string) (*bytes.Buffer, error) {
	config := Dockerfile{
		HasYarn: false,
	}

	node, err := extraVersion(root)
	if err != nil {
		return nil, err
	}
	config.NodeVersion = node

	templatePath, err := utils.GetAssetFilePath("assets/templates/node/Dockerfile.template")
	if err != nil {
		return nil, err
	}

	if checkYarn(root) {
		config.HasYarn = true
	}

	dockerfileTemplate, err := template.ParseFiles(templatePath)
	if err != nil {
		return nil, err
	}

	fileBuffer := new(bytes.Buffer)
	fileWriter := bufio.NewWriter(fileBuffer)

	if err := dockerfileTemplate.Execute(fileWriter, config); err != nil {
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
