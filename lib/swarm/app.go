package swarm

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jysperm/deploybeta/config"
	"github.com/jysperm/deploybeta/lib/db"
	"github.com/jysperm/deploybeta/lib/models"

	"github.com/docker/docker/api/types/swarm"
)

type UpstreamConfig struct {
	Port uint32 `json:"port"`
}

var ErrNetworkJoined = errors.New("Had joined the network")
var ErrNetworkNoUnlinkable = errors.New("No network could be unlinking")

func UpdateAppService(app *models.Application) error {
	if app.VersionTag == "" {
		return nil
	}

	networkConfigs := []swarm.NetworkAttachmentConfig{}
	environments := []string{}

	dataSources := make([]models.DataSource, 0)

	err := app.DataSources().FetchAll(&dataSources)

	if err != nil {
		return err
	}

	for _, dataSource := range dataSources {
		dataSourceService, _, err := swarmClient.ServiceInspectWithRaw(context.TODO(), dataSource.SwarmServiceName())

		if err != nil {
			return err
		}

		networkConfigs = append(networkConfigs, swarm.NetworkAttachmentConfig{
			Target: dataSource.SwarmNetworkName(),
		})

		key := fmt.Sprintf("DATA_SOURCE_%s", strings.Replace(strings.ToUpper(dataSource.Name), "-", "_", -1))
		value := fmt.Sprintf("%s:%d", config.HostPrivateAddress, getServicePort(&dataSourceService))

		environments = append(environments, fmt.Sprintf("%s=%s", key, value))
	}

	version := models.Version{}
	err = app.Version().Fetch(&version)

	if err != nil {
		return err
	}

	if err := UpdateService(app, []swarm.PortConfig{}, networkConfigs, version.DockerImageName(), environments); err != nil {
		return err
	}

	appService, _, err := swarmClient.ServiceInspectWithRaw(context.TODO(), app.SwarmServiceName())

	if err != nil {
		return err
	}

	upstreams := []UpstreamConfig{
		UpstreamConfig{
			Port: getServicePort(&appService),
		},
	}

	if err := db.PutJSON(fmt.Sprintf("/upstreams/%s", app.Name), upstreams); err != nil {
		return err
	}

	return nil
}

func RemoveApp(app *models.Application) error {
	if err := db.DeleteKey(fmt.Sprintf("/upstreams/%s", app.Name)); err != nil {
		return err
	}

	return RemoveService(app)
}

func ListNodes(app *models.Application) ([]Container, error) {
	containers, err := ListContainers(app)
	if err != nil {
		return nil, err
	}

	if len(containers) == 0 {
		return []Container{}, nil
	}

	for i := 0; i < len(containers); i++ {
		containers[i].Image = ""
		containers[i].VersionTag = app.VersionTag
	}

	return containers, nil
}
