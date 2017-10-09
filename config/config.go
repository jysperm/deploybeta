package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"

	"github.com/jysperm/deploying/lib/utils"
)

var Listen string
var EtcdEndpoints []string
var DefaultRegistry string

func init() {
	err := godotenv.Load()

	if !strings.Contains(err.Error(), "no such file or directory") {
		panic("Load .env failed: " + err.Error())
	}

	err = godotenv.Load(utils.GetAssetFilePath("defaults.env"))

	if err != nil {
		panic("Load defaults.env failed: " + err.Error())
	}

	Listen = os.Getenv("LISTEN")
	EtcdEndpoints = strings.Split(os.Getenv("ETCD_ENDPOINTS"), ",")
	DefaultRegistry = os.Getenv("DEFAULT_REGISTRY")

	fmt.Println(Listen, EtcdEndpoints, DefaultRegistry)
}
