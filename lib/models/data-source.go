package models

import (
	"errors"
	"fmt"

	"github.com/hashicorp/errwrap"

	"github.com/jysperm/deploybeta/config"
	"github.com/jysperm/deploybeta/lib/db"
	"github.com/jysperm/deploybeta/lib/utils"
)

const (
	ROLE_MASTER  = "master"
	ROLE_SLAVE   = "slave"
	ROLE_UNKNOWN = ""

	COMMAND_CHANGE_ROLE   = "change-role"
	COMMAND_REPORT_STATUS = "report-status"
)

var ErrInvalidDataSourceType = errors.New("invalid datasource type")
var ErrDataSourceNotFound = errors.New("dataSource not found")
var ErrDataSourceNodeNotFound = errors.New("dataSource node not found")

type DataSource struct {
	db.ResourceMeta

	Name          string `json:"name"`
	OwnerUsername string `json:"ownerUsername"`
	Type          string `json:"type"`
	Instances     int    `json:"instances"`

	// Current master node
	MasterNodeHost string `json:"masterNodeHost"`
	// HTTP API token scoped to this dataSource
	AgentToken string `json:"agentToken"`
}

func (dataSource *DataSource) ResourceKey() string {
	return fmt.Sprintf("/data-sources/%s", dataSource.Name)
}

func (dataSource *DataSource) Associations() []db.Association {
	return []db.Association{
		dataSource.Nodes(),
		dataSource.Apps(),
		dataSource.Owner(),
		dataSource.MasterNode(),
	}
}

func (dataSource *DataSource) MasterNode() db.HasOneAssociation {
	return db.HasOne(fmt.Sprintf("/data-sources/%s/nodes/%s", dataSource.Name, dataSource.MasterNodeHost))
}

func (dataSource *DataSource) Nodes() db.HasManyAssociation {
	return db.HasManyPrefix(fmt.Sprintf("/data-sources/%s/nodes/", dataSource.Name))
}

func (dataSource *DataSource) Apps() db.HasManyAssociation {
	return db.HasManyThrough(fmt.Sprintf("/data-sources/%s/apps", dataSource.Name))
}

func (dataSource *DataSource) Owner() db.BelongsToAssociation {
	return db.BelongsTo(
		(&Account{Username: dataSource.OwnerUsername}).ResourceKey(),
		fmt.Sprintf("/accounts/%s/data-sources", dataSource.OwnerUsername),
	)
}

type DataSourceNode struct {
	db.ResourceMeta

	DataSourceName string `json:"dataSourceName"`

	// Reported address and port, like `10.0.1.1:6380`
	Host string `json:"host"`
	// Reported Role, `master` or `slave`
	Role string `json:"role"`
	// Reported master host, like `10.0.1.1:6380`
	MasterHost string `json:"masterHost"`
}

func (node *DataSourceNode) ResourceKey() string {
	return fmt.Sprintf("/data-sources/%s/nodes/%s", node.DataSourceName, node.Host)
}

func (node *DataSourceNode) Associations() []db.Association {
	return []db.Association{
		node.DataSource(),
	}
}

func (node *DataSourceNode) DataSource() db.BelongsToAssociation {
	return db.BelongsTo((&DataSource{Name: node.DataSourceName}).ResourceKey())
}

type DataSourceNodeCommand struct {
	Command    string `json:"command"`
	Role       string `json:"role"`
	MasterHost string `json:"masterHost"`
}

var availableTypes = []string{"mongodb", "mysql", "redis"}

func CreateDataSource(dataSource *DataSource) error {
	if !validName.MatchString(dataSource.Name) {
		return ErrInvalidName
	}

	if !utils.StringInSlice(availableTypes, dataSource.Type) {
		return ErrInvalidDataSourceType
	}

	dataSource.AgentToken = utils.RandomString(32)

	_, err := db.StartTransaction(func(tran db.Transaction) {
		tran.Create(dataSource)
	})

	return err
}

func FindDataSourceByName(name string) (*DataSource, error) {
	dataSource := &DataSource{
		Name: name,
	}

	err := db.Fetch(dataSource)

	if err == db.ErrResourceNotFound {
		return nil, errwrap.Wrap(ErrDataSourceNotFound, err)
	}

	return dataSource, err
}

func GetDataSourceOfAccount(dataSourceName string, account *Account) (*DataSource, error) {
	dataSource, err := FindDataSourceByName(dataSourceName)

	if err != nil {
		return nil, err
	}

	if dataSource.OwnerUsername != account.Username {
		return nil, ErrDataSourceNotFound
	}

	return dataSource, nil
}

func (dataSource *DataSource) UpdateInstances(instances int) error {
	_, err := db.StartTransaction(func(tran db.Transaction) {
		err := db.Fetch(dataSource)

		if err != nil {
			tran.SetError(err)
			return
		}

		dataSource.Instances = instances

		tran.Update(dataSource)
	})

	return err
}

func (dataSource *DataSource) LinkApp(app *Application) error {
	_, err := db.StartTransaction(func(tran db.Transaction) {
		app.DataSources().Attach(tran, dataSource)
		dataSource.Apps().Attach(tran, app)
	})

	return err
}

func (dataSource *DataSource) UnlinkApp(app *Application) error {
	_, err := db.StartTransaction(func(tran db.Transaction) {
		app.DataSources().Detach(tran, dataSource)
		dataSource.Apps().Detach(tran, app)
	})

	return err
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

func (dataSource *DataSource) Destroy() error {
	_, err := db.StartTransaction(func(tran db.Transaction) {
		tran.Delete(dataSource)
	})

	return err
}

func (dataSource *DataSource) FindNodeByHost(host string) (*DataSourceNode, error) {
	node := &DataSourceNode{
		DataSourceName: dataSource.Name,
		Host:           host,
	}

	err := db.Fetch(node)

	if err == db.ErrResourceNotFound {
		return nil, errwrap.Wrap(ErrDataSourceNodeNotFound, err)
	}

	return node, err
}

func (dataSource *DataSource) CreateNode(node *DataSourceNode) (*DataSourceNodeCommand, error) {
	node.DataSourceName = dataSource.Name

	command := &DataSourceNodeCommand{
		Command: COMMAND_CHANGE_ROLE,
	}

	_, err := db.StartTransaction(func(tran db.Transaction) {
		err := db.Fetch(dataSource)

		if err != nil {
			tran.SetError(err)
			return
		}

		if dataSource.MasterNodeHost == ROLE_UNKNOWN {
			dataSource.MasterNodeHost = node.Host
			command.Role = ROLE_MASTER
		} else {
			command.Role = ROLE_SLAVE
			command.MasterHost = dataSource.MasterNodeHost
		}

		tran.Update(dataSource)
		tran.Create(node)
	})

	if err != nil {
		return nil, errwrap.Wrapf("create dataSource node: {{err}}", err)
	}

	return command, nil
}

func (node *DataSourceNode) SetMaster() error {
	dataSource := &DataSource{}

	_, err := db.StartTransaction(func(tran db.Transaction) {
		err := node.DataSource().Fetch(dataSource)

		if err != nil {
			tran.SetError(err)
			return
		}

		dataSource.MasterNodeHost = node.Host

		tran.Update(dataSource)
	})

	return err
}

func (node *DataSourceNode) Update(updates *DataSourceNode) error {
	_, err := db.StartTransaction(func(tran db.Transaction) {
		err := db.Fetch(node)

		if err != nil {
			tran.SetError(err)
			return
		}

		node.Role = updates.Role
		node.MasterHost = updates.MasterHost

		tran.Update(node)
	})

	if err != nil {
		return errwrap.Wrapf("update dataSource node: {{err}}", err)
	}

	return nil
}

func (node *DataSourceNode) Destroy() error {
	_, err := db.StartTransaction(func(tran db.Transaction) {
		tran.Delete(node)
	})

	return err
}

func (node *DataSourceNode) WaitForCommand() (*DataSourceNodeCommand, error) {
	dataSource := &DataSource{}
	err := node.DataSource().Fetch(dataSource)

	if err != nil {
		return nil, err
	}

	checkNewCommand := func() *DataSourceNodeCommand {
		if node.Role == ROLE_UNKNOWN {
			return &DataSourceNodeCommand{
				Command: COMMAND_REPORT_STATUS,
			}
		} else if dataSource.MasterNodeHost == node.Host && node.Role != ROLE_MASTER {
			return &DataSourceNodeCommand{
				Command: COMMAND_CHANGE_ROLE,
				Role:    ROLE_MASTER,
			}
		} else if dataSource.MasterNodeHost != node.Host && node.Role == ROLE_MASTER {
			return &DataSourceNodeCommand{
				Command:    COMMAND_CHANGE_ROLE,
				Role:       ROLE_SLAVE,
				MasterHost: dataSource.MasterNodeHost,
			}
		} else {
			return nil
		}
	}

	if command := checkNewCommand(); command != nil {
		return command, nil
	}

	cancelWatchDataSource, dataSourceUpdates, dataSourceErrs := db.WatchUpdates(dataSource)
	cancelWatchNode, nodeUpdates, nodeErrs := db.WatchUpdates(node)

	defer cancelWatchDataSource()
	defer cancelWatchNode()

	for {
		select {
		case updated := <-dataSourceUpdates:
			db.Assign(dataSource, updated)

			if command := checkNewCommand(); command != nil {
				return command, nil
			}
		case updated := <-nodeUpdates:
			db.Assign(node, updated)

			if command := checkNewCommand(); command != nil {
				return command, nil
			}
		case err := <-dataSourceErrs:
			return nil, err
		case err := <-nodeErrs:
			return nil, err
		}
	}
}
