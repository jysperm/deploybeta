package runtimes

import (
	"errors"
	"net"

	"github.com/docker/docker/api/types/swarm"
)

var ErrInvalidDataSourceType = errors.New("invalid dataSource type")

type DataSourceRuntime interface {
	DockerImageName() string
	ExposeProtocol() swarm.PortConfigProtocol
	ExposePort() uint16
	CheckNodeAvailability(host string) error
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

func checkTcpPort(host string) error {
	conn, err := net.Dial("tcp", host)

	if err != nil {
		return err
	} else {
		defer conn.Close()
		return nil
	}
}
