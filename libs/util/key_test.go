package util

import (
	"encoding/hex"
	"fmt"
	"hash/crc32"
	"testing"
)

func TestEncryptKey(t *testing.T) {
	key := "/dachu/zhuanpan99"
	key1 := HashKey(key)
	t.Log("key:" + key1)
}

func TestCrc32Key(t *testing.T) {
	var castagnoliTable = crc32.MakeTable(crc32.Castagnoli)
	key := "/dachu/zhuanpan99"
	crc := crc32.New(castagnoliTable)
	crc.Write([]byte(key))

	fmt.Printf("Sum32 : %x \n", crc.Sum32())
}

func TestHexKey(t *testing.T) {
	key := "/dachu/zhuanpan99"
	src := []byte(key)
	dst := make([]byte, hex.EncodedLen(len(src)))
	hex.Encode(dst, src)
	fmt.Printf("%s\n", dst)
}
