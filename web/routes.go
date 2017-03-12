package web

import (
	"github.com/labstack/echo"

	"github.com/jysperm/deploying/web/handlers"
	"github.com/jysperm/deploying/web/handlers/helpers"
)

func CreateWebServer() *echo.Echo {
	app := echo.New()

	app.File("/", "./web/frontend/public/index.html")
	app.Static("/assets", "./web/frontend/public")

	app.POST("/accounts", handlers.RegisterAccount)
	app.POST("/sessions", handlers.CreateSession)

	app.GET("/session/account", handlers.CurrentAccount, helpers.AuthenticateMiddleware)

	app.GET("/apps", handlers.GetMyApps, helpers.AuthenticateMiddleware)
	app.POST("/apps", handlers.CreateApp, helpers.AuthenticateMiddleware)
	app.PATCH("/apps/:name", handlers.UpdateApp, helpers.AuthenticateMiddleware)
	app.DELETE("/apps/:name", handlers.DeleteApp, helpers.AuthenticateMiddleware)
	app.POST("/apps/:name/images", handlers.CreateImage, helpers.AuthenticateMiddleware)

	return app
}
