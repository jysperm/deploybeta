package web

import "github.com/kataras/iris"

import "github.com/jysperm/deploying/web/handlers"

func init() {
  app := iris.New()

  app.Get("/", func(ctx *iris.Context) {
    ctx.ServeFile("./web/frontend/public/index.html", true)
  })

  app.Get("/hello", handlers.Index)

  app.StaticWeb("/assets", "./web/frontend/public")

  app.Listen(":7000")
}
