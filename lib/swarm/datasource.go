package swarm

import (
	"errors"

	"github.com/docker/docker/api/types/swarm"
	"github.com/jysperm/deploybeta/config"
	"github.com/jysperm/deploybeta/lib/models"
	"github.com/jysperm/deploybeta/lib/runtimes"
)

var ErrNetworkNotFound = errors.New("Network not found")

func UpdateDataSource(dataSource *models.DataSource) error {
	networkID, err := FindNetworkByName(dataSource.SwarmNetworkName())
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

	runtime := runtimes.NewDataSourceRuntime(dataSource.Type)

	portConfig := swarm.PortConfig{
		Protocol:   runtime.ExposeProtocol(),
		TargetPort: uint32(runtime.ExposePort()),
	}

	environments := []string{
		"AGENT_TOKEN=" + dataSource.AgentToken,
		"DATASOURCE_NAME=" + dataSource.Name,
		"DEPLOYBETA_URL=http://" + config.HostPrivateAddress + config.Listen,
	}

	return UpdateService(dataSource, []swarm.PortConfig{portConfig}, []swarm.NetworkAttachmentConfig{networkOpts}, runtime.DockerImageName(), environments)
}

func RemoveDataSource(datasource *models.DataSource) error {
	if err := RemoveService(datasource); err != nil {
		return err
	}

	return RemoveOverlay(datasource)
}

func ListDataSources() []models.DataSource {
	//TODO
	return []models.DataSource{}
}
