// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db

import (
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
	ExtraJson   string `db:"extra"` // JSON map

	CreateTime string `db:"create_time"`
	UpdateTime string `db:"update_time"`
	StatusTime string `db:"status_time"`
}

type Group struct {
	Gid       string `db:"gid"`
	GidParent string `db:"gid_parent"`

	Name        string `db:"name"`
	Email       string `db:"email"`
	Description string `db:"description"`
	Status      string `db:"status"`
	ExtraJson   string `db:"extra"` // JSON map

	CreateTime string `db:"create_time"`
	UpdateTime string `db:"update_time"`
	StatusTime string `db:"status_time"`
}

func NewUserFrom(src *pbim.User) *User {
	var dst = new(User)
	copyutil.MustDeepCopy(dst, src)
	dst.ExtraJson = string(jsonutil.Encode(src.Extra))
	return dst
}

func (src *User) ToPbUser() *pbim.User {
	var dst = new(pbim.User)
	copyutil.MustDeepCopy(dst, src)
	dst.Extra = make(map[string]string)
	jsonutil.Decode([]byte(src.ExtraJson), dst.Extra)
	return dst
}

func NewGroupFrom(src *pbim.Group) *Group {
	var dst = new(Group)
	copyutil.MustDeepCopy(dst, src)
	dst.ExtraJson = string(jsonutil.Encode(src.Extra))
	return dst
}

func (src *Group) ToPbGroup() *pbim.Group {
	var dst = new(pbim.Group)
	copyutil.MustDeepCopy(dst, src)
	dst.Extra = make(map[string]string)
	jsonutil.Decode([]byte(src.ExtraJson), dst.Extra)
	return dst
}
