// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db

import (
	"context"
	"time"

	"github.com/chai2010/template"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	idpkg "openpitrix.io/iam/pkg/id"
	"openpitrix.io/iam/pkg/internal/funcutil"
	"openpitrix.io/iam/pkg/internal/strutil"
	pbam "openpitrix.io/iam/pkg/pb/am"
	"openpitrix.io/iam/pkg/service/am/db_spec"
	"openpitrix.io/iam/pkg/validator"
	"openpitrix.io/logger"
)

func (p *Database) GetRoleModule(ctx context.Context, req *pbam.RoleId) (*pbam.RoleModule, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	req.RoleId = strutil.SimplifyString(req.RoleId)

	if !validator.IsValidId(req.RoleId) {
		err := status.Errorf(codes.InvalidArgument, "invalid role_id: %q", req.RoleId)
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	// 0. get role
	var role = &db_spec.Role{RoleId: req.RoleId}
	if err := p.DB.Model(db_spec.Role{}).Take(role).Error; err != nil {
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	var sActionBundleVisibility string
	switch role.Portal {
	case db_spec.Portal_Admin:
		sActionBundleVisibility = "global_admin_action_bundle_visibility=1"
	case db_spec.Portal_Isv:
		sActionBundleVisibility = "isv_action_bundle_visibility=1"

	case db_spec.Portal_Dev:
		sActionBundleVisibility = "user_action_bundle_visibility=1"
	case db_spec.Portal_User:
		sActionBundleVisibility = "user_action_bundle_visibility=1"

	default:
		sActionBundleVisibility = "user_action_bundle_visibility=1"
	}

	// 1. query roleModuleBindList
	query, err := template.Render(`
		select distinct * from role_module_binding where 1=1
			and role_id='{{.RoleId}}'
		`,
		req,
	)
	if err != nil {
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	query = strutil.SimplifyString(query)
	logger.Infof(ctx, "query: %s", query)

	var roleModuleBindList []db_spec.RoleModuleBinding
	p.DB.Raw(query).Find(&roleModuleBindList)
	if err := p.DB.Error; err != nil {
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}
	if len(roleModuleBindList) == 0 {
		reply := &pbam.RoleModule{RoleId: req.RoleId}
		return reply, nil
	}

	// 2. query moduleApiList
	query, err = template.Render(`
		select distinct
			module_api.*
		from
			role_module_binding, module_api
		where 1=1
			and role_module_binding.module_id=module_api.module_id

			and {{sActionBundleVisibility}}

			and role_module_binding.role_id='{{.RoleId}}'
		`,
		req, template.FuncMap{
			"sActionBundleVisibility": func() string { return sActionBundleVisibility },
		},
	)
	if err != nil {
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	query = strutil.SimplifyString(query)
	logger.Infof(ctx, "query: %s", query)

	var moduleApiList []db_spec.ModuleApi
	p.DB.Raw(query).Find(&moduleApiList)
	if err := p.DB.Error; err != nil {
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}
	if len(moduleApiList) == 0 {
		logger.Warnf(ctx, "no module_api, req: %v", req)
		reply := &pbam.RoleModule{RoleId: req.RoleId}
		return reply, nil
	}

	// 3. query enableActionList
	query, err = template.Render(`
		select distinct enable_action_bundle.* from
			enable_action_bundle, role_module_binding, module_api
		where 1=1
			and enable_action_bundle.bind_id=role_module_binding.bind_id
			and enable_action_bundle.action_bundle_id=module_api.action_bundle_id
			and module_api.module_id=role_module_binding.module_id

			and {{sActionBundleVisibility}}

			and role_module_binding.role_id='{{.RoleId}}'
		`,
		req, template.FuncMap{
			"sActionBundleVisibility": func() string { return sActionBundleVisibility },
		},
	)
	if err != nil {
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	query = strutil.SimplifyString(query)
	logger.Infof(ctx, "query: %s", query)

	var enableActionList []db_spec.EnableActionBundle
	p.DB.Raw(query).Find(&enableActionList)
	if err := p.DB.Error; err != nil {
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	// 4. build module tree
	return p.buildRoleModuleTree(ctx, role,
		roleModuleBindList,
		moduleApiList,
		enableActionList,
	)
}

func (p *Database) ModifyRoleModule(ctx context.Context, req *pbam.RoleModule) (*pbam.RoleModule, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	// check role_id
	if !validator.IsValidId(req.RoleId) {
		err := status.Errorf(codes.InvalidArgument, "invalid RoleId: %v", req.RoleId)
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}
	if len(req.Module) == 0 {
		return req, nil
	}

	// get RoleModuleBindingList
	var roleModuleBindingList []db_spec.RoleModuleBinding
	for _, v := range req.Module {
		var isCheckAll = 0
		if v.IsCheckAll {
			isCheckAll = 1
		}
		roleModuleBindingList = append(
			roleModuleBindingList,
			db_spec.RoleModuleBinding{
				BindId:     idpkg.GenId("xid-"),
				RoleId:     req.RoleId,
				ModuleId:   v.ModuleId,
				DataLevel:  v.DataLevel,
				IsCheckAll: isCheckAll,
				CreateTime: time.Now(), // todo: query from DB
				UpdateTime: time.Now(),
			},
		)
	}

	// get EnableActionList
	var enableActionList []db_spec.EnableActionBundle
	for i, module := range req.Module {
		for _, feature := range module.Feature {
			for _, action := range feature.ActionBundle {
				if action.ActionBundleEnabled || strutil.Contains(feature.CheckedActionBundleId, action.ActionBundleId) {
					enableActionList = append(
						enableActionList,
						db_spec.EnableActionBundle{
							EnableId:       idpkg.GenId("xid-"),
							BindId:         roleModuleBindingList[i].BindId,
							ActionBundleId: action.ActionBundleId,
						},
					)
				}
			}
		}
	}

	// TODO: check all ModuleId exists
	// TODO: check all ActionId exists

	tx := p.DB.Begin()
	{
		// delete old RoleModuleBindingList
		for _, v := range roleModuleBindingList {
			tx.Exec(
				`DELETE from role_module_binding where role_id=? and module_id=?`,
				req.RoleId, v.ModuleId,
			)
			if err := tx.Error; err != nil {
				tx.Rollback()
				logger.Warnf(ctx, "%+v", err)
				return nil, err
			}
		}

		// delete old EnableActionList
		for _, v := range enableActionList {
			tx.Exec(
				`DELETE from enable_action_bundle where bind_id=? and action_bundle_id=?`,
				v.BindId, v.ActionBundleId,
			)
			if err := tx.Error; err != nil {
				tx.Rollback()
				logger.Warnf(ctx, "%+v", err)
				return nil, err
			}
		}

		// insert new RoleModuleBindingList
		for _, v := range roleModuleBindingList {
			tx.NewRecord(v)
			if err := tx.Error; err != nil {
				tx.Rollback()
				logger.Warnf(ctx, "%+v", err)
				return nil, err
			}
		}
		// insert new EnableActionList
		for _, v := range enableActionList {
			tx.NewRecord(v)
			if err := tx.Error; err != nil {
				tx.Rollback()
				logger.Warnf(ctx, "%+v", err)
				return nil, err
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	return req, nil
}
