package swarm

import (
	"context"
	"errors"

	"github.com/docker/docker/api/types/filters"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
)

var ErrServiceNotFound = errors.New("Not found service")
var swarmClient *client.Client

type SwarmService interface {
	SwarmServiceName() string
}

type Container struct {
	State      string `json:"state"`
	Image      string `json:"image,omitempty"`
	VersionTag string `json:"versionTag, omitempty"`
	CreatedAt  string `json:"createdAt"`
}

func init() {
	var err error
	swarmClient, err = client.NewEnvClient()
	if err != nil {
		panic(err)
	}
}

func UpdateService(service SwarmService, instances uint64, portConfig []swarm.PortConfig, networkConfig []swarm.NetworkAttachmentConfig, image string, envs []string) error {
	var serviceName = service.SwarmServiceName()

	var create bool
	var err error

	serviceID, err := RetrieveServiceID(serviceName)
	if err == ErrServiceNotFound {
		create = true
	} else {
		create = false
	}

	if err != nil && err != ErrServiceNotFound {
		return err
	}

	containerSpec := swarm.ContainerSpec{
		Image: image,
		Labels: map[string]string{
			"deploying.name": serviceName,
		},
	}
	if len(envs) != 0 {
		containerSpec.Env = envs
	}

	taskSpec := swarm.TaskSpec{
		ContainerSpec: containerSpec,
		LogDriver: &swarm.Driver{
			Name: "json-file",
			Options: map[string]string{
				"labels": "deploying.name",
			},
		},
		Networks: networkConfig,
	}

	replicatedService := swarm.ReplicatedService{Replicas: &instances}
	serviceMode := swarm.ServiceMode{Replicated: &replicatedService}
	var endpointSpec swarm.EndpointSpec
	if len(portConfig) != 0 {
		endpointSpec.Mode = "vip"
		endpointSpec.Ports = portConfig
	} else {
		port := swarm.PortConfig{
			Protocol:    swarm.PortConfigProtocolTCP,
			TargetPort:  3000,
			PublishMode: swarm.PortConfigPublishModeIngress,
		}
		endpointSpec.Mode = "vip"
		endpointSpec.Ports = []swarm.PortConfig{port}
	}

	serviceSpec := swarm.ServiceSpec{
		Annotations:  swarm.Annotations{Name: serviceName},
		TaskTemplate: taskSpec,
		Mode:         serviceMode,
		EndpointSpec: &endpointSpec,
		Networks:     networkConfig,
	}

	if create {
		serviceRes, err := swarmClient.ServiceCreate(context.Background(), serviceSpec, types.ServiceCreateOptions{})
		if err != nil {
			return err
		}
		serviceID = serviceRes.ID
	} else {
		version, err := RetrieveServiceVersion(serviceName)
		if err != nil {
			return err
		}

		if _, err := swarmClient.ServiceUpdate(context.Background(), serviceID, *version, serviceSpec, types.ServiceUpdateOptions{}); err != nil {
			return err
		}
	}

	return nil
}

func RemoveService(service SwarmService) error {
	serviceID, err := RetrieveServiceID(service.SwarmServiceName())
	if err != nil {
		return err
	}

	if err := swarmClient.ServiceRemove(context.Background(), serviceID); err != nil {
		return err
	}

	return nil
}

func ListContainers(service SwarmService) ([]Container, error) {
	filter := filters.NewArgs()
	filter.Add("service", service.SwarmServiceName())
	listOpts := types.TaskListOptions{
		Filters: filter,
	}

	tasks, err := swarmClient.TaskList(context.Background(), listOpts)
	if err != nil {
		return nil, err
	}
	if len(tasks) == 0 {
		return []Container{}, nil
	}

	var containers []Container
	for _, v := range tasks {
		c := Container{
			State:     string(v.Status.State),
			Image:     v.Spec.ContainerSpec.Image,
			CreatedAt: v.Status.Timestamp.String(),
		}
		containers = append(containers, c)
	}

	return containers, nil
}
