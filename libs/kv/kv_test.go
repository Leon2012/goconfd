package kv

import (
	"testing"
)

func TestJsonEncode(t *testing.T) {
	kv := &Kv{
		Revision: 0,
		Event:    0,
		Key:      "dachu/dachu99_actid",
		Value:    "48800",
	}
	b, err := kv.Encode(JsonEncode)
	if err != nil {
		t.Error(err)
	}
	t.Log(string(b))
}

func TestJsonDecode(t *testing.T) {
	json := "{\"rev\":0,\"type\":0,\"key\":\"dachu/dachu99_actid\",\"value\":\"48800\"}"
	data := []byte(json)
	kv, err := Decode(data, JsonDecode)
	if err != nil {
		t.Error(err)
	}
	t.Log(kv.String())
}

func TestPhpEncode(t *testing.T) {
	kv := &Kv{
		Revision: 0,
		Event:    0,
		Key:      "dachu/dachu99_actid",
		Value:    "48800",
	}
	b, err := kv.Encode(PhpEncode)
	if err != nil {
		t.Error(err)
	}
	t.Log(string(b))
}

func TestPhpDecode(t *testing.T) {
	php := `<?php
		$_64616368752f646163687539395f6163746964_rev='0';
		$_64616368752f646163687539395f6163746964_type='0';
		$_64616368752f646163687539395f6163746964_key='dachu/dachu99_actid';
		$_64616368752f646163687539395f6163746964_value='48800';`

	data := []byte(php)
	kv, err := Decode(data, PhpDecode)
	if err != nil {
		t.Error(err)
	}
	t.Log(kv.String())

}
