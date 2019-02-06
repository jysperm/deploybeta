package datasources

import (
	"fmt"

	"github.com/docker/docker/api/types/swarm"

	"github.com/jysperm/deploybeta/config"
)

type RedisRuntime struct {
}

func (runtime *RedisRuntime) DockerImageName() string {
	return fmt.Sprintf("%s/%sdatasource-redis", config.DefaultRegistry, config.DockerPrefix)
}

func (runtime *RedisRuntime) ExposeProtocol() swarm.PortConfigProtocol {
	return swarm.PortConfigProtocolTCP
}

func (runtime *RedisRuntime) ExposePort() uint16 {
	return 6379
}

func (runtime *RedisRuntime) CheckNodeAvailability(host string) error {
  return checkTcpPort(host)
}
