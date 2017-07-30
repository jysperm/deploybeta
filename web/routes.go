package web

import (
	"github.com/labstack/echo"

	"github.com/jysperm/deploying/web/handlers"
	"github.com/jysperm/deploying/web/handlers/helpers"
)

func CreateWebServer() *echo.Echo {
	app := echo.New()

	app.File("/", "./frontend/public/index.html")
	app.Static("/assets", "./frontend/public")

	app.POST("/accounts", handlers.RegisterAccount)
	app.POST("/sessions", handlers.CreateSession)

	app.GET("/session/account", handlers.CurrentAccount, helpers.AuthenticateMiddleware)

	app.POST("/apps/:name/versions", handlers.CreateVersion, helpers.AuthenticateMiddleware, helpers.AppOwnerMiddleware)
	app.POST("/apps/:name/version", handlers.CreateAndDeploy, helpers.AuthenticateMiddleware, helpers.AppOwnerMiddleware)
	app.PUT("/apps/:name/version", handlers.DeployVersion, helpers.AuthenticateMiddleware, helpers.AppOwnerMiddleware)

	app.GET("/apps", handlers.GetMyApps, helpers.AuthenticateMiddleware)
	app.POST("/apps", handlers.CreateApp, helpers.AuthenticateMiddleware)
	app.PATCH("/apps/:name", handlers.UpdateApp, helpers.AuthenticateMiddleware, helpers.AppOwnerMiddleware)
	app.DELETE("/apps/:name", handlers.DeleteApp, helpers.AuthenticateMiddleware)

	return app
}
