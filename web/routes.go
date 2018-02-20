package web

import (
	"github.com/labstack/echo"

	"github.com/jysperm/deploying/web/handlers"
	"github.com/jysperm/deploying/web/handlers/helpers"
)

func CreateWebServer() *echo.Echo {
	app := echo.New()

	auth := helpers.AuthenticateMiddleware
	appOwner := helpers.AppOwnerMiddleware
	dataSource := helpers.DataSourceMiddleware
	dataSourceAgent := helpers.DataSourceAgentMiddleware

	app.File("/", "./frontend/public/index.html")
	app.Static("/assets", "./frontend/public")

	app.POST("/accounts", handlers.RegisterAccount)
	app.POST("/sessions", handlers.CreateSession)

	app.GET("/session/account", handlers.CurrentAccount, auth)

	app.POST("/apps/:name/versions", handlers.CreateVersion, auth, appOwner)
	app.PUT("/apps/:name/version", handlers.DeployVersion, auth, appOwner)
	app.GET("/apps/:name/versions/:tag/progress", handlers.PushProgress, auth, appOwner)

	app.GET("/apps", handlers.GetMyApps, auth)
	app.POST("/apps", handlers.CreateApp, auth)
	app.PATCH("/apps/:name", handlers.UpdateApp, auth, appOwner)
	app.DELETE("/apps/:name", handlers.DeleteApp, auth, appOwner)

	app.GET("/data-sources", handlers.ListDataSources, auth)
	app.POST("/data-sources", handlers.CreateDataSource, auth)
	app.PATCH("/data-sources/:name", handlers.UpdateDataSource, auth, dataSource)
	app.DELETE("/data-sources/:name", handlers.DeleteDataSource, auth, dataSource)

	app.POST("/data-sources/:name/links/:appName", handlers.LinkDataSource, auth, dataSource)
	app.PUT("/data-sources/:name/links/:appName", handlers.UnlinkDataSource, auth, dataSource)

	app.POST("/data-sources/:name/agents", handlers.CreateDataSourceNode, dataSourceAgent)
	app.PUT("/data-sources/:name/agents/:host", handlers.UpdateDataSourceNode, dataSourceAgent)
	app.GET("/data-sources/:name/agents/:host/commands", handlers.PollDataSourceNodeCommands, dataSourceAgent)

	return app
}
