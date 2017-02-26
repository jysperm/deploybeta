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

	app.Use(helpers.AuthenticateMiddleware)

	app.GET("/session/account", handlers.CurrentAccount)

	app.GET("/apps", handlers.GetMyApps)
	app.POST("/apps", handlers.CreateApp)
	app.PATCH("/apps/:name", handlers.UpdateApp)
	app.DELETE("/apps/:name", handlers.DeleteApp)

	return app
}
