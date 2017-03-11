package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo"

	"github.com/docker/docker/api/types"

	appModel "github.com/jysperm/deploying/lib/models/app"
	"github.com/jysperm/deploying/lib/services"
	"github.com/jysperm/deploying/lib/services/builder"
	. "github.com/jysperm/deploying/web/handlers/helpers"

	"golang.org/x/net/context"
)

func generateVersion() string {
	now := time.Now()
	return fmt.Sprintf("%d%d%d-%d%d%d", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())
}

func CreateImage(ctx echo.Context) error {
	params := new(appModel.Application)
	if err := ctx.Bind(params); err != nil {
		return NewHTTPError(http.StatusBadRequest, err)
	}

	version := generateVersion()
	buildOpts := types.ImageBuildOptions{
		Tags: []string{version},
	}
	shasum, err := builder.BuildImage(buildOpts, params.GitRepository)
	if err != nil {
		return NewHTTPError(http.StatusInternalServerError, err)
	}

	imageKey := fmt.Sprintf("/apps/%s/versions/%s", params.Name, version)
	if _, err := services.EtcdClient.Put(context.Background(), imageKey, version); err != nil {
		return NewHTTPError(http.StatusInternalServerError, err)
	}
	return ctx.JSON(http.StatusCreated, NewImageResponse(shasum, version))
}
