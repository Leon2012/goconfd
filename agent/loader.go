package agent

import (
	"errors"
	_ "time"

	"github.com/Leon2012/goconfd/libs/kv"
	"github.com/Leon2012/goconfd/libs/util"
	_ "golang.org/x/net/context"
)

//更新本地数据
func UpdateLocalVales(ctx *Context) error {
	files, err := util.GetFileList(ctx.Agent.opts.SavePath, "php")
	if err != nil {
		return err
	}
	for _, f := range files {
		encodeKey := util.GetName(f)
		decodeKey, err := util.UnHexKey(encodeKey)
		if err != nil {
			ctx.Agent.logf("decode hex key : %s faile", err.Error())
			break
		}
		k, err := LoadValueByKey(ctx, decodeKey)
		if err != nil {
			ctx.Agent.logf("load local key :%s faile", encodeKey)
			break
		} else {
			//ctx.Agent.local.Save(k)
			ctx.Agent.watchKVChan <- k
			ctx.Agent.logf("load local key : %s success", encodeKey)
		}
	}
	return nil
}

//监控数据
func WatchValuesByPrefix(ctx *Context) {
	//var err error
	keyPrefix := ctx.Agent.opts.KeyPrefix
	ctx.Agent.idc.WatchWithPrefix(keyPrefix, func(k *kv.Kv, prefix string) bool {
		// kp := &KvPair{
		// 	KeyPrefix: keyPrefix,
		// 	Kv:        k,
		// }
		ctx.Agent.watchKVChan <- k
		return true
	})
}

//监控多前缀数据
func WatchValuesByPrefixs(ctx *Context, prefixs []string) {
	ctx.Agent.idc.WatchWithPrefixs(prefixs, func(k *kv.Kv, prefix string) bool {
		ctx.Agent.watchKVChan <- k
		return true
	})
}

//监控单key数据
func WatchValuesByKey(ctx *Context, key string) {
	ctx.Agent.idc.WatchWithKey(key, func(k *kv.Kv, prefix string) bool {
		if key == MONITOR_RPC_HOST_KEY {
			ctx.Agent.rpcClient.ReloadAddrs()
		} else {
			ctx.Agent.watchKVChan <- k
		}
		return true
	})
}

//加载远端数据
func LoadValueByKey(ctx *Context, key string) (*kv.Kv, error) {
	var err error
	kvs, err := ctx.Agent.idc.Get(key)
	if err != nil {
		return nil, err
	} else {
		if len(kvs) > 0 {
			k := kvs[0]
			ctx.Agent.watchKVChan <- k
			return k, nil
		} else {
			return nil, errors.New("Not Found")
		}
	}
}
