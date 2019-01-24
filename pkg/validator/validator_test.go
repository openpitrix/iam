// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package validator_test

import (
	"testing"

	idpkg "openpitrix.io/iam/pkg/id"
	. "openpitrix.io/iam/pkg/internal/assert"
	"openpitrix.io/iam/pkg/validator"
)

func TestIsValidId(t *testing.T) {
	Assert(t, validator.IsValidId(
		idpkg.GenId("uid-", 12),
		idpkg.GenId("gid-", 12),
		idpkg.GenId("xid-", 12),
	))
}

func TestIsValidName(t *testing.T) {
	Assert(t, validator.IsValidName(
		"a", "a b",
	))
	Assert(t, !validator.IsValidName(
		" a ", " ", "",
	))
}

func TestIsValidEmail(t *testing.T) {
	Assert(t, validator.IsValidEmail(
		"admin@kubesphere.io",
		"admin@op.com",
		"dev@op.com",
		"user@op.com",
		"kubernetes-admin@kubernetes",
		"admin@kubernetes",
		"admin@openpitrix.io",
		"dev@openpitrix.io",
		"user@openpitrix.io",
		"1@qq.com",
	))

	Assert(t, !validator.IsValidEmail(
		"123", "aaa",
	))
}

func TestIsValidPhoneNumber(t *testing.T) {
	Assert(t, validator.IsValidPhoneNumber(
		"123",
		"110",
		"123456",
	))
	Assert(t, !validator.IsValidPhoneNumber(
		"admin@openpitrix.io",
		"dev@openpitrix.io",
		"user@openpitrix.io",
		"1@qq.com",
		"1 3",
		" 1 ",
	))
}
