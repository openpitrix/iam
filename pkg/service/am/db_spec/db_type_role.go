// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db_spec

import (
	"time"

	"github.com/golang/protobuf/ptypes"

	pbam "openpitrix.io/iam/pkg/pb/am"
)

type DBRole struct {
	RoleId      string    `db:"role_id"`
	RoleName    string    `db:"role_name"`
	Description string    `db:"description"`
	Portal      string    `db:"portal"`
	Owner       string    `db:"owner"`
	OwnerPath   string    `db:"owner_path"`
	CreateTime  time.Time `db:"create_time"`
	UpdateTime  time.Time `db:"update_time"`
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

func (p *DBRole) ValidateForInsert() error {
	return nil
}
func (p *DBRole) ValidateForUpdate() error {
	return nil
}
