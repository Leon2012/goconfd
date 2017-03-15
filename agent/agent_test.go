package agent

import (
	"testing"
)

func TestAgent(t *testing.T) {
	opts := NewOptions()
	opts.Hosts = "localhost:2379"
	opts.KeyPrefix = "/develop/activity/"
	opts.SavePath = "/home/vagrant"

	agent := NewAgent(opts)
	defer agent.Exit()

	agent.Main()

}
