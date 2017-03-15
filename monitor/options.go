package monitor

import (
	"errors"
	"fmt"
	"log"
	"os"
)

type Options struct {
	Logger     *log.Logger
	RpcAddress string `flag:"rpc_address"`
	DBUrl      string `flag:"db_url" cfg:"db_url"`
	DBName     string `flag:"db_name" cfg:"db_name"`
	DBUser     string `flag:"db_user" cfg:"db_user"`
	DBPass     string `flag:"db_pass" cfg:"db_pass"`
	DBTimeout  int    `flag:"db_timeout" cfg:"db_timeout"`
}

func NewOptions() *Options {
	return &Options{
		RpcAddress: "0.0.0.0:3002",
		Logger:     log.New(os.Stderr, "[monitor]", log.Ldate|log.Ltime|log.Lmicroseconds),
		DBUrl:      "127.0.0.1:27017",
		DBName:     "goconfd",
		DBTimeout:  5,
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
