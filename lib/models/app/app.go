package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"

	etcd "github.com/coreos/etcd/clientv3"
	etcdpb "github.com/coreos/etcd/mvcc/mvccpb"
	accountModel "github.com/jysperm/deploying/lib/models/account"
	"github.com/jysperm/deploying/lib/services"
	"golang.org/x/net/context"
)

var ErrInvalidName = errors.New("invalid app name")
var ErrUpdateConflict = errors.New("update version conflict")

type Application struct {
	Name          string `json:"name"`
	Owner         string `json:"owner"`
	GitRepository string `json:"gitRepository"`
	Instances     int    `json:"instances"`
	Version       string `json:"version"`
}

var validName = regexp.MustCompile(`^[A-Za-z0-9_-]+$`)

func CreateApp(app *Application) error {
	if !validName.MatchString(app.Name) {
		return ErrInvalidName
	}

	appKey := fmt.Sprint("/apps/", app.Name)
	accountAppsKey := fmt.Sprintf("/account/%s/apps", app.Owner)

	appBytes, err := json.Marshal(app)

	if err != nil {
		return err
	}

	resp, err := services.EtcdClient.Get(context.Background(), accountAppsKey)

	if err != nil {
		return err
	}

	var accountAppsKeyValue *etcdpb.KeyValue
	var accountApps []string

	compares := []etcd.Cmp{
		etcd.Compare(etcd.CreateRevision(appKey), "=", 0),
	}

	ops := []etcd.Op{
		etcd.OpPut(appKey, string(appBytes)),
	}

	if len(resp.Kvs) > 0 {
		accountAppsKeyValue = resp.Kvs[0]
		err = json.Unmarshal([]byte(resp.Kvs[0].Value), &accountApps)

		if err != nil {
			return err
		}

		compares = append(compares, etcd.Compare(etcd.Version(accountAppsKey), "=", accountAppsKeyValue.Version))

		accountAppsBytes, err := json.Marshal(append(accountApps, app.Name))

		if err != nil {
			return err
		}

		ops = append(ops, etcd.OpPut(accountAppsKey, string(accountAppsBytes)))
	} else {
		compares = append(compares, etcd.Compare(etcd.CreateRevision(appKey), "=", 0))

		accountAppsBytes, err := json.Marshal([]string{app.Name})

		if err != nil {
			return err
		}

		ops = append(ops, etcd.OpPut(accountAppsKey, string(accountAppsBytes)))
	}

	txnResp, err := services.EtcdClient.Txn(context.Background()).If(compares...).Then(ops...).Commit()

	if err != nil {
		return err
	}

	if txnResp.Succeeded == false {
		return ErrUpdateConflict
	}

	return nil
}

// TODO: Delete app name from `/account/:name/apps`
func DeleteByName(name string) error {
	appKey := fmt.Sprint("/apps/", name)

	_, err := services.EtcdClient.Delete(context.Background(), appKey)

	return err
}

func GetAppsOfAccount(account *accountModel.Account) (result []Application, err error) {
	accountAppsKey := fmt.Sprintf("/account/%s/apps", account.Username)
	resp, err := services.EtcdClient.Get(context.Background(), accountAppsKey)

	if err != nil {
		return nil, err
	}

	if len(resp.Kvs) == 0 {
		return result, nil
	}

	accountApps := []string{}
	err = json.Unmarshal([]byte(resp.Kvs[0].Value), &accountApps)

	if err != nil {
		return nil, err
	}

	for _, appName := range accountApps {
		appKey := fmt.Sprint("/apps/", appName)
		resp, err = services.EtcdClient.Get(context.Background(), appKey)

		if err != nil {
			return result, err
		}

		app := Application{}

		if len(resp.Kvs) != 0 {
			err = json.Unmarshal([]byte(resp.Kvs[0].Value), &app)

			if err != nil {
				return result, err
			}

			result = append(result, app)
		}
	}

	return result, nil
}

func (app *Application) UpdateGitRepository(GgtRepository string) error {
	return nil
}

func (app *Application) UpdateInstances(instances int) error {
	return nil
}

func (app *Application) UpdateVersion(version string) error {
	return nil
}

func FindByName(name string) (Application, error) {
	appKey := fmt.Sprintf("/apps/%s", name)
	resp, err := services.EtcdClient.Get(context.Background(), appKey)
	if err != nil {
		return Application{}, err
	}
	var appFound Application
	if err := json.Unmarshal(resp.Kvs[0].Value, &appFound); err != nil {
		return Application{}, err
	}

	return appFound, nil
}
