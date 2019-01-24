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

	"openpitrix.io/iam/pkg/internal/funcutil"
	pbam "openpitrix.io/iam/pkg/pb/am"
	"openpitrix.io/logger"
)

func (p *Database) GetRoleModule(ctx context.Context, req *pbam.RoleId) (*pbam.RoleModule, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	if !isValidIds(req.RoleId) {
		err := status.Errorf(codes.InvalidArgument, "invalid role_id: %q", req.RoleId)
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	query, err := template.Render(sqlGetRoleModule_by_roleId, req)
	if err != nil {
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	var rows = []ModuleApiInfo{}
	if err := p.DB.Raw(query).Scan(&rows).Error; err != nil {
		logger.Warnf(nil, "%v", query)
		logger.Warnf(nil, "%+v", err)
		return nil, err
	}

	roleModuleMap := ModuleApiInfoList(rows).ToRoleModuleMap()
	return roleModuleMap[req.RoleId], nil
}

func (p *Database) ModifyRoleModule(ctx context.Context, req *pbam.RoleModule) (*pbam.RoleModule, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	var moduleIdList = func() []string {
		var ss []string
		for _, mod := range req.Module {
			ss = append(ss, mod.ModuleId)
		}
		return ss
	}()

	tx := p.DB.Begin()

	tx.Raw(
		`delete from role_module_binding where role_id=? and module_id in (?)`,
		req.RoleId, moduleIdList,
	)
	if err := tx.Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	for _, mod := range req.Module {
		tx.Raw(template.MustRender(
			`delete from enable_action where bind_id in (
				select bind_id from role_module_binding where
					role_id='{{.RoleId}}' and module_id='{{.ModuleId}}'
			);`,
			struct{ RoleId, ModuleId string }{
				RoleId:   req.RoleId,
				ModuleId: mod.ModuleId,
			},
		))
		if err := tx.Error; err != nil {
			tx.Rollback()
			return nil, err
		}

		bindId := genId("bind-", 12)
		tx.NewRecord(RoleModuleBinding{
			BindId:     bindId,
			RoleId:     req.RoleId,
			ModuleId:   mod.ModuleId,
			DataLevel:  mod.DataLevel,
			CreateTime: time.Now(),
			UpdateTime: time.Now(),
			Owner:      mod.Owner,
			IsCheckAll: btoi(mod.IsCheckAll),
		})
		if err := tx.Error; err != nil {
			tx.Rollback()
			return nil, err
		}

		// for chekced actions
		for _, feature := range mod.Feature {
			for _, action := range feature.Action {
				tx.NewRecord(EnableAction{
					EnableId: genId("id-", 12),
					BindId:   bindId,
					ActionId: action.ActionId,
				})
				if err := tx.Error; err != nil {
					tx.Rollback()
					return nil, err
				}
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	return req, nil
}

const sqlGetRoleModule_by_roleId = `
select distinct
	'{{.RoleId}}' as role_id,
	t.module_id,
	t.module_name,
	t.data_level,
	t.owner,
	(case when isnull(t.is_check_all)=0 then 1 else 0 end) as is_check_all,
	t.feature_id,
	t.feature_name,
	t.action_id,
	t.action_name ,
	(case when isnull(tt.action_id)=0 then 1 else 0 end) as action_enabled
FROM (
		select distinct
			t3.role_id,
			t3.role_name,
			t3.portal,
			t1.module_id,
			t1.module_name,
			t2.data_level,
			t2.is_check_all,
			t1.feature_id,
			t1.feature_name,
			t1.action_id,
			t1.action_name,
			t2.bind_id
		from module_api t1
			left join role_module_binding t2 on t1.module_id=t2.module_id  and t2.role_id='{{.RoleId}}'
		left join role t3 on t2.role_id=t3.role_id   and t3.role_id='{{.RoleId}}'
	)t
	left join enable_action tt on t.action_id=tt.action_id and t.bind_id=tt.bind_id
order by t.action_id;
`
