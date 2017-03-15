package util

/*
#include <stdlib.h>
#include <stdio.h>
char **EMPTY = NULL;
*/
import "C" //此行和上面的注释之间不能有空行，否则会报错
import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"regexp"
	"unsafe"
)

/**
 * @brief str_hash Hash algorithm of processing a md5 string.
 *
 * @param str The md5 string.
 *
 * @return The number less than 1024.
 */
func str_hash(str string) int {
	b := []byte(str)
	c := b[0:3]

	cc := C.CString(string(c))
	defer C.free(unsafe.Pointer(cc))

	d := C.strtol(cc, C.EMPTY, 16)
	d = d / 4
	return int(d)
}

func GenMd5Str(data []byte) string {
	fmt.Println("Begin to Caculate MD5...")
	m := md5.New()
	m.Write(data)
	return hex.EncodeToString(m.Sum(nil))
}

func IsMd5(str string) bool {
	regular := `^([0-9a-zA-Z]){32}$`
	regx := regexp.MustCompile(regular)
	return regx.MatchString(str)
}

func HashKey(key string) string {
	data := []byte(key)
	md5Sum := GenMd5Str(data)
	lvl1 := str_hash(md5Sum)
	lvl2 := str_hash(string(md5Sum[3:]))
	return fmt.Sprintf("%d/%d/%s", lvl1, lvl2, md5Sum)
}

func HexKey(key string) string {
	src := []byte(key)
	dst := make([]byte, hex.EncodedLen(len(src)))
	hex.Encode(dst, src)
	return string(dst)
}

func UnHexKey(key string) (string, error) {
	decoded, err := hex.DecodeString(key)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}
