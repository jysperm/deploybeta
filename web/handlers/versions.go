package handlers

import (
	"net/http"

	"github.com/labstack/echo"

	appModel "github.com/jysperm/deploying/lib/models/app"
	versionModel "github.com/jysperm/deploying/lib/models/version"
	"github.com/jysperm/deploying/lib/swarm"
	. "github.com/jysperm/deploying/web/handlers/helpers"
)

func CreateVersion(ctx echo.Context) error {
	app := ctx.Get("app").(appModel.Application)

	params := map[string]string{}
	if err := ctx.Bind(&params); err != nil {
		return NewHTTPError(http.StatusBadRequest, err)
	}

	version, err := versionModel.CreateVersion(&app, "", params["gitTag"])
	if err != nil {
		return NewHTTPError(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusCreated, NewVersionResponse(&version))
}

func DeployVersion(ctx echo.Context) error {
	app := ctx.Get("app").(appModel.Application)

	params := map[string]string{}
	if err := ctx.Bind(&params); err != nil {
		return NewHTTPError(http.StatusBadRequest, err)
	}

	version, err := versionModel.FindByTag(app, params["tag"])
	if err != nil || version == nil {
		return NewHTTPError(http.StatusBadRequest, err)
	}

	if err := app.UpdateVersion(version.Tag); err != nil {
		return NewHTTPError(http.StatusInternalServerError, err)
	}

	if err := swarm.UpdateService(app); err != nil {
		return NewHTTPError(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusOK, NewVersionResponse(version))
}

func CreateAndDeploy(ctx echo.Context) error {
	app := ctx.Get("app").(appModel.Application)

	params := map[string]string{}
	if err := ctx.Bind(&params); err != nil {
		return NewHTTPError(http.StatusBadRequest, err)
	}

	version, err := versionModel.CreateVersion(&app, "", params["gitTag"])
	if err != nil {
		return NewHTTPError(http.StatusInternalServerError, err)
	}

	if err := app.UpdateVersion(version.Tag); err != nil {
		return NewHTTPError(http.StatusInternalServerError, err)
	}

	if err := swarm.UpdateService(app); err != nil {
		return NewHTTPError(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusCreated, NewVersionResponse(&version))
}

// TODO:
func DeleteVersion(ctx echo.Context) error {
	return nil
}
