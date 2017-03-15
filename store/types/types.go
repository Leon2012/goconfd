package types

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type ObjHeader struct {
	Id        bson.ObjectId `bson:"_id,omitempty"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Agent struct {
	ObjHeader     `bson:",inline"`
	HostName      string
	KeyPrefix     string
	IpAddress     string
	Port          int
	Status        int //0-离线, 1-上线,
	HeartbeatTime time.Time
}

type Heartbeat struct {
	ObjHeader   `bson:",inline"`
	HostName    string
	KeyPrefix   string
	LatestKey   string
	LatestValue string
	LatestTime  time.Time
}
