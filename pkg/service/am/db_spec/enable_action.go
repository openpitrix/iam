// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db_spec

type EnableAction struct {
	EnableId string `gorm:"type:varchar(50);primary_key"`
	BindId   string `gorm:"type:varchar(50);not null"`
	ActionId string `gorm:"type:varchar(50);not null"`
}
