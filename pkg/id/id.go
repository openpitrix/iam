// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package id

import (
	"crypto/rand"
	"strings"

	"openpitrix.io/iam/pkg/internal/base58"
)

const DefaultMaxLength = 16

func GenId(prefix string) string {
	return GenIdWithMaxLen(prefix, DefaultMaxLength)
}

func GenIdWithMaxLen(prefix string, maxLen int) string {
	if prefix == "" {
		prefix = "xid-"
	}
	if maxLen <= 0 {
		maxLen = DefaultMaxLength
	}

	if maxLen <= len(prefix) {
		maxLen += len(prefix)
	}

	buf := make([]byte, maxLen-len(prefix))
	rand.Read(buf)
	s := strings.ToLower(string(base58.EncodeBase58(buf)))
	return prefix + s[:len(buf)]
}
