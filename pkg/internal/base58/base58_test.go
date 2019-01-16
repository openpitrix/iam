// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package base58

import (
	"testing"
)

type base58Test struct {
	en, de string
}

var base58Golden = []base58Test{
	{"", ""},
	{"\x61", "2g"},
	{"\x62\x62\x62", "a3gV"},
	{"\x63\x63\x63", "aPEr"},
	{"\x73\x69\x6d\x70\x6c\x79\x20\x61\x20\x6c\x6f\x6e\x67\x20\x73\x74\x72\x69\x6e\x67", "2cFupjhnEsSn59qHXstmK2ffpLv2"},
	{"\x00\xeb\x15\x23\x1d\xfc\xeb\x60\x92\x58\x86\xb6\x7d\x06\x52\x99\x92\x59\x15\xae\xb1\x72\xc0\x66\x47", "1NS17iag9jJgTHD1VXjvLCEnZuQ3rJDE9L"},
	{"\x51\x6b\x6f\xcd\x0f", "ABnLTmg"},
	{"\xbf\x4f\x89\x00\x1e\x67\x02\x74\xdd", "3SEo3LWLoPntC"},
	{"\x57\x2e\x47\x94", "3EFU7m"},
	{"\xec\xac\x89\xca\xd9\x39\x23\xc0\x23\x21", "EJDM8drfXA6uyA"},
	{"\x10\xc8\x51\x1e", "Rt5zm"},
	{"\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00", "1111111111"},
}

func TestEncodeBase58(t *testing.T) {
	for _, g := range base58Golden {
		s := string(EncodeBase58([]byte(g.en)))
		if s != g.de {
			t.Errorf("Bad EncodeBase58. Need=%v, Got=%v", g.de, s)
		}
	}
}

func TestDecodeBase58(t *testing.T) {
	for _, g := range base58Golden {
		s := string(DecodeBase58([]byte(g.de)))
		if s != g.en {
			t.Errorf("Bad DecodeBase58. Need=%v, Got=%v", g.en, s)
		}
	}
}

func TestBase58Check(t *testing.T) {
	ba := []byte("Bitcoin")
	ba = EncodeBase58Check(ba)
	if !DecodeBase58Check(ba) {
		t.Errorf("TestBase58Check. Got=%v", ba)
	}
}
