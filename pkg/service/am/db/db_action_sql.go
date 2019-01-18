// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db

const sqlDescribeActionsBy_RoleId_Protal = `
-- argument[0]: role_id
-- argument[1]: portal
select distinct
	t1.role_id,
	t1.role_name,
	t1.portal,

	t3.module_id,
	t3.module_name,
	t2.data_level,

	t3.feature_id,
	t3.feature_name,

	t3.action_id,
	t3.action_name,
	(case when isnull(t4.action_id)=0 then 'true' else 'false' end) as action_enabled,

	t3.api_id,
	t3.api_method,
	t3.api_description,

	t3.url_method,
	t3.url
from
	role t1,
	role_module_binding t2,
	action2 t3 left join enable_action t4 on t4.action_id=t3.action_id
where
	t1.role_id=t2.role_id and
	t2.module_id=t3.module_id and
	t1.role_id=? and
	t1.portal=?
order by
	t3.module_id,
	t3.feature_id,
	t3.action_id
`
