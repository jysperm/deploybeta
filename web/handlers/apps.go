package handlers

import (
	"io"
	"net/http"
	"strconv"

	"github.com/buger/jsonparser"
	"github.com/labstack/echo"

	"github.com/jysperm/deploying/lib/models"
	"github.com/jysperm/deploying/lib/swarm"
	. "github.com/jysperm/deploying/web/handlers/helpers"
)

func GetMyApps(ctx echo.Context) error {
	account := GetSessionAccount(ctx)
	apps, err := models.GetAppsOfAccount(account)

	if err != nil {
		return NewHTTPError(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusOK, NewAppsResponse(apps))
}

func CreateApp(ctx echo.Context) error {
	params := map[string]string{}
	err := ctx.Bind(&params)

	if err != nil {
		return NewHTTPError(http.StatusBadRequest, err)
	}

	app := &models.Application{
		Name:          params["name"],
		Owner:         GetSessionAccount(ctx).Username,
		GitRepository: params["gitRepository"],
		Instances:     1,
		Version:       params["version"],
	}

	err = models.CreateApp(app)

	if err != nil && err == models.ErrUpdateConflict {
		return NewHTTPError(http.StatusConflict, err)
	} else if err != nil && err == models.ErrInvalidName {
		return NewHTTPError(http.StatusBadRequest, err)
	} else if err != nil {
		return NewHTTPError(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusCreated, NewAppResponse(app))
}

func UpdateApp(ctx echo.Context) error {
	app := ctx.Get("app").(models.Application)

	jsonBuf := make([]byte, 1024)

	update := models.Application{
		Name:          app.Name,
		Owner:         app.Owner,
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

	if err := swarm.UpdateService(&update); err != nil {
		return NewHTTPError(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusOK, NewAppResponse(&app))
}

func DeleteApp(ctx echo.Context) error {
	appName := ctx.Param("name")
	app, err := models.FindAppByName(appName)
	if err != nil {
		return NewHTTPError(http.StatusBadRequest, err)
	}
	if err := swarm.RemoveService(app); err != nil {
		return NewHTTPError(http.StatusInternalServerError, err)
	}

	return ctx.String(http.StatusOK, "")
}
