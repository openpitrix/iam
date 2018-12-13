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
	_ gorpHooker = (*Group)(nil)
)

type Group struct {
	Gid       string `db:"gid, primarykey"`
	GidParent string `db:"gid_parent"`

	Name        string            `db:"name"`
	Email       string            `db:"email"`
	Description string            `db:"description"`
	Status      string            `db:"status"`
	Extra       map[string]string `db:"-"`     // JSON map
	Extra_XXX   string            `db:"extra"` // JSON map

	CreateTime     *timestamp.Timestamp `db:"-"`
	CreateTime_XXX *time.Time           `db:"create_time"`
	UpdateTime     *timestamp.Timestamp `db:"-"`
	UpdateTime_XXX *time.Time           `db:"update_time"`
	StatusTime     *timestamp.Timestamp `db:"-"`
	StatusTime_XXX *time.Time           `db:"status_time"`
}

func NewGroupFrom(src *pbim.Group) *Group {
	var dst = new(Group)
	copyutil.MustDeepCopy(dst, src)
	return dst
}

func (p *Group) PostGet(s gorp.SqlExecutor) error   { return p.hookPostRead(s) }
func (p *Group) PreInsert(s gorp.SqlExecutor) error { return p.hookPreWrite(s) }
func (p *Group) PreUpdate(s gorp.SqlExecutor) error { return p.hookPreWrite(s) }
func (p *Group) PreDelete(s gorp.SqlExecutor) error { return p.hookPreWrite(s) }

func (p *Group) hookPostRead(s gorp.SqlExecutor) error {
	p.Extra = make(map[string]string)
	jsonutil.Decode([]byte(p.Extra_XXX), p.Extra)
	return nil
}
func (p *Group) hookPreWrite(s gorp.SqlExecutor) error {
	p.Extra_XXX = string(jsonutil.Encode(p.Extra))

	p.CreateTime_XXX = new(time.Time)
	p.UpdateTime_XXX = new(time.Time)
	p.StatusTime_XXX = new(time.Time)

	*p.CreateTime_XXX, _ = ptypes.Timestamp(p.CreateTime)
	*p.UpdateTime_XXX, _ = ptypes.Timestamp(p.UpdateTime)
	*p.StatusTime_XXX, _ = ptypes.Timestamp(p.StatusTime)

	return nil
}

func (src *Group) ToPbGroup() *pbim.Group {
	var dst = new(pbim.Group)
	copyutil.MustDeepCopy(dst, src)
	return dst
}
