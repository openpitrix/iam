// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package service

import (
	"openpitrix.io/iam/openpitrix/pkg/pb"
)

type DBAction struct {
	ActionId    string `pb:"action_id"`
	ActionName  string `pb:"action_name"`
	Method      string `pb:"method"`
	Description string `pb:"description"`
	FeatureId   string `pb:"feature_id"`
	FeatureName string `pb:"feature_name"`
	ModuleId    string `pb:"module_id"`
	ModuleName  string `pb:"module_name"`
}

func pbActionToDB(p *pb.Action) *DBAction {
	if p == nil {
		return new(DBAction)
	}
	var q = &DBAction{
		ActionId:    p.ActionId,
		ActionName:  p.ActionName,
		Method:      p.Method,
		Description: p.Description,
		FeatureId:   p.FeatureId,
		FeatureName: p.FeatureName,
		ModuleId:    p.ModuleId,
		ModuleName:  p.ModuleName,
	}

	return q
}

func (p *DBAction) ToPb() *pb.Action {
	if p == nil {
		return new(pb.Action)
	}
	var q = &pb.Action{
		ActionId:    p.ActionId,
		ActionName:  p.ActionName,
		Method:      p.Method,
		Description: p.Description,
		FeatureId:   p.FeatureId,
		FeatureName: p.FeatureName,
		ModuleId:    p.ModuleId,
		ModuleName:  p.ModuleName,
	}

	return q
}

func (p *DBAction) ValidateForInsert() error {
	return nil
}
func (p *DBAction) ValidateForUpdate() error {
	return nil
}
