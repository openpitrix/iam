// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db

import (
	"time"

	"github.com/golang/protobuf/ptypes"

	"openpitrix.io/iam/pkg/internal/copyutil"
	"openpitrix.io/iam/pkg/internal/jsonutil"
	"openpitrix.io/iam/pkg/pb/im"
)

type User struct {
	Uid string `db:"uid"`
	Gid string `db:"gid"`

	Name        string `db:"name"`
	Email       string `db:"email"`
	Description string `db:"description"`
	Password    string `db:"password"`
	Status      string `db:"status"`
	XXX_Extra   string `db:"extra"` // JSON map

	XXX_CreateTime time.Time `db:"create_time"`
	XXX_UpdateTime time.Time `db:"update_time"`
	XXX_StatusTime time.Time `db:"status_time"`
}

type Group struct {
	Gid       string `db:"gid"`
	GidParent string `db:"gid_parent"`

	Name        string `db:"name"`
	Email       string `db:"email"`
	Description string `db:"description"`
	Status      string `db:"status"`
	XXX_Extra   string `db:"extra"` // JSON map

	XXX_CreateTime time.Time `db:"create_time"`
	XXX_UpdateTime time.Time `db:"update_time"`
	XXX_StatusTime time.Time `db:"status_time"`
}

func NewUserFrom(src *pbim.User) *User {
	var dst = new(User)

	copyutil.MustDeepCopy(dst, src)

	dst.XXX_Extra = string(jsonutil.Encode(src.Extra))

	dst.XXX_CreateTime, _ = ptypes.Timestamp(src.CreateTime)
	dst.XXX_UpdateTime, _ = ptypes.Timestamp(src.UpdateTime)
	dst.XXX_StatusTime, _ = ptypes.Timestamp(src.StatusTime)

	return dst
}

func (src *User) ToPbUser() *pbim.User {
	var dst = new(pbim.User)

	copyutil.MustDeepCopy(dst, src)

	dst.Extra = make(map[string]string)
	jsonutil.Decode([]byte(src.XXX_Extra), dst.Extra)

	dst.CreateTime, _ = ptypes.TimestampProto(src.XXX_CreateTime)
	dst.UpdateTime, _ = ptypes.TimestampProto(src.XXX_UpdateTime)
	dst.StatusTime, _ = ptypes.TimestampProto(src.XXX_StatusTime)

	return dst
}

func NewGroupFrom(src *pbim.Group) *Group {
	var dst = new(Group)

	copyutil.MustDeepCopy(dst, src)

	dst.XXX_Extra = string(jsonutil.Encode(src.Extra))

	dst.XXX_CreateTime, _ = ptypes.Timestamp(src.CreateTime)
	dst.XXX_UpdateTime, _ = ptypes.Timestamp(src.UpdateTime)
	dst.XXX_StatusTime, _ = ptypes.Timestamp(src.StatusTime)

	return dst
}

func (src *Group) ToPbGroup() *pbim.Group {
	var dst = new(pbim.Group)

	copyutil.MustDeepCopy(dst, src)

	dst.Extra = make(map[string]string)
	jsonutil.Decode([]byte(src.XXX_Extra), dst.Extra)

	dst.CreateTime, _ = ptypes.TimestampProto(src.XXX_CreateTime)
	dst.UpdateTime, _ = ptypes.TimestampProto(src.XXX_UpdateTime)
	dst.StatusTime, _ = ptypes.TimestampProto(src.XXX_StatusTime)

	return dst
}
