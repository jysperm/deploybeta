package models

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jysperm/deploying/lib/etcd"
)

type DataSource struct {
	Name      string `json:"name"`
	Owner     string `json:"owner"`
	Type      string `json:"type"`
	Instances int    `json:"instances"`
}

type DataSourceNode struct {
	Host string `json:"host"`
	Role string `json:"role"`
}

func CreateDataSource(dataSource *DataSource) error {
	if !validName.MatchString(dataSource.Name) {
		return ErrInvalidName
	}

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

func DeleteDataSourceByName(name string) error {
	_, err := etcd.Client.Delete(context.Background(), fmt.Sprint("/data-sources/", name))

	return err
}

func CreateDataSourceNode(dataSource *DataSource, dataSourceNode *DataSourceNode) error {
	tran := etcd.NewTransaction()

	tran.CreateJSON(fmt.Sprintf("/data-sources/%s/nodes/%s", dataSource.Name, dataSourceNode.Host), dataSourceNode)

	resp, err := tran.Execute()

	if err != nil {
		return err
	}

	if resp.Succeeded == false {
		return ErrUpdateConflict
	}

	return nil
}
