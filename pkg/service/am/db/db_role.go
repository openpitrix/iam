// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db

import (
	"context"
	"strings"
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

	tx := p.DB.Begin()
	{
		// select modules by portal
		var moduleApiList []db_spec.ModuleApi
		if err := tx.Find(&moduleApiList).Error; err != nil {
			logger.Warnf(ctx, "%+v, %v", err, dbRole)
			return nil, err
		}
		var moduleIdMap = make(map[string]string)
		for _, v := range moduleApiList {
			moduleIdMap[v.ModuleId] = v.ModuleId
		}

		// create new record
		if err := p.DB.Create(dbRole).Error; err != nil {
			tx.Rollback()
			logger.Warnf(ctx, "%+v, %v", err, dbRole)
			return nil, err
		}

		// bind modules with no check
		var now = time.Now()
		for moduleId, _ := range moduleIdMap {
			var dbRoleModuleBinding = &db_spec.RoleModuleBinding{
				BindId:     idpkg.GenId("bind-"),
				RoleId:     dbRole.RoleId,
				ModuleId:   moduleId,
				DataLevel:  "self",
				CreateTime: now,
				UpdateTime: now,
			}
			if err := p.DB.Create(dbRoleModuleBinding).Error; err != nil {
				tx.Rollback()
				logger.Warnf(ctx, "%+v, %v", err, dbRole)
				return nil, err
			}
		}
	}
	if err := tx.Commit().Error; err != nil {
		logger.Warnf(ctx, "%+v", err)
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
			logger.Warnf(ctx, "%+v", err)
			return nil, err
		}

		tx.Raw("delete from user_role_binding where role_id in (?)", req.RoleId)
		if err := tx.Error; err != nil {
			tx.Rollback()
			logger.Warnf(ctx, "%+v", err)
			return nil, err
		}

		tx.Raw("delete from role_module_binding where role_id in (?)", req.RoleId)
		if err := tx.Error; err != nil {
			tx.Rollback()
			logger.Warnf(ctx, "%+v", err)
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

	var role = db_spec.Role{RoleId: req.RoleId}
	if err := p.DB.Model(db_spec.Role{}).Take(&role).Error; err != nil {
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	// find users
	var userRoleBindList []db_spec.UserRoleBinding
	p.DB.Where(&db_spec.UserRoleBinding{RoleId: req.RoleId}).Find(&userRoleBindList)
	if err := p.DB.Error; err != nil {
		logger.Warnf(ctx, "%+v", err)
		// ignore err
	}

	pbRole := role.ToPB()
	for _, v := range userRoleBindList {
		pbRole.UserId = append(pbRole.UserId, v.UserId)
	}

	// OK
	return pbRole, nil
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

	if len(req.RoleId) > 0 {
		if !validator.IsValidId(req.RoleId...) {
			err := status.Errorf(codes.InvalidArgument, "invalid RoleId: %v", req.RoleId)
			logger.Warnf(ctx, "%+v", err)
			return nil, err
		}
	}
	if len(req.RoleName) > 0 {
		if !validator.IsValidName(req.RoleName...) {
			err := status.Errorf(codes.InvalidArgument, "invalid RoleName: %v", req.RoleId)
			logger.Warnf(ctx, "%+v", err)
			return nil, err
		}
	}
	if len(req.Portal) > 0 {
		if !validator.IsValidPortal(req.Portal...) {
			err := status.Errorf(codes.InvalidArgument, "invalid Portal: %v", req.Portal)
			logger.Warnf(ctx, "%+v", err)
			return nil, err
		}
	}
	if len(req.UserId) > 0 {
		if !validator.IsValidId(req.UserId...) {
			err := status.Errorf(codes.InvalidArgument, "invalid UserId: %v", req.UserId)
			logger.Warnf(ctx, "%+v", err)
			return nil, err
		}
	}

	const sqlTmpl = `
		{{if not .UserId}}
			select distinct * from role where 1=1
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
			select distinct role.* from
				role, user_role_binding
			where 1=1
				and user_role_binding.role_id=role.role_id
				and user_role_binding.user_id in (
					{{range $i, $v := .UserId}}
						{{if eq $i 0}} '{{$v}}' {{else}} ,'{{$v}}' {{end}}
					{{end}}
				)
				{{if .RoleId}}
					and role.role_id in (
						{{range $i, $v := .RoleId}}
							{{if eq $i 0}} '{{$v}}' {{else}} ,'{{$v}}' {{end}}
						{{end}}
					)
				{{end}}
				{{if .RoleName}}
					and role.role_name in (
						{{range $i, $v := .RoleName}}
							{{if eq $i 0}} '{{$v}}' {{else}} ,'{{$v}}' {{end}}
						{{end}}
					)
				{{end}}
				{{if .Portal}}
					and role.portal in (
						{{range $i, $v := .Portal}}
							{{if eq $i 0}} '{{$v}}' {{else}} ,'{{$v}}' {{end}}
						{{end}}
					)
				{{end}}
		{{end}}
	`

	var query, err = template.Render(sqlTmpl, req)
	if err != nil {
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	query = strutil.SimplifyString(query)
	logger.Infof(ctx, "query: %s", query)

	var rows []db_spec.Role
	p.DB.Raw(query).Find(&rows)
	if err := p.DB.Error; err != nil {
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}
	if len(rows) == 0 {
		return &pbam.RoleList{}, nil
	}

	// query user_id
	query, err = template.Render(`
		SELECT DISTINCT * FROM user_role_binding WHERE 1=0
			{{range $i, $v := .}}
				OR role_id='{{$v.RoleId}}'
			{{end}}
		`, rows,
	)
	if err != nil {
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	query = strutil.SimplifyString(query)
	logger.Infof(ctx, "query: %s", query)

	var bindRows []db_spec.UserRoleBinding
	p.DB.Raw(query).Find(&bindRows)
	if err := p.DB.Error; err != nil {
		logger.Warnf(ctx, "%+v", err)
		// ignore err
	}

	// convert to pb type
	var sets []*pbam.Role
	for _, v := range rows {
		sets = append(sets, v.ToPB())
	}

	// save user_id
	for _, v := range bindRows {
		for j, vj := range sets {
			if v.RoleId == vj.RoleId {
				sets[j].UserId = append(sets[j].UserId, v.UserId)
			}
		}
	}

	reply := &pbam.RoleList{
		Value: sets,
	}

	return reply, nil
}
