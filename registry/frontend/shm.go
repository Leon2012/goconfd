package frontend

import (
	"sync"

	"path/filepath"

	"errors"

	"github.com/Leon2012/goconfd/libs/kv"
	"github.com/Leon2012/goconfd/libs/shm"
	"github.com/Leon2012/goconfd/libs/util"
)

const DEFAULT_SHM_PATH = "/dev/shm"

type ShmSaver struct {
	sync.RWMutex
	EncodeFunc kv.EncodeFunc
	DecodeFunc kv.DecodeFunc
	shmPath    string
}

func NewShmSaver(shmPath string) (*ShmSaver, error) {
	s := &ShmSaver{}
	s.EncodeFunc = kv.JsonEncode
	s.DecodeFunc = kv.JsonDecode
	s.shmPath = shmPath
	return s, nil
}

func (s *ShmSaver) Save(k *kv.Kv) error {
	data, err := k.Encode(s.EncodeFunc)
	if err != nil {
		return err
	}
	key := k.Key
	hexKey := util.HexKey(key)
	fileName := filepath.Join(s.shmPath, hexKey)
	shmid, err := shm.Open(fileName)
	if err != nil {
		return err
	}
	err = shm.Write(shmid, string(data))
	if err != nil {
		return err
	}
	return nil
}

func (s *ShmSaver) Get(k string) (*kv.Kv, error) {
	hexKey := util.HexKey(k)
	fileName := filepath.Join(s.shmPath, hexKey)
	data := shm.Read(fileName)
	if data == "" {
		return nil, errors.New("read faile")
	}
	nkv, err := kv.Decode([]byte(data), s.DecodeFunc)
	if err != nil {
		return nil, err
	}
	return nkv, nil
}

func (s *ShmSaver) Keys() []string {
	return nil
}
