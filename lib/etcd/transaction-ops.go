package etcd

func (tran *Transaction) AppendStringArray(key string, values ...string) {
	tran.WatchJSON(key, &[]string{}, func(watchedKey interface{}) error {
		keyValues := *watchedKey.(*[]string)

		tran.PutJSON(key, append(keyValues, values...))

		return nil
	})
}
