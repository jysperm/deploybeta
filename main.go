package main

import (
	"github.com/jysperm/deploybeta/config"
	_ "github.com/jysperm/deploybeta/lib/datasource"
	"github.com/jysperm/deploybeta/web"
)

func main() {
	web.CreateWebServer().Start(config.Listen)
}
