// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db

import (
	"context"
	"sort"

	pbam "openpitrix.io/iam/pkg/pb/am"
	"openpitrix.io/iam/pkg/service/am/db_spec"
)

func (p *Database) buildRoleModuleTree(
	ctx context.Context, role *db_spec.Role,
	roleModuleBindList []db_spec.RoleModuleBinding,
	moduleApiList []db_spec.ModuleApi,
	enableActionList []db_spec.EnableAction,
) (*pbam.RoleModule, error) {
	var (
		featureMap = make(map[string]*pbam.Feature)
		moduleMap  = make(map[string]*pbam.Module)
		roleModule = &pbam.RoleModule{RoleId: role.RoleId}
	)

	// 1. moduleApiList => actionList
	var actionList []*pbam.Action
	for _, v := range moduleApiList {
		actionList = append(actionList,
			p.buildActionFromModuleApi(
				&v, roleModuleBindList, enableActionList, role,
			),
		)
	}

	// 2. actionList => feature map
	for _, v := range actionList {
		m := featureMap[v.FeatureId]
		if m == nil {
			m = new(pbam.Feature)
		}

		m.FeatureId = v.FeatureId
		m.FeatureName = v.FeatureName
		m.Action = append(m.Action, v)
		if v.ActionEnabled {
			m.CheckedActionId = append(m.CheckedActionId, v.ActionId)
		}

		featureMap[m.FeatureId] = m
	}
	for _, m := range featureMap {
		sort.Slice(m.Action, func(i, j int) bool {
			return m.Action[i].ActionId < m.Action[j].ActionId
		})
	}

	// 3. feature map => module map
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
		m.Owner = action.Owner

		moduleMap[m.ModuleId] = m
	}
	for _, m := range moduleMap {
		sort.Slice(m.Feature, func(i, j int) bool {
			return m.Feature[i].FeatureId < m.Feature[j].FeatureId
		})
	}

	// module map => role module
	for _, v := range moduleMap {
		roleModule.Module = append(roleModule.Module, v)
	}

	// OK
	return roleModule, nil
}

func (p *Database) buildActionFromModuleApi(
	actionApi *db_spec.ModuleApi,
	roleModuleBindList []db_spec.RoleModuleBinding,
	enableActionList []db_spec.EnableAction,
	role *db_spec.Role,
) *pbam.Action {
	var (
		roleBind     db_spec.RoleModuleBinding
		enableAction db_spec.EnableAction
	)
	for _, v := range roleModuleBindList {
		if v.ModuleId == actionApi.ModuleId {
			roleBind = v
			break
		}
	}
	for _, v := range enableActionList {
		if v.ActionId == actionApi.ActionId {
			enableAction = v
			break
		}
	}

	q := &pbam.Action{
		RoleId:         role.RoleId,
		RoleName:       role.RoleName,
		Portal:         role.Portal,
		ModuleId:       actionApi.ModuleId,
		ModuleName:     actionApi.ModuleName,
		DataLevel:      roleBind.DataLevel,
		Owner:          role.Owner,
		FeatureId:      actionApi.FeatureId,
		FeatureName:    actionApi.FeatureName,
		ActionId:       actionApi.ActionId,
		ActionName:     actionApi.ActionName,
		ActionEnabled:  enableAction != (db_spec.EnableAction{}),
		ApiId:          actionApi.ApiId,
		ApiMethod:      actionApi.ApiId,
		ApiDescription: actionApi.ApiDescription,
		Url:            actionApi.Url,
		UrlMethod:      actionApi.UrlMethod,
	}

	return q
}
