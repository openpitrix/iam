// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db_spec

import (
	"fmt"
	"time"

	"openpitrix.io/iam/pkg/validator"
)

type UserGroupBinding struct {
	Id         string    `gorm:"type:varchar(50);primary_key"`
	GroupId    string    `gorm:"type:varchar(50);not null"`
	UserId     string    `gorm:"type:varchar(50);not null"`
	CreateTime time.Time `gorm:"default CURRENT_TIMESTAMP"`
}

func (p *UserGroupBinding) IsValidForCreate() error {
	if !validator.IsValidId(p.Id) {
		return fmt.Errorf("UserGroupBinding.IsValidForCreate: invalid Id %q", p.Id)
	}
	if !validator.IsValidId(p.GroupId) {
		return fmt.Errorf("UserGroupBinding.IsValidForCreate: invalid GroupId %q", p.GroupId)
	}
	if !validator.IsValidId(p.UserId) {
		return fmt.Errorf("UserGroupBinding.IsValidForCreate: invalid UserId %q", p.UserId)
	}

	return nil
}
