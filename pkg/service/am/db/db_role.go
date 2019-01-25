// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db

import (
	"context"
	"strings"

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

func (p *Database) CreateRole(ctx context.Context, req *pbam.Role) (*pbam.Role, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	// must generate new id
	req.RoleId = idpkg.GenId("role-")

	var dbRole = db_spec.NewRoleFromPB(req).AdjustForCreate()
	if err := dbRole.IsValidForCreate(); err != nil {
		err = status.Errorf(codes.InvalidArgument, "%v", err)
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	// create new record
	if err := p.DB.Create(dbRole).Error; err != nil {
		logger.Warnf(ctx, "%+v, %v", err, dbRole)
		return nil, err
	}

	return p.GetRole(ctx, &pbam.RoleId{RoleId: req.RoleId})
}
func (p *Database) DeleteRoles(ctx context.Context, req *pbam.RoleIdList) (*pbam.Empty, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	if len(req.RoleId) == 0 {
		err := status.Errorf(codes.InvalidArgument, "empty RoleId")
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}
	if !validator.IsValidId(req.RoleId...) {
		err := status.Errorf(codes.InvalidArgument, "invalid RoleId: %v", req.RoleId)
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	tx := p.DB.Begin()
	{
		tx.Raw("delete from role where role_id in (?)", req.RoleId)
		if err := tx.Error; err != nil {
			tx.Rollback()
			return nil, err
		}

		tx.Raw("delete from user_role_binding where role_id in (?)", req.RoleId)
		if err := tx.Error; err != nil {
			tx.Rollback()
			return nil, err
		}

		tx.Raw("delete from role_module_binding where role_id in (?)", req.RoleId)
		if err := tx.Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	if err := tx.Commit().Error; err != nil {
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	return &pbam.Empty{}, nil
}

func (p *Database) ModifyRole(ctx context.Context, req *pbam.Role) (*pbam.Role, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	var dbRole = db_spec.NewRoleFromPB(req).AdjustForUpdate()
	if err := dbRole.IsValidForUpdate(); err != nil {
		err = status.Errorf(codes.InvalidArgument, "%v", err)
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	if err := p.DB.Model(dbRole).Updates(dbRole).Error; err != nil {
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	return p.GetRole(ctx, &pbam.RoleId{RoleId: req.RoleId})
}

func (p *Database) GetRole(ctx context.Context, req *pbam.RoleId) (*pbam.Role, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	var v = db_spec.Role{RoleId: req.RoleId}
	if err := p.DB.Model(db_spec.Role{}).Take(&v).Error; err != nil {
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	return v.ToPB(), nil
}
func (p *Database) GetRoleListByUserId(ctx context.Context, req *pbam.UserId) (*pbam.RoleList, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	var query = `
		select
			role.*
		from
			role,
			user_role_binding
		where
			role.role_id=user_role_binding.role_id and
			user_role_binding.user_id=?
	`

	var rows []db_spec.Role
	if err := p.DB.Raw(query, req.UserId).Find(&rows).Error; err != nil {
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	var sets []*pbam.Role
	for _, v := range rows {
		sets = append(sets, v.ToPB())
	}

	reply := &pbam.RoleList{
		Value: sets,
	}
	return reply, nil
}

func (p *Database) DescribeRoles(ctx context.Context, req *pbam.DescribeRolesRequest) (*pbam.RoleList, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	// fix repeated fileds
	// GET donot support repeated type in grpc-gateway
	if len(req.RoleId) == 1 && strings.Contains(req.RoleId[0], ",") {
		req.RoleId = strings.Split(req.RoleId[0], ",")
	}
	if len(req.RoleName) == 1 && strings.Contains(req.RoleName[0], ",") {
		req.RoleName = strings.Split(req.RoleName[0], ",")
	}
	if len(req.Portal) == 1 && strings.Contains(req.Portal[0], ",") {
		req.Portal = strings.Split(req.Portal[0], ",")
	}
	if len(req.UserId) == 1 && strings.Contains(req.UserId[0], ",") {
		req.UserId = strings.Split(req.UserId[0], ",")
	}

	req.RoleId = strutil.SimplifyStringList(req.RoleId)
	req.RoleName = strutil.SimplifyStringList(req.RoleName)
	req.Portal = strutil.SimplifyStringList(req.Portal)
	req.UserId = strutil.SimplifyStringList(req.UserId)

	// all list
	if len(req.RoleId)+len(req.RoleName)+len(req.Portal)+len(req.UserId) == 0 {
		var query = `select * from role`
		logger.Infof(ctx, "query: %s", strutil.SimplifyString(query))

		var rows []db_spec.Role
		p.DB.Raw(query).Find(&rows)

		var sets []*pbam.Role
		for _, v := range rows {
			sets = append(sets, v.ToPB())
		}

		reply := &pbam.RoleList{Value: sets}
		return reply, nil
	}

	// no user_id:
	// select * from role where 1=1 and role_id in ('a','b')
	//
	// user bind:
	// select distinct t1.* from role t1 where 1=1 and t1.role_id in (
	//	select t2.role_id from user_role_binding t1, role t2
	//	where t1.role_id=t2.role_id and user_id in ('a', 'b')
	// )

	var query = template.MustRender(`
			{{if not .UserId}}
				select * from role where 1=1
				{{if .RoleId}}
					and role_id in (
						{{range $i, $v := .RoleId}}
							{{if eq $i 0}} '{{$v}}' {{else}} ,'{{$v}}' {{end}}
						{{end}}
					)
				{{end}}
				{{if .RoleName}}
					and role_name in (
						{{range $i, $v := .RoleName}}
							{{if eq $i 0}} '{{$v}}' {{else}} ,'{{$v}}' {{end}}
						{{end}}
					)
				{{end}}
				{{if .Portal}}
					and portal in (
						{{range $i, $v := .Portal}}
							{{if eq $i 0}} '{{$v}}' {{else}} ,'{{$v}}' {{end}}
						{{end}}
					)
				{{end}}
			{{else}}
				select distinct t1.* from role t1 where 1=1
				{{if .RoleId}}
					and role_id in (
						{{range $i, $v := .RoleId}}
							{{if eq $i 0}} '{{$v}}' {{else}} ,'{{$v}}' {{end}}
						{{end}}
					)
				{{end}}
				{{if .RoleName}}
					and role_name in (
						{{range $i, $v := .RoleName}}
							{{if eq $i 0}} '{{$v}}' {{else}} ,'{{$v}}' {{end}}
						{{end}}
					)
				{{end}}
				{{if .Portal}}
					and portal in (
						{{range $i, $v := .Portal}}
							{{if eq $i 0}} '{{$v}}' {{else}} ,'{{$v}}' {{end}}
						{{end}}
					)
				{{end}}

				and t1.role_id in (
					select t2.role_id from
						user_role_binding t1,
						role t2
					where
						t1.role_id=t2.role_id and
						user_id in (
							{{range $i, $v := .UserId}}
								{{if eq $i 0}} '{{$v}}' {{else}} ,'{{$v}}' {{end}}
							{{end}}
						)
				)
			{{end}}
		`, req,
	)

	query = strutil.SimplifyString(query)
	logger.Infof(ctx, "query: %s", query)

	var rows []db_spec.Role
	p.DB.Raw(query).Find(&rows)

	var sets []*pbam.Role
	for _, v := range rows {
		sets = append(sets, v.ToPB())
	}

	reply := &pbam.RoleList{Value: sets}
	return reply, nil

}
