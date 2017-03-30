package shm

import (
	"fmt"
	_ "os"
	"testing"
)

var SHM_FILE = "/dev/shm/confd6"

func TestWrite(t *testing.T) {
	id, err := Open(SHM_FILE)
	if err != nil {
		t.Error(err)
	}
	//defer Close(id)
	t.Log("shm id" + fmt.Sprintf("%d", id))
	err = Write(id, "111111111hello world11111")
	if err != nil {
		t.Error(err)
	}
}

func TestDel(t *testing.T) {
	Del(SHM_FILE)

}

func TestRead(t *testing.T) {
	str := Read(SHM_FILE)
	t.Log("result:" + str)
}
