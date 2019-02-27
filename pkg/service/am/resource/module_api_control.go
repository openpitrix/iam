// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package resource

import (
	"context"
	"strings"

	"openpitrix.io/iam/pkg/constants"
	"openpitrix.io/iam/pkg/gerr"
	"openpitrix.io/iam/pkg/global"
	"openpitrix.io/iam/pkg/models"
	"openpitrix.io/logger"
)

func GetModuleIds(ctx context.Context) ([]string, error) {
	query := "select module_id from module_api group by module_id"
	rows, err := global.Global().Database.Raw(query).Rows()
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}
	var moduleIds []string
	for rows.Next() {
		var moduleId string
		rows.Scan(&moduleId)
		moduleIds = append(moduleIds, moduleId)
	}
	return moduleIds, nil
}

func GetVisibilityModuleIds(ctx context.Context, roleId string) ([]string, error) {
	role, err := GetRole(ctx, roleId)
	if err != nil {
		return nil, err
	}

	columnActionBundleVisibility := role.Portal + constants.ColumnActionBundleVisibilitySuffix

	query := `
		select module_api.module_id
		from
			role_module_binding, module_api
		where role_module_binding.module_id=module_api.module_id
			and role_module_binding.role_id=? and module_api.` +
		columnActionBundleVisibility + "=1 group by module_api.module_id"
	var moduleIds []string
	rows, err := global.Global().Database.Raw(query, roleId).Rows()
	if err != nil {
		logger.Errorf(ctx, "Get visibility module ids by role id [%s] failed: %+v", roleId, err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}
	for rows.Next() {
		var moduleId string
		rows.Scan(&moduleId)
		moduleIds = append(moduleIds, moduleId)
	}
	return moduleIds, nil
}

func getEnableModuleApis(ctx context.Context, roleIds []string) ([]*models.ModuleApi, error) {
	const query = `
		select module_api.* from
			enable_action_bundle, role_module_binding, module_api
		where enable_action_bundle.bind_id=role_module_binding.bind_id
			and enable_action_bundle.action_bundle_id=module_api.action_bundle_id
			and module_api.module_id=role_module_binding.module_id
			and role_module_binding.role_id in (?)
		`
	var enableModuleApis []*models.ModuleApi
	if err := global.Global().Database.Raw(query, roleIds).Scan(&enableModuleApis).Error; err != nil {
		logger.Errorf(ctx, "Get enable module apis by role id [%s] failed: %+v", strings.Join(roleIds, ","), err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}

	return enableModuleApis, nil
}

func GetEnableModuleApis(ctx context.Context, roleIds []string) ([]*models.ModuleApi, error) {
	roleModuleBindings, err := GetRoleModuleBindingsByRoleIds(ctx, roleIds)
	if err != nil {
		return nil, err
	}

	var isCheckAllModuleIds []string
	for _, roleModuleBinding := range roleModuleBindings {
		if roleModuleBinding.IsCheckAll {
			isCheckAllModuleIds = append(isCheckAllModuleIds, roleModuleBinding.ModuleId)
		}
	}

	allEnableModuleApis, err := GetModuleApisByModuleIds(ctx, isCheckAllModuleIds)
	if err != nil {
		return nil, err
	}

	enableModuleApis, err := getEnableModuleApis(ctx, roleIds)
	if err != nil {
		return nil, err
	}

	allEnableModuleApis = append(allEnableModuleApis, enableModuleApis...)

	return models.Unique(allEnableModuleApis), nil
}

func GetVisibilityModuleApis(ctx context.Context, roleId string) ([]*models.ModuleApi, error) {
	role, err := GetRole(ctx, roleId)
	if err != nil {
		return nil, err
	}

	columnActionBundleVisibility := role.Portal + constants.ColumnActionBundleVisibilitySuffix

	query := `
		select module_api.*
		from
			role_module_binding, module_api
		where role_module_binding.module_id=module_api.module_id
			and role_module_binding.role_id=? and module_api.` + columnActionBundleVisibility + "=1"
	var moduleApis []*models.ModuleApi
	if err := global.Global().Database.Raw(query, roleId).Scan(&moduleApis).Error; err != nil {
		logger.Errorf(ctx, "Get module apis by role id [%s] failed: %+v", roleId, err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}

	return moduleApis, nil
}

func GetModuleApisByModuleIds(ctx context.Context, moduleIds []string) ([]*models.ModuleApi, error) {
	var moduleApis []*models.ModuleApi
	if err := global.Global().Database.Table(constants.TableModuleApi).
		Where(constants.ColumnModuleId+" in (?)", moduleIds).
		Find(&moduleApis).Error; err != nil {
		logger.Errorf(ctx, "Get module apis by module ids [%s] failed: %+v", strings.Join(moduleIds, ","), err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}

	return moduleApis, nil
}
