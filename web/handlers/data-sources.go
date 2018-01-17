package handlers

import (
	"net/http"

	"github.com/labstack/echo"

	"github.com/jysperm/deploying/lib/models"
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

	return ctx.JSON(http.StatusCreated, helpers.NewDataSourceResponse(dataSource))
}

func UpdateDataSource(ctx echo.Context) error {
	return nil
}

func DeleteDataSource(ctx echo.Context) error {
	return nil
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
