// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db

import (
	"bytes"
	"encoding/gob"

	pbam "openpitrix.io/iam/pkg/pb/am"
)

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

type DBAction struct {
	ApiId          string
	ApiMethod      string
	ApiDescription string

	ModuleId   string
	ModuleName string

	FeatureId   string
	FeatureName string

	ActionId   string
	ActionName string

	Url       string
	UrlMethod string

	// in other tables

	RoleId   string
	RoleName string
	Portal   string

	DataLevel string

	ActionEnabled string
}

func (p *DBAction) ToPB() *pbam.Action {
	var buf bytes.Buffer
	gob.NewEncoder(&buf).Encode(p)

	var q = new(pbam.Action)
	gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(q)

	return q
}
