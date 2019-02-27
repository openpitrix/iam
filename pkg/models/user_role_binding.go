// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

import (
	"time"

	"openpitrix.io/iam/pkg/constants"
	"openpitrix.io/iam/pkg/util/idutil"
)

type UserRoleBinding struct {
	Id         string    `gorm:"type:varchar(50);primary_key"`
	UserId     string    `gorm:"type:varchar(50);"`
	RoleId     string    `gorm:"type:varchar(50);"`
	CreateTime time.Time `gorm:"default CURRENT_TIMESTAMP"`
}

func NewUserRoleBinding(userId, roleId string) *UserRoleBinding {
	return &UserRoleBinding{
		Id:         idutil.GetUuid(constants.PrefixUserRoleBindingId),
		RoleId:     roleId,
		UserId:     userId,
		CreateTime: time.Now(),
	}
}
