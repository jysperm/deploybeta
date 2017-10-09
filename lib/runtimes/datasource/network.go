package datasource

import (
	"errors"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"

	"golang.org/x/net/context"
)

var swarmClient *client.Client

func init() {
	var err error
	swarmClient, err = client.NewEnvClient()
	if err != nil {
		panic(err)
	}
}

func CreateOverlay(name string) (string, error) {
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
		Attachable:     false,
		IPAM:           &net,
	}

	res, err := swarmClient.NetworkCreate(context.Background(), name, options)
	if err != nil {
		return "", err
	}

	return res.ID, nil
}

func RemoveOverlay(name string) error {
	id, err := findByName(name)
	if err != nil {
		return err
	}
	if err == nil && id == "" {
		return errors.New("Not found overlay network")
	}
	if err := swarmClient.NetworkRemove(context.Background(), id); err != nil {
		return err
	}
	return nil
}

func ListOverlay() ([]types.NetworkResource, error) {
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

func findByName(name string) (string, error) {
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
