package local

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/Leon2012/goconfd/libs/kv"
	"github.com/Leon2012/goconfd/libs/util"
)

const (
	FILE_SAVE_TYPE_JSON = iota
	FILE_SAVE_TYPE_PHP
)

type FileSaver struct {
	sync.RWMutex
	Path       string
	Ext        string
	EncodeFunc kv.EncodeFunc
	DecodeFunc kv.DecodeFunc
}

func NewFileSaver(path string, ext string) (*FileSaver, error) {
	err := util.IsWritable(path)
	if err != nil {
		return nil, err
	}
	s := &FileSaver{
		Path: path,
	}
	if !strings.HasPrefix(".", ext) {
		s.Ext = "." + ext
	}
	if ext == "json" {
		s.EncodeFunc = kv.JsonEncode
		s.DecodeFunc = kv.JsonDecode
	} else if ext == "php" {
		s.EncodeFunc = kv.PhpEncode
		s.DecodeFunc = kv.PhpDecode
	}
	return s, nil
}

func (s *FileSaver) Save(k *kv.Kv) error {
	s.Lock()
	defer s.Unlock()
	var err error
	key := k.Key
	newKey := s.safeKey(key)
	var fullFilePath string
	lastIndex := strings.LastIndex(newKey, "/")
	if lastIndex != -1 {
		path := newKey[0:(lastIndex + 1)]
		file := newKey[(lastIndex + 1):len(newKey)]
		fullPath := filepath.Join(s.Path, path)
		if _, err = os.Stat(fullPath); os.IsNotExist(err) {
			err = os.MkdirAll(fullPath, 0777)
			if err != nil {
				return err
			}
		}
		fullFilePath = filepath.Join(fullPath, file)
	} else {
		fullFilePath = filepath.Join(s.Path, newKey)
	}
	fullFilePath += s.Ext
	//fmt.Println(fullFilePath)
	if k.Event == 0 { //KV_EVENT_PUT
		err = s.save(fullFilePath, k)
	} else if k.Event == 1 { //KV_EVENT_DELETE
		err = s.delete(fullFilePath, k)
	} else { //KV_EVENT_NONE
		err = s.save(fullFilePath, k)
	}
	return err
}

func (s *FileSaver) Get(k string) (*kv.Kv, error) {
	s.Lock()
	defer s.Unlock()
	return s.get(k)
}

func (s *FileSaver) Keys() []string {
	return nil
}

func (s *FileSaver) save(fullFilePath string, k *kv.Kv) error {
	data, err := k.Encode(s.EncodeFunc)
	if err != nil {
		return err
	}
	f, err := os.Create(fullFilePath)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func (s *FileSaver) delete(fullFilePath string, k *kv.Kv) error {
	if util.IsExist(fullFilePath) {
		return os.Remove(fullFilePath)
	} else {
		return errors.New("key: " + k.Key + " not exist in local")
	}
}

func (s *FileSaver) get(k string) (*kv.Kv, error) {
	fullFilePath := filepath.Join(s.Path, k)
	f, err := os.Open(fullFilePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}
	fileSize := fi.Size()
	buf := make([]byte, fileSize)
	_, err = f.Read(buf)
	if err != nil {
		return nil, err
	}
	nkv, err := kv.Decode(buf, s.DecodeFunc)
	if err != nil {
		return nil, err
	}
	return nkv, nil
}

func (s *FileSaver) safeKey(key string) string {
	var newKey string
	newKey = util.HexKey(key)
	return newKey
}
