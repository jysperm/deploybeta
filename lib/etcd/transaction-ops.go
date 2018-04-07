package etcd

import (
	"github.com/jysperm/deploying/lib/utils"
)

func (tran *Transaction) AppendStringArray(key string, values ...string) {
	tran.WatchJSON(key, &[]string{}, func(watchedKey interface{}) error {
		keyValues := *watchedKey.(*[]string)

		tran.PutJSON(key, append(keyValues, values...))

		return nil
	})
}

func (tran *Transaction) PullStringArray(key string, value string) {
	tran.WatchJSON(key, &[]string{}, func(watchedKey interface{}) error {
		keyValues := *watchedKey.(*[]string)

		tran.PutJSON(key, utils.PullStringFromSlice(keyValues, value))

		return nil
	})
}
