package node

import (
	"github.com/jysperm/deploying/lib/builder/runtimes"
	"github.com/jysperm/deploying/lib/utils"
)

type Dockerfile struct {
	NodeVersion string
	HasYarn     bool
}

func GenerateDockerfile(vNode string, root string) error {
	if vNode == "" {
		vNode = `'lts/*'`
	}
	config := Dockerfile{
		NodeVersion: vNode,
		HasYarn:     false,
	}

	templatePath, err := utils.GetAssetFilePath("lib/builder/runtimes/node/Dockerfile.template")
	if err != nil {
		return err
	}

	if CheckYarn(root) {
		config.HasYarn = true
	}

	if runtimes.GenerateDockerfile(templatePath, root, config); err != nil {
		return err
	}

	return nil

}
