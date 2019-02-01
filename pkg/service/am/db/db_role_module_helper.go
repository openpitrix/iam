// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db

import (
	"context"
	"sort"

	"openpitrix.io/iam/pkg/internal/strutil"
	pbam "openpitrix.io/iam/pkg/pb/am"
	"openpitrix.io/iam/pkg/service/am/db_spec"
)

func (p *Database) buildRoleModuleTree(
	ctx context.Context, role *db_spec.Role,
	roleModuleBindList []db_spec.RoleModuleBinding,
	moduleApiList []db_spec.ModuleApi,
	enableActionList []db_spec.EnableActionBundle,
) (*pbam.RoleModule, error) {
	var (
		featureMap = make(map[string]*pbam.ModuleFeature)
		moduleMap  = make(map[string]*pbam.RoleModuleElem)
		roleModule = &pbam.RoleModule{RoleId: role.RoleId}
	)

	// 1. moduleApiList => actionList
	var actionList []*pbam.ModuleFeatureActionBundle
	for _, v := range moduleApiList {
		actionList = append(actionList,
			p.buildActionFromModuleApi(
				&v, roleModuleBindList, enableActionList, role,
			),
		)
	}
	sort.Slice(actionList, func(i, j int) bool {
		return actionList[i].ActionBundleId < actionList[j].ActionBundleId
	})

	// 2. actionList => feature map
	for _, v := range actionList {
		m := featureMap[v.FeatureId]
		if m == nil {
			m = new(pbam.ModuleFeature)
		}

		m.FeatureId = v.FeatureId
		m.FeatureName = v.FeatureName
		m.ActionBundle = append(m.ActionBundle, v)
		if v.ActionBundleEnabled && !strutil.Contains(m.CheckedActionBundleId, v.ActionBundleId) {
			m.CheckedActionBundleId = append(m.CheckedActionBundleId, v.ActionBundleId)
		}

		featureMap[m.FeatureId] = m
	}
	for _, m := range featureMap {
		sort.Slice(m.ActionBundle, func(i, j int) bool {
			return m.ActionBundle[i].ActionBundleId < m.ActionBundle[j].ActionBundleId
		})
	}

	// 3. feature map => module map
	for _, v := range featureMap {
		action := v.ActionBundle[0]

		m := moduleMap[action.ModuleId]
		if m == nil {
			m = new(pbam.RoleModuleElem)
		}

		m.ModuleId = action.ModuleId
		m.ModuleName = action.ModuleName
		m.Feature = append(m.Feature, v)
		m.Owner = action.Owner
		m.DataLevel = action.DataLevel

		moduleMap[m.ModuleId] = m
	}
	for _, m := range moduleMap {
		sort.Slice(m.Feature, func(i, j int) bool {
			return m.Feature[i].FeatureId < m.Feature[j].FeatureId
		})
		for _, bind := range roleModuleBindList {
			if m.ModuleId == bind.ModuleId && bind.RoleId == role.RoleId {
				m.IsCheckAll = bind.IsCheckAll != 0
			}
		}
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
	enableActionList []db_spec.EnableActionBundle,
	role *db_spec.Role,
) *pbam.ModuleFeatureActionBundle {
	var (
		roleBind     db_spec.RoleModuleBinding
		enableAction db_spec.EnableActionBundle
	)
	for _, v := range roleModuleBindList {
		if v.ModuleId == actionApi.ModuleId && v.RoleId == role.RoleId {
			roleBind = v
			break
		}
	}
	for _, v := range enableActionList {
		if v.ActionBundleId == actionApi.ActionBundleId && v.BindId == roleBind.BindId {
			enableAction = v
			break
		}
	}

	q := &pbam.ModuleFeatureActionBundle{
		RoleId:              role.RoleId,
		RoleName:            role.RoleName,
		Portal:              role.Portal,
		ModuleId:            actionApi.ModuleId,
		ModuleName:          actionApi.ModuleName,
		DataLevel:           roleBind.DataLevel,
		Owner:               role.Owner,
		FeatureId:           actionApi.FeatureId,
		FeatureName:         actionApi.FeatureName,
		ActionBundleId:      actionApi.ActionBundleId,
		ActionBundleName:    actionApi.ActionBundleName,
		ActionBundleEnabled: enableAction != (db_spec.EnableActionBundle{}),
		ApiId:               actionApi.ApiId,
		ApiMethod:           actionApi.ApiId,
		ApiDescription:      actionApi.ApiDescription,
		Url:                 actionApi.Url,
		UrlMethod:           actionApi.UrlMethod,

		GlobalAdminActionBundleVisibility: actionApi.GlobalAdminActionBundleVisibility,
		IsvActionBundleVisibility:         actionApi.IsvActionBundleVisibility,
		UserActionBundleVisibility:        actionApi.UserActionBundleVisibility,
	}

	return q
}
