package monitor

import (
	"net"
	"net/rpc"

	"github.com/Leon2012/goconfd/libs/protocol"
	"github.com/Leon2012/goconfd/store/types"
)

type MonitorRpc struct {
	ctx *Context
}

func NewRpcServer(c *Context) *MonitorRpc {
	r := &MonitorRpc{
		ctx: c,
	}
	return r
}

func (r *MonitorRpc) Ping(args *protocol.NoArg, reply *protocol.NoReply) error {
	return nil
}

func (r *MonitorRpc) Online(args *protocol.OnlineArg, reply *protocol.Ack) error {
	if args.HostName == "" || args.KeyPrefix == "" {
		reply.Code = 1001
		reply.Message = "hostname or keyprefix is empty"
		return nil
	}

	agent := &types.Agent{}
	agent.HostName = args.HostName
	agent.KeyPrefix = args.KeyPrefix
	agent.IpAddress = args.IpAddress
	r.ctx.Monitor.logf("INFO: call Online func, hostName-keyPrefix(%s-%s)", agent.HostName, agent.KeyPrefix)
	err := r.ctx.Monitor.db.Online(agent)
	if err != nil {
		return err
	}
	reply.Code = 1000
	reply.Message = ""
	return nil
}

func (r *MonitorRpc) Offline(args *protocol.OfflineArg, reply *protocol.Ack) error {
	if args.HostName == "" || args.KeyPrefix == "" {
		reply.Code = 1001
		reply.Message = "hostname or keyprefix is empty"
		return nil
	}
	agent := &types.Agent{}
	agent.HostName = args.HostName
	agent.KeyPrefix = args.KeyPrefix
	r.ctx.Monitor.logf("INFO: call Online func, hostName-keyPrefix(%s-%s)", agent.HostName, agent.KeyPrefix)
	err := r.ctx.Monitor.db.Offline(agent)
	if err != nil {
		return err
	}
	reply.Code = 1000
	reply.Message = ""
	return nil
}

func (r *MonitorRpc) Heartbeat(args *protocol.HeartbeatArg, reply *protocol.HeartbeatReply) error {
	r.ctx.Monitor.logf("INFO: call Heartbeat func, hostName-keyPrefix(%s-%s)", args.HostName, args.KeyPrefix)

	heartbeat := &types.Heartbeat{}
	heartbeat.HostName = args.HostName
	heartbeat.KeyPrefix = args.KeyPrefix
	heartbeat.LatestKey = args.Key
	heartbeat.LatestValue = args.Value
	heartbeat.LatestTime = args.Time
	err := r.ctx.Monitor.db.Heartbeat(heartbeat)
	if err != nil {
		r.ctx.Monitor.logf("ERROR: call heartbeat faile - %s", err)
		return err
	}
	return nil
}

func (r *MonitorRpc) serve(lis net.Listener) error {
	rpc.Register(r)
	go func() {
		for {
			conn, err := lis.Accept()
			if err != nil {
				continue
			}
			go rpc.ServeConn(conn)
		}
	}()

	return nil
}
