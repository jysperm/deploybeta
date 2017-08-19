package swarm

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"

	"github.com/jysperm/deploying/lib/builder"
	"github.com/jysperm/deploying/lib/etcd"
	"github.com/jysperm/deploying/lib/models/app"
	"github.com/jysperm/deploying/lib/models/version"
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

//UpdateService will update or create a app
func UpdateService(app app.Application) error {
	var create bool
	var err error
	serviceID, err := extractServiceID(app.Name)
	if err != nil && err.Error() == "Not found service" {
		create = true
	} else {
		create = false
	}
	if err != nil && err.Error() != "Not found service" {
		return err
	}

	currentVersion, err := version.FindByTag(app, app.Version)
	if err != nil {
		return err
	}

	repoTag, err := builder.LookupRepoTag(app.Name, currentVersion.Shasum)
	if err != nil {
		return err
	}

	var upstreamConfig UpstreamConfig
	uint64Instances := uint64(app.Instances)

	containerSpec := swarm.ContainerSpec{
		Image: repoTag,
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
		serviceResponse, err := swarmClient.ServiceCreate(context.Background(), serviceSpec, types.ServiceCreateOptions{})
		if err != nil {
			return err
		}
		serviceID = serviceResponse.ID
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
	if _, err := etcd.Client.Delete(context.Background(), upstreamKey); err != nil {
		return err
	}

	return nil
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
		return "", errors.New("Not found service")
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
