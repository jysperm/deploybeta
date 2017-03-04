package swarm

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"

	"github.com/jysperm/deploying/lib/models/app"
	"github.com/jysperm/deploying/lib/services"
	"github.com/jysperm/deploying/lib/services/builder"

	"golang.org/x/net/context"
)

//UpstreamConfig define the structure of upstream config
type UpstreamConfig struct {
	Port uint32 `json:"port"`
}

var swarmClient *client.Client

func init() {
	var err error
	swarmClient, err = client.NewEnvClient()
	if err != nil {
		panic(err)
	}
}

//CreateService will build a image from Dockerfile and deploy a service
func CreateService(app app.Application) error {
	versionTag := fmt.Sprintf("%s:%s", app.Name, app.Version)
	buildOpts := types.ImageBuildOptions{
		Tags: []string{versionTag},
	}
	shasum, err := builder.BuildImage(buildOpts, app.GitRepository)
	if err != nil {
		return err
	}

	uint64Instances := uint64(app.Instances)
	imageName := fmt.Sprintf("%s@sha256:%s", versionTag, shasum)
	containerSpec := swarm.ContainerSpec{Image: imageName}
	taskSpec := swarm.TaskSpec{ContainerSpec: containerSpec}
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
	serviceResponse, err := swarmClient.ServiceCreate(context.Background(), serviceSpec, types.ServiceCreateOptions{})
	if err != nil {
		return err
	}

	publishedPort, err := extractPort(serviceResponse.ID)
	if err != nil {
		return err
	}
	upstream, err := json.Marshal([]UpstreamConfig{UpstreamConfig{Port: publishedPort}})
	if err != nil {
		return err
	}
	upstreamKey := fmt.Sprintf("/upstream/%s", app.Name)
	if _, err := services.EtcdClient.Put(context.Background(), upstreamKey, string(upstream)); err != nil {
		return err
	}

	return nil
}

//UpdateService will updaate the config of given app
func UpdateService(app app.Application) error {
	serviceID, err := extractServiceID(app.Name)
	if err != nil {
		return err
	}

	versionTag := fmt.Sprintf("%s:%s", app.Name, app.Version)
	buildOpts := types.ImageBuildOptions{
		Tags: []string{versionTag},
	}
	shasum, err := builder.BuildImage(buildOpts, app.GitRepository)
	if err != nil {
		return err
	}

	uint64Instances := uint64(app.Instances)
	imageName := fmt.Sprintf("%s@sha256:%s", versionTag, shasum)
	containerSpec := swarm.ContainerSpec{Image: imageName}
	taskSpec := swarm.TaskSpec{ContainerSpec: containerSpec}
	replcatedService := swarm.ReplicatedService{Replicas: &uint64Instances}
	serviceMode := swarm.ServiceMode{Replicated: &replcatedService}
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
	_, err = swarmClient.ServiceUpdate(context.Background(), serviceID, swarm.Version{}, serviceSpec, types.ServiceUpdateOptions{})
	if err != nil {
		return err
	}

	publishedPort, err := extractPort(serviceID)
	if err != nil {
		return err
	}
	upstream, err := json.Marshal([]UpstreamConfig{UpstreamConfig{Port: publishedPort}})
	if err != nil {
		return err
	}
	upstreamKey := fmt.Sprintf("/upstream/%s", app.Name)
	if _, err := services.EtcdClient.Put(context.Background(), upstreamKey, string(upstream)); err != nil {
		return err
	}

	return nil
}

//RemoveService will remove the given service
func RemoveService(app app.Application) error {
	serviceID, err := extractServiceID(app.Name)
	if err != nil {
		return err
	}
	if err := swarmClient.ServiceRemove(context.Background(), serviceID); err != nil {
		return err
	}

	upstreamKey := fmt.Sprintf("/upstream/%s", app.Name)
	if _, err := services.EtcdClient.Delete(context.Background(), upstreamKey); err != nil {
		return err
	}

	return nil
}

func extractServiceID(name string) (string, error) {
	var serviceID string
	services, err := swarmClient.ServiceList(context.Background(), types.ServiceListOptions{})
	if err != nil {
		return "", err
	}

	for _, i := range services {
		if i.Spec.Annotations.Name == name {
			serviceID = i.ID
		}
	}
	if serviceID == "" {
		return "", errors.New("Not found service")
	}
	return serviceID, nil
}

func extractPort(serviceID string) (uint32, error) {
	service, _, err := swarmClient.ServiceInspectWithRaw(context.Background(), serviceID)
	if err != nil {
		return 0, nil
	}
	portConfig := service.Endpoint.Ports[0]
	return portConfig.PublishedPort, nil
}
