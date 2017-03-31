package registry

import (
	_ "time"

	"github.com/Leon2012/goconfd/libs/kv"
)

type Backend interface {
	Del(k string) error
	Put(k, v string) error
	Get(k string) ([]*kv.Kv, error)
	BatchGetByPrefix(prefix string) ([]*kv.Kv, error)
	WatchWithPrefix(prefix string, f kv.KvFunc)
	WatchWithPrefixs(prefixs []string, f kv.KvFunc)
	WatchWithKey(key string, f kv.KvFunc)
	WatchWithKeys(keys []string, f kv.KvFunc)
	PutWithTTL(k, v string, ttl int64) error
	SetTTL(k string) error
	Close()
}

var DefaultBackend Backend
