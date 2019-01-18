// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db

import (
	"time"

	"github.com/golang/protobuf/ptypes"

	pbam "openpitrix.io/iam/pkg/pb/am"
)

type DBRole struct {
	RoleId      string `gorm:"primary_key"`
	RoleName    string
	Description string
	Portal      string
	Owner       string
	OwnerPath   string
	CreateTime  time.Time
	UpdateTime  time.Time
}

func PBRoleToDB(p *pbam.Role) *DBRole {
	if p == nil {
		return new(DBRole)
	}
	var q = &DBRole{
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

func (p *DBRole) ToPB() *pbam.Role {
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
