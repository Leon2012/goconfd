package monitor

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/Leon2012/goconfd/libs/net2"
)

type Options struct {
	Logger            *log.Logger
	RpcAddress        string `flag:"rpc_address"`
	DBUrl             string `flag:"db_url" cfg:"db_url"`
	DBName            string `flag:"db_name" cfg:"db_name"`
	DBUser            string `flag:"db_user" cfg:"db_user"`
	DBPass            string `flag:"db_pass" cfg:"db_pass"`
	DBTimeout         int    `flag:"db_timeout" cfg:"db_timeout"`
	Hosts             string `flag:"hosts" cfg:"hosts"`
	DialTimeout       int    `flag:"dial_timeout"`
	RequestTimeout    int    `flag:"request_timeout"`
	HeartbeatInterval int    `flag:"heartbeat_interval" cfg:"heartbeat_interval"`
}

func NewOptions() *Options {
	return &Options{
		RpcAddress:        "0.0.0.0:3002",
		Logger:            log.New(os.Stderr, "[monitor]", log.Ldate|log.Ltime|log.Lmicroseconds),
		DBUrl:             "127.0.0.1:27017",
		DBName:            "goconfd",
		DBTimeout:         5,
		DialTimeout:       5,
		RequestTimeout:    5,
		HeartbeatInterval: 5,
	}
}

func (o *Options) Valid() error {
	if o.DBUrl == "" {
		return errors.New("db url is require")
	}
	if o.DBName == "" {
		return errors.New("db name is require")
	}
	return nil
}

func (o *Options) String() string {
	return fmt.Sprintf("rpc_address:%s, db_url:%s, db_name:%s \n", o.RpcAddress, o.DBUrl, o.DBName)

}

func (o *Options) ParseHosts() []string {
	return net2.ParseHosts(o.Hosts)
}
