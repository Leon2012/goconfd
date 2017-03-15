package idc

import (
	"fmt"
	"testing"
	"time"

	"github.com/Leon2012/goconfd/libs/kv"
	"github.com/coreos/etcd/clientv3"
	"golang.org/x/net/context"
)

const ETCD_HOST = "localhost:2379"
const ETCD_REQUEST_TIMEOUT = 5 * time.Second
const ETCD_DIAL_TIMEOUT = 5 * time.Second

func getEtcdClient() (*clientv3.Client, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{ETCD_HOST},
		DialTimeout: ETCD_DIAL_TIMEOUT,
	})
	return cli, err
}

func TestGet(t *testing.T) {
	cli, err := getEtcdClient()
	if err != nil {
		t.Error(err)
	}
	defer cli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), ETCD_REQUEST_TIMEOUT)
	cancel()

	k := "test1"
	resp, err := cli.Get(ctx, k)
	if err != nil {
		t.Error(err)
	}

	t.Log(resp)
}

func TestPut(t *testing.T) {
	cli, err := getEtcdClient()
	if err != nil {
		t.Error(err)
	}
	defer cli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), ETCD_REQUEST_TIMEOUT)
	defer cancel()

	resp, err := cli.Put(ctx, "foo1", "bar1")
	if err != nil {
		t.Error(err)
	}

	fmt.Println(resp)
}

func TestWatch(t *testing.T) {
	cli, err := getEtcdClient()
	if err != nil {
		t.Error(err)
	}
	defer cli.Close()
	rch := cli.Watch(context.Background(), "/develop/activity/", clientv3.WithPrefix())
	for wresp := range rch {
		for _, ev := range wresp.Events {
			fmt.Printf("%s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
		}
	}
}

func TestBatchGet(t *testing.T) {
	cli, err := getEtcdClient()
	if err != nil {
		t.Error(err)
	}
	defer cli.Close()
	ctx, cancel := context.WithTimeout(context.Background(), ETCD_REQUEST_TIMEOUT)
	resp, err := cli.Get(ctx, "/develop/activity/", clientv3.WithPrefix(), clientv3.WithSort(clientv3.SortByKey, clientv3.SortAscend))
	cancel()
	if err != nil {
		t.Error(err)
	}
	for _, ev := range resp.Kvs {
		fmt.Printf("%s : %s\n", ev.Key, ev.Value)
	}
}

func TestEtcdPut(t *testing.T) {
	hosts := []string{"localhost:2379"}
	etcd, err := NewEtcdAdpater(hosts, 5, 5)
	if err != nil {
		t.Error(err)
	}
	defer etcd.Close()

	err = etcd.Put("foo3", "bar3")
	if err != nil {
		t.Error(err)
	}
}

func TestWatchWithPrefixs(t *testing.T) {
	//hosts := []string{"127.0.0.1:2379", "192.168.174.114:2379"}
	hosts := []string{"localhost:2379"}
	prefixs := []string{"develop.activity", "develop.user"}
	etcd, err := NewEtcdAdpater(hosts, 5, 5)
	if err != nil {
		t.Error(err)
	}
	defer etcd.Close()
	var f kv.KvFunc
	f = func(k *kv.Kv, prefix string) bool {
		str := fmt.Sprintf("watch prefix %s change", prefix)
		fmt.Println(str)
		return true
	}
	fmt.Println(f)
	etcd.WatchWithPrefixs(prefixs, f)
	for {

	}
}

func TestWatchPrefix(t *testing.T) {
	//hosts := []string{"127.0.0.1:2379", "192.168.174.114:2379"}
	hosts := []string{"localhost:2379"}
	prefix := "develop.activity"
	etcd, err := NewEtcdAdpater(hosts, 5, 5)
	if err != nil {
		t.Error(err)
	}
	defer etcd.Close()
	var f kv.KvFunc
	f = func(k *kv.Kv, prefix string) bool {
		str := fmt.Sprintf("watch prefix %s change", prefix)
		fmt.Println(str)
		return true
	}
	fmt.Println(f)
	etcd.WatchWithPrefix(prefix, f)
}
