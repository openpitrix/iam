// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package resource

import (
	"context"
	"strings"

	"openpitrix.io/iam/pkg/gerr"
	"openpitrix.io/iam/pkg/global"
	"openpitrix.io/iam/pkg/models"
	"openpitrix.io/logger"
)

func GetRoleModuleBindingsByRoleIds(ctx context.Context, roleIds []string) ([]*models.RoleModuleBinding, error) {
	const query = `
		select role_module_binding.*
		from
			role_module_binding
		where role_module_binding.role_id in (?)
	`
	var roleModuleBindings []*models.RoleModuleBinding
	if err := global.Global().Database.Raw(query, roleIds).Scan(&roleModuleBindings).Error; err != nil {
		logger.Errorf(ctx, "Get role module bindings by role [%s] failed: %+v", strings.Join(roleIds, ","), err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}

	return roleModuleBindings, nil
}

func GetRoleModuleBindingsByRoleIdsAndModuleIds(ctx context.Context, roleIds, moduleIds []string) ([]*models.RoleModuleBinding, error) {
	const query = `
		select role_module_binding.*
		from
			role_module_binding
		where role_id in (?) and module_id in (?)
	`
	var roleModuleBindings []*models.RoleModuleBinding
	if err := global.Global().Database.Raw(query, roleIds, moduleIds).Scan(&roleModuleBindings).Error; err != nil {
		logger.Errorf(ctx, "Get role module bindings by role [%s] module [%s] failed: %+v",
			strings.Join(roleIds, ","), strings.Join(moduleIds, ","), err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}

	return roleModuleBindings, nil
}
