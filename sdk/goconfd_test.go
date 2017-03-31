package sdk

import (
	"testing"
)

func TestGet(t *testing.T) {
	gconfd, err := NewGoconfd("/dev/shm")
	if err != nil {
		t.Error(err)
	}
	k, err := gconfd.Get("develop.activity.k6")
	if err != nil {
		t.Error(err)
	}
	t.Log(k.String())
}

func TestGetFromAgent(t *testing.T) {
	gconfd, err := NewGoconfd("/dev/shm")
	if err != nil {
		t.Error(err)
	}
	k, err := gconfd.Get("develop.activity.k7")
	if err != nil {
		t.Error(err)
	}
	t.Log(k.Value)
}
