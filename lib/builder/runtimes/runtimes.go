package runtimes

import (
	"errors"

	"github.com/jysperm/deploying/lib/builder/runtimes/golang"
	"github.com/jysperm/deploying/lib/builder/runtimes/node"
)

var ErrUnknowType = errors.New("unknown type of project")

func Dockerlize(root string, extra interface{}) error {
	if err := golang.Check(root); err == nil {
		err := golang.GenerateDockerfile(root, (extra).(string))
		if err != nil {
			return err
		}
		return nil
	}

	if err := node.Check(root); err == nil {
		err := node.GenerateDockerfile(root)
		if err != nil {
			return err
		}
		return nil
	}

	return ErrUnknowType
}
