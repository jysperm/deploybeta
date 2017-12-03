package models

import (
	"encoding/json"
	"fmt"

	"github.com/jysperm/deploying/config"

	etcdv3 "github.com/coreos/etcd/clientv3"
	"github.com/jysperm/deploying/lib/etcd"
	"golang.org/x/net/context"
)

type Version struct {
	Tag      string `json:"tag"`
	Registry string `json:"registry"`
	Status   string `json:"status"`
}

func CreateVersion(app *Application, tag string) (*Version, error) {

	versionKey := fmt.Sprintf("/apps/%s/versions/%s", app.Name, tag)

	newVersion := new(Version)
	newVersion.Registry = config.DefaultRegistry
	newVersion.Tag = tag
	newVersion.Status = "building"

	jsonVersion, err := json.Marshal(newVersion)
	if err != nil {
		return nil, err
	}

	if _, err := etcd.Client.Put(context.Background(), versionKey, string(jsonVersion)); err != nil {
		return nil, err
	}

	return newVersion, nil
}

func DeleteVersionByTag(app *Application, tag string) error {
	versionKey := fmt.Sprintf("/apps/%s/versions/%s", app.Name, tag)

	if _, err := etcd.Client.Delete(context.Background(), versionKey); err != nil {
		return err
	}

	return nil
}

func FindVersionByTag(app *Application, tag string) (*Version, error) {
	versionKey := fmt.Sprintf("/apps/%s/versions/%s", app.Name, tag)

	res, err := etcd.Client.Get(context.Background(), versionKey)
	if err != nil {
		return nil, err
	}

	if len(res.Kvs[0].Value) == 0 {
		return nil, nil
	}

	var version Version
	if err := json.Unmarshal(res.Kvs[0].Value, &version); err != nil {
		return nil, err
	}

	return &version, nil

}

func ListVersions(app *Application) (*[]Version, error) {
	versionPrefix := fmt.Sprintf("/apps/%s/versions/", app.Name)

	res, err := etcd.Client.Get(context.Background(), versionPrefix, etcdv3.WithPrefix())
	if err != nil {
		return &[]Version{}, err
	}

	var versionArray []Version

	for _, v := range res.Kvs {
		var t Version
		_ = json.Unmarshal(v.Value, &t)
		versionArray = append(versionArray, t)
	}

	if len(versionArray) == 0 {
		return &[]Version{}, nil
	}

	return &versionArray, nil
}

func DeleteAllVersion(app *Application) error {
	versionPrefix := fmt.Sprintf("/apps/%s/versions/", app.Name)

	_, err := etcd.Client.Delete(context.Background(), versionPrefix, etcdv3.WithPrefix())
	if err != nil {
		return err
	}

	return nil
}

func (v *Version) UpdateStatus(app *Application, status string) error {
	versionKey := fmt.Sprintf("/apps/%s/versions/%s", app.Name, v.Tag)

	tran := etcd.NewTransaction()

	tran.WatchJSON(versionKey, &Version{})

	resp, err := tran.Execute(func(watchedKeys map[string]interface{}) error {
		version := *watchedKeys[versionKey].(*Version)

		version.Status = status

		tran.PutJSONOnSuccess(versionKey, version)

		return nil
	})

	if err != nil {
		return err
	}

	if resp.Succeeded == false {
		return ErrUpdateConflict
	}

	v.Status = status

	return nil
}
