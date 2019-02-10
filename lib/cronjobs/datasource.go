package cronjobs

import (
	"fmt"
	"log"

	"github.com/hashicorp/errwrap"

	"github.com/jysperm/deploybeta/lib/datasources"
	"github.com/jysperm/deploybeta/lib/db"
	"github.com/jysperm/deploybeta/lib/models"
)

func checkDataSourceNodes() {
	dataSources := make([]models.DataSource, 0)

	err := db.FetchAllFrom("/data-sources", &dataSources)

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

		runtime := datasources.NewDataSourceRuntime(dataSource.Type)

		for _, node := range nodes {
			err := runtime.CheckNodeAvailability(node.Host)

			if err != nil {
				log.Println(errwrap.Wrapf("check datasource nodes: {{err}}", err))
			}
		}
	}
}
