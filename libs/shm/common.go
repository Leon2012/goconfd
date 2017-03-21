package shm

/*
#include <stdlib.h>
#include <sys/types.h>
#include <sys/ipc.h>
#include <sys/msg.h>
#include <sys/shm.h>
key_t ftok(const char *pathname, int proj_id);
*/
import "C"
import (
	"errors"
	"unsafe"
)

func Ftok(pathname string, projID uint8) (int64, error) {
	if projID == 0 {
		return -1, errors.New("sysvipc: projID must be nonzero")
	}
	cpath := C.CString(pathname)
	defer C.free(unsafe.Pointer(cpath))
	rckey, err := C.ftok(cpath, C.int(projID))
	rc := int64(rckey)
	if rc == -1 {
		return -1, err
	}
	return rc, nil
}
