package handlers

import (
	"io"
	"net/http"
	"strconv"

	"github.com/buger/jsonparser"

	"github.com/labstack/echo"

	"github.com/jysperm/deploying/lib/models"
	"github.com/jysperm/deploying/lib/swarm"
	"github.com/jysperm/deploying/web/handlers/helpers"
)

func ListDataSources(ctx echo.Context) error {
	account := helpers.GetSessionAccount(ctx)

	dataSources, err := models.GetDataSourcesOfAccount(account)

	if err != nil {
		return helpers.NewHTTPError(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusOK, helpers.NewDataSourcesResponse(dataSources))
}

func CreateDataSource(ctx echo.Context) error {
	params := map[string]string{}
	err := ctx.Bind(&params)

	if err != nil {
		return helpers.NewHTTPError(http.StatusBadRequest, err)
	}

	dataSource := &models.DataSource{
		Name:      params["name"],
		Type:      params["type"],
		Owner:     helpers.GetSessionAccount(ctx).Username,
		Instances: 1,
	}

	err = models.CreateDataSource(dataSource)

	if err != nil && err == models.ErrUpdateConflict {
		return helpers.NewHTTPError(http.StatusConflict, err)
	} else if err != nil && err == models.ErrInvalidName {
		return helpers.NewHTTPError(http.StatusBadRequest, err)
	} else if err != nil {
		return helpers.NewHTTPError(http.StatusInternalServerError, err)
	}

	if err := swarm.UpdateDataSource(dataSource, uint64(dataSource.Instances)); err != nil {
		return helpers.NewHTTPError(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusCreated, helpers.NewDataSourceResponse(dataSource))
}

func UpdateDataSource(ctx echo.Context) error {
	dataSource := ctx.Get("datasource").(models.DataSource)
	jsonBuf := make([]byte, 1024)

	if _, err := ctx.Request().Body.Read(jsonBuf); err != nil && err != io.EOF {
		return helpers.NewHTTPError(http.StatusBadRequest, err)
	}

	instances, valueType, _, err := jsonparser.Get(jsonBuf, "instances")
	if err != jsonparser.KeyPathNotFoundError && err != nil {
		return helpers.NewHTTPError(http.StatusBadRequest, err)
	}

	if valueType != jsonparser.NotExist {
		realValue, err := strconv.Atoi(string(instances))
		if err != nil {
			return helpers.NewHTTPError(http.StatusBadRequest, err)
		}
		dataSource.Instances = realValue
	}

	if err := dataSource.UpdateInstances(dataSource.Instances); err != nil {
		return helpers.NewHTTPError(http.StatusInternalServerError, err)
	}

	if err := swarm.UpdateDataSource(&dataSource, uint64(dataSource.Instances)); err != nil {
		return helpers.NewHTTPError(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusOK, helpers.NewDataSourceResponse(&dataSource))
}

func LinkDataSource(ctx echo.Context) error {
	appName := ctx.Param("appName")
	dataSource := ctx.Get("datasource").(models.DataSource)

	app, err := models.FindAppByName(appName)
	if err != nil {
		return helpers.NewHTTPError(http.StatusBadRequest, err)
	}

	if err := swarm.LinkDataSource(app, &dataSource); err != nil {
		return helpers.NewHTTPError(http.StatusInternalServerError, err)
	}

	if err := models.LinkDataSource(&dataSource, app); err != nil {
		return helpers.NewHTTPError(http.StatusInternalServerError, err)
	}

	return ctx.String(http.StatusOK, "")
}

func UnlinkDataSource(ctx echo.Context) error {
	appName := ctx.Param("appName")
	dataSource := helpers.GetDataSourceInfo(ctx)

	app, err := models.FindAppByName(appName)
	if err != nil {
		return helpers.NewHTTPError(http.StatusBadRequest, err)
	}

	if err := swarm.UnlinkDataSource(app, dataSource); err != nil {
		return helpers.NewHTTPError(http.StatusInternalServerError, err)
	}

	if err := models.UnlinkDataSource(dataSource, app); err != nil {
		return helpers.NewHTTPError(http.StatusInternalServerError, err)
	}

	return ctx.String(http.StatusOK, "")
}

func DeleteDataSource(ctx echo.Context) error {
	dataSource := helpers.GetDataSourceInfo(ctx)

	if err := swarm.RemoveDataSource(dataSource); err != nil {
		return helpers.NewHTTPError(http.StatusInternalServerError, err)
	}

	if err := models.DeleteDataSourceByName(dataSource.Name); err != nil {
		return helpers.NewHTTPError(http.StatusInternalServerError, err)
	}

	return ctx.String(http.StatusOK, "")
}

func CreateDataSourceNode(ctx echo.Context) error {
	params := map[string]string{}
	err := ctx.Bind(&params)

	if err != nil {
		return helpers.NewHTTPError(http.StatusBadRequest, err)
	}

	dataSourceNode := &models.DataSourceNode{
		Host: params["host"],
		Role: "master",
	}

	err = models.CreateDataSourceNode(helpers.GetDataSourceInfo(ctx), dataSourceNode)

	if err != nil {
		return helpers.NewHTTPError(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusCreated, helpers.NewDataSourceNodeResponse(dataSourceNode))
}
