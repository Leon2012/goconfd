package kv

import "encoding/json"

func JsonEncode(kv *Kv) ([]byte, error) {
	b, err := json.Marshal(kv)
	return b, err
}

func JsonDecode(data []byte) (*Kv, error) {
	var k Kv
	err := json.Unmarshal(data, &k)
	if err != nil {
		return nil, err
	} else {
		return &k, nil
	}
}
