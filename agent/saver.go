package agent

import (
	"github.com/Leon2012/goconfd/libs/kv"
)

func SaveKV(ctx *Context) {
	var kv *kv.Kv
	for {
		select {
		case kv = <-ctx.Agent.watchKVChan:
			err := ctx.Agent.local.Save(kv)
			if err != nil {
				ctx.Agent.logf(err.Error())
			} else {
				ctx.Agent.setLastHeartbeat(kv)
				ctx.Agent.logf("key:" + kv.Key + " save success")
			}
		}
	}
}
