package etcd

import (
	etcdv3 "github.com/coreos/etcd/clientv3"

	"github.com/jysperm/deploying/config"
)

var etcdConfig = etcdv3.Config{
	Endpoints: config.EtcdEndpoints,
}

var Client *etcdv3.Client

func init() {
	var err error

	Client, err = etcdv3.New(etcdConfig)

	if err != nil {
		panic(err)
	}
}
