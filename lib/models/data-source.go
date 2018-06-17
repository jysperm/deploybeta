package models

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	etcdv3 "github.com/coreos/etcd/clientv3"
	"github.com/hashicorp/errwrap"

	"github.com/jysperm/deploybeta/config"
	"github.com/jysperm/deploybeta/lib/etcd"
	"github.com/jysperm/deploybeta/lib/utils"
)

var ErrInvalidDataSourceType = errors.New("invalid datasource type")

// Serialize to /data-sources/:name
type DataSource struct {
	Name      string `json:"name"`
	Owner     string `json:"owner"`
	Type      string `json:"type"`
	Instances int    `json:"instances"`

	// HTTP API token scoped to this dataSource
	AgentToken string `json:"agentToken"`
}

// Serialize to /data-sources/:name/nodes/:host
type DataSourceNode struct {
	// Reference to DataSource.Name
	DataSourceName string `json:"dataSourceName"`
	// Reported address and port, like `10.0.1.1:6380`
	Host string `json:"host"`
	// Reported Role, `master` or `slave`
	Role string `json:"role"`
	// Reported master host, like `10.0.1.1:6380`
	MasterHost string `json:"masterHost"`

	ExpectedRole string `json:"expectedRole"`
}

type DataSourceNodeCommand struct {
	Command string `json:"command"`
	Role    string `json:"role"`
}

var availableTypes = []string{"mongodb", "redis"}

func CreateDataSource(dataSource *DataSource) error {
	if !validName.MatchString(dataSource.Name) {
		return ErrInvalidName
	}

	if !utils.StringInSlice(availableTypes, dataSource.Type) {
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

func (dataSource *DataSource) UpdateInstances(instances int) error {
	dataSourceKey := fmt.Sprintf("/data-sources/%s", dataSource.Name)

	tran := etcd.NewTransaction()

	tran.WatchJSON(dataSourceKey, &DataSource{}, func(watchedKey interface{}) error {
		ds := *watchedKey.(*DataSource)

		ds.Instances = instances

		tran.PutJSON(dataSourceKey, ds)

		return nil
	})

	resp, err := tran.Execute()

	if err != nil {
		return err
	}

	if resp.Succeeded == false {
		return ErrUpdateConflict
	}

	dataSource.Instances = instances

	return nil
}

func (dataSource *DataSource) LinkApp(app *Application) error {
	tran := etcd.NewTransaction()

	tran.AppendStringArray(fmt.Sprintf("/data-sources/%s/links", dataSource.Name), app.Name)
	tran.AppendStringArray(fmt.Sprintf("/apps/%s/data-sources", app.Name), dataSource.Name)

	err := tran.ExecuteMustSuccess()

	if err != nil {
		return err
	}

	return nil
}

func (dataSource *DataSource) UnlinkApp(app *Application) error {
	tran := etcd.NewTransaction()

	tran.PullStringArray(fmt.Sprintf("/data-sources/%s/links", dataSource.Name), app.Name)
	tran.PullStringArray(fmt.Sprintf("/apps/%s/data-sources", app.Name), dataSource.Name)

	err := tran.ExecuteMustSuccess()

	if err != nil {
		return err
	}

	return nil
}

func (dataSource *DataSource) GetLinkedAppNames() ([]string, error) {
	linkedApps := make([]string, 0)

	_, err := etcd.LoadKey(fmt.Sprintf("/data-sources/%s/links", dataSource.Name), linkedApps)

	if err != nil {
		return linkedApps, err
	}

	return linkedApps, nil
}

func (dataSource *DataSource) SwarmServiceName() string {
	return fmt.Sprintf("%s%s", config.DockerPrefix, dataSource.Name)
}

func (dataSource *DataSource) SwarmNetworkName() string {
	return fmt.Sprintf("%s%s", config.DockerPrefix, dataSource.Name)
}

func (dataSource *DataSource) SwarmInstances() int {
	return dataSource.Instances
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

func GetDataSourcesOfApp(app *Application) ([]DataSource, error) {
	dataSources := []DataSource{}

	dataSourceNames := make([]string, 0)
	_, err := etcd.LoadKey(fmt.Sprintf("/apps/%s/data-sources", app.Name), &dataSourceNames)

	if err != nil {
		return dataSources, err
	}

	for _, dataSourceName := range dataSourceNames {
		dataSource := DataSource{}
		found, err := etcd.LoadKey(fmt.Sprintf("/data-sources/%s", dataSourceName), &dataSource)

		if err != nil {
			return dataSources, err
		}

		if found {
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
	found, err := etcd.LoadKey(fmt.Sprintf("/data-sources/%s/nodes/%s", dataSource.Name, host), &dataSourceNode)

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

func (dataSource *DataSource) ListNodes() (nodes []DataSourceNode, err error) {
	resp, err := etcd.Client.Get(context.Background(), fmt.Sprintf("/data-sources/%s/nodes/", dataSource.Name), etcdv3.WithPrefix())

	if err != nil {
		return nodes, err
	}

	for _, v := range resp.Kvs {
		node := DataSourceNode{}

		err = json.Unmarshal(v.Value, &node)

		if err != nil {
			return nodes, err
		}

		nodes = append(nodes, node)
	}

	return nodes, nil
}

func (node *DataSourceNode) SetMaster() error {
	nodeKey := fmt.Sprintf("/data-sources/%s/nodes/%s", node.DataSourceName, node.Host)

	tran := etcd.NewTransaction()

	tran.WatchJSON(nodeKey, &DataSourceNode{}, func(watchedKey interface{}) error {
		node := *watchedKey.(*DataSourceNode)

		node.ExpectedRole = "master"

		tran.PutJSON(nodeKey, node)

		return nil
	})

	err := tran.ExecuteMustSuccess()

	if err != nil {
		return err
	}

	node.ExpectedRole = "master"

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

func (node *DataSourceNode) WaitForCommand() (*DataSourceNodeCommand, error) {
	nodeKey := fmt.Sprintf("/data-sources/%s/nodes/%s", node.DataSourceName, node.Host)

	checkNewCommand := func(node DataSourceNode) *DataSourceNodeCommand {
		if node.ExpectedRole != "" && node.Role != node.ExpectedRole {
			return &DataSourceNodeCommand{
				Command: "change-role",
				Role:    node.ExpectedRole,
			}
		} else {
			return nil
		}
	}

	watcher := etcd.Client.Watch(context.TODO(), nodeKey)

	latestNode := DataSourceNode{}

	found, err := etcd.LoadKey(nodeKey, &latestNode)

	if err != nil {
		return nil, err
	}

	if found {
		if command := checkNewCommand(latestNode); command != nil {
			return command, nil
		}
	}

	for w := range watcher {
		for _, ev := range w.Events {
			latestNode := DataSourceNode{}
			err = json.Unmarshal([]byte(ev.Kv.Value), &latestNode)

			if err != nil {
				log.Panicln(err)
			}

			if command := checkNewCommand(latestNode); command != nil {
				return command, nil
			}
		}
	}

	return nil, nil
}
