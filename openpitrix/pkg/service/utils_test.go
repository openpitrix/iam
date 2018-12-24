// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package service

import (
	"testing"

	"openpitrix.io/iam/openpitrix/pkg/pb"
)

func TestGetStructNonZeroFiledNames(t *testing.T) {
	var names, values = pkgGetTableFiledNamesAndValues(pb.Role{
		RoleId:   " ",
		RoleName: "chai",
	})

	tAssert(t, len(names) == 2)

	tAssert(t, names[0] == "role_id")
	tAssert(t, names[1] == "role_name")

	tAssert(t, values[0].(string) == "")
	tAssert(t, values[1].(string) == "chai")
}
