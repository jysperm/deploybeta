package runtimes

import (
	"bytes"
	"errors"

	"github.com/docker/docker/api/types/swarm"

	"github.com/jysperm/deploybeta/lib/runtimes/datasource"
	"github.com/jysperm/deploybeta/lib/runtimes/golang"
	"github.com/jysperm/deploybeta/lib/runtimes/node"
)

var ErrUnknowType = errors.New("unknown type of project")
var ErrInvalidDataSourceType = errors.New("invalid dataSource type")

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

type DataSourceRuntime interface {
	DockerImageName() string
	ExposeProtocol() swarm.PortConfigProtocol
	ExposePort() uint16
}

func NewDataSourceRuntime(dataSourceType string) DataSourceRuntime {
	if dataSourceType == "mongodb" {
		return &datasource.MongoDBRuntime{}
	} else if dataSourceType == "redis" {
		return &datasource.RedisRuntime{}
	} else {
		panic(ErrInvalidDataSourceType)
	}
}
