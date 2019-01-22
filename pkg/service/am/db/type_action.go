// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db

type Action2 struct {
	ApiId string `gorm:"type:varchar(50);primary_key"`

	ModuleId   string `gorm:"type:varchar(50);"`
	ModuleName string `gorm:"type:varchar(50);"`

	FeatureId   string `gorm:"type:varchar(50);"`
	FeatureName string `gorm:"type:varchar(50);"`

	ActionId   string `gorm:"type:varchar(50);"`
	ActionName string `gorm:"type:varchar(50);"`

	ApiMethod      string `gorm:"type:varchar(50);"`
	ApiDescription string `gorm:"type:varchar(100);"`

	UrlMethod string `gorm:"type:varchar(20);"`
	Url       string `gorm:"type:varchar(500);"`
}

type EnableAction struct {
	EnableId string `gorm:"type:varchar(50);primary_key"`
	BindId   string `gorm:"type:varchar(50);"`
	ActionId string `gorm:"type:varchar(50);"`
}
