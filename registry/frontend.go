package registry

import (
	_ "time"

	"github.com/Leon2012/goconfd/libs/kv"
)

type Frontend interface {
	Save(k *kv.Kv) error
	Get(k string) (*kv.Kv, error)
	Keys() []string
}

var DefaultFrontend Frontend
