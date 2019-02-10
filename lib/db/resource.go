package db

import (
	"encoding/json"
	"errors"
	"reflect"
	"bytes"

	"golang.org/x/net/context"
	etcdv3 "github.com/coreos/etcd/clientv3"
)

var ErrResourceNotFound = errors.New("resource not found")

// Resource represents a Golang struct can store in Etcd, in JSON serialization.
type Resource interface {
	ResourceKey() string
	Associations() []Association
	GetResourceMeta() *ResourceMeta
}

type ResourceMeta struct {
	EtcdVersion int64 `json:"-"`
}

func (meta *ResourceMeta) GetResourceMeta() *ResourceMeta {
	return meta
}

func Assign(dest Resource, src Resource) {
	reflect.ValueOf(dest).Set(reflect.ValueOf(src))
}

func Fetch(resource Resource) error {
	return FetchFrom(resource.ResourceKey(), resource)
}

func FetchFrom(key string, resource Resource) error {
	resp, err := client.Get(context.Background(), key)

	if err != nil {
		return err
	}

	if len(resp.Kvs) == 0 {
		return ErrResourceNotFound
	}

	keyValue := resp.Kvs[0]

	err = json.Unmarshal([]byte(keyValue.Value), resource)

	resource.GetResourceMeta().EtcdVersion = keyValue.Version

	if err != nil {
		return err
	}

	return nil
}

func FetchAllFrom(prefix string, resources interface{}) error {
	resp, err := client.Get(context.Background(), prefix, etcdv3.WithPrefix())

	if err != nil {
		return err
	}

	resourceBytesList := [][]byte{}

	for _, keyValue := range resp.Kvs {
		resourceBytesList = append(resourceBytesList, keyValue.Value)
	}

	resourcesBytes := bytes.Buffer{}
	resourcesBytes.Write([]byte("["))
	resourcesBytes.Write(bytes.Join(resourceBytesList, []byte(",")))
	resourcesBytes.Write([]byte("]"))

	return json.Unmarshal(resourcesBytes.Bytes(), resources)
}

func clone(dest Resource, src Resource) {
	reflect.ValueOf(&dest).Set(reflect.New(reflect.TypeOf(reflect.ValueOf(src))))
}
