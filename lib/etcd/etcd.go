package etcd

import (
	"encoding/json"

	etcdv3 "github.com/coreos/etcd/clientv3"

	"github.com/jysperm/deploybeta/config"

	"golang.org/x/net/context"
)

var etcdConfig = etcdv3.Config{
	Endpoints: config.EtcdEndpoints,
}

var Client *etcdv3.Client

func LoadKey(key string, modelStruct interface{}) (bool, error) {
	resp, err := Client.Get(context.Background(), key)

	if err != nil || len(resp.Kvs) == 0 {
		return false, err
	}

	err = json.Unmarshal([]byte(resp.Kvs[0].Value), modelStruct)

	if err != nil {
		return false, err
	}

	return true, nil
}

func PutKey(key string, data interface{}) error {
	dataBytes, err := json.Marshal(data)

	if err != nil {
		return err
	}

	_, err = Client.Put(context.Background(), key, string(dataBytes))

	return err
}

func init() {
	var err error

	Client, err = etcdv3.New(etcdConfig)

	if err != nil {
		panic(err)
	}
}
