package kv

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"regexp"
	"strconv"
	"strings"

	"github.com/Leon2012/goconfd/libs/reflect"
	"github.com/Leon2012/goconfd/libs/util"
)

const PHP_CODE_PREFIX = "<?php"
const PHP_VAR_PREFIX = "$_"
const PHP_VAR_SUFFIX = ";"
const PHP_CODE_LINE_BREAK = "\n"

func PhpEncode(kv *Kv) ([]byte, error) {
	var buffer bytes.Buffer
	buffer.WriteString(PHP_CODE_PREFIX)
	buffer.WriteString(PHP_CODE_LINE_BREAK)

	revTag := reflect.GetTag(kv, "Revision", "php")
	revStr := key2phpVar(kv, revTag, strconv.FormatInt(kv.Revision, 10))
	buffer.WriteString(revStr)

	eventTag := reflect.GetTag(kv, "Event", "php")
	eventStr := key2phpVar(kv, eventTag, strconv.Itoa(int(kv.Event)))
	buffer.WriteString(eventStr)

	keyTag := reflect.GetTag(kv, "Key", "php")
	keyStr := key2phpVar(kv, keyTag, kv.Key)
	buffer.WriteString(keyStr)

	valTag := reflect.GetTag(kv, "Value", "php")
	valStr := key2phpVar(kv, valTag, kv.Value)
	buffer.WriteString(valStr)

	return buffer.Bytes(), nil
}

func PhpDecode(data []byte) (*Kv, error) {
	var err error
	str := string(data)
	strArr := strings.Split(str, PHP_CODE_LINE_BREAK)
	if strArr[0] != PHP_CODE_PREFIX {
		return nil, errors.New("no php code")
	}
	kv := &Kv{}
	for i := 1; i < len(strArr); i++ {
		err = nil
		line := strArr[i]
		err = phpVar2Key(kv, line)
		if err != nil {
			break
		}
	}
	if err != nil {
		return nil, err
	} else {
		return kv, nil
	}
}

func key2phpVar(kv *Kv, key, val string) string {
	var s, prefix string
	prefix = PHP_VAR_PREFIX + util.HexKey(kv.Key) + "_"
	s = prefix + key + "=" + "'" + val + "'" + PHP_VAR_SUFFIX + PHP_CODE_LINE_BREAK
	return s
}

func phpVar2Key(kv *Kv, line string) error {
	r, err := regexp.Compile("\\${1}(_([a-z0-9]+)){2}=\\'([a-zA-Z0-9\\/\\._]+)\\';")
	if err != nil {
		return err
	}
	strs := r.FindAllStringSubmatch(line, -1)
	if len(strs) == 0 {
		return errors.New("parse php error")
	}

	revTag := reflect.GetTag(kv, "Revision", "php")
	eventTag := reflect.GetTag(kv, "Event", "php")
	keyTag := reflect.GetTag(kv, "Key", "php")
	valTag := reflect.GetTag(kv, "Value", "php")

	strss := strs[0]
	if len(strss) == 4 {
		key := strss[2]
		val := strss[3]
		if key == revTag {
			i, _ := strconv.ParseInt(val, 10, 64)
			kv.Revision = i
		} else if key == eventTag {
			b, _ := strconv.Atoi(val)
			kv.Event = int32(b)
		} else if key == keyTag {
			kv.Key = val
		} else if key == valTag {
			kv.Value = val
		}
		return nil
	} else {
		return errors.New("parse php error")
	}
}

func md5Sum(str string) string {
	data := []byte(str)
	m := md5.New()
	m.Write(data)
	return hex.EncodeToString(m.Sum(nil))
}

func safeVal(val string) string {
	// var newVal string
	// newVal = strings.Replace(val, "\'", "\\'", -1);
	return val
}
