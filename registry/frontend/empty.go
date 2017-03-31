package frontend

import (
	"fmt"
	"github.com/Leon2012/goconfd/libs/kv"
)

type Empty struct{}

func (e *Empty) Save(k *kv.Kv) error {
	//fmt.Printf("key : %s, value : %s \n", k.Key, k.Value)
	fmt.Println(k)
	return nil
}

func (e *Empty) Get(k string) (*kv.Kv, error) {
	fmt.Printf("get key : %s \n", k)
	return nil, nil
}
