package dashboard

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/Leon2012/goconfd/libs/net2"
)

type Options struct {
	Logger         *log.Logger
	HttpAddress    string `flag:"http_address"`
	DBUrl          string `flag:"db_url" cfg:"db_url"`
	DBName         string `flag:"db_name" cfg:"db_name"`
	DBUser         string `flag:"db_user" cfg:"db_user"`
	DBPass         string `flag:"db_pass" cfg:"db_pass"`
	DBTimeout      int    `flag:"db_timeout" cfg:"db_timeout"`
	Hosts          string `flag:"hosts" cfg:"hosts"`
	DialTimeout    int    `flag:"dial_timeout"`
	RequestTimeout int    `flag:"request_timeout"`
	TemplatePath   string `flag:"template_path" cfg:"template_path"`
}

func NewOptions() *Options {
	return &Options{
		HttpAddress:    "0.0.0.0:3003",
		Logger:         log.New(os.Stderr, "[dashboard]", log.Ldate|log.Ltime|log.Lmicroseconds),
		DBUrl:          "127.0.0.1:27017",
		DBName:         "goconfd",
		DBTimeout:      5,
		DialTimeout:    5,
		RequestTimeout: 5,
		Hosts:          "localhost:2379",
		TemplatePath:   "./",
	}
}

func (o *Options) String() string {
	return fmt.Sprintf("http_address:%s, etcd_hosts:%s, db_url:%s \n", o.HttpAddress, o.Hosts, o.DBUrl)
}

func (o *Options) Valid() error {
	if o.DBUrl == "" {
		return errors.New("db url is require")
	}
	if o.DBName == "" {
		return errors.New("db name is require")
	}
	if o.Hosts == "" {
		return errors.New("etcd hosts is require")
	}
	return nil
}

func (o *Options) ParseHosts() []string {
	return net2.ParseHosts(o.Hosts)
}
