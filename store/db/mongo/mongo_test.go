package mongo

import (
	"os"
	"testing"
	"time"

	"github.com/Leon2012/goconfd/store/types"
)

var config MongoConfig
var adapter *MongoAdapter

func init() {
	config = MongoConfig{
		Url:      "127.0.0.1:27017",
		DbName:   "goconfd",
		Timeout:  1,
		Username: "",
		Password: "",
	}
	adapter = NewMongoAdapter()
}

func TestOnline(t *testing.T) {
	hostName, _ := os.Hostname()
	a := &types.Agent{
		HostName:  hostName,
		KeyPrefix: "developer.activity",
		IpAddress: "127.0.0.1",
		Port:      30,
	}
	adapter.Open(config)
	defer adapter.Close()
	err := adapter.Online(a)
	if err != nil {
		t.Error(err)
	}
}

func TestOffline(t *testing.T) {
	hostName, _ := os.Hostname()
	a := &types.Agent{
		HostName: hostName,
	}
	adapter.Open(config)
	defer adapter.Close()
	err := adapter.Offline(a)
	if err != nil {
		t.Error(err)
	}
}

func TestHeartbeat(t *testing.T) {
	hostName, _ := os.Hostname()
	l := &types.Heartbeat{
		HostName:    hostName,
		KeyPrefix:   "developer.activity",
		LatestKey:   "a",
		LatestValue: "b",
		LatestTime:  time.Now(),
	}
	adapter.Open(config)
	defer adapter.Close()
	err := adapter.Heartbeat(l)
	if err != nil {
		t.Error(err)
	}
}

func TestGetAgents(t *testing.T) {
	adapter.Open(config)
	defer adapter.Close()
	agents, err := adapter.GetAgents()
	if err != nil {
		t.Error(err)
	} else {
		t.Log(agents)
	}
}
