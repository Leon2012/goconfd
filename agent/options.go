package agent

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/Leon2012/goconfd/libs/net2"
	"github.com/Leon2012/goconfd/libs/util"
)

type Options struct {
	Logger            *log.Logger
	HttpAddress       string `flag:"http_address"`
	KeyPrefix         string `flag:"key_prefix" cfg:"key_prefix"`
	Hosts             string `flag:"hosts" cfg:"hosts"`
	DialTimeout       int    `flag:"dial_timeout"`
	RequestTimeout    int    `flag:"request_timeout"`
	SavePath          string `flag:"save_path" cfg:"save_path"`
	SaveType          int    `flag:"save_type" cfg:"save_type"`
	FileExt           string `flag:"file_ext" cfg:"file_ext"`
	AutoLoad          bool   `flag:"auto_load"`
	HeartbeatInterval int    `flag:"heartbeat_interval" cfg:"heartbeat_interval"`
}

func NewOptions() *Options {
	return &Options{
		HttpAddress:       "0.0.0.0:3001",
		Logger:            log.New(os.Stderr, "[agent]", log.Ldate|log.Ltime|log.Lmicroseconds),
		DialTimeout:       5,
		RequestTimeout:    5,
		HeartbeatInterval: 5,
	}
}

func (o *Options) String() string {
	return fmt.Sprintf("http_address:%s, key_prefix:%s, etcd_hosts:%s, dial_timeout:%d, request_timeout:%d, save_path:%s \n", o.HttpAddress, o.KeyPrefix, o.Hosts, o.DialTimeout, o.RequestTimeout, o.SavePath)
}

func (o *Options) ParseHosts() []string {
	return net2.ParseHosts(o.Hosts)
}

func (o *Options) Valid() error {
	if o.KeyPrefix == "" {
		return errors.New("key prefix is require")
	}
	if o.Hosts == "" {
		return errors.New("hosts is require")
	}
	if o.SavePath == "" {
		return errors.New("save path is require")
	}
	exist := util.IsExist(o.SavePath)
	if !exist {
		return errors.New("save path is not exist")
	}
	return nil
}
