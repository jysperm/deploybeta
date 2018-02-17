package models

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jysperm/deploying/lib/etcd"
	"github.com/jysperm/deploying/lib/utils"
)

var ErrInvalidDataSourceType = errors.New("invalid datasource type")

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

var availableTypes = []string{"mongodb", "redis", "mysql"}

func CreateDataSource(dataSource *DataSource) error {
	if !validName.MatchString(dataSource.Name) {
		return ErrInvalidName
	}

	if !utils.StringInSlice(dataSource.Type, availableTypes) {
		return ErrInvalidDataSourceType
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
