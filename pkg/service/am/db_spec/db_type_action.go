// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db_spec

import (
	"openpitrix.io/iam/pkg/pb/am"
)

type DBAction struct {
	ActionId    string `db:"action_id" gorm:"primary_key"`
	ActionName  string `db:"action_name"`
	Method      string `db:"method"`
	Description string `db:"description"`
	FeatureId   string `db:"feature_id"`
	FeatureName string `db:"feature_name"`
	ModuleId    string `db:"module_id"`
	ModuleName  string `db:"module_name"`
	Url         string `db:"url"`
	UrlMethod   string `db:"url_method"`
	Api         string `db:"api"`
	ApiMethod   string `db:"api_method"`
}

func (p *DBAction) ToPB() *pbam.Action {
	return &pbam.Action{
		ActionId:    p.ActionId,
		ActionName:  p.ActionName,
		Method:      p.Method,
		Description: p.Description,
		FeatureId:   p.FeatureId,
		FeatureName: p.FeatureName,
		ModuleId:    p.ModuleId,
		ModuleName:  p.ModuleName,
		Url:         p.Url,
		UrlMethod:   p.UrlMethod,
	}
}
