package monitor

import (
	"fmt"
	"runtime"
	"time"

	"errors"

	"github.com/Leon2012/goconfd/libs/net2"
	"github.com/Leon2012/goconfd/libs/node"
	"github.com/Leon2012/goconfd/registry"
	"github.com/Leon2012/goconfd/registry/backend"
)

const (
	defaultSystemMonitorNodeKey = "/system/monitor/nodes/"
)

type Service struct {
	ctx               *Context
	info              *node.Node
	heartbeatInterval int64
	ttlMax            int
	ttlFailCnt        int
	heartbeatTicker   *time.Ticker
	backend           registry.Backend
	exitChan          chan int
	exit              bool
}

func NewService(ctx *Context) (*Service, error) {
	s := &Service{
		ctx: ctx,
		info: &node.Node{
			Name:           ctx.Monitor.hostName,
			ServiceAddress: ctx.Monitor.opts.RpcAddress,
			CPU:            runtime.NumCPU(),
			IP:             net2.LocalIPString(),
		},
		heartbeatInterval: int64(ctx.Monitor.opts.HeartbeatInterval),
		ttlMax:            3,
		ttlFailCnt:        0,
		exitChan:          make(chan int),
		exit:              false,
	}
	hosts := ctx.Monitor.opts.ParseHosts()
	cli, err := backend.NewEtcdAdpater(hosts, ctx.Monitor.opts.DialTimeout, ctx.Monitor.opts.RequestTimeout)
	if err != nil {
		return nil, err
	}
	s.backend = cli
	return s, nil
}

func (a *Service) key() string {
	return fmt.Sprintf("%s/%s/%s", defaultSystemMonitorNodeKey, a.ctx.Monitor.hostName, a.ctx.Monitor.opts.RpcAddress)
}

//心跳
func (a *Service) doHeartbeat() {
	key := a.key()
	err := a.backend.SetTTL(key)
	if err != nil {
		a.ttlFailCnt++
	} else {
		a.ttlFailCnt = 0
	}
	if a.ttlFailCnt >= a.ttlMax {
		a.Deregister()
	}
}

func (a *Service) Heartbeat() {
	d := time.Duration(a.heartbeatInterval) * time.Second
	a.heartbeatTicker = time.NewTicker(d)
	defer func() {
		a.heartbeatTicker.Stop()
	}()
	for {
		select {
		case <-a.exitChan:
			return
		case <-a.heartbeatTicker.C:
			a.doHeartbeat()
		}
	}
}

//服务注册
func (a *Service) Register() error {
	key := a.key()
	err := a.backend.PutWithTTL(key, a.info.String(), (a.heartbeatInterval * int64(a.ttlMax)))
	return err
}

//服务注销
func (a *Service) Deregister() error {
	if !a.exit {
		a.exit = true
		key := a.key()
		err := a.backend.Del(key)
		close(a.exitChan)
		return err
	}
	return errors.New("service is exit")
}

func (a *Service) WatchNodes() {

}
