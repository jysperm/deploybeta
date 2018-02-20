package swarm

import (
	"errors"
	"fmt"

	"github.com/docker/docker/api/types/swarm"
	"github.com/jysperm/deploying/config"
	"github.com/jysperm/deploying/lib/models"
)

var ErrNetworkNotFound = errors.New("Network not found")

func UpdateDataSource(dataSource *models.DataSource, instances uint64) error {
	image := fmt.Sprintf("%s/%s:latest", config.DefaultRegistry, dataSource.Type)

	networkID, err := FindNetworkByName(dataSource.Name)
	if err != nil {
		return err
	}
	if networkID == "" {
		networkID, err = CreateOverlay(dataSource)
		if err != nil {
			return err
		}
	}

	networkOpts := swarm.NetworkAttachmentConfig{
		Target: networkID,
	}

	var portConfig swarm.PortConfig
	if dataSource.Type == "redis" {
		portConfig.Protocol = swarm.PortConfigProtocolTCP
		portConfig.TargetPort = uint32(config.DefaultRedisPort)
	} else if dataSource.Type == "mongodb" {
		portConfig.Protocol = swarm.PortConfigProtocolTCP
		portConfig.TargetPort = uint32(config.DefaultMongoDBPort)
	}

	environments := []string{
		"DATASOURCE_NAME=" + dataSource.Name,
		"DEPLOYING_URL=http://" + config.HostPrivateAddress + config.Listen,
	}

	return UpdateService(dataSource.Name, instances, []swarm.PortConfig{portConfig}, []swarm.NetworkAttachmentConfig{networkOpts}, image, environments)

}

func RemoveDataSource(datasource *models.DataSource) error {
	if err := RemoveService(datasource.Name); err != nil {
		return err
	}

	return RemoveOverlay(datasource)
}

func ListDataSources() []models.DataSource {
	//TODO
	return []models.DataSource{}
}
