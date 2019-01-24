// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package id

import (
	"crypto/rand"

	"openpitrix.io/iam/pkg/internal/base58"
)

func GenId(prefix string, maxLen int) string {
	if prefix == "" {
		prefix = "xid-"
	}
	if maxLen <= 0 {
		maxLen = 12
	}

	if maxLen <= len(prefix) {
		maxLen += len(prefix)
	}

	buf := make([]byte, maxLen-len(prefix))
	rand.Read(buf)
	s := string(base58.EncodeBase58(buf))
	return prefix + s[:len(buf)]
}
