// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db

import (
	"reflect"
	"testing"

	. "openpitrix.io/iam/pkg/internal/assert"
)

func _TestModuleApiInfo(t *testing.T) {
	u0 := &ModuleApiInfo{
		RoleId:   "role-001",
		RoleName: "role-name",
		Portal:   "portal-xxx",

		ModuleId:   "mod-001",
		ModuleName: "mod-name",
		DataLevel:  "data-level-001",
		//IsFeatureAllChecked: "1",

		FeatureId:   "feature-001",
		FeatureName: "feature-name",

		ActionId:   "action-001",
		ActionName: "action-name",
		//ActionEnabled: "true",

		ApiId:          "api-001",
		ApiMethod:      "api.method/001",
		ApiDescription: "api-001-desc",

		Url:       "/api/method/001",
		UrlMethod: "GET",
	}

	pbU0, err := u0.ToProtoMessage()
	Assert(t, err == nil, err)
	Assert(t, pbU0.RoleId == "role-001")

	u1 := NewModuleApiInfoFromPB(pbU0)

	pbU1, err := u1.ToProtoMessage()
	Assert(t, err == nil, err)
	u2 := NewModuleApiInfoFromPB(pbU1)

	Assertf(t, reflect.DeepEqual(u1, u0),
		"u1 = %v, u0 = %v", u1, u0,
	)
	Assertf(t, reflect.DeepEqual(u1, u2),
		"u1 = %v, u2 = %v", u1, u2,
	)
}
