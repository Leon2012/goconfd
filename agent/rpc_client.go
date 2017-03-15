package agent

import (
	"strings"
	"time"

	"errors"

	"github.com/Leon2012/goconfd/libs/client"
	"github.com/Leon2012/goconfd/libs/net2"
	"github.com/Leon2012/goconfd/libs/protocol"
)

type RpcClient struct {
	client            *client.RPCClient
	ctx               *Context
	maxErrCnt, errCnt int
}

func NewRpcClient(c *Context) *RpcClient {
	key := c.Agent.hostName + "_" + c.Agent.opts.KeyPrefix
	selector := client.NewRingSelector(key)
	cl := &RpcClient{
		ctx:       c,
		maxErrCnt: 10,
		errCnt:    0,
	}
	cl.client = client.NewRPCClient(selector)
	return cl
}

func (r *RpcClient) Online() error {
	var err error
	args := &protocol.OnlineArg{}
	args.HostName = r.ctx.Agent.hostName
	args.IpAddress = r.ctx.Agent.localIP
	args.KeyPrefix = r.ctx.Agent.opts.KeyPrefix
	var reply protocol.Ack
	err = r.doCall("MonitorRpc.Online", args, reply)
	if err != nil {
		return err
	}
	if reply.Code != 1000 {
		return errors.New(reply.Message)
	}
	return nil
}

func (r *RpcClient) Offline() error {
	var err error
	args := &protocol.OfflineArg{}
	args.HostName = r.ctx.Agent.hostName
	args.KeyPrefix = r.ctx.Agent.opts.KeyPrefix
	var reply protocol.Ack
	err = r.doCall("MonitorRpc.Offline", args, reply)
	if err != nil {
		return err
	}
	if reply.Code != 1000 {
		return errors.New(reply.Message)
	}
	return nil
}

func (r *RpcClient) Heartbeat(interval int) {
	var pp time.Duration
	pp = ((time.Duration(interval) * time.Second) * 9) / 10
	ticker := time.NewTicker(pp) //心跳包
	defer func() {
		ticker.Stop()
	}()
	for {
		select {
		case <-r.ctx.Agent.exitChan:
			r.ctx.Agent.logf("agent exit!!!!!!")
			return
		case <-ticker.C: //心跳
			r.doHeartbeat()
		}
	}
}

func (r *RpcClient) ReloadAddrs() bool {
	if r.ctx.Agent.monitorKv == nil {
		return false
	}
	addrStr := r.ctx.Agent.monitorKv.Value
	if addrStr == "" {
		return false
	}
	addrs := []string{}
	if strings.Index(addrStr, ";") != -1 {
		addrArr := strings.Split(addrStr, ";")
		for _, addr := range addrArr {
			addr = strings.TrimSpace(addr)
			if net2.CheckIp(addr) {
				addrs = append(addrs, addr)
			}
		}
	} else {
		addrs = append(addrs, addrStr)
	}
	r.client.SetAddrs(addrs)
	r.ctx.Agent.logf("set monitor rpc address : %s", addrs)
	return true
}

func (r *RpcClient) doHeartbeat() {
	r.ctx.Agent.logf("call doHeartbeat func!")
	args := &protocol.HeartbeatArg{}
	args.HostName = r.ctx.Agent.hostName
	args.KeyPrefix = r.ctx.Agent.opts.KeyPrefix
	if r.ctx.Agent.lastHeartbeat != nil {
		if r.ctx.Agent.lastHeartbeat.Kv != nil {
			args.Key = r.ctx.Agent.lastHeartbeat.Kv.Key
			args.Value = r.ctx.Agent.lastHeartbeat.Kv.Value
		}
		args.Time = r.ctx.Agent.lastHeartbeat.UpdateTime
	} else {
		args.Key = ""
		args.Value = ""
	}
	var reply protocol.HeartbeatReply
	err := r.doCall("MonitorRpc.Heartbeat", args, reply)
	if err != nil {
		r.ctx.Agent.logf("doHeartbeat error : %s", err.Error())
		return
	}
}

func (r *RpcClient) doCall(method string, args, reply interface{}) error {
	var err error
	if r.errCnt > r.maxErrCnt {
		return errors.New("error cnt more than the max error cnt, require ignore!")
	}
	err = r.client.Open()
	if err != nil {
		r.errCnt++
		r.ctx.Agent.logf("rpc call %s error : %s", method, err.Error())
		return err
	}
	defer r.client.Close()
	err = r.client.Call(method, args, &reply)
	if err != nil {
		r.errCnt++
		return err
	}
	return nil
}
