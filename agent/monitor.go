package agent

import (
	"net"

	"fmt"

	"sync"

	"github.com/Leon2012/goconfd/libs/kv"
	"github.com/Leon2012/goconfd/libs/node"
)

const (
	defaultSystemMonitorNodeKey = "/system/monitor/nodes/"
)

type Monitor struct {
	nodes []*node.Node
	ctx   *Context
	sync.Mutex
}

func NewMonitor(ctx *Context) *Monitor {
	return &Monitor{
		ctx:   ctx,
		nodes: []*node.Node{},
	}
}

func (m *Monitor) WatchNodes() {
	m.ctx.Agent.idc.WatchWithPrefix(defaultSystemMonitorNodeKey, func(k *kv.Kv, prefix string) bool {
		nd, err := node.NewNode(k.Value)
		if err == nil {
			if k.Event == kv.KV_EVENT_PUT {
				m.AddNode(nd)
			} else if k.Event == kv.KV_EVENT_DELETE {
				m.DelNode(nd)
			}
			m.ctx.Agent.rpcClient.ReloadAddrs()
			return true
		} else {
			return false
		}
	})
}

func (m *Monitor) AddNode(node *node.Node) {
	m.Lock()
	for idx, nd := range m.nodes {
		if nd.Name == node.Name {
			m.nodes = append(m.nodes[:idx], m.nodes[idx+1:]...)
			break
		}
	}
	m.nodes = append(m.nodes, node)
	m.Unlock()
}

func (m *Monitor) DelNode(node *node.Node) {
	m.Lock()
	for idx, nd := range m.nodes {
		if nd.Name == node.Name {
			m.nodes = append(m.nodes[:idx], m.nodes[idx+1:]...)
			break
		}
	}
	m.Unlock()
}

func (m *Monitor) UpdateNode(node *node.Node) {
	m.Lock()
	for _, nd := range m.nodes {
		if nd.Name == node.Name {
			nd.CPU = node.CPU
			nd.IP = node.IP
			nd.ServiceAddress = node.ServiceAddress
			break
		}
	}
	m.Unlock()
}

func (m *Monitor) GetNodesAddr() []string {
	addrs := []string{}
	for _, node := range m.nodes {
		host, port, err := net.SplitHostPort(node.ServiceAddress)
		if err == nil {
			if host == "0.0.0.0" {
				host = node.IP
			}
			addr := fmt.Sprintf("%s:%s", host, port)
			addrs = append(addrs, addr)
		}
	}
	return addrs
}
