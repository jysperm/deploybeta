package swarm

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/docker/docker/api/types/swarm"
	"github.com/hashicorp/errwrap"

	"github.com/jysperm/deploybeta/config"
	"github.com/jysperm/deploybeta/lib/models"
	"github.com/jysperm/deploybeta/lib/runtimes"
)

var ErrNetworkNotFound = errors.New("Network not found")

type DataSourceNode struct {
	// NodeID in swarm cluster
	NodeID string
	// overlay address, like `10.0.1.1`
	Addr string
}

func UpdateDataSource(dataSource *models.DataSource) error {
	networkID, err := FindNetworkByName(dataSource.SwarmNetworkName())
	if err != nil {
		return err
	}
	if networkID == "" {
		networkID, err = CreateOverlay(dataSource)
		if err != nil {
			return err
		}
	}

	networkOpts := swarm.NetworkAttachmentConfig{
		Target: networkID,
	}

	runtime := runtimes.NewDataSourceRuntime(dataSource.Type)

	portConfig := swarm.PortConfig{
		Protocol:   runtime.ExposeProtocol(),
		TargetPort: uint32(runtime.ExposePort()),
	}

	environments := []string{
		"AGENT_TOKEN=" + dataSource.AgentToken,
		"DATASOURCE_NAME=" + dataSource.Name,
		"DEPLOYBETA_URL=http://" + config.HostPrivateAddress + config.Listen,
	}

	return UpdateService(dataSource, []swarm.PortConfig{portConfig}, []swarm.NetworkAttachmentConfig{networkOpts}, runtime.DockerImageName(), environments)
}

func RemoveDataSource(dataSource *models.DataSource) error {
	if err := RemoveService(dataSource); err != nil {
		return err
	}

	return RemoveOverlay(dataSource)
}

func ListDataSourceNodes(dataSource *models.DataSource) ([]DataSourceNode, error) {
	tasks, err := swarmClient.TaskList(context.Background(), getSerciceTasksFilter(dataSource))

	if err != nil {
		return nil, err
	}

	nodes := make([]DataSourceNode, 0)

	for _, task := range tasks {
		addr, err := findOverlayAddress(dataSource, &task)

		if err != nil {
			fmt.Println(errwrap.Wrapf("list dataSource nodes: {{err}}", err))
			continue
		}

		node := DataSourceNode{
			NodeID: task.NodeID,
			Addr:   addr,
		}

		nodes = append(nodes, node)
	}

	return nodes, nil
}

func findOverlayAddress(dataSource *models.DataSource, task *swarm.Task) (string, error) {
	for _, attachment := range task.NetworksAttachments {
		if attachment.Network.Spec.Name == dataSource.SwarmNetworkName() && len(attachment.Addresses) > 0 {
			addr, _, err := net.ParseCIDR(attachment.Addresses[0])

			if err != nil {
				return "", err
			}

			return addr.String(), nil
		}
	}

	return "", fmt.Errorf("network `%s` not found", dataSource.SwarmNetworkName())
}
