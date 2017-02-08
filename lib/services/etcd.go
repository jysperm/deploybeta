package services

import (
	etcdv3 "github.com/coreos/etcd/clientv3"

	"github.com/jysperm/deploying/config"
)

var etcdConfig = etcdv3.Config{
	Endpoints: config.EtcdEndpoints,
}

var EtcdClient *etcdv3.Client

func init() {
	var err error

	EtcdClient, err = etcdv3.New(etcdConfig)

	if err != nil {
		panic(err)
	}
}
