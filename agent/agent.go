package agent

import (
	"fmt"
	"net"
	"os"
	"sync"
	"time"

	"github.com/Leon2012/goconfd/agent/idc"
	"github.com/Leon2012/goconfd/agent/local"
	"github.com/Leon2012/goconfd/libs/kv"
	"github.com/Leon2012/goconfd/libs/net2"
	"github.com/Leon2012/goconfd/libs/util"
	"github.com/Leon2012/goconfd/libs/version"
	_ "github.com/coreos/etcd/clientv3"
)

const MONITOR_RPC_HOST_KEY = "system.monitor.rpc.host"

type KvPair struct {
	KeyPrefix string
	Kv        *kv.Kv
}

type Heartbeat struct {
	Kv         *kv.Kv
	UpdateTime time.Time
}

type Agent struct {
	sync.RWMutex
	opts          *Options
	httpListener  net.Listener
	waitGroup     util.WaitGroupWrapper
	idc           IdcInterface
	local         LocalInterface
	keyCache      map[string]string
	monitorKv     *kv.Kv
	lastHeartbeat *Heartbeat
	watchKVChan   chan *kv.Kv
	hostName      string
	localIP       string
	rpcClient     *RpcClient
	exitChan      chan bool
}

func NewAgent(opts *Options) *Agent {
	a := &Agent{
		opts:          opts,
		keyCache:      make(map[string]string),
		watchKVChan:   make(chan *kv.Kv),
		lastHeartbeat: nil,
		rpcClient:     nil,
		exitChan:      make(chan bool),
	}
	a.logf(version.String("agent"))
	return a
}

func (a *Agent) setLastHeartbeat(kv *kv.Kv) {
	if a.lastHeartbeat == nil {
		a.lastHeartbeat = &Heartbeat{}
	}
	a.lastHeartbeat.Kv = kv
	a.lastHeartbeat.UpdateTime = time.Now()
}

func (a *Agent) logf(f string, args ...interface{}) {
	if a.opts.Logger == nil {
		return
	}
	a.opts.Logger.Output(2, fmt.Sprintf(f, args...))
}

func (a *Agent) Main() {
	ctx := &Context{a}

	hostName, err := util.GetHostName()
	if err != nil {
		a.logf("FATAL: get host name faile")
		os.Exit(1)
	}

	localIP, err := net2.GetLocalIPv4()
	if err != nil {
		a.logf("FATAL: get local ipv4 faile")
		os.Exit(1)
	}

	a.Lock()
	a.hostName = hostName
	a.localIP = localIP
	a.Unlock()

	//初始化etcd client
	if len(a.opts.Hosts) == 0 {
		a.logf("FATAL: etcd host empty")
		os.Exit(1)
	}
	hosts := a.opts.ParseHosts()
	cli, err := idc.NewEtcdAdpater(hosts, a.opts.DialTimeout, a.opts.RequestTimeout)
	if err != nil {
		a.logf("FATAL: create etcd client failed - %s, hosts : %s", err.Error(), hosts)
		os.Exit(1)
	}
	a.Lock()
	a.idc = cli
	a.Unlock()

	var l LocalInterface
	if a.opts.SaveType == 1 {
		l, err = local.NewFileSaver(a.opts.SavePath, a.opts.FileExt)
	} else if a.opts.SaveType == 2 {
		l, err = local.NewShmSaver(a.opts.SavePath)
	}
	if err != nil {
		a.logf("FATAL: create local save failed - %s", err.Error())
		os.Exit(1)
	}
	a.Lock()
	a.local = l
	a.Unlock()

	httpListener, err := net.Listen("tcp", a.opts.HttpAddress)
	if err != nil {
		a.logf("FATAL: listen (%s) failed - %s", a.opts.HttpAddress, err)
		os.Exit(1)
	}
	a.Lock()
	a.httpListener = httpListener
	a.Unlock()
	httpServer := newHttpServer(ctx)
	a.waitGroup.Wrap(func() {
		httpServer.serve(a.httpListener)
	})

	go SaveKV(ctx)

	//go LoadValuesByPrefix(ctx)
	if a.opts.AutoLoad {
		UpdateLocalVales(ctx)
	}
	//获取monitor host
	mrhKv, err := LoadValueByKey(ctx, MONITOR_RPC_HOST_KEY)
	if err == nil {
		a.monitorKv = mrhKv
	} else {
		a.monitorKv = nil
	}
	go WatchValuesByKey(ctx, MONITOR_RPC_HOST_KEY)

	rpcClient := NewRpcClient(ctx)
	a.Lock()
	a.rpcClient = rpcClient
	a.Unlock()
	a.rpcClient.ReloadAddrs()
	err = a.rpcClient.Online()
	if err != nil {
		a.logf("agent call online faile : %s", err.Error())
	}

	go WatchValuesByPrefix(ctx)
	go rpcClient.Heartbeat(a.opts.HeartbeatInterval)
}

func (a *Agent) Exit() {
	if a.httpListener != nil {
		a.httpListener.Close()
	}
	if a.idc != nil {
		a.idc.Close()
	}
	err := a.rpcClient.Offline()
	if err != nil {
		a.logf("agent call offline faile : %s", err.Error())
	}
	close(a.exitChan)
	a.waitGroup.Wait()
}
