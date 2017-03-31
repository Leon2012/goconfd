package dashboard

import (
	"fmt"
	"net"
	"os"

	"sync"

	"github.com/Leon2012/goconfd/libs/util"
	"github.com/Leon2012/goconfd/libs/version"
	"github.com/Leon2012/goconfd/registry"
	"github.com/Leon2012/goconfd/registry/backend"
	"github.com/Leon2012/goconfd/store/db"
	"github.com/Leon2012/goconfd/store/db/mongo"
)

type Dashboard struct {
	sync.RWMutex
	opts         *Options
	httpListener net.Listener
	waitGroup    util.WaitGroupWrapper
	dbConfig     mongo.MongoConfig
	db           db.Adapter
	idc          registry.Backend
}

func NewDashboard(o *Options) *Dashboard {
	d := &Dashboard{
		opts: o,
	}
	d.logf(version.String("dashboard"))
	return d
}

func (a *Dashboard) logf(f string, args ...interface{}) {
	if a.opts.Logger == nil {
		return
	}
	a.opts.Logger.Output(2, fmt.Sprintf(f, args...))
}

func (a *Dashboard) Main() {
	ctx := &Context{
		Dashboard: a,
	}
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
	err := a.db.Open(a.dbConfig)
	if err != nil {
		a.logf("FATAL: open db failed - %s", err)
		os.Exit(1)
	}

	//初始化etcd client
	if len(a.opts.Hosts) == 0 {
		a.logf("FATAL: etcd host empty")
		os.Exit(1)
	}
	cli, err := backend.NewEtcdAdpater(a.opts.ParseHosts(), a.opts.DialTimeout, a.opts.RequestTimeout)
	if err != nil {
		a.logf("FATAL: create etcd client failed - %s", err.Error())
		os.Exit(1)
	}
	a.Lock()
	a.idc = cli
	a.Unlock()

	httpListener, err := net.Listen("tcp", a.opts.HttpAddress)
	if err != nil {
		a.logf("FATAL: listen (%s) failed - %s", a.opts.HttpAddress, err)
		os.Exit(1)
	}
	a.Lock()
	a.httpListener = httpListener
	a.Unlock()
	httpServer := newHttpServer(ctx, a.opts.TemplatePath)
	a.waitGroup.Wrap(func() {
		httpServer.serve(a.httpListener)
	})
}

func (a *Dashboard) Exit() {
	if a.db != nil {
		a.db.Close()
	}
	if a.httpListener != nil {
		a.httpListener.Close()
	}
	a.waitGroup.Wait()
}
