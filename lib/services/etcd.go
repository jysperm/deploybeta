package services

import (
	"context"
	"encoding/json"

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

type EtcdTransaction struct {
	watchedKeys map[string]interface{}
	compares    []etcdv3.Cmp
	successOps  []etcdv3.Op
	failedOps   []etcdv3.Op
	err         error
}

func NewEtcdTransaction() *EtcdTransaction {
	return &EtcdTransaction{
		watchedKeys: make(map[string]interface{}),
	}
}

func (tran *EtcdTransaction) WatchJSON(key string, schema interface{}) {
	tran.watchedKeys[key] = schema
}

func (tran *EtcdTransaction) CreateJSON(key string, data interface{}) {
	dataBytes, err := json.Marshal(data)

	if err != nil {
		tran.err = err
	} else {
		compare := etcdv3.Compare(etcdv3.CreateRevision(key), "=", 0)
		successOp := etcdv3.OpPut(key, string(dataBytes))

		tran.compares = append(tran.compares, compare)
		tran.successOps = append(tran.successOps, successOp)
	}
}

func (tran *EtcdTransaction) PutJSONOnSuccess(key string, data interface{}) {
	dataBytes, err := json.Marshal(data)

	if err != nil {
		tran.err = err
	} else {
		tran.successOps = append(tran.successOps, etcdv3.OpPut(key, string(dataBytes)))
	}
}

func (tran *EtcdTransaction) Execute(resolvers ...func(map[string]interface{}) error) (*etcdv3.TxnResponse, error) {
	if tran.err != nil {
		return nil, tran.err
	}

	for key, schema := range tran.watchedKeys {
		resp, err := EtcdClient.Get(context.Background(), key)

		if err != nil {
			return nil, err
		}

		if len(resp.Kvs) > 0 {
			keyValue := resp.Kvs[0]

			tran.compares = append(tran.compares, etcdv3.Compare(etcdv3.Version(key), "=", keyValue.Version))

			err = json.Unmarshal([]byte(keyValue.Value), schema)

			if err != nil {
				return nil, err
			}
		}
	}

	for _, resolver := range resolvers {
		err := resolver(tran.watchedKeys)

		if err != nil {
			return nil, err
		}
	}

	return EtcdClient.Txn(context.Background()).
		If(tran.compares...).
		Then(tran.successOps...).
		Else(tran.failedOps...).
		Commit()
}
