// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db_spec

import (
	"time"
)

type RoleModuleBinding struct {
	BindId     string `gorm:"type:varchar(50);primary_key"`
	RoleId     string `gorm:"type:varchar(50);not null"`
	ModuleId   string `gorm:"type:varchar(50);not null"`
	DataLevel  string `gorm:"type:varchar(50);not null"`
	IsCheckAll int    `gorm:"type:tinyint;not null"`

	CreateTime time.Time
	UpdateTime time.Time
}
