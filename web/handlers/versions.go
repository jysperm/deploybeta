package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/coreos/etcd/clientv3"
	"github.com/labstack/echo"

	"github.com/jysperm/deploybeta/lib/builder"
	"github.com/jysperm/deploybeta/lib/etcd"
	"github.com/jysperm/deploybeta/lib/models"
	"github.com/jysperm/deploybeta/lib/swarm"
	. "github.com/jysperm/deploybeta/web/handlers/helpers"
)

func CreateVersion(ctx echo.Context) error {
	app := ctx.Get("app").(models.Application)

	params := map[string]string{}
	if err := ctx.Bind(&params); err != nil {
		return NewHTTPError(http.StatusBadRequest, err)
	}

	version, err := builder.BuildVersion(&app, params["gitTag"])
	if err != nil {
		return NewHTTPError(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusCreated, NewVersionResponse(version))
}

func DeployVersion(ctx echo.Context) error {
	app := ctx.Get("app").(models.Application)

	params := map[string]string{}
	if err := ctx.Bind(&params); err != nil {
		return NewHTTPError(http.StatusBadRequest, err)
	}

	version, err := models.FindVersionByTag(&app, params["tag"])
	if err != nil {
		return NewHTTPError(http.StatusBadRequest, err)
	}

	if version.Status != "success" {
		return NewHTTPError(http.StatusBadRequest, errors.New("Version hadn't been built or had failed building"))
	}

	if err := app.UpdateVersion(version.Tag); err != nil {
		return NewHTTPError(http.StatusInternalServerError, err)
	}

	if err := swarm.UpdateAppService(&app); err != nil {
		return NewHTTPError(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusOK, NewVersionResponse(&version))
}

func PushProgress(ctx echo.Context) error {
	app := ctx.Get("app").(models.Application)
	tag := ctx.Param("tag")
	finishied := false

	watchPrefix := fmt.Sprintf("/progress/%s/%s/", app.Name, tag)
	watcher := etcd.Client.Watch(context.Background(), watchPrefix, clientv3.WithPrefix())

	rw := ctx.Response().Writer
	flusher, ok := rw.(http.Flusher)
	if !ok {
		return NewHTTPError(http.StatusInternalServerError, errors.New("Streaming unsupported!"))
	}

	ctx.Response().Header().Set("Content-Type", "text/event-stream")
	ctx.Response().WriteHeader(http.StatusOK)

	resp, err := etcd.Client.Get(context.Background(), watchPrefix, clientv3.WithPrefix())
	if err != nil {
		return NewHTTPError(http.StatusInternalServerError, err)
	}
	for _, ev := range resp.Kvs {
		fmt.Fprintf(rw, "data: %s\n\n", string(ev.Value))
		flusher.Flush()
		if strings.Contains(string(ev.Value), "Deploybeta: Building finished.") {
			finishied = true
			break
		}
	}
	if !finishied {
		for w := range watcher {
			for _, ev := range w.Events {
				fmt.Fprintf(rw, "data: %s\n\n", string(ev.Kv.Value))
				flusher.Flush()
				if strings.Contains(string(ev.Kv.Value), "Deploybeta: Building finished.") {
					break
				}
				break
			}
		}
	}
	return ctx.NoContent(http.StatusOK)
}

// TODO:
func DeleteVersion(ctx echo.Context) error {
	return nil
}
