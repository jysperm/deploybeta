package main

import (
	"github.com/jysperm/deploying/config"
	"github.com/jysperm/deploying/web"
)

func main() {
	web.Listen(config.Port)
}