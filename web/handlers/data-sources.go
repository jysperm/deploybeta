package handlers

import (
	"net/http"

	"github.com/labstack/echo"

	"github.com/jysperm/deploying/lib/models"
	. "github.com/jysperm/deploying/web/handlers/helpers"
)

func ListDataSources(ctx echo.Context) error {
	account := GetSessionAccount(ctx)

	dataSources, err := models.GetDataSourcesOfAccount(account)

	if err != nil {
		return NewHTTPError(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusOK, NewDataSourcesResponse(dataSources))
}

func CreateDataSource(ctx echo.Context) error {
	params := map[string]string{}
	err := ctx.Bind(&params)

	if err != nil {
		return NewHTTPError(http.StatusBadRequest, err)
	}

	dataSource := &models.DataSource{
		Name:      params["name"],
		Type:      params["type"],
		Owner:     GetSessionAccount(ctx).Username,
		Instances: 1,
	}

	err = models.CreateDataSource(dataSource)

	if err != nil && err == models.ErrUpdateConflict {
		return NewHTTPError(http.StatusConflict, err)
	} else if err != nil && err == models.ErrInvalidName {
		return NewHTTPError(http.StatusBadRequest, err)
	} else if err != nil {
		return NewHTTPError(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusCreated, NewDataSourceResponse(dataSource))
}
