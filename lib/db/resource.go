package db

import (
	"encoding/json"
	"errors"
	"reflect"

	"golang.org/x/net/context"
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

func clone(dest Resource, src Resource) {
	reflect.ValueOf(&dest).Set(reflect.New(reflect.TypeOf(reflect.ValueOf(src))))
}
