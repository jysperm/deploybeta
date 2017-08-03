package runtimes

import (
	"bytes"
	"errors"

	"github.com/jysperm/deploying/lib/builder/runtimes/golang"
	"github.com/jysperm/deploying/lib/builder/runtimes/node"
)

var ErrUnknowType = errors.New("unknown type of project")

func Dockerlize(root string, extra interface{}) (*bytes.Buffer, error) {
	if err := golang.Check(root); err == nil {
		buf, err := golang.GenerateDockerfile(root, (extra).(string))
		if err != nil {
			return nil, err
		}
		return buf, nil
	}

	if err := node.Check(root); err == nil {
		buf, err := node.GenerateDockerfile(root)
		if err != nil {
			return nil, err
		}
		return buf, nil
	}

	return nil, ErrUnknowType
}
