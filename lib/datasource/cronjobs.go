package datasources

import (
	"fmt"
	"strings"
	"time"

	"github.com/jysperm/deploybeta/lib/db"
	"github.com/jysperm/deploybeta/lib/models"
	"github.com/jysperm/deploybeta/lib/swarm"
)

func init() {
	ticker := time.NewTicker(time.Minute)

	go func() {
		for {
			select {
			case <-ticker.C:
				checkDataSourceNodes()
			}
		}
	}()
}

func checkDataSourceNodes() {
	dataSources := make([]models.DataSource, 0)

	err := db.FetchAllFrom("/data-sources", &dataSources, 2)

	if err != nil {
		fmt.Println(err)
		return
	}

	for _, dataSource := range dataSources {
		nodes := make([]models.DataSourceNode, 0)
		err := dataSource.Nodes().FetchAll(&nodes)

		if err != nil {
			fmt.Println(err)
			return
		}

		swarmNodes, err := swarm.ListDataSourceNodes(&dataSource)

		if err != nil {
			fmt.Println(err)
			return
		}

		for _, node := range nodes {
			if !nodeExistsInSwarm(swarmNodes, &node) {
				fmt.Printf("%+v need to be remote from etcd\n", node)
			}
		}
	}
}

func nodeExistsInSwarm(swarmNodes []swarm.DataSourceNode, node *models.DataSourceNode) bool {
	for _, swarmNode := range swarmNodes {
		if strings.HasPrefix(node.Host, swarmNode.Addr) {
			return true
		}
	}

	return false
}
