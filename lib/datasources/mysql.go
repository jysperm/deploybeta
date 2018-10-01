package datasources

import (
	"fmt"

	"github.com/docker/docker/api/types/swarm"

	"github.com/jysperm/deploybeta/config"
)

type MySQLRuntime struct {
}

func (runtime *MySQLRuntime) DockerImageName() string {
	return fmt.Sprintf("%s/%sdatasource-mysql", config.DefaultRegistry, config.DockerPrefix)
}

func (runtime *MySQLRuntime) ExposeProtocol() swarm.PortConfigProtocol {
	return swarm.PortConfigProtocolTCP
}

func (runtime *MySQLRuntime) ExposePort() uint16 {
	return 3306
}
