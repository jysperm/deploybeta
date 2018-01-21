package swarm

import (
	"errors"
	"fmt"

	"github.com/docker/docker/api/types/swarm"
	"github.com/jysperm/deploying/config"
	"github.com/jysperm/deploying/lib/models"
)

var ErrNetworkNotFound = errors.New("Network not found")

func UpdateDataSource(datasource *models.DataSource, instances uint64) error {
	image := fmt.Sprintf("%s/%s:latest", config.DefaultRegistry, datasource.Type)

	networkID, err := FindNetworkByName(datasource.Name)
	if err != nil {
		return err
	}
	if networkID == "" {
		networkID, err = CreateOverlay(datasource)
		if err != nil {
			return err
		}
	}

	networkOpts := swarm.NetworkAttachmentConfig{
		Target: networkID,
	}

	return UpdateService(datasource.Name, instances, []swarm.PortConfig{}, []swarm.NetworkAttachmentConfig{networkOpts}, image)

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
