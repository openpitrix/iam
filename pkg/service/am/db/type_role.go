// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db

import (
	"time"

	"github.com/golang/protobuf/ptypes"

	pbam "openpitrix.io/iam/pkg/pb/am"
)

type Role struct {
	RoleId      string `gorm:"primary_key"`
	RoleName    string
	Description string
	Portal      string
	CreateTime  time.Time
	UpdateTime  time.Time
	Owner       string
	OwnerPath   string
}

func PBRoleToDB(p *pbam.Role) *Role {
	if p == nil {
		return new(Role)
	}
	var q = &Role{
		RoleId:      p.RoleId,
		RoleName:    p.RoleName,
		Description: p.Description,
		Portal:      p.Portal,
		Owner:       p.Owner,
		OwnerPath:   p.OwnerPath,
	}

	q.CreateTime, _ = ptypes.Timestamp(p.CreateTime)
	q.UpdateTime, _ = ptypes.Timestamp(p.UpdateTime)

	return q
}

func (p *Role) ToPB() *pbam.Role {
	if p == nil {
		return new(pbam.Role)
	}
	var q = &pbam.Role{
		RoleId:      p.RoleId,
		RoleName:    p.RoleName,
		Description: p.Description,
		Portal:      p.Portal,
		Owner:       p.Owner,
		OwnerPath:   p.OwnerPath,
	}

	q.CreateTime, _ = ptypes.TimestampProto(p.CreateTime)
	q.UpdateTime, _ = ptypes.TimestampProto(p.UpdateTime)

	return q
}
