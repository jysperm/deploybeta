package web

import (
	"github.com/kataras/iris"

	"github.com/jysperm/deploying/web/handlers"
	"github.com/jysperm/deploying/web/handlers/helpers"
)

var app = iris.New()

func init() {
	app.Get("/", func(ctx *iris.Context) {
		ctx.ServeFile("./web/frontend/public/index.html", true)
	})

	app.StaticWeb("/assets", "./web/frontend/public")

	app.Post("/accounts", handlers.RegisterAccount)
	app.Post("/sessions", handlers.CreateSession)

	app.Use(&helpers.AuthenticateMiddleware{})

	app.Get("/session", handlers.CurrentAccount)
}

func Listen(port string) {
	app.Listen(port)
}
