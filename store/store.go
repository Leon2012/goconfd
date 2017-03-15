package store

import (
	"github.com/Leon2012/goconfd/store/db"
)

var adaptr db.Adapter

func Register(name string, adapter db.Adapter) {
	if adapter == nil {
		panic("store: Register adapter is nil")
	}
	if adaptr != nil {
		panic("store: Adapter already registered")
	}
	adaptr = adapter
}
