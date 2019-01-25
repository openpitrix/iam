// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db_spec

type UserRoleBinding struct {
	Id     string `gorm:"type:varchar(50);primary_key"`
	UserId string `gorm:"type:varchar(50);"`
	RoleId string `gorm:"type:varchar(50);"`
}
