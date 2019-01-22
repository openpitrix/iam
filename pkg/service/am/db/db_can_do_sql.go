// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db

const sqlCanDo = `
select distinct
	t3.url,
	t3.url_method
from
	role t1,
	role_module_binding t2,
	action2 t3,
	enable_action t4
where
	t1.role_id=t2.role_id and
	t2.module_id=t3.module_id and
	t1.role_id in (
		select
			t2.role_id
		from
			user_role_binding t1,
			role t2
		where
			t1.role_id=t2.role_id and
			t1.user_id=?
	) and
	t4.action_id=t3.action_id and
	right(t3.url,length(t3.url)-locate( '/',t3.url, 2))=? and
	t3.url_method=?
`

type sqlGetAccessPath_args struct {
	UserId    string
	OwnerPath string
	Url       string
	UrlMethod string
}

const sqlGetAccessPath_tmpl = `
select distinct
	(case
		when t2.data_level='all'
			then substring_index('{{.OwnerPath}}', '.', 1)
		when t2.data_level='department'
			then replace('{{.OwnerPath}}', '{{.UserId}}.', '')
		else '{{.OwnerPath}}'
	end) as access_path
from
	role t1,
	role_module_binding t2,
	action2 t3,
	enable_action t4
where
	t1.role_id=t2.role_id and t2.module_id=t3.module_id
	and t1.role_id in
	(select t2.role_id from user_role_binding t1, role t2 where  t1.role_id=t2.role_id and t1.user_id='uid-PYu7bdqa')
	and t4.action_id=t3.action_id
	and t3.module_id=(
		select module_id from action2 where url_method='{{.UrlMethod}}' and
		right(url,length(url)-locate( '/',url, 2))='{{.Url}}'
	)
`
