// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db

import (
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"gopkg.in/gorp.v2"

	"openpitrix.io/iam/pkg/internal/copyutil"
	"openpitrix.io/iam/pkg/internal/jsonutil"
	"openpitrix.io/iam/pkg/pb/im"
)

var (
	_ gorpHooker = (*User)(nil)
)

type User struct {
	Uid string `db:"uid, primarykey"`
	Gid string `db:"gid"`

	Name        string            `db:"name"`
	Email       string            `db:"email"`
	Description string            `db:"description"`
	Password    string            `db:"password"`
	Status      string            `db:"status"`
	Extra       map[string]string `db:"-"`     // JSON map
	Extra_XXX   string            `db:"extra"` // JSON map

	CreateTime     *timestamp.Timestamp `db:"-"`
	CreateTime_XXX time.Time            `db:"create_time"`
	UpdateTime     *timestamp.Timestamp `db:"-"`
	UpdateTime_XXX time.Time            `db:"update_time"`
	StatusTime     *timestamp.Timestamp `db:"-"`
	StatusTime_XXX time.Time            `db:"status_time"`
}

func NewUserFrom(src *pbim.User) *User {
	var dst = new(User)
	copyutil.MustDeepCopy(dst, src)
	return dst
}

func (p *User) PostGet(s gorp.SqlExecutor) error   { return p.hookPostRead(s) }
func (p *User) PreInsert(s gorp.SqlExecutor) error { return p.hookPreWrite(s) }
func (p *User) PreUpdate(s gorp.SqlExecutor) error { return p.hookPreWrite(s) }
func (p *User) PreDelete(s gorp.SqlExecutor) error { return p.hookPreWrite(s) }

func (p *User) hookPostRead(s gorp.SqlExecutor) error {
	p.Extra = make(map[string]string)
	jsonutil.Decode([]byte(p.Extra_XXX), p.Extra)
	return nil
}
func (p *User) hookPreWrite(s gorp.SqlExecutor) error {
	p.Extra_XXX = string(jsonutil.Encode(p.Extra))

	p.CreateTime_XXX, _ = ptypes.Timestamp(p.CreateTime)
	p.UpdateTime_XXX, _ = ptypes.Timestamp(p.UpdateTime)
	p.StatusTime_XXX, _ = ptypes.Timestamp(p.StatusTime)

	return nil
}

func (src *User) ToPbUser() *pbim.User {
	var dst = new(pbim.User)
	copyutil.MustDeepCopy(dst, src)
	return dst
}
