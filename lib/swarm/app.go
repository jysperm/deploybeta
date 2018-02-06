package swarm

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jysperm/deploying/config"
	"github.com/jysperm/deploying/lib/etcd"
	"github.com/jysperm/deploying/lib/models"

	"github.com/docker/docker/api/types/swarm"
)

type UpstreamConfig struct {
	Port uint32 `json:"port"`
}

var ErrNetworkJoined = errors.New("Had joined the network")
var ErrNetworkNoUnlinkable = errors.New("No network could be unlinking")

func UpdateApp(app *models.Application) error {
	currentVersion, err := models.FindVersionByTag(app, app.Version)
	if err != nil {
		return err
	}
	nameVersion := fmt.Sprintf("%s/%s:%s", config.DefaultRegistry, app.Name, currentVersion.Tag)

	var upstreamConfig UpstreamConfig

	if err := UpdateService(app.Name, uint64(app.Instances), []swarm.PortConfig{}, []swarm.NetworkAttachmentConfig{}, nameVersion, []string{}); err != nil {
		return err
	}

	serviceID, _ := RetrieveServiceID(app.Name)
	upstreamConfig.Port, err = RetrievePort(serviceID)
	if err != nil {
		return err
	}

	upstream, err := json.Marshal([]UpstreamConfig{upstreamConfig})
	if err != nil {
		return err
	}

	upstreamKey := fmt.Sprintf("/upstream/%s", app.Name)
	if _, err := etcd.Client.Put(context.Background(), upstreamKey, string(upstream)); err != nil {
		return err
	}

	if err := app.Update(app); err != nil {
		return err
	}

	if err := app.UpdateVersion(currentVersion.Tag); err != nil {
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

	return RemoveService(app.Name)
}

func LinkDataSource(app *models.Application, datasource *models.DataSource) error {
	networkID, err := FindNetworkByName(datasource.Name)
	if err != nil {
		return err
	}
	if networkID == "" {
		return ErrNetworkNotFound
	}

	serviceID, err := RetrieveServiceID(app.Name)
	if err != nil {
		return err
	}
	if serviceID == "" {
		return ErrServiceNotFound
	}

	service, _, err := swarmClient.ServiceInspectWithRaw(context.Background(), serviceID)
	if err != nil {
		return err
	}

	for _, v := range service.Spec.TaskTemplate.Networks {
		if v.Target == networkID {
			return ErrNetworkJoined
		}
	}

	networkOpts := swarm.NetworkAttachmentConfig{
		Target: networkID,
	}
	service.Spec.TaskTemplate.Networks = append(service.Spec.TaskTemplate.Networks, networkOpts)

	envs := service.Spec.TaskTemplate.ContainerSpec.Env

	currentVersion, err := models.FindVersionByTag(app, app.Version)
	if err != nil {
		return err
	}
	nameVersion := fmt.Sprintf("%s/%s:%s", config.DefaultRegistry, app.Name, currentVersion.Tag)

	port, err := RetrievePort(serviceID)
	if err != nil {
		return err
	}
	dataSourceEnv := fmt.Sprintf("%s=%s:%d", datasource.Name, config.HostPrivateAddress, port)
	envs = append(envs, dataSourceEnv)

	return UpdateService(app.Name, uint64(app.Instances), []swarm.PortConfig{}, service.Spec.TaskTemplate.Networks, nameVersion, envs)

}

func UnlinkDataSource(app *models.Application, dataSource *models.DataSource) error {
	networkID, err := FindNetworkByName(dataSource.Name)
	if err != nil {
		return err
	}

	if networkID == "" {
		return ErrNetworkNotFound
	}

	serviceID, err := RetrieveServiceID(app.Name)
	if err != nil {
		return err
	}

	if serviceID == "" {
		return ErrServiceNotFound
	}

	service, _, err := swarmClient.ServiceInspectWithRaw(context.Background(), serviceID)
	if err != nil {
		return err
	}

	if len(service.Spec.TaskTemplate.Networks) == 0 {
		return ErrNetworkNoUnlinkable
	}

	port, err := RetrievePort(serviceID)
	if err != nil {
		return err
	}

	dataSourceEnv := fmt.Sprintf("%s=%s:%d", dataSource.Name, config.HostPrivateAddress, port)

	for i := 0; i < len(service.Spec.TaskTemplate.Networks); i++ {
		if service.Spec.TaskTemplate.Networks[i].Target == networkID {
			service.Spec.TaskTemplate.Networks = append(service.Spec.TaskTemplate.Networks[:i], service.Spec.TaskTemplate.Networks[i+1:]...)
			for j := 0; j < len(service.Spec.TaskTemplate.ContainerSpec.Env); j++ {
				if service.Spec.TaskTemplate.ContainerSpec.Env[i] == dataSourceEnv {
					service.Spec.TaskTemplate.ContainerSpec.Env = append(service.Spec.TaskTemplate.ContainerSpec.Env[:i], service.Spec.TaskTemplate.ContainerSpec.Env[i+1:]...)
					currentVersion, err := models.FindVersionByTag(app, app.Version)
					if err != nil {
						return err
					}
					nameVersion := fmt.Sprintf("%s/%s:%s", config.DefaultRegistry, app.Name, currentVersion.Tag)
					return UpdateService(app.Name, uint64(app.Instances), []swarm.PortConfig{}, service.Spec.TaskTemplate.Networks, nameVersion, service.Spec.TaskTemplate.ContainerSpec.Env)
				}
			}
		}
	}

	return ErrNetworkNoUnlinkable
}

func ListNodes(app *models.Application) ([]Container, error) {
	containers, err := ListContainers(app.Name)
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
