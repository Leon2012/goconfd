package protocol

import "time"

type Ack struct {
	Code    int
	Message string
}

type NoArg struct {
}

type NoReply struct {
}

type OnlineArg struct {
	HostName  string
	KeyPrefix string
	IpAddress string
}

type OnlineReply struct {
	Status bool
}

type OfflineArg struct {
	HostName  string
	KeyPrefix string
}

type OfflineReply struct {
}

type HeartbeatArg struct {
	HostName  string
	KeyPrefix string
	Key       string
	Value     string
	Time      time.Time
}
type HeartbeatReply struct {
}
