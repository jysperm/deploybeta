package swarm

import (
	"errors"

	"github.com/jysperm/deploying/lib/models"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"

	"golang.org/x/net/context"
)

func CreateOverlay(datasource *models.DataSource) (string, error) {
	net := network.IPAM{
		Driver:  "default",
		Options: map[string]string{},
		Config:  []network.IPAMConfig{},
	}
	options := types.NetworkCreate{
		CheckDuplicate: true,
		Driver:         "overlay",
		EnableIPv6:     false,
		Internal:       false,
		Attachable:     true,
		IPAM:           &net,
	}

	res, err := swarmClient.NetworkCreate(context.Background(), datasource.SwarmNetworkName(), options)
	if err != nil {
		return "", err
	}

	return res.ID, nil
}

func RemoveOverlay(datasource *models.DataSource) error {
	id, err := FindNetworkByName(datasource.SwarmNetworkName())
	if err != nil {
		return err
	}
	if err == nil && id == "" {
		return errors.New("Not found overlay network")
	}
	return swarmClient.NetworkRemove(context.Background(), id)
}

func ListOverlays() ([]types.NetworkResource, error) {
	filter := filters.NewArgs()
	filter.Add("driver", "overlay")
	options := types.NetworkListOptions{
		Filters: filter,
	}
	list, err := swarmClient.NetworkList(context.Background(), options)
	if err != nil {
		return []types.NetworkResource{}, err
	}
	return list, nil
}

func FindNetworkByName(name string) (string, error) {
	filter := filters.NewArgs()
	filter.Add("driver", "overlay")
	options := types.NetworkListOptions{
		Filters: filter,
	}
	list, err := swarmClient.NetworkList(context.Background(), options)
	if err != nil {
		return "", err
	}
	for _, i := range list {
		if i.Name == name {
			return i.ID, nil
		}
	}

	return "", nil
}
