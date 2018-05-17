package main

import (
	"github.com/jysperm/deploybeta/config"
	"github.com/jysperm/deploybeta/web"
)

func main() {
	web.CreateWebServer().Start(config.Listen)
}
