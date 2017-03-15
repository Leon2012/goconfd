package agent

import (
	"testing"

	"github.com/Leon2012/goconfd/libs/client"
	"github.com/Leon2012/goconfd/libs/protocol"
)

var rpcAddress = "0.0.0.0:3002"
var addrs []string
var rpcClient *client.RPCClient

func init() {
	addrs = []string{}
	addrs = append(addrs, rpcAddress)
}

func TestOnline(t *testing.T) {
	var err error
	getClient()
	err = rpcClient.Open()
	if err != nil {
		t.Error(err)
		return
	}
	args := &protocol.OnlineArg{}
	args.HostName = "hostname"
	args.IpAddress = "127.0.0.1"
	args.KeyPrefix = "developer.usergroup"
	var reply protocol.Ack

	err = rpcClient.Call("MonitorRpc.Online", args, &reply)
	if err != nil {
		t.Error(err)
		return
	}

	if reply.Code == 1000 {
		t.Log("online success")
	} else {
		t.Log(reply.Message)
	}

	rpcClient.Close()
}

func getClient() {
	s := &client.DefaultSelector{}
	c := client.NewRPCClient(s)
	c.SetAddrs(addrs)
	rpcClient = c
}
