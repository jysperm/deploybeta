package version

import (
	"fmt"
	"time"

	"github.com/docker/docker/api/types"

	"golang.org/x/net/context"

	appModel "github.com/jysperm/deploying/lib/models/app"
	"github.com/jysperm/deploying/lib/services"
	"github.com/jysperm/deploying/lib/services/builder"
)

type Version struct {
	Shasum string `json:"shasum"`
	Tag    string `json:"tag"`
}

func CreateVersion(app *appModel.Application) (Version, error) {
	tag := generateTag()
	tagKey := fmt.Sprintf("/apps/%s/versions/%s", app.Name, tag)

	buildOpts := types.ImageBuildOptions{
		Tags: []string{tag},
	}
	shasum, err := builder.BuildImage(buildOpts, app.GitRepository)
	if err != nil {
		return Version{}, err
	}

	if _, err := services.EtcdClient.Put(context.Background(), tagKey, shasum); err != nil {
		return Version{}, err
	}

	return Version{Shasum: shasum, Tag: tag}, nil
}

func DeleteVersion(app appModel.Application, version string) error {
	versionKey := fmt.Sprintf("/apps/%s/versions/%s", app.Name, version)

	if _, err := services.EtcdClient.Delete(context.Background(), versionKey); err != nil {
		return err
	}

	return nil
}

func generateTag() string {
	now := time.Now()
	return fmt.Sprintf("%d%d%d-%d%d%d", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())
}
