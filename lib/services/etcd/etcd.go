package etcd

import (
	etcd "github.com/coreos/etcd/client"

	"github.com/jysperm/deploying/config"
)

var etcdConfig = etcd.Config{
	Endpoints: config.EtcdEndpoints,
}

var Keys etcd.KeysAPI

func init() {
	connection, err := etcd.New(etcdConfig)

	if err != nil {
		panic(err)
	}

	Keys = etcd.NewKeysAPI(connection)
}
