package runtimes

import (
	"errors"

	"github.com/jysperm/deploying/lib/builder/runtimes/golang"
	"github.com/jysperm/deploying/lib/builder/runtimes/node"
	"github.com/jysperm/deploying/lib/utils"
)

var ErrUnknowType = errors.New("unknown type of project")

func Dockerlize(root string, extra interface{}) error {
	if err := checkGo(root); err == nil {
		err := golang.GenerateDockerfile(root, (extra).(string))
		if err != nil {
			return err
		}
		return nil
	}

	if err := checkNodejs(root); err == nil {
		err := node.GenerateDockerfile(root)
		if err != nil {
			return err
		}
		return nil
	}

	return ErrUnknowType
}

func checkGo(root string) error {
	if utils.CheckDep(root) || utils.CheckGlide(root) {
		return nil
	}
	return ErrUnknowType
}

func checkNodejs(root string) error {
	if utils.ExistsInRoot("package.json", root) {
		return nil
	}
	return ErrUnknowType
}
