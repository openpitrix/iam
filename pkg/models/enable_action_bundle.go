// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

import (
	"time"

	"openpitrix.io/iam/pkg/constants"
	"openpitrix.io/iam/pkg/util/idutil"
)

type EnableActionBundle struct {
	EnableId       string `gorm:"type:varchar(50);primary_key"`
	BindId         string `gorm:"type:varchar(50);not null"`
	ActionBundleId string `gorm:"type:varchar(50);not null"`

	CreateTime time.Time
}

func NewEnableActionBundle(bindId, actionBundleId string) *EnableActionBundle {
	now := time.Now()
	return &EnableActionBundle{
		EnableId:       idutil.GetUuid(constants.PrefixEnableActionBundleId),
		BindId:         bindId,
		ActionBundleId: actionBundleId,
		CreateTime:     now,
	}
}
