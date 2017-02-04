package web

import "github.com/kataras/iris"

import "github.com/jysperm/deploying/web/handlers"

var app = iris.New()

func init() {
	app.Get("/", func(ctx *iris.Context) {
		ctx.ServeFile("./web/frontend/public/index.html", true)
	})

	app.Post("/accounts", handlers.RegisterAccount)
	app.Post("/sessions", handlers.CreateSession)

	app.StaticWeb("/assets", "./web/frontend/public")
}

func Listen(port string) {
	app.Listen(port)
}
