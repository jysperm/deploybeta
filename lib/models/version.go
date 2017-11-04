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
	Shasum   string `json:"shasum"`
	Tag      string `json:"tag"`
	Registry string `json:"registry"`
}

func CreateVersion(app *Application, gitTag string, tag string, shasum string) (*Version, error) {

	versionKey := fmt.Sprintf("/apps/%s/versions/%s", app.Name, tag)

	newVersion := new(Version)
	newVersion.Registry = config.DefaultRegistry
	newVersion.Shasum = shasum
	newVersion.Tag = tag

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

	version := new(Version)
	if err := json.Unmarshal(res.Kvs[0].Value, version); err != nil {
		return nil, err
	}

	return version, nil

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
