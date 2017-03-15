package local

import (
	"github.com/Leon2012/goconfd/libs/kv"
	"strings"
	"testing"
)

func TestSplitKey(t *testing.T) {
	key := "/develop/activity/dachu99/actid2"
	lastIndex := strings.LastIndex(key, "/")
	path := key[0:(lastIndex + 1)]
	file := key[(lastIndex + 1):len(key)]

	t.Log("file:" + file)
	t.Log("path:" + path)
}

func TestKeySave(t *testing.T) {
	k := kv.NewKv(1, "/develop/activity/dachu99/actid2", "40086", 1)
	s, err := NewFileSaver("/home/vagrant")
	if err != nil {
		t.Error(err)
	}
	err = s.Save(k)
	if err != nil {
		t.Error(err)
	}
}

func TestKeyGet(t *testing.T) {
	key := "/develop/activity/dachu99/actid2"
	s, err := NewFileSaver("/home/vagrant")
	if err != nil {
		t.Error(err)
	}
	kv, err := s.Get(key)
	if err != nil {
		t.Error(err)
	}
	t.Log(kv)
}
