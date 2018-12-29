// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

// Package jsonutil provides json helper functions.
package jsonutil

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
)

func Encode(m interface{}) []byte {
	data, err := json.MarshalIndent(m, "", "\t")
	if err != nil {
		return nil
	}
	data = bytes.Replace(data, []byte("\n"), []byte("\r\n"), -1)
	return data
}

func Decode(data []byte, m interface{}) error {
	return json.Unmarshal(data, m)
}

func Load(filename string, m interface{}) (err error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, m)
	if err != nil {
		return
	}
	return
}

func Save(filename string, m interface{}) (err error) {
	if err := ioutil.WriteFile(filename, Encode(m), 0666); err != nil {
		return err
	}
	return nil
}
