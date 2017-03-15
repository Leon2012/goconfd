package kv

import "fmt"

const (
	KV_EVENT_NONE   = -1
	KV_EVENT_PUT    = 0
	KV_EVENT_DELETE = 1
)

type Kv struct {
	Revision int64  `json:"rev" php:"rev"`
	Event    int32  `json:"type" php:"type"`
	Key      string `json:"key" php:"key"`
	Value    string `json:"value" php:"value"`
}

type KvFunc func(kv *Kv, prefix string) bool
type EncodeFunc func(kv *Kv) ([]byte, error)
type DecodeFunc func([]byte) (*Kv, error)

func NewKv(event int32, key, value string, rev int64) *Kv {
	k := &Kv{
		Event:    event,
		Key:      key,
		Value:    value,
		Revision: rev,
	}
	return k
}

func (e *Kv) String() string {
	return fmt.Sprintf("key : %s, value : %s, event : %d , rev : %d \n", e.Key, e.Value, int(e.Event), int(e.Revision))
}

func (e *Kv) Encode(f EncodeFunc) ([]byte, error) {
	b, err := f(e)
	return b, err
}

func Decode(data []byte, f DecodeFunc) (*Kv, error) {
	kv, err := f(data)
	return kv, err
}
