package agent

import (
	"github.com/Leon2012/goconfd/libs/kv"
	_ "time"
)

type LocalInterface interface {
	Save(k *kv.Kv) error
	Get(k string) (*kv.Kv, error)
	Keys() []string
}
