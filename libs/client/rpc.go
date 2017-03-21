package client

import (
	"errors"
	"net"
	"net/rpc"
)

type RPCClient struct {
	client   *rpc.Client
	conn     *net.TCPConn
	addrs    []string
	selector Selector
}

func NewRPCClient(sel Selector) *RPCClient {
	if sel == nil {
		sel = &DefaultSelector{}
	}
	return &RPCClient{
		selector: sel,
	}
}

func (c *RPCClient) SetAddrs(addrs []string) {
	c.addrs = addrs
}

func (c *RPCClient) Select() string {
	s, err := c.selector.Select(c.addrs)
	if err != nil {
		return ""
	}
	return s
}

func (c *RPCClient) Open() error {
	addr := c.Select()
	if addr == "" {
		return errors.New("no addr select")
	}
	address, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return err
	}
	conn, err := net.DialTCP("tcp", nil, address)
	if err != nil {
		return err
	}
	c.conn = conn
	client := rpc.NewClient(conn)
	c.client = client
	return nil
}

func (c *RPCClient) Call(method string, args interface{}, reply interface{}) error {
	err := c.Open()
	if err != nil {
		return err
	}
	defer c.Close()
	err = c.client.Call(method, args, reply)
	return err
}

func (c *RPCClient) Close() {
	if c.conn != nil {
		c.conn.Close()
		c.conn = nil
	}
	if c.client != nil {
		c.client.Close()

		c.client = nil
	}
}
