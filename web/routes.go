package web

import (
	"github.com/labstack/echo"

	"github.com/jysperm/deploying/web/handlers"
	"github.com/jysperm/deploying/web/handlers/helpers"
)

var app = echo.New()

func init() {
	app.File("/", "./web/frontend/public/index.html")
	app.Static("/assets", "./web/frontend/public")

	app.POST("/accounts", handlers.RegisterAccount)
	app.POST("/sessions", handlers.CreateSession)

	app.Use(helpers.AuthenticateMiddleware)

	app.GET("/session/account", handlers.CurrentAccount)
}

func Listen(port string) {
	app.Start(port)
}
