package db

import (
	"context"
	"encoding/json"

	etcdv3 "github.com/coreos/etcd/clientv3"
)

func WatchUpdates(resource Resource) (func(), <-chan Resource, <-chan error) {
	var canceled bool
	var watcher etcdv3.Watcher
	var updated Resource
	var err error

	cancel := func() {
		canceled = true

		if watcher != nil {
			watcher.Close()
			watcher = nil
		}
	}

	updateds := make(chan Resource)
	errs := make(chan error)

	watcher = etcdv3.NewWatcher(client)

	go func() {
		for response := range watcher.Watch(context.TODO(), resource.ResourceKey()) {
			for _, event := range response.Events {
				clone(updated, resource)

				err = json.Unmarshal(event.Kv.Value, updated)

				if err != nil {
					errs <- err
				} else {
					updateds <- updated
				}

				if canceled {
					return
				}
			}
		}
	}()

	return cancel, updateds, errs
}
