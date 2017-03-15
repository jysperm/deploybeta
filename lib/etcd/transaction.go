package etcd

import (
	"context"
	"encoding/json"

	etcdv3 "github.com/coreos/etcd/clientv3"
)

type Transaction struct {
	watchedKeys map[string]interface{}
	compares    []etcdv3.Cmp
	successOps  []etcdv3.Op
	failedOps   []etcdv3.Op
	err         error
}

func NewTransaction() *Transaction {
	return &Transaction{
		watchedKeys: make(map[string]interface{}),
	}
}

func (tran *Transaction) WatchJSON(key string, schema interface{}) {
	tran.watchedKeys[key] = schema
}

func (tran *Transaction) CreateJSON(key string, data interface{}) {
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

func (tran *Transaction) PutJSONOnSuccess(key string, data interface{}) {
	dataBytes, err := json.Marshal(data)

	if err != nil {
		tran.err = err
	} else {
		tran.successOps = append(tran.successOps, etcdv3.OpPut(key, string(dataBytes)))
	}
}

func (tran *Transaction) Execute(resolvers ...func(map[string]interface{}) error) (*etcdv3.TxnResponse, error) {
	if tran.err != nil {
		return nil, tran.err
	}

	for key, schema := range tran.watchedKeys {
		resp, err := Client.Get(context.Background(), key)

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

	return Client.Txn(context.Background()).
		If(tran.compares...).
		Then(tran.successOps...).
		Else(tran.failedOps...).
		Commit()
}
