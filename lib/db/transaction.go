package db

import (
	"context"
	"encoding/json"
	"errors"

	etcdv3 "github.com/coreos/etcd/clientv3"

	"github.com/jysperm/deploybeta/lib/utils"
)

var ErrEtcdTransactionFailed = errors.New("etcd transaction failed")

type Transaction struct {
	compares   []etcdv3.Cmp
	successOps []etcdv3.Op
	failedOps  []etcdv3.Op
	err        error
}

func RetryTransaction(updater func(tran *Transaction)) (*etcdv3.TxnResponse, error) {
	for attempts := 1; attempts <= 3; attempts++ {
		resp, err := StartTransaction(updater)

		if err == ErrEtcdTransactionFailed {
			continue
		} else if err != nil {
			return nil, err
		} else {
			return resp, nil
		}
	}

	return nil, ErrEtcdTransactionFailed
}

func StartTransaction(updater func(tran *Transaction)) (*etcdv3.TxnResponse, error) {
	tran := &Transaction{}

	updater(tran)

	resp, err := tran.Execute()

	if err != nil {
		return nil, err
	}

	if resp.Succeeded == false {
		return nil, ErrEtcdTransactionFailed
	}

	return resp, nil
}

func (tran *Transaction) Create(resource Resource) {
	key := resource.ResourceKey()

	for _, association := range resource.Associations() {
		association.onCreate(tran, resource)
	}

	dataBytes, err := json.Marshal(resource)

	if err != nil {
		tran.err = err
	} else {
		compare := etcdv3.Compare(etcdv3.CreateRevision(key), "=", 0)
		successOp := etcdv3.OpPut(key, string(dataBytes))

		tran.compares = append(tran.compares, compare)
		tran.successOps = append(tran.successOps, successOp)
	}
}

func (tran *Transaction) AddToStringSet(key string, value string) {
	resp, err := client.Get(context.Background(), key)

	if err != nil {
		tran.SetError(err)
		return
	}

	values := []string{}

	if len(resp.Kvs) > 0 {
		keyValue := resp.Kvs[0]

		err = json.Unmarshal([]byte(keyValue.Value), &values)

		if err != nil {
			tran.SetError(err)
			return
		}

		tran.compares = append(tran.compares, etcdv3.Compare(etcdv3.Version(key), "=", keyValue.Version))
	} else {
		tran.compares = append(tran.compares, etcdv3.Compare(etcdv3.CreateRevision(key), "=", 0))
	}

	values = utils.AddStringToUniqueSlice(values, value)

	dataBytes, err := json.Marshal(values)

	if err != nil {
		tran.SetError(err)
	} else {
		tran.successOps = append(tran.successOps, etcdv3.OpPut(key, string(dataBytes)))
	}
}

func (tran *Transaction) PullfromStringSet(key string, value string) {
	resp, err := client.Get(context.Background(), key)

	if err != nil {
		tran.SetError(err)
		return
	}

	values := []string{}

	if len(resp.Kvs) > 0 {
		keyValue := resp.Kvs[0]

		err = json.Unmarshal([]byte(keyValue.Value), &values)

		if err != nil {
			tran.SetError(err)
			return
		}

		tran.compares = append(tran.compares, etcdv3.Compare(etcdv3.Version(key), "=", keyValue.Version))
	} else {
		tran.compares = append(tran.compares, etcdv3.Compare(etcdv3.CreateRevision(key), "=", 0))
	}

	values = utils.PullStringFromSlice(values, value)

	dataBytes, err := json.Marshal(values)

	if err != nil {
		tran.SetError(err)
	} else {
		tran.successOps = append(tran.successOps, etcdv3.OpPut(key, string(dataBytes)))
	}
}

func (tran *Transaction) SetError(err error) {
	if tran.err == nil {
		tran.err = err
	}
}

func (tran *Transaction) Put(resource Resource) {
	key := resource.ResourceKey()

	for _, association := range resource.Associations() {
		association.onCreate(tran, resource)
	}

	dataBytes, err := json.Marshal(resource)

	if err != nil {
		tran.SetError(err)
	} else {
		successOp := etcdv3.OpPut(key, string(dataBytes))

		tran.successOps = append(tran.successOps, successOp)
	}
}

func (tran *Transaction) Update(resource Resource) {
	key := resource.ResourceKey()
	dataBytes, err := json.Marshal(resource)

	if err != nil {
		tran.SetError(err)
	} else {
		compare := etcdv3.Compare(etcdv3.Version(key), "=", resource.GetResourceMeta().EtcdVersion)
		successOp := etcdv3.OpPut(key, string(dataBytes))

		tran.compares = append(tran.compares, compare)
		tran.successOps = append(tran.successOps, successOp)
	}
}

func (tran *Transaction) Remove(resource Resource) {
	key := resource.ResourceKey()

	for _, association := range resource.Associations() {
		association.onDelete(tran, resource)
	}

	compare := etcdv3.Compare(etcdv3.CreateRevision(key), "!=", 0)
	successOp := etcdv3.OpDelete(key)

	tran.compares = append(tran.compares, compare)
	tran.successOps = append(tran.successOps, successOp)
}

func (tran *Transaction) Delete(resource Resource) {
	key := resource.ResourceKey()

	for _, association := range resource.Associations() {
		association.onDelete(tran, resource)
	}

	successOp := etcdv3.OpDelete(key)

	tran.successOps = append(tran.successOps, successOp)
}

func (tran *Transaction) Execute(updaters ...func() error) (*etcdv3.TxnResponse, error) {
	if tran.err != nil {
		return nil, tran.err
	}

	for _, updater := range updaters {
		err := updater()

		if err != nil {
			return nil, err
		}
	}

	return client.Txn(context.Background()).
		If(tran.compares...).
		Then(tran.successOps...).
		Else(tran.failedOps...).
		Commit()
}

func (tran *Transaction) ExecuteMustSuccess() error {
	resp, err := tran.Execute()

	if err != nil {
		return err
	}

	if resp.Succeeded == false {
		return ErrEtcdTransactionFailed
	}

	return nil
}
