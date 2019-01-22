// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db

import (
	"testing"
)

func TestIsValidIds(t *testing.T) {
	Assert(t, isValidIds(
		genId("uid-", 12),
		genId("gid-", 12),
		genId("xid-", 12),
	))
}

func TestIsValidEmails(t *testing.T) {
	Assert(t, isValidEmails(
		"admin@openpitrix.io",
		"dev@openpitrix.io",
		"user@openpitrix.io",
		"1@qq.com",
	))
	Assert(t, !isValidEmails(
		"123", "aaa",
	))
}

func TestIsValidPhoneNumbers(t *testing.T) {
	Assert(t, isValidPhoneNumbers(
		"123",
		"110",
		"123456",
	))
	Assert(t, !isValidPhoneNumbers(
		"admin@openpitrix.io",
		"dev@openpitrix.io",
		"user@openpitrix.io",
		"1@qq.com",
		"1 3",
		" 1 ",
	))
}

func TestSimplifyStringList(t *testing.T) {
	s0 := []string{"a", "", "c", "  ", " d "}
	s1 := simplifyStringList(s0)

	Assert(t, len(s1) == 3)
	Assert(t, s1[0] == "a")
	Assert(t, s1[1] == "c")
	Assert(t, s1[2] == "d")
}

func TestSimplifyString(t *testing.T) {
	Assert(t, simplifyString("\ta  b  c") == "a b c")
}
