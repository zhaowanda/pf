package nosql

import "container/list"

type KV interface {
	Put(key, value string) (bool, error)
	Get(key string) (string, error)
	Del(key string) (bool, error)
	BatchPut(mapKey map[string]string) (bool, error)
	BatchDel(keys list.List) (bool, error)
}
