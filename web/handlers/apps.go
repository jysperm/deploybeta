package handlers

import (
	"io"
	"net/http"
	"strconv"

	"github.com/buger/jsonparser"
	"github.com/labstack/echo"

	appModel "github.com/jysperm/deploying/lib/models/app"
	"github.com/jysperm/deploying/lib/swarm"
	. "github.com/jysperm/deploying/web/handlers/helpers"
)

func GetMyApps(ctx echo.Context) error {
	account := GetSessionAccount(ctx)
	apps, err := appModel.GetAppsOfAccount(account)

	if err != nil {
		return NewHTTPError(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusOK, apps)
}

func CreateApp(ctx echo.Context) error {
	params := map[string]string{}
	err := ctx.Bind(&params)

	if err != nil {
		return NewHTTPError(http.StatusBadRequest, err)
	}

	app := &appModel.Application{
		Name:          params["name"],
		Owner:         GetSessionAccount(ctx).Username,
		GitRepository: "",
		Instances:     1,
		Version:       "",
	}

	err = appModel.CreateApp(app)

	if err != nil && err == appModel.ErrUpdateConflict {
		return NewHTTPError(http.StatusConflict, err)
	} else if err != nil && err == appModel.ErrInvalidName {
		return NewHTTPError(http.StatusBadRequest, err)
	} else if err != nil {
		return NewHTTPError(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusCreated, NewAppResponse(app))
}

func UpdateApp(ctx echo.Context) error {
	app := ctx.Get("app").(appModel.Application)

	jsonBuf := make([]byte, 1024)

	update := appModel.Application{
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

	if valueType != jsonparser.Null {
		realValue, err := strconv.Atoi(string(instances))
		if err != nil {
			return err
		}
		update.Instances = realValue
	}

	if err := app.Update(&update); err != nil {
		return NewHTTPError(http.StatusConflict, err)
	}

	if app.Instances != update.Instances {
		err := swarm.UpdateService(app)
		if err != nil {
			return NewHTTPError(http.StatusInternalServerError, err)
		}
	}

	return ctx.JSON(http.StatusOK, NewAppResponse(&app))
}

func DeleteApp(ctx echo.Context) error {
	appName := ctx.Param("name")
	if err := appModel.DeleteByName(appName); err != nil {
		return NewHTTPError(http.StatusInternalServerError, err)
	}

	return ctx.String(http.StatusOK, "")
}
