package monitor

import (
	"fmt"
	"os"
	"sync"

	"net"

	"github.com/Leon2012/goconfd/libs/util"
	"github.com/Leon2012/goconfd/libs/version"
	"github.com/Leon2012/goconfd/store/db"
	"github.com/Leon2012/goconfd/store/db/mongo"
)

type Monitor struct {
	sync.RWMutex
	opts        *Options
	tcpListener net.Listener
	waitGroup   util.WaitGroupWrapper
	dbConfig    mongo.MongoConfig
	db          db.Adapter
	hostName    string
	service     *Service
}

func NewMonitor(o *Options) *Monitor {
	m := &Monitor{
		opts: o,
	}
	m.logf(version.String("monitor"))
	return m
}

func (a *Monitor) logf(f string, args ...interface{}) {
	if a.opts.Logger == nil {
		return
	}
	a.opts.Logger.Output(2, fmt.Sprintf(f, args...))
}

func (a *Monitor) Main() {
	ctx := &Context{a}

	hostName, err := util.GetHostName()
	if err != nil {
		a.logf("FATAL: get host name faile")
		os.Exit(1)
	}
	a.Lock()
	a.hostName = hostName
	a.Unlock()

	a.Lock()
	a.dbConfig = mongo.MongoConfig{
		Url:      a.opts.DBUrl,
		DbName:   a.opts.DBName,
		Timeout:  a.opts.DBTimeout,
		Username: a.opts.DBUser,
		Password: a.opts.DBPass,
	}
	a.db = mongo.NewMongoAdapter()
	a.Unlock()
	err = a.db.Open(a.dbConfig)
	if err != nil {
		a.logf("FATAL: open db failed - %s", err)
		os.Exit(1)
	}

	service, err := NewService(ctx)
	if err != nil {
		a.logf("FATAL: create service failed - %s", err.Error())
		os.Exit(1)
	}
	err = service.Register()
	if err != nil {
		a.logf("FATAL: register service failed - %s", err.Error())
		os.Exit(1)
	}
	a.Lock()
	a.service = service
	a.Unlock()
	a.waitGroup.Wrap(func() {
		a.service.Heartbeat()
	})

	tcpListener, err := net.Listen("tcp", a.opts.RpcAddress)
	if err != nil {
		a.logf("FATAL: listen (%s) failed - %s", a.opts.RpcAddress, err)
		os.Exit(1)
	}
	a.Lock()
	a.tcpListener = tcpListener
	a.Unlock()
	rpcServer := NewRpcServer(ctx)
	a.waitGroup.Wrap(func() {
		rpcServer.serve(a.tcpListener)
		a.logf("INFO: rpc server listen(%s) success", a.opts.RpcAddress)
	})
}

func (a *Monitor) Exit() {
	a.service.Deregister()
	if a.db != nil {
		a.db.Close()
	}
	if a.tcpListener != nil {
		a.tcpListener.Close()
	}
	a.waitGroup.Wait()
}
