package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"

	etcd "github.com/coreos/etcd/clientv3"
	etcdpb "github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/jysperm/deploying/lib/services"
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

func (app *Application) UpdateGitRepository(GgtRepository string) error {
	return nil
}

func (app *Application) UpdateInstances(instances int) error {
	return nil
}

func (app *Application) UpdateVersion(version string) error {
	return nil
}
