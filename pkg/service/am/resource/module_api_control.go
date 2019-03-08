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

func GetModuleApisByActionBundleIds(ctx context.Context, actionBundleIds []string) ([]*models.ModuleApi, error) {
	var moduleApis []*models.ModuleApi
	if err := global.Global().Database.
		Table(constants.TableModuleApi).
		Where(constants.ColumnActionBundleId+" in (?)", actionBundleIds).
		Find(&moduleApis).Error; err != nil {
		return nil, gerr.New(ctx, gerr.Internal, gerr.ErrorInternalError)
	}

	return moduleApis, nil
}

func GetModuleIds(ctx context.Context) ([]string, error) {
	query := "select module_id from `module_api` group by module_id"
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

	var moduleIds []string
	rows, err := global.Global().Database.
		Table(constants.TableModuleApi).
		Select(constants.TableModuleApi+"."+constants.ColumnModuleId).
		Where(constants.TableModuleApi+"."+columnActionBundleVisibility+" = 1").
		Joins("JOIN "+constants.TableRoleModuleBinding+" on "+constants.TableRoleModuleBinding+"."+constants.ColumnRoleId+" = ?", roleId).
		Group(constants.TableModuleApi + "." + constants.ColumnModuleId).Rows()
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
	var enableModuleApis []*models.ModuleApi
	if err := global.Global().Database.
		Table(constants.TableModuleApi).
		Select(constants.TableModuleApi+".*").
		Joins("JOIN "+constants.TableRoleModuleBinding+" on "+
			constants.TableRoleModuleBinding+"."+constants.ColumnModuleId+" = "+constants.TableModuleApi+"."+constants.ColumnModuleId).
		Where(constants.TableRoleModuleBinding+"."+constants.ColumnRoleId+" in (?)", roleIds).
		Joins("JOIN " + constants.TableEnableActionBundle + " on " +
			constants.TableEnableActionBundle + "." + constants.ColumnActionBundleId + " = " + constants.TableModuleApi + "." + constants.ColumnActionBundleId).
		Scan(&enableModuleApis).Error; err != nil {
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

	var moduleApis []*models.ModuleApi
	if err := global.Global().Database.
		Table(constants.TableModuleApi).
		Select(constants.TableModuleApi+".*").
		Joins("JOIN "+constants.TableRoleModuleBinding+" on "+
			constants.TableRoleModuleBinding+"."+constants.ColumnModuleId+" = "+constants.TableModuleApi+"."+constants.ColumnModuleId).
		Where(constants.TableRoleModuleBinding+"."+constants.ColumnRoleId+" = ?", roleId).
		Where(constants.TableModuleApi + "." + columnActionBundleVisibility + " = 1").
		Scan(&moduleApis).Error; err != nil {
		logger.Errorf(ctx, "Get module apis by role id [%s] failed: %+v", roleId, err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}

	return moduleApis, nil
}

func GetModuleApisByModuleIds(ctx context.Context, moduleIds []string) ([]*models.ModuleApi, error) {
	var moduleApis []*models.ModuleApi
	if err := global.Global().Database.
		Table(constants.TableModuleApi).
		Where(constants.ColumnModuleId+" in (?)", moduleIds).
		Find(&moduleApis).Error; err != nil {
		logger.Errorf(ctx, "Get module apis by module ids [%s] failed: %+v", strings.Join(moduleIds, ","), err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}

	return moduleApis, nil
}
