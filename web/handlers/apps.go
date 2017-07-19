package handlers

import (
	"net/http"

	"github.com/labstack/echo"

	appModel "github.com/jysperm/deploying/lib/models/app"
	versionMdel "github.com/ jysperm/deploying/lib/models/version"
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
	appName := ctx.Param("name")
	app, err := appModel.FindByName(appName)
	if err != nil {
		return NewHTTPError(http.StatusBadRequest, err)
	}

	update := new(appModel.Application)
	if err := ctx.Bind(update); err != nil {
		return NewHTTPError(http.StatusBadRequest, err)
	}

	if app.GitRepository != update.GitRepository {
		if err := app.UpdateGitRepository(update.GitRepository); err != nil {
			return NewHTTPError(http.StatusInternalServerError, err)
		}
	}

	if app.Instances != update.Instances || app.Version != update.Version {

		if app.Instances != update.Instances {
			if err := app.UpdateInstances(update.Instances); err != nil {
				return NewHTTPError(http.StatusInternalServerError, err)
			}
		}

		if app.Version != update.Version {
			if version, err := versionModel.FindByTag(*app, update.Version); version == nil && err != nil {
				return NewHTTPError(http.StatusBadRequest, err)
			}
			if err := app.UpdateVersion(update.Version); err != nil {
				return NewHTTPError(http.StatusInternalServerError, err)
			}
		}

		if err := swarm.UpdateService(*app); err != nil {
			return NewHTTPError(http.StatusInternalServerError, err)
		}
	}

	return ctx.JSON(http.StatusCreated, NewAppResponse(app))
}

func DeleteApp(ctx echo.Context) error {
	appName := ctx.Param("name")
	if err := appModel.DeleteByName(appName); err != nil {
		return NewHTTPError(http.StatusInternalServerError, err)
	}

	return ctx.String(http.StatusOK, "")
}
