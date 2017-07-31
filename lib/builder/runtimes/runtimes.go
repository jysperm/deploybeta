package runtimes

import (
	"errors"

	"github.com/jysperm/deploying/lib/builder/runtimes/golang"
	gohelpers "github.com/jysperm/deploying/lib/builder/runtimes/golang/helpers"
	"github.com/jysperm/deploying/lib/builder/runtimes/node"
	nodehelpers "github.com/jysperm/deploying/lib/builder/runtimes/node/helpers"
)

var ErrUnknowType = errors.New("unknown type of project")

func Dockerlize(root string, extra interface{}) error {
	if err := gohelpers.CheckGo(root); err == nil {
		err := golang.GenerateDockerfile(root, (extra).(string))
		if err != nil {
			return err
		}
		return nil
	}

	if err := nodehelpers.CheckNodejs(root); err == nil {
		err := node.GenerateDockerfile(root)
		if err != nil {
			return err
		}
		return nil
	}

	return ErrUnknowType
}
