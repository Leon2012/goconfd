package client

import (
	"testing"
)

var addrs []string

func init() {
	addrs = []string{"127.0.0.1:1000", "127.0.0.1:1001", "127.0.0.1:1002", "127.0.0.1:1003", "127.0.0.1:1004", "127.0.0.1:1005"}
}

func TestDefaultSelector(t *testing.T) {
	selector := &DefaultSelector{}
	addr := selector.Select(addrs)
	t.Log(addr)
}

func TestRingSelector(t *testing.T) {
	selector := &RingSelector{}
	selector.Key = "1"
	addr := selector.Select(addrs)
	t.Log(addr)
}
