package datasources

import (
	"errors"

	"github.com/docker/docker/api/types/swarm"
)

var ErrInvalidDataSourceType = errors.New("invalid dataSource type")

type DataSourceRuntime interface {
	DockerImageName() string
	ExposeProtocol() swarm.PortConfigProtocol
	ExposePort() uint16
}

func NewDataSourceRuntime(dataSourceType string) DataSourceRuntime {
	if dataSourceType == "mongodb" {
		return &MongoDBRuntime{}
	} else if dataSourceType == "mysql" {
		return &MySQLRuntime{}
	} else if dataSourceType == "redis" {
		return &RedisRuntime{}
	} else {
		panic(ErrInvalidDataSourceType)
	}
}
