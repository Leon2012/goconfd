package idc

import (
	"fmt"
	"time"

	"sync"

	"github.com/Leon2012/goconfd/libs/kv"
	"github.com/coreos/etcd/clientv3"
	"golang.org/x/net/context"
)

type prefixChan struct {
	p string
	c clientv3.WatchChan
}

type EtcdAdpater struct {
	client         *clientv3.Client
	hosts          []string
	requestTimeout time.Duration
	dialTimeout    time.Duration
	wg             sync.WaitGroup
}

func NewEtcdAdpater(hosts []string, request, dial int) (*EtcdAdpater, error) {
	etcd := &EtcdAdpater{
		hosts:          hosts,
		requestTimeout: time.Duration(request) * time.Second,
		dialTimeout:    time.Duration(dial) * time.Second,
	}
	err := etcd.connect()
	if err != nil {
		return nil, err
	}
	return etcd, nil
}

func (e *EtcdAdpater) connect() error {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   e.hosts,
		DialTimeout: e.dialTimeout,
	})
	if err != nil {
		return err
	}
	e.client = cli
	return nil
}

func (e *EtcdAdpater) Put(k, v string) error {
	ctx, cancel := context.WithTimeout(context.Background(), e.requestTimeout)
	_, err := e.client.Put(ctx, k, v)
	cancel()
	if err != nil {
		return err
	}
	return nil
}
func (e *EtcdAdpater) Get(k string) ([]*kv.Kv, error) {
	fmt.Println("get key:" + k)
	kvs := []*kv.Kv{}
	ctx, cancel := context.WithTimeout(context.Background(), e.requestTimeout)
	resp, err := e.client.Get(ctx, k)
	cancel()
	if err != nil {
		return kvs, err
	}
	revision := resp.Header.Revision
	for _, ev := range resp.Kvs {
		kv := kv.NewKv(0, string(ev.Key), string(ev.Value), revision)
		kvs = append(kvs, kv)
	}
	return kvs, nil
}
func (e *EtcdAdpater) BatchGetByPrefix(prefix string) ([]*kv.Kv, error) {
	kvs := []*kv.Kv{}
	ctx, cancel := context.WithTimeout(context.Background(), e.requestTimeout)
	resp, err := e.client.Get(ctx, prefix, clientv3.WithPrefix(), clientv3.WithSort(clientv3.SortByKey, clientv3.SortAscend))
	cancel()
	if err != nil {
		return kvs, err
	}
	revision := resp.Header.Revision
	for _, ev := range resp.Kvs {
		kv := kv.NewKv(0, string(ev.Key), string(ev.Value), revision)
		kvs = append(kvs, kv)
	}
	return kvs, nil
}

func (e *EtcdAdpater) WatchWithKey(key string, f kv.KvFunc) {
	rch := e.client.Watch(context.Background(), key)
	for wresp := range rch {
		revision := wresp.Header.Revision
		for _, ev := range wresp.Events {
			kv := kv.NewKv(int32(ev.Type), string(ev.Kv.Key), string(ev.Kv.Value), revision)
			f(kv, key)
		}
	}
}

func (e *EtcdAdpater) WatchWithPrefix(prefix string, f kv.KvFunc) {
	rch := e.client.Watch(context.Background(), prefix, clientv3.WithPrefix())
	for wresp := range rch {
		revision := wresp.Header.Revision
		for _, ev := range wresp.Events {
			kv := kv.NewKv(int32(ev.Type), string(ev.Kv.Key), string(ev.Kv.Value), revision)
			f(kv, prefix)
		}
	}
}

func (e *EtcdAdpater) WatchWithKeys(keys []string, f kv.KvFunc) {
	pcs := []*prefixChan{}
	for _, key := range keys {
		rch := e.client.Watch(context.Background(), key)
		pc := &prefixChan{
			p: key,
			c: rch,
		}
		pcs = append(pcs, pc)
	}
	e.processChans(pcs, f)
}

func (e *EtcdAdpater) WatchWithPrefixs(prefixs []string, f kv.KvFunc) {
	pcs := []*prefixChan{}
	for _, prefix := range prefixs {
		rch := e.client.Watch(context.Background(), prefix, clientv3.WithPrefix())
		pc := &prefixChan{
			p: prefix,
			c: rch,
		}
		pcs = append(pcs, pc)
	}
	e.processChans(pcs, f)
}

func (e *EtcdAdpater) processChans(pcs []*prefixChan, f kv.KvFunc) {
	n := len(pcs)
	e.wg.Add(n)
	for i, pc := range pcs {
		p := pc.p
		c := pc.c
		go e.processChan(i, c, p, f)
	}
}

func (e *EtcdAdpater) processChan(i int, c clientv3.WatchChan, p string, f kv.KvFunc) {
	for wresp := range c {
		revision := wresp.Header.Revision
		for _, ev := range wresp.Events {
			kv := kv.NewKv(int32(ev.Type), string(ev.Kv.Key), string(ev.Kv.Value), revision)
			f(kv, p)
		}
	}
	e.wg.Done()
}

func (e *EtcdAdpater) Close() {
	if e.client != nil {
		e.client.Close()
	}
	e.wg.Wait()
}
