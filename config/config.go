package config

import (
	"os"
	"strings"

	"github.com/joho/godotenv"

	"strconv"

	"github.com/jysperm/deploying/lib/utils"
)

var Listen string
var EtcdEndpoints []string
var DefaultRegistry string
var HttpProxy string
var HttpsProxy string
var AptCnMirror string
var NpmCnMirror string
var NvmCnMirror string
var DefaultRedisPort int
var DefaultMongoDBPort int
var HostPrivateAddress string

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
	HttpProxy = os.Getenv("PROXY_HTTP")
	HttpsProxy = os.Getenv("PROXY_HTTPS")
	AptCnMirror = os.Getenv("APT_CN_MIRROR")
	NpmCnMirror = os.Getenv("NPM_CN_MIRROR")
	NvmCnMirror = os.Getenv("NVM_CN_MIRROR")
	StrDefaultRedisPort := os.Getenv("DEFAULT_REDIS_PORT")
	StrDefaultMongoDBPort := os.Getenv("DEFAULT_MONGODB_PORT")
	DefaultRedisPort, _ = strconv.Atoi(StrDefaultRedisPort)
	DefaultMongoDBPort, _ = strconv.Atoi(StrDefaultMongoDBPort)
	HostPrivateAddress = os.Getenv("HOST_PRIVATE_ADDRESS")
}
