package swarm

import (
	"github.com/docker/docker/api/types"
	dockerSwarm "github.com/docker/docker/api/types/swarm"
	"golang.org/x/net/context"
)

//CreateService will create a new app and return its id and warnings
func CreateService(spec dockerSwarm.ServiceSpec, options types.ServiceCreateOptions) (string, []string, error) {
	ctx := context.Background()

	response, err := swarmClient.ServiceCreate(ctx, spec, options)
	if err != nil {
		return "", []string{}, err
	}

	return response.ID, response.Warnings, nil
}

//RemoveService will remove a existed app
func RemoveService(serviceID string) error {
	ctx := context.Background()

	if err := swarmClient.ServiceRemove(ctx, serviceID); err != nil {
		return err
	}

	return nil
}

//UpdateService will update the config of a existed app and return warnings of updating opterations
func UpdateService(serviceID string, version dockerSwarm.Version, spec dockerSwarm.ServiceSpec, options types.ServiceUpdateOptions) ([]string, error) {
	ctx := context.Background()

	response, err := swarmClient.ServiceUpdate(ctx, serviceID, version, spec, options)
	if err != nil {
		return []string{}, err
	}

	return response.Warnings, nil
}

//ListServices will list all services of the swarm
func ListServices(options types.ServiceListOptions) ([]dockerSwarm.Service, error) {
	ctx := context.Background()

	services, err := swarmClient.ServiceList(ctx, options)
	if err != nil {
		return []dockerSwarm.Service{}, err
	}

	return services, nil
}
