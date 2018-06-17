package swarm

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jysperm/deploybeta/config"
	"github.com/jysperm/deploybeta/lib/etcd"
	"github.com/jysperm/deploybeta/lib/models"

	"github.com/docker/docker/api/types/swarm"
)

// Serialize to /upstreams/:appName
type UpstreamConfig struct {
	Port uint32 `json:"port"`
}

var ErrNetworkJoined = errors.New("Had joined the network")
var ErrNetworkNoUnlinkable = errors.New("No network could be unlinking")

func UpdateAppService(app *models.Application) error {
	networkConfigs := []swarm.NetworkAttachmentConfig{}
	environments := []string{}

	dataSources, err := models.GetDataSourcesOfApp(app)

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

		key := fmt.Sprintf("DATA_SOURCE_%s", strings.ToUpper(dataSource.Name))
		value := fmt.Sprintf("%s:%d", config.HostPrivateAddress, getServicePort(&dataSourceService))

		environments = append(environments, fmt.Sprintf("%s=%s", key, value))
	}

	version, err := models.FindVersionByTag(app, app.Version)

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

	if err := etcd.PutKey(fmt.Sprintf("/upstreams/%s", app.Name), upstreams); err != nil {
		return err
	}

	return nil
}

func RemoveApp(app *models.Application) error {
	upstreamKey := fmt.Sprintf("/upstreams/%s", app.Name)
	if _, err := etcd.Client.Delete(context.Background(), upstreamKey); err != nil {
		return err
	}

	if err := models.DeleteAppByName(app.Name); err != nil {
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
		containers[i].VersionTag = app.Version
	}

	return containers, nil
}
