package config

import (
	"os"
	"strings"

	"github.com/joho/godotenv"

	"github.com/jysperm/deploybeta/lib/utils"
)

var Listen string
var EtcdEndpoints []string
var DefaultRegistry string
var HttpProxy string
var AptMirror string
var NpmMirror string
var NvmMirror string
var HostPrivateAddress string
var DockerPrefix string
var WildcardDomain string

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
	AptMirror = os.Getenv("APT_MIRROR")
	NpmMirror = os.Getenv("NPM_MIRROR")
	NvmMirror = os.Getenv("NVM_MIRROR")
	HostPrivateAddress = os.Getenv("HOST_PRIVATE_ADDRESS")
	DockerPrefix = os.Getenv("DOCKER_PREFIX")
	WildcardDomain = os.Getenv("WILDCARD_DOMAIN")
}
