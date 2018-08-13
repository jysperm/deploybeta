package db

import (
	"context"
	"encoding/json"
	"errors"

	etcdv3 "github.com/coreos/etcd/clientv3"

	"github.com/jysperm/deploybeta/lib/utils"
)

var ErrEtcdTransactionFailed = errors.New("etcd transaction failed")

// Transaction represents a Txn operation of Etcd.
type Transaction interface {
	Put(resource Resource)
	Create(resource Resource)
	Update(resource Resource)
	Remove(resource Resource)
	Delete(resource Resource)
	DeleteKey(key string)
	DeletePrefix(prefix string)

	AddToStringSet(key string, value string)
	PullfromStringSet(key string, value string)

	SetError(err error)
	Execute(updaters ...func() error) (*etcdv3.TxnResponse, error)
}

func NewTransaction() Transaction {
	return &transaction{}
}

func StartTransaction(updater func(tran Transaction)) (*etcdv3.TxnResponse, error) {
	tran := NewTransaction()

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

func RetryTransaction(updater func(tran Transaction)) (*etcdv3.TxnResponse, error) {
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

type transaction struct {
	compares   []etcdv3.Cmp
	successOps []etcdv3.Op
	failedOps  []etcdv3.Op
	err        error
}

func (tran *transaction) Put(resource Resource) {
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

func (tran *transaction) Create(resource Resource) {
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

func (tran *transaction) Update(resource Resource) {
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

func (tran *transaction) Remove(resource Resource) {
	key := resource.ResourceKey()

	for _, association := range resource.Associations() {
		association.onDelete(tran, resource)
	}

	compare := etcdv3.Compare(etcdv3.CreateRevision(key), "!=", 0)
	successOp := etcdv3.OpDelete(key)

	tran.compares = append(tran.compares, compare)
	tran.successOps = append(tran.successOps, successOp)
}

func (tran *transaction) Delete(resource Resource) {
	for _, association := range resource.Associations() {
		association.onDelete(tran, resource)
	}

	DeleteKey(resource.ResourceKey())
}

func (tran *transaction) DeleteKey(key string) {
	successOp := etcdv3.OpDelete(key)

	tran.successOps = append(tran.successOps, successOp)
}

func (tran *transaction) DeletePrefix(prefix string) {
	successOp := etcdv3.OpDelete(prefix, etcdv3.WithPrefix())

	tran.successOps = append(tran.successOps, successOp)
}

func (tran *transaction) AddToStringSet(key string, value string) {
	values := []string{}
	kv, err := FetchJSON(key, &values)

	if err != nil {
		tran.SetError(err)
		return
	}

	if kv != nil {
		tran.compares = append(tran.compares, etcdv3.Compare(etcdv3.Version(key), "=", kv.Version))
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

func (tran *transaction) PullfromStringSet(key string, value string) {
	values := []string{}
	kv, err := FetchJSON(key, &values)

	if err != nil {
		tran.SetError(err)
		return
	}

	if kv != nil {
		tran.compares = append(tran.compares, etcdv3.Compare(etcdv3.Version(key), "=", kv.Version))

		values = utils.PullStringFromSlice(values, value)

		dataBytes, err := json.Marshal(values)

		if err != nil {
			tran.SetError(err)
		} else {
			tran.successOps = append(tran.successOps, etcdv3.OpPut(key, string(dataBytes)))
		}

	} else {
		tran.compares = append(tran.compares, etcdv3.Compare(etcdv3.CreateRevision(key), "=", 0))
	}
}

func (tran *transaction) SetError(err error) {
	if tran.err == nil {
		tran.err = err
	}
}

func (tran *transaction) Execute(updaters ...func() error) (*etcdv3.TxnResponse, error) {
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
