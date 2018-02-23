package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jysperm/deploying/config"

	etcdv3 "github.com/coreos/etcd/clientv3"
	"github.com/jysperm/deploying/lib/etcd"
	"golang.org/x/net/context"
)

// Serialize to /apps/:appName/versions/:tag
type Version struct {
	AppName  string `json:"appName"`
	Tag      string `json:"tag"`
	Registry string `json:"registry"`
	Status   string `json:"status"`
}

func NewVersion(app *Application) Version {
	return Version{
		AppName:  app.Name,
		Tag:      newVersionTag(),
		Registry: config.DefaultRegistry,
		Status:   "building",
	}
}

func (version *Version) Save() error {
	tran := etcd.NewTransaction()
	tran.CreateJSON(version.etcdKey(), version)
	return tran.ExecuteMustSuccess()
}

func DeleteVersionByTag(app *Application, tag string) error {
	versionKey := fmt.Sprintf("/apps/%s/versions/%s", app.Name, tag)

	if _, err := etcd.Client.Delete(context.Background(), versionKey); err != nil {
		return err
	}

	return nil
}

func FindVersionByTag(app *Application, tag string) (version Version, err error) {
	found, err := etcd.LoadKey(fmt.Sprintf("/apps/%s/versions/%s", app.Name, tag), &version)

	if err != nil {
		return version, err
	} else if !found {
		return version, errors.New("version not found")
	} else {
		return version, nil
	}
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

	tran.WatchJSON(versionKey, &Version{}, func(watchedKey interface{}) error {
		version := *watchedKey.(*Version)

		version.Status = status

		tran.PutJSON(versionKey, version)

		return nil
	})

	resp, err := tran.Execute()

	if err != nil {
		return err
	}

	if resp.Succeeded == false {
		return ErrUpdateConflict
	}

	v.Status = status

	return nil
}

func (version *Version) DockerImageName() string {
	return fmt.Sprintf("%s/%s%s:%s", version.Registry, config.DockerPrefix, version.AppName, version.Tag)
}

func (version *Version) etcdKey() string {
	return fmt.Sprintf("/apps/%s/versions/%s", version.AppName, version.Tag)
}

func newVersionTag() string {
	now := time.Now()
	return fmt.Sprintf("%d%d%d-%d%d%d", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())
}
