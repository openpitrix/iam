// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

import (
	"time"

	"openpitrix.io/iam/pkg/constants"
	"openpitrix.io/iam/pkg/util/idutil"
)

type RoleModuleBinding struct {
	BindId     string `gorm:"type:varchar(50);primary_key"`
	RoleId     string `gorm:"type:varchar(50);not null"`
	ModuleId   string `gorm:"type:varchar(50);not null"`
	DataLevel  string `gorm:"type:varchar(50);not null"`
	IsCheckAll bool

	CreateTime time.Time
}

func NewRoleModuleBinding(roleId, moduleId, dataLevel string, isCheckAll bool) *RoleModuleBinding {
	now := time.Now()
	return &RoleModuleBinding{
		BindId:     idutil.GetUuid(constants.PrefixRoleModuleBindingId),
		RoleId:     roleId,
		ModuleId:   moduleId,
		DataLevel:  dataLevel,
		IsCheckAll: isCheckAll,
		CreateTime: now,
	}
}
