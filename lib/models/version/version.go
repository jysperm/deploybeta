package version

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/docker/docker/api/types"

	"golang.org/x/net/context"

	"github.com/jysperm/deploying/lib/etcd"
	appModel "github.com/jysperm/deploying/lib/models/app"
	"github.com/jysperm/deploying/lib/services/"
	"github.com/jysperm/deploying/lib/services/builder"
)

type Version struct {
	Shasum string `json:"shasum"`
	Tag    string `json:"tag"`
}

func CreateVersion(app *appModel.Application) (*Version, error) {
	version := generateTag()
	nameVersion := fmt.Sprintf("%s:%s", app.Name, version)
	versionKey := fmt.Sprintf("/apps/%s/versions/%s", app.Name, version)

	buildOpts := types.ImageBuildOptions{
		Tags: []string{nameVersion},
	}
	shasum, err := builder.BuildImage(buildOpts, app.GitRepository)
	if err != nil {
		return nil, err
	}

	newVersion := Version{Shasum: shasum, Tag: version}
	jsonVersion, err := json.Marshal(newVersion)
	if err != nil {
		return nil, err
	}
	if _, err := services.EtcdClient.Put(context.Background(), versionKey, string(jsonVersion)); err != nil {
		return nil, err
	}

	return &Version{Shasum: shasum, Tag: version}, nil
}

func DeleteVersion(app appModel.Application, version string) error {
	versionKey := fmt.Sprintf("/apps/%s/versions/%s", app.Name, version)

	if _, err := etcd.Client.Delete(context.Background(), versionKey); err != nil {
		return err
	}

	return nil
}

func LookupVersion(app appModel.Application, version string) (*Version, error) {
	versionKey := fmt.Sprintf("/apps/%s/versions/%s", app.Name, version)

	resp, err := services.EtcdClient.Get(context.Background(), versionKey)
	if err != nil {
		return nil, err
	}

	var v Version
	if err := json.Unmarshal(resp.Kvs[0].Value, &v); err != nil {
		return nil, err
	}

	return &v, nil
}
func generateTag() string {
	now := time.Now()
	return fmt.Sprintf("%d%d%d-%d%d%d", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())
}
