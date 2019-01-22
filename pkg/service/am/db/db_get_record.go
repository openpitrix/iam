// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db

import (
	"openpitrix.io/iam/pkg/internal/funcutil"
	"openpitrix.io/logger"
)

func (p *Database) getRecordsByRoleId(roleId string) ([]DBRecord, error) {
	logger.Infof(nil, funcutil.CallerName(1))

	var rows = []DBRecord{}
	err := p.DB.Raw(sqlGetAllRecords_by_roleId, roleId).Scan(&rows).Error
	if err != nil {
		logger.Warnf(nil, "%v", sqlGetAllRecords_by_roleId)
		logger.Warnf(nil, "%+v", err)
		return nil, err
	}

	return rows, nil
}

const sqlGetAllRecords_by_roleId = `
-- argument[0]: role_id
-- argument[1]: portal
select distinct
	t3.role_id,
	t3.role_name,
	t3.portal,
	t1.module_id,
	t1.module_name,
	t1.feature_id,
	t1.feature_name,
	t2.data_level,
	t1.action_id,
	(case when isnull(t4.action_id)=0 then 'true' else 'false' end) as action_enabled,
	t1.action_name,
	t1.api_id,
	t1.api_method,
	t1.url_method,
	t1.url
from
	action2 t1
	left join role_module_binding t2 on t1.module_id=t2.module_id
	left join role t3 on t2.role_id=t3.role_id and t3.role_id=?
	left join enable_action t4 on t4.action_id= t1.action_id
`
