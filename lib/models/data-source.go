package models

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hashicorp/errwrap"

	"github.com/jysperm/deploying/config"
	"github.com/jysperm/deploying/lib/etcd"
	"github.com/jysperm/deploying/lib/utils"
)

var ErrInvalidDataSourceType = errors.New("invalid datasource type")

// Serialize to /data-source/:name
type DataSource struct {
	Name      string `json:"name"`
	Owner     string `json:"owner"`
	Type      string `json:"type"`
	Instances int    `json:"instances"`

	// HTTP API token scoped to this dataSource
	AgentToken string `json:"agentToken"`
}

// Serialize to /data-source/:name/nodes/:host
type DataSourceNode struct {
	// Reference to DataSource.Name
	DataSourceName string `json:"dataSourceName"`
	// Reported address and port, like `10.0.1.1:6380`
	Host string `json:"host"`
	// Reported Role, `master` or `slave`
	Role string `json:"role"`
	// Reported master host, like `10.0.1.1:6380`
	MasterHost string `json:"masterHost"`
}

var availableTypes = []string{"mongodb", "redis"}

func CreateDataSource(dataSource *DataSource) error {
	if !validName.MatchString(dataSource.Name) {
		return ErrInvalidName
	}

	if !utils.StringInSlice(dataSource.Type, availableTypes) {
		return ErrInvalidDataSourceType
	}

	dataSource.AgentToken = utils.RandomString(32)

	tran := etcd.NewTransaction()

	tran.AppendStringArray(fmt.Sprintf("/accounts/%s/data-sources", dataSource.Owner), dataSource.Name)
	tran.CreateJSON(fmt.Sprint("/data-sources/", dataSource.Name), dataSource)

	resp, err := tran.Execute()

	if err != nil {
		return err
	}

	if resp.Succeeded == false {
		return ErrUpdateConflict
	}

	return nil
}

func (datasource *DataSource) UpdateInstances(instances int) error {
	datasourceKey := fmt.Sprintf("/data-source/%s", datasource.Name)

	tran := etcd.NewTransaction()

	tran.WatchJSON(datasourceKey, &DataSource{}, func(watchedKey interface{}) error {
		ds := *watchedKey.(*DataSource)

		ds.Instances = instances

		tran.PutJSON(datasourceKey, ds)

		return nil
	})

	resp, err := tran.Execute()

	if err != nil {
		return err
	}

	if resp.Succeeded == false {
		return ErrUpdateConflict
	}

	datasource.Instances = instances

	return nil
}

func (datasource *DataSource) SwarmServiceName() string {
	return fmt.Sprintf("%s%s", config.DockerPrefix, datasource.Name)
}

func (datasource *DataSource) SwarmNetworkName() string {
	return fmt.Sprintf("%s%s", config.DockerPrefix, datasource.Name)
}

func LinkDataSource(dataSource *DataSource, app *Application) error {
	linksKey := fmt.Sprintf("/data-source/%s/links", dataSource.Name)

	tran := etcd.NewTransaction()

	tran.WatchJSON(linksKey, &[]string{}, func(watchedKey interface{}) error {
		apps := *watchedKey.(*[]string)

		for _, v := range apps {
			if v == app.Name {
				return errors.New("DataSource has been attached")
			}
		}

		apps = append(apps, app.Name)

		tran.PutJSON(linksKey, apps)

		return nil
	})

	resp, err := tran.Execute()

	if err != nil {
		return err
	}

	if resp.Succeeded == false {
		return ErrUpdateConflict
	}

	if err := UpdateDataSourceLinks(dataSource, app); err != nil {
		return err
	}

	return nil
}

func UnlinkDataSource(dataSource *DataSource, app *Application) error {
	linksKey := fmt.Sprintf("/data-source/%s/links", dataSource.Name)

	tran := etcd.NewTransaction()

	tran.WatchJSON(linksKey, &[]string{}, func(watchedKey interface{}) error {
		links := *watchedKey.(*[]string)

		for i := 0; i < len(links); i++ {
			if links[i] == app.Name {
				links = append(links[:i], links[i+1:]...)
				tran.PutJSON(linksKey, links)
				return nil
			}
		}

		return errors.New("Not found link")
	})

	resp, err := tran.Execute()

	if err != nil {
		return err
	}

	if resp.Succeeded == false {
		return ErrUpdateConflict
	}

	if err := UpdateDataSourceLinks(dataSource, app); err != nil {
		return err
	}

	return nil
}

func UpdateDataSourceLinks(dataSource *DataSource, app *Application) error {
	appLinksKey := fmt.Sprintf("/app/%s/data-sources", app.Name)

	tran := etcd.NewTransaction()

	tran.WatchJSON(appLinksKey, &[]string{}, func(watchedKey interface{}) error {
		appLinks := *watchedKey.(*[]string)

		for i := 0; i < len(appLinks); i++ {
			if appLinks[i] == dataSource.Name {
				appLinks = append(appLinks[:i], appLinks[i+1:]...)
				tran.PutJSON(appLinksKey, appLinks)
				return nil
			}
		}

		appLinks = append(appLinks, dataSource.Name)

		return nil
	})

	resp, err := tran.Execute()

	if err != nil {
		return err
	}

	if resp.Succeeded == false {
		return ErrUpdateConflict
	}

	return nil
}

func GetDataSourcesOfAccount(account *Account) (dataSources []DataSource, err error) {
	dataSources = make([]DataSource, 0)

	resp, err := etcd.Client.Get(context.Background(), fmt.Sprintf("/accounts/%s/data-sources", account.Username))

	if err != nil {
		return nil, err
	}

	if len(resp.Kvs) == 0 {
		return dataSources, nil
	}

	accountDataSources := []string{}

	err = json.Unmarshal([]byte(resp.Kvs[0].Value), &accountDataSources)

	if err != nil {
		return nil, err
	}

	for _, dataSourceName := range accountDataSources {
		resp, err = etcd.Client.Get(context.Background(), fmt.Sprint("/data-sources/", dataSourceName))

		if err != nil {
			return nil, err
		}

		dataSource := DataSource{}

		if len(resp.Kvs) != 0 {
			err = json.Unmarshal([]byte(resp.Kvs[0].Value), &dataSource)

			if err != nil {
				return nil, err
			}

			dataSources = append(dataSources, dataSource)
		}
	}

	return dataSources, nil
}

func GetDataSourceOfAccount(dataSourceName string, account *Account) (*DataSource, error) {
	resp, err := etcd.Client.Get(context.Background(), fmt.Sprint("/data-sources/", dataSourceName))

	if err != nil {
		return nil, err
	}

	dataSource := &DataSource{}

	if len(resp.Kvs) != 0 {
		err = json.Unmarshal([]byte(resp.Kvs[0].Value), &dataSource)

		if err != nil {
			return nil, err
		}
	}

	return dataSource, nil
}

func FindDataSourceByName(name string) (dataSource DataSource, err error) {
	found, err := etcd.LoadKey(fmt.Sprintf("/data-sources/%s", name), &dataSource)

	if err != nil {
		return dataSource, err
	} else if !found {
		return dataSource, errors.New("dataSource not found")
	} else {
		return dataSource, nil
	}
}

func DeleteDataSourceByName(name string) error {
	_, err := etcd.Client.Delete(context.Background(), fmt.Sprint("/data-sources/", name))

	return err
}

func (dataSource *DataSource) FindNodeByHost(host string) (dataSourceNode DataSourceNode, err error) {
	found, err := etcd.LoadKey(fmt.Sprintf("/data-source/%s/nodes/%s", dataSource.Name, host), &dataSourceNode)

	if err != nil {
		return dataSourceNode, err
	} else if !found {
		return dataSourceNode, errors.New("dataSource node not found")
	} else {
		return dataSourceNode, nil
	}
}

func (dataSource *DataSource) CreateNode(node *DataSourceNode) error {
	node.DataSourceName = dataSource.Name

	tran := etcd.NewTransaction()

	tran.CreateJSON(fmt.Sprintf("/data-sources/%s/nodes/%s", dataSource.Name, node.Host), node)

	resp, err := tran.Execute()

	if err != nil {
		return err
	}

	if resp.Succeeded == false {
		return errwrap.Wrapf("create dataSource node: {{err}}", ErrUpdateConflict)
	}

	return nil
}

func (node *DataSourceNode) Update(updates *DataSourceNode) error {
	nodeKey := fmt.Sprintf("/data-sources/%s/nodes/%s", node.DataSourceName, node.Host)

	tran := etcd.NewTransaction()

	tran.WatchJSON(nodeKey, &DataSourceNode{}, func(watchedKey interface{}) error {
		node := *watchedKey.(*DataSourceNode)

		node.Role = updates.Role
		node.MasterHost = updates.MasterHost

		tran.PutJSON(nodeKey, node)

		return nil
	})

	resp, err := tran.Execute()

	if err != nil {
		return err
	}

	if resp.Succeeded == false {
		return errwrap.Wrapf("update dataSource node: {{err}}", ErrUpdateConflict)
	}

	node.Role = updates.Role
	node.MasterHost = updates.MasterHost

	return nil
}
