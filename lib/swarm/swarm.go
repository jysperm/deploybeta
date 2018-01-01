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

type Container struct {
	State      string `json:"state"`
	Image      string `json:"image,omitempty"`
	VersionTag string `json:"versionTag,omitempty`
	CreatedAt  string `json:"createdAt"`
}

func init() {
	var err error
	swarmClient, err = client.NewEnvClient()
	if err != nil {
		panic(err)
	}
}

func UpdateService(name string, instances uint64, portConfig []swarm.PortConfig, networkConfig []swarm.NetworkAttachmentConfig, image string) error {
	var create bool
	var err error

	serviceID, err := RetrieveServiceID(name)
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
			"deploying.name": name,
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
		Annotations:  swarm.Annotations{Name: name},
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
		version, err := RetrieveServiceVersion(serviceID)
		if err != nil {
			return err
		}

		if _, err := swarmClient.ServiceUpdate(context.Background(), serviceID, *version, serviceSpec, types.ServiceUpdateOptions{}); err != nil {
			return err
		}
	}

	return nil
}

func RemoveService(name string) error {
	serviceID, err := RetrieveServiceID(name)
	if err != nil {
		return err
	}

	if err := swarmClient.ServiceRemove(context.Background(), serviceID); err != nil {
		return err
	}

	return nil
}

func ListContainers(name string) (*[]Container, error) {
	filter := filters.NewArgs()
	filter.Add("service", name)
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
			State:     string(v.Status.State),
			Image:     v.Spec.ContainerSpec.Image,
			CreatedAt: v.Status.Timestamp.String(),
		}
		containers = append(containers, c)
	}

	return &containers, nil
}
