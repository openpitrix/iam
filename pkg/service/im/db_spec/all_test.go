// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db_spec

import (
	"reflect"
	"testing"
	"time"

	idpkg "openpitrix.io/iam/pkg/id"
	. "openpitrix.io/iam/pkg/internal/assert"
	"openpitrix.io/iam/pkg/internal/strutil"
)

func TestUser(t *testing.T) {
	u0 := &User{
		UserId:      idpkg.GenId("fuckid-"),
		UserName:    "fuck-name",
		Email:       "120@qq.com",
		PhoneNumber: "119",
		Description: "Description",
		Password:    "123456",
		Status:      "failed",
		CreateTime:  time.Date(2019, 1, 2, 3, 4, 5, 0, time.UTC),
		UpdateTime:  time.Date(2019, 1, 2, 3, 4, 5, 0, time.UTC),
		StatusTime:  time.Date(2019, 1, 2, 3, 4, 5, 0, time.UTC),
		Extra:       strutil.NewString(`{"abc":"abc-value"}`),
	}

	pbU0, err := u0.ToProtoMessage()
	Assert(t, err == nil, err)
	u1 := NewUserFromPB(pbU0)

	pbU1, err := u1.ToProtoMessage()
	Assert(t, err == nil, err)
	u2 := NewUserFromPB(pbU1)

	Assertf(t, reflect.DeepEqual(u1, u2),
		"u1 = %v, u2 = %v", u1, u2,
	)
}

func TestUserGroup(t *testing.T) {
	g0 := &UserGroup{
		ParentGroupId:  "",
		GroupId:        "gid-0001",
		GroupPath:      "gid-0001.",
		GroupName:      "play-dev",
		Description:    "Description",
		Status:         "ok",
		CreateTime:     time.Date(2019, 1, 2, 3, 4, 5, 0, time.UTC),
		UpdateTime:     time.Date(2019, 1, 2, 3, 4, 5, 0, time.UTC),
		StatusTime:     time.Date(2019, 1, 2, 3, 4, 5, 0, time.UTC),
		Extra:          strutil.NewString(`{"abc":"abc-value"}`),
		GroupPathLevel: 0,
	}

	pbG0, err := g0.ToProtoMessage()
	Assert(t, err == nil, err)
	g1 := NewUserGroupFromPB(pbG0)

	pbG1, err := g1.ToProtoMessage()
	Assert(t, err == nil, err)
	g2 := NewUserGroupFromPB(pbG1)

	Assertf(t, reflect.DeepEqual(g1, g2),
		"g1 = %v, g2 = %v", g1, g2,
	)
}
