// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

// Package copyutil provides deep copy for value type.
package copyutil

import (
	"bytes"
	"encoding/gob"
)

func DeepCopy(dst, src interface{}) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(src); err != nil {
		return err
	}
	return gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(dst)
}

func MustDeepCopy(dst, src interface{}) {
	if err := DeepCopy(dst, src); err != nil {
		panic(err)
	}
}
