package helpers

import (
	"errors"
	"net/http"

	"github.com/hashicorp/errwrap"
	"github.com/labstack/echo"

	"github.com/jysperm/deploying/lib/models"
)

func AuthenticateMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		session, err := models.FindSessionByToken(ctx.Request().Header.Get("Authorization"))

		if err != nil {
			return NewHTTPError(http.StatusUnauthorized, err)
		}

		account, err := models.FindAccountByName(session.Username)

		if err != nil {
			return NewHTTPError(http.StatusUnauthorized, err)
		}

		ctx.Set("account", &account)

		return next(ctx)
	}
}

func AppOwnerMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		account := GetSessionAccount(ctx)
		appName := ctx.Param("name")
		apps, err := models.GetAppsOfAccount(account)

		if err != nil {
			return NewHTTPError(http.StatusInternalServerError, err)
		}

		if err == nil && len(apps) == 0 {
			return NewHTTPError(http.StatusUnauthorized, errors.New("Not found application"))
		}

		for _, app := range apps {
			if app.Name == appName {
				ctx.Set("app", app)
				return next(ctx)
			}
		}

		return NewHTTPError(http.StatusBadRequest, errors.New("Not found application"))
	}
}

func DataSourceMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		account := GetSessionAccount(ctx)
		dataSourceName := ctx.Param("name")

		dataSource, err := models.GetDataSourceOfAccount(dataSourceName, account)

		if err != nil {
			return NewHTTPError(http.StatusInternalServerError, err)
		}

		if dataSource == nil {
			return NewHTTPError(http.StatusUnauthorized, errors.New("Not found datasource"))
		}

		ctx.Set("dataSource", &dataSource)

		return next(ctx)
	}
}

func DataSourceAgentMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		dataSource, err := models.FindDataSourceByName(ctx.Param("name"))

		if err != nil {
			return NewHTTPError(http.StatusBadRequest, errwrap.Wrapf("find dataSource: {{err}}", err))
		}

		if dataSource.AgentToken != ctx.Request().Header.Get("Authorization") {
			return NewHTTPError(http.StatusUnauthorized, errors.New("invalid agent token"))
		}

		dataSouceNode, err := dataSource.FindNodeByHost(ctx.Param("host"))

		if err != nil {
			return NewHTTPError(http.StatusBadRequest, errwrap.Wrapf("find dataSource node: {{err}}", err))
		}

		ctx.Set("dataSource", &dataSource)
		ctx.Set("dataSourceNode", &dataSouceNode)

		return next(ctx)
	}
}

func GetSessionAccount(ctx echo.Context) *models.Account {
	return ctx.Get("account").(*models.Account)
}

func GetDataSourceInfo(ctx echo.Context) *models.DataSource {
	return ctx.Get("dataSource").(*models.DataSource)
}

func GetDataSourceNodeInfo(ctx echo.Context) *models.DataSourceNode {
	return ctx.Get("dataSourceNode").(*models.DataSourceNode)
}
