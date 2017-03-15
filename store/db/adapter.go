package db

import (
	"time"

	"github.com/Leon2012/goconfd/store/types"
)

type Adapter interface {
	Open(c interface{}) error
	Close() error
	IsOpen() bool
	Online(agent *types.Agent) error
	Offline(agent *types.Agent) error
	Heartbeat(pack *types.Heartbeat) error
	GetAgents() ([]*types.Agent, error)
	GetHeartbeats() ([]*types.Heartbeat, error)
	GetHeartbeatsByAgent(agent *types.Agent) ([]*types.Heartbeat, error)
	GetHeartbeatsByTime(time time.Time) ([]*types.Heartbeat, error)
}
