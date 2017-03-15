package reflect

import (
	"testing"
)

type Kv struct {
	KeyPrefix string `json:"prefix" php:"prefix"`
	Revision  int64  `json:"rev" php:"rev"`
	Event     int32  `json:"type" php:"type"`
	Key       string `json:"key" php:"key"`
	Value     string `json:"value" php:"value"`
}

func TestGetTag(t *testing.T) {
	kv := &Kv{}
	tag := GetTag(kv, "Revision", "php")
	t.Log(tag)
}
