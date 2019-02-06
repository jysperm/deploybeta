package datasources

import (
	"fmt"

	"github.com/docker/docker/api/types/swarm"

	"github.com/jysperm/deploybeta/config"
)

type MongoDBRuntime struct {
}

func (runtime *MongoDBRuntime) DockerImageName() string {
	return fmt.Sprintf("%s/%sdatasource-mongodb", config.DefaultRegistry, config.DockerPrefix)
}

func (runtime *MongoDBRuntime) ExposeProtocol() swarm.PortConfigProtocol {
	return swarm.PortConfigProtocolTCP
}

func (runtime *MongoDBRuntime) ExposePort() uint16 {
	return 27017
}

func (runtime *MongoDBRuntime) CheckNodeAvailability(host string) error {
  return checkTcpPort(host)
}
