package db

import (
	"encoding/json"

	etcdv3 "github.com/coreos/etcd/clientv3"
	etcdv3pb "github.com/coreos/etcd/mvcc/mvccpb"
	"golang.org/x/net/context"

	"github.com/jysperm/deploybeta/config"
)

var client *etcdv3.Client
var Client *etcdv3.Client

func FetchJSON(key string, data interface{}) (*etcdv3pb.KeyValue, error) {
	resp, err := client.Get(context.Background(), key)

	if err != nil {
		return nil, err
	}

	if resp.Count > 0 {
		keyValue := resp.Kvs[0]

		return keyValue, json.Unmarshal(keyValue.Value, data)
	} else {
		return nil, nil
	}
}

func PutJSON(key string, data interface{}) error {
	dataBytes, err := json.Marshal(data)

	if err != nil {
		return err
	}

	_, err = client.Put(context.Background(), key, string(dataBytes))

	return err
}

func DeleteKey(key string) error {
	_, err := client.Delete(context.Background(), key)

	return err
}

func init() {
	var err error

	client, err = etcdv3.New(etcdv3.Config{
		Endpoints: config.EtcdEndpoints,
	})

	Client = client

	if err != nil {
		panic(err)
	}
}
