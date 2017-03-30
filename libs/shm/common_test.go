package shm

import (
	"testing"
)

func TestFtok(t *testing.T) {
	pathname := "/tmp/queue1"
	projID := uint8(0x01)
	key, err := Ftok(pathname, projID)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(key)
	}
}
