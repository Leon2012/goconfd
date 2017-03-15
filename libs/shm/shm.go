package shm

// #include <stdlib.h>
// #include <string.h>
// #include <sys/shm.h>
// #include <sys/types.h>
/*
int my_shm_open(char* filename, int open_flag){
    int shm_id;
    key_t key;
    key = ftok(filename, 0x01);
    if(key == -1){
        return -1;
    }
    if(open_flag)
        shm_id = shmget(key, 4096, IPC_CREAT|IPC_EXCL|0666);
    else
        shm_id = shmget(key, 0, 0);
    if(shm_id == -1){
        return -1;
    }
    return shm_id;
}

int my_shm_update(int shm_id, char* content){
    char* addr;
    addr = (char*)shmat(shm_id, NULL, 0);
    if(addr == (char*)-1){
        return -1;
    }
    if(strlen(content) > 4095)
        return -1;
    strcpy(addr, content);
    shmdt(addr);
    return 0;
}

int my_shm_del(char* filename){
    int shm_id;
    char* addr;
    char* s;
    shm_id = my_shm_open(filename, 0);
    if(shm_id == -1){
        return -1;
    }
    shmctl(shm_id, IPC_RMID, NULL);
    return 0;
}

char* my_shm_read(char* filename){
    int shm_id;
    char* addr;
    char* s;
    shm_id = my_shm_open(filename, 0);
    if(shm_id == -1)
        return NULL;
    addr = (char*)shmat(shm_id, NULL, 0);
    if(addr == (char*)-1){
        return NULL;
    }
    s = (char*)malloc(strlen(addr) + 1);
    strcpy(s, addr);
    shmdt(addr);
    return s;
}
*/
import "C"

import "unsafe"

import (
	_ "log"
	"os"
	_ "time"
)

type errorString struct {
	s string
}

func (e *errorString) Error() string {
	return e.s
}

func Open(filename string) (int, error) {
	//filename := filepath.Join("/tmp", file)
	fp, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return 0, err
	}
	defer fp.Close()
	f := C.CString(filename)
	defer C.free(unsafe.Pointer(f))
	r := int(C.my_shm_open(f, C.int(1)))
	if r == -1 {
		return 0, &errorString{"Open error"}
	}
	return r, nil
}

func Write(shm_id int, content string) error {
	c := C.CString(content)
	defer C.free(unsafe.Pointer(c))
	r := int(C.my_shm_update(C.int(shm_id), c))
	if r == -1 {
		return &errorString{"Write error"}
	}
	return nil
}

func Read(filename string) string {
	//filename := filepath.Join("/tmp", file)
	f := C.CString(filename)
	defer C.free(unsafe.Pointer(f))
	s := C.my_shm_read(f)
	defer C.free(unsafe.Pointer(s))
	return C.GoString(s)
}

func Del(filename string) error {
	//C.my_shm_del(C.int(shm_id))
	//filename := filepath.Join("/tmp", file)
	f := C.CString(filename)
	defer C.free(unsafe.Pointer(f))
	C.my_shm_del(f)
	err := os.Remove(filename)
	return err
}
