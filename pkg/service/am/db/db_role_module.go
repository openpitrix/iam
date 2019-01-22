// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db

import (
	"context"
	"sort"
	"strings"

	"openpitrix.io/iam/pkg/internal/funcutil"
	pbam "openpitrix.io/iam/pkg/pb/am"
	"openpitrix.io/logger"
)

func (p *Database) GetRoleModule(ctx context.Context, req *pbam.RoleId) (*pbam.RoleModule, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	records, err := p.getRecordsByRoleId(req.RoleId)
	if err != nil {
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	var (
		featureMap = make(map[string]*pbam.Feature)
		moduleMap  = make(map[string]*pbam.Module)
	)

	// action => feature map
	for _, v := range records {
		m := featureMap[v.FeatureId]
		if m == nil {
			m = new(pbam.Feature)
		}

		m.FeatureId = v.FeatureId
		m.FeatureName = v.FeatureName
		m.Action = append(m.Action, v.ToPB())
		if v.ActionEnabled == "1" || strings.EqualFold(v.ActionEnabled, "true") {
			m.CheckedActionId = append(m.CheckedActionId, v.ActionId)
		}

		featureMap[m.FeatureId] = m
	}
	for _, m := range featureMap {
		sort.Slice(m.Action, func(i, j int) bool {
			return m.Action[i].ActionId < m.Action[j].ActionId
		})
	}

	// feature map => module map
	for _, v := range featureMap {
		action := v.Action[0]

		m := moduleMap[action.ModuleId]
		if m == nil {
			m = new(pbam.Module)
		}

		m.ModuleId = action.ModuleId
		m.ModuleName = action.ModuleName
		m.Feature = append(m.Feature, v)
		m.DataLevel = action.DataLevel
		m.CheckAll = action.IsFeatureCheckAll == "1" || strings.EqualFold(action.IsFeatureCheckAll, "true")

		moduleMap[m.ModuleId] = m
	}
	for _, m := range moduleMap {
		sort.Slice(m.Feature, func(i, j int) bool {
			return m.Feature[i].FeatureId < m.Feature[j].FeatureId
		})
	}

	// module map => role module
	roleModule := new(pbam.RoleModule)
	for _, v := range moduleMap {
		action := v.Feature[0].Action[0]
		if action.RoleId != req.RoleId {
			continue
		}

		roleModule.RoleId = action.RoleId
		roleModule.RoleName = action.RoleName
		roleModule.Module = append(roleModule.Module, v)
	}

	// OK
	return roleModule, nil
}

func (p *Database) ModifyRoleModule(ctx context.Context, req *pbam.RoleModule) (*pbam.RoleModule, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	panic("todo")
}
