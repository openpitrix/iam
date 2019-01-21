// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db

import (
	"context"
	"time"

	"openpitrix.io/iam/pkg/internal/funcutil"
	pbam "openpitrix.io/iam/pkg/pb/am"
	"openpitrix.io/logger"
)

func (p *Database) ModifyRoleModuleBinding(ctx context.Context, req *pbam.ModifyRoleModuleBindingRequest) (*pbam.ActionList, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	tx := p.DB.Begin()
	{
		tx.Raw(
			`
			-- argument[0]: role_id
			-- argument[1]: module_id list
			delete from enable_action where bind_id in(
				select distinct
				bind_id
			from
				role_module_binding
			where
				role_id=? and
				module_id in (?)
			);
			`,
			req.RoleId,
			req.ModuleId,
		)
		if err := tx.Error; err != nil {
			tx.Rollback()
			return nil, err
		}

		tx.Raw(
			`
			-- argument[0]: role_id
			-- argument[1]: module_id list
			delete from
				role_module_binding
			where
				role_id=? and
				module_id in (?);
			`,
			req.RoleId,
			req.ModuleId,
		)
		if err := tx.Error; err != nil {
			tx.Rollback()
			return nil, err
		}

		for i, v := range req.Module {

			enable_id := "" // TODO: ?
			bind_id := genId("bind-", 12)
			action_id := "" // TODO: ?

			createTime := time.Now()
			updateTime := time.Now()
			owner := "" // TODO: ?

			tx.Raw(
				`
				insert into role_module_binding(
					bind_id,
					role_id,
					module_id,
					data_level,
					create_time,
					update_time,
					owner
				) values(
					?, -- bind_id
					?, -- role_id
					?, -- module_id
					?, -- data_level
					?, -- create_time
					?, -- update_time
					?  -- owner
				)
				`,
				bind_id,
				req.RoleId,
				req.ModuleId[i],
				v.DataLevel,
				createTime,
				updateTime,
				owner,
			)
			if err := tx.Error; err != nil {
				tx.Rollback()
				return nil, err
			}

			tx.Raw(
				`
				insert into enable_action (
					enable_id,
					bind_id,
					action_id
				) values('','','')
				`,
				enable_id,
				bind_id,
				action_id,
			)
		}
		if err := tx.Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	logger.Infof(ctx, "TODO")
	return nil, nil
}
