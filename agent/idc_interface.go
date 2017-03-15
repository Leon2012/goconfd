package agent

import (
	_ "time"

	"github.com/Leon2012/goconfd/libs/kv"
)

type IdcInterface interface {
	Put(k, v string) error
	Get(k string) ([]*kv.Kv, error)
	BatchGetByPrefix(prefix string) ([]*kv.Kv, error)
	WatchWithPrefix(prefix string, f kv.KvFunc)
	WatchWithPrefixs(prefixs []string, f kv.KvFunc)
	WatchWithKey(key string, f kv.KvFunc)
	WatchWithKeys(keys []string, f kv.KvFunc)
	Close()
}
