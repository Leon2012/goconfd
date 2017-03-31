package sdk

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Leon2012/goconfd/libs/kv"
	"github.com/Leon2012/goconfd/registry"
	"github.com/Leon2012/goconfd/registry/frontend"
)

type Goconfd struct {
	local    registry.Frontend
	agentUrl string
}

func NewGoconfd(shmPath string) (*Goconfd, error) {
	g := &Goconfd{}
	g.agentUrl = "http://127.0.0.1:3001/"
	local, err := frontend.NewShmSaver(shmPath)
	if err != nil {
		return nil, err
	}
	g.local = local
	return g, nil
}

func (g *Goconfd) SetAgentUrl(url string) {
	g.agentUrl = url
}

func (g *Goconfd) Get(key string) (*kv.Kv, error) {
	return g.local.Get(key)
}

func (g *Goconfd) GetFromAgent(key string) (*kv.Kv, error) {
	url := fmt.Sprintf("%s/get/%s", g.agentUrl, key)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	nkv, err := kv.Decode(data, kv.JsonDecode)
	if err != nil {
		return nil, err
	}
	return nkv, nil
}
