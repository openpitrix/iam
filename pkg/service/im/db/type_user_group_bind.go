// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db

import (
	"time"
)

type UserGroupBinding struct {
	Id         string    `gorm:"type:varchar(50);primary_key"`
	GroupId    string    `gorm:"type:varchar(50);not null"`
	UserId     string    `gorm:"type:varchar(50);not null"`
	CreateTime time.Time `gorm:"default CURRENT_TIMESTAMP"`
}
