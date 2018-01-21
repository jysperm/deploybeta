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

func UpdateApp(app *models.Application) error {
	currentVersion, err := models.FindVersionByTag(app, app.Version)
	if err != nil {
		return err
	}
	nameVersion := fmt.Sprintf("%s/%s:%s", config.DefaultRegistry, app.Name, currentVersion.Tag)

	var upstreamConfig UpstreamConfig

	if err := UpdateService(app.Name, uint64(app.Instances), []swarm.PortConfig{}, []swarm.NetworkAttachmentConfig{}, nameVersion); err != nil {
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

func JoinDataSource(app *models.Application, datasource *models.DataSource) error {
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

	currentVersion, err := models.FindVersionByTag(app, app.Version)
	if err != nil {
		return err
	}
	nameVersion := fmt.Sprintf("%s/%s:%s", config.DefaultRegistry, app.Name, currentVersion.Tag)

	return UpdateService(app.Name, uint64(app.Instances), []swarm.PortConfig{}, service.Spec.TaskTemplate.Networks, nameVersion)

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
