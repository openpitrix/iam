// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db_spec

import (
	"time"

	"openpitrix.io/iam/pkg/pb/am"
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

func (p *DBRole) ToPB() *pbam.Role {
	return nil
}
