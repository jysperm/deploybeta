package swarm

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
)

func RetrieveServiceID(name string) (string, error) {
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
		return "", ErrServiceNotFound
	}
	return serviceID, nil
}

func RetrievePort(serviceID string) (uint32, error) {
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

func RetrieveServiceVersion(serviceID string) (*swarm.Version, error) {
	service, _, err := swarmClient.ServiceInspectWithRaw(context.Background(), serviceID)
	if err != nil {
		return nil, err
	}
	return &service.Meta.Version, nil
}
