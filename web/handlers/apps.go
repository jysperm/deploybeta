package handlers

import (
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/buger/jsonparser"
	"github.com/hashicorp/errwrap"
	"github.com/labstack/echo"

	"github.com/jysperm/deploybeta/lib/models"
	"github.com/jysperm/deploybeta/lib/swarm"
	. "github.com/jysperm/deploybeta/web/handlers/helpers"
)

func GetMyApps(ctx echo.Context) error {
	account := GetSessionAccount(ctx)

	apps := make([]models.Application, 0)
	err := account.Apps().FetchAll(&apps)

	if err != nil {
		return NewHTTPError(http.StatusInternalServerError, err)
	}

	responses := make([]ApplicationResponse, len(apps))

	for i, app := range apps {
		versions := make([]models.Version, 0)
		err := app.Versions().FetchAll(&versions)

		if err != nil {
			return err
		}

		upstreams := make([]models.Upstream, 0)
		err = app.Upstreams().FetchAll(&upstreams)

		if err != nil {
			return err
		}

		nodes, err := swarm.ListNodes(&app)

		if err != nil {
			log.Println(errwrap.Wrapf("list swarm nodes: {{err}}", err))
		}

		responses[i] = *NewApplicationResponse(&app, versions, nodes, upstreams)
	}

	return ctx.JSON(http.StatusOK, responses)
}

func CreateApp(ctx echo.Context) error {
	params := map[string]string{}
	err := ctx.Bind(&params)

	if err != nil {
		return NewHTTPError(http.StatusBadRequest, err)
	}

	app := &models.Application{
		Name:          params["name"],
		OwnerUsername: GetSessionAccount(ctx).Username,
		GitRepository: params["gitRepository"],
		Instances:     1,
		VersionTag:    params["versionTag"],
	}

	err = models.CreateApp(app)

	if err != nil && err == models.ErrUpdateConflict {
		return NewHTTPError(http.StatusConflict, err)
	} else if err != nil && err == models.ErrInvalidName {
		return NewHTTPError(http.StatusBadRequest, err)
	} else if err != nil {
		return NewHTTPError(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusCreated, NewApplicationResponse(app, nil, nil, nil))
}

func UpdateApp(ctx echo.Context) error {
	app := ctx.Get("app").(models.Application)

	jsonBuf := make([]byte, 1024)

	update := models.Application{
		Name:          app.Name,
		OwnerUsername: app.OwnerUsername,
		GitRepository: "",
	}

	if _, err := ctx.Request().Body.Read(jsonBuf); err != nil && err != io.EOF {
		return NewHTTPError(http.StatusInternalServerError, err)
	}

	gitRepository, valueType, _, err := jsonparser.Get(jsonBuf, "gitRepository")
	if err != jsonparser.KeyPathNotFoundError && err != nil {
		return NewHTTPError(http.StatusBadRequest, err)
	}

	if valueType != jsonparser.Null {
		update.GitRepository = string(gitRepository)
	}

	instances, valueType, _, err := jsonparser.Get(jsonBuf, "instances")
	if err != jsonparser.KeyPathNotFoundError && err != nil {
		return NewHTTPError(http.StatusBadRequest, err)
	}

	if valueType != jsonparser.NotExist {
		realValue, err := strconv.Atoi(string(instances))
		if err != nil {
			return NewHTTPError(http.StatusInternalServerError, err)
		}
		update.Instances = realValue
	}

	if err := app.Update(&update); err != nil {
		return NewHTTPError(http.StatusInternalServerError, err)
	}

	if err := swarm.UpdateAppService(&app); err != nil {
		return NewHTTPError(http.StatusInternalServerError, errwrap.Wrapf("apply changes to swarm: {{err}}", err))
	}

	return ctx.JSON(http.StatusOK, NewApplicationResponse(&app, nil, nil, nil))
}

func AddAppDomain(ctx echo.Context) error {
	app := ctx.Get("app").(models.Application)

	if err := app.AddUpstream(ctx.Param("domain")); err != nil {
		return NewHTTPError(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusOK, NewApplicationResponse(&app, nil, nil, nil))
}

func RemoveAppDomain(ctx echo.Context) error {
	app := ctx.Get("app").(models.Application)

	if err := app.RemoveUpstream(ctx.Param("domain")); err != nil {
		return NewHTTPError(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusOK, NewApplicationResponse(&app, nil, nil, nil))
}

func DeleteApp(ctx echo.Context) error {
	appName := ctx.Param("name")
	app, err := models.FindAppByName(appName)
	if err != nil {
		return NewHTTPError(http.StatusBadRequest, err)
	}

	if err := app.Destroy(); err != nil {
		return NewHTTPError(http.StatusInternalServerError, err)
	}

	if err := swarm.RemoveApp(app); err != nil {
		return NewHTTPError(http.StatusInternalServerError, err)
	}

	return ctx.String(http.StatusOK, "")
}
