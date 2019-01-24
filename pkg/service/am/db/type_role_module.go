// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db

import (
	"time"
)

type ModuleApi struct {
	ApiId string `gorm:"type:varchar(50);primary_key"`

	ModuleId   string `gorm:"type:varchar(50);not null"`
	ModuleName string `gorm:"type:varchar(50);not null"`

	FeatureId   string `gorm:"type:varchar(50);not null"`
	FeatureName string `gorm:"type:varchar(50);not null"`

	ActionId   string `gorm:"type:varchar(50);not null"`
	ActionName string `gorm:"type:varchar(50);not null"`

	ApiMethod      string `gorm:"type:varchar(50);not null"`
	ApiDescription string `gorm:"type:varchar(100);not null"`

	UrlMethod string `gorm:"type:varchar(100);not null"`
	Url       string `gorm:"type:varchar(255);not null"`
}

type RoleModuleBinding struct {
	BindId     string `gorm:"type:varchar(50);primary_key"`
	RoleId     string `gorm:"type:varchar(50);not null"`
	ModuleId   string `gorm:"type:varchar(50);not null"`
	DataLevel  string `gorm:"type:varchar(50);not null"`
	CreateTime time.Time
	UpdateTime time.Time
	Owner      string `gorm:"type:varchar(50);not null"`
	IsCheckAll int    `gorm:"type:tinyint;not null"`
}

type EnableAction struct {
	EnableId string `gorm:"type:varchar(50);primary_key"`
	BindId   string `gorm:"type:varchar(50);not null"`
	ActionId string `gorm:"type:varchar(50);not null"`
}
