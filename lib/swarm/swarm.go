package swarm

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jysperm/deploying/config"
	"github.com/jysperm/deploying/lib/etcd"
	"github.com/jysperm/deploying/lib/models"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
)

type UpstreamConfig struct {
	Port uint32 `json:"port"`
}

type Container struct {
	State      string `json:"state"`
	VersionTag string `json:"versionTag"`
	CreatedAt  string `json:"createdAt:`
}

var ErrNotFoundService = errors.New("Not found service")
var swarmClient *client.Client

func init() {
	var err error
	swarmClient, err = client.NewEnvClient()
	if err != nil {
		panic(err)
	}
}

func UpdateService(app *models.Application) error {
	var create bool
	var err error
	serviceID, err := extractServiceID(app.Name)
	if err == ErrNotFoundService {
		create = true
	} else {
		create = false
	}

	if err != nil && err != ErrNotFoundService {
		return err
	}

	currentVersion, err := models.FindVersionByTag(app, app.Version)
	if err != nil {
		return err
	}
	nameVersion := fmt.Sprintf("%s/%s:%s", config.DefaultRegistry, app.Name, currentVersion.Tag)

	var upstreamConfig UpstreamConfig
	uint64Instances := uint64(app.Instances)

	containerSpec := swarm.ContainerSpec{
		Image: nameVersion,
		Labels: map[string]string{
			"deploying.name": app.Name,
		},
	}

	taskSpec := swarm.TaskSpec{
		ContainerSpec: containerSpec,
		LogDriver: &swarm.Driver{
			Name: "json-file",
			Options: map[string]string{
				"labels": "deploying.name",
			},
		},
	}

	replicatedService := swarm.ReplicatedService{Replicas: &uint64Instances}
	serviceMode := swarm.ServiceMode{Replicated: &replicatedService}
	portConfig := swarm.PortConfig{
		Protocol:    swarm.PortConfigProtocolTCP,
		TargetPort:  3000,
		PublishMode: swarm.PortConfigPublishModeIngress,
	}
	endpointSpec := swarm.EndpointSpec{
		Mode:  "vip",
		Ports: []swarm.PortConfig{portConfig},
	}
	serviceSpec := swarm.ServiceSpec{
		Annotations:  swarm.Annotations{Name: app.Name},
		TaskTemplate: taskSpec,
		Mode:         serviceMode,
		EndpointSpec: &endpointSpec,
	}

	if create {
		serviceRes, err := swarmClient.ServiceCreate(context.Background(), serviceSpec, types.ServiceCreateOptions{})
		if err != nil {
			return err
		}
		serviceID = serviceRes.ID
	} else {
		internalVersion, err := extractInternalVersion(serviceID)
		if err != nil {
			return err
		}

		if _, err := swarmClient.ServiceUpdate(context.Background(), serviceID, internalVersion, serviceSpec, types.ServiceUpdateOptions{}); err != nil {
			return err
		}
	}

	upstreamConfig.Port, err = extractPort(serviceID)
	if err != nil {
		return err
	}

	upstream, err := json.Marshal([]UpstreamConfig{upstreamConfig})
	if err != nil {
		return err
	}
	upstreamKey := fmt.Sprintf("/upstreams/%s", app.Name)
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

func RemoveService(app *models.Application) error {
	serviceID, err := extractServiceID(app.Name)
	if err == ErrNotFoundService {
		return nil
	}

	if err != nil {
		return err
	}

	if err := swarmClient.ServiceRemove(context.Background(), serviceID); err != nil {
		return err
	}

	upstreamKey := fmt.Sprintf("/upstreams/%s", app.Name)
	if _, err := etcd.Client.Delete(context.Background(), upstreamKey); err != nil {
		return err
	}

	if err := models.DeleteAppByName(app.Name); err != nil {
		return err
	}

	if err := models.DeleteAllVersion(app); err != nil {
		return err
	}

	return nil
}

func ListContainers(app *models.Application) (*[]Container, error) {
	filter := filters.NewArgs()
	filter.Add("service", app.Name)
	listOpts := types.TaskListOptions{
		Filters: filter,
	}

	tasks, err := swarmClient.TaskList(context.Background(), listOpts)
	if err != nil {
		return nil, err
	}
	if len(tasks) == 0 {
		return nil, nil
	}

	var containers []Container
	for _, v := range tasks {
		c := Container{
			State:      string(v.Status.State),
			VersionTag: app.Version,
			CreatedAt:  v.Status.Timestamp.String(),
		}
		containers = append(containers, c)
	}

	return &containers, nil
}

func extractServiceID(name string) (string, error) {
	var serviceID string
	query := filters.NewArgs()
	query.Add("name", name)
	listOpts := types.ServiceListOptions{
		Filters: query,
	}
	services, err := swarmClient.ServiceList(context.Background(), listOpts)
	if err != nil {
		return "", err
	}
	for _, i := range services {
		if i.Spec.Annotations.Name == name {
			serviceID = i.ID
			break
		}
	}
	if serviceID == "" {
		return "", ErrNotFoundService
	}
	return serviceID, nil
}

func extractPort(serviceID string) (uint32, error) {
	var srv swarm.Service
	var err error
	var portConfig swarm.PortConfig
	for {
		srv, _, err = swarmClient.ServiceInspectWithRaw(context.Background(), serviceID)
		if err != nil {
			return 0, err
		}
		if len(srv.Endpoint.Ports) != 0 {
			portConfig = srv.Endpoint.Ports[0]
			break
		}
	}
	return portConfig.PublishedPort, nil
}

func extractInternalVersion(serviceID string) (swarm.Version, error) {
	service, _, err := swarmClient.ServiceInspectWithRaw(context.Background(), serviceID)
	if err != nil {
		return swarm.Version{}, err
	}
	return service.Meta.Version, nil
}
