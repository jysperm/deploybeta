package node

import (
	"os"
	"path/filepath"
	"text/template"

	"github.com/jysperm/deploying/lib/utils"
)

type Dockerfile struct {
	NodeVersion string
	HasYarn     bool
}

func GenerateDockerfile(root string, version string) error {
	config := Dockerfile{
		NodeVersion: version,
		HasYarn:     false,
	}

	templatePath, err := utils.GetAssetFilePath("lib/builder/runtimes/node/Dockerfile.template")
	if err != nil {
		return err
	}

	if utils.CheckYarn(root) {
		config.HasYarn = true
	}

	dockerfileTemplate, err := template.ParseFiles(templatePath)
	if err != nil {
		return err
	}

	dockerfilePath := filepath.Join(root, "Dockerfile")
	Dockerfile, err := os.OpenFile(dockerfilePath, os.O_WRONLY|os.O_CREATE|os.O_EXCL|os.O_TRUNC, 0666)
	defer Dockerfile.Close()

	if err := dockerfileTemplate.Execute(Dockerfile, config); err != nil {
		return err
	}

	return nil

}
