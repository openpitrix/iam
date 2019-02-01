// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db_spec

type ModuleApi struct {
	ApiId string `gorm:"type:varchar(50);primary_key"`

	ModuleId   string `gorm:"type:varchar(50);not null"`
	ModuleName string `gorm:"type:varchar(50);not null"`

	FeatureId   string `gorm:"type:varchar(50);not null"`
	FeatureName string `gorm:"type:varchar(50);not null"`

	ActionBundleId   string `gorm:"type:varchar(50);not null"`
	ActionBundleName string `gorm:"type:varchar(50);not null"`

	ApiMethod      string `gorm:"type:varchar(50);not null"`
	ApiDescription string `gorm:"type:varchar(100);not null"`

	UrlMethod string `gorm:"type:varchar(100);not null"`
	Url       string `gorm:"type:varchar(255);not null"`

	GlobalAdminActionBundleVisibility bool
	IsvActionBundleVisibility         bool
	UserActionBundleVisibility        bool
}
