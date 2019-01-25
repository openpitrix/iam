// Copyright 2018 The OpenPitrix Authors. All rights reserved.
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
	"openpitrix.io/iam/pkg/pb/im"
	"openpitrix.io/iam/pkg/service/im/db_spec"
	"openpitrix.io/iam/pkg/validator"
	"openpitrix.io/logger"
)

func (p *Database) CreateUser(ctx context.Context, req *pbim.User) (*pbim.User, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	// must generate new id
	req.UserId = idpkg.GenId("uid-")

	var dbUser = db_spec.NewUserFromPB(req).AdjustForCreate()
	if err := dbUser.IsValidForCreate(); err != nil {
		err = status.Errorf(codes.InvalidArgument, "%v", err)
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	// create new record
	if err := p.DB.Create(dbUser).Error; err != nil {
		logger.Warnf(ctx, "%+v, %v", err, dbUser)
		return nil, err
	}

	return p.GetUser(ctx, &pbim.UserId{UserId: req.UserId})
}

func (p *Database) DeleteUsers(ctx context.Context, req *pbim.UserIdList) (*pbim.Empty, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	if len(req.UserId) == 0 {
		err := status.Errorf(codes.InvalidArgument, "empty UserId")
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}
	if !validator.IsValidId(req.UserId...) {
		err := status.Errorf(codes.InvalidArgument, "invalid UserId: %v", req.UserId)
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	tx := p.DB.Begin()
	{
		tx.Delete(db_spec.User{}, `user_id in (?)`, req.UserId)
		if err := tx.Error; err != nil {
			tx.Rollback()
			return nil, err
		}

		tx.Delete(db_spec.User{}, `user_id in (?)`, req.UserId)
		if err := tx.Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	reply := &pbim.Empty{}
	return reply, nil
}

func (p *Database) ModifyUser(ctx context.Context, req *pbim.User) (*pbim.User, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	// ignore password
	req.Password = ""

	var dbUser = db_spec.NewUserFromPB(req).AdjustForUpdate()
	if err := dbUser.IsValidForUpdate(); err != nil {
		err = status.Errorf(codes.InvalidArgument, "%v", err)
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	if err := p.DB.Model(dbUser).Updates(dbUser).Error; err != nil {
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	return p.GetUser(ctx, &pbim.UserId{UserId: req.UserId})
}

func (p *Database) GetUser(ctx context.Context, req *pbim.UserId) (*pbim.User, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	var v = db_spec.User{UserId: req.UserId}
	if err := p.DB.Model(db_spec.User{}).Take(&v).Error; err != nil {
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	// ignore Password
	v.Password = ""

	return v.ToPB(), nil
}

func (p *Database) ListUsers(ctx context.Context, req *pbim.ListUsersRequest) (*pbim.ListUsersResponse, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	// fix repeated fileds
	if len(req.GroupId) == 1 && strings.Contains(req.GroupId[0], ",") {
		req.GroupId = strings.Split(req.GroupId[0], ",")
	}
	if len(req.UserId) == 1 && strings.Contains(req.UserId[0], ",") {
		req.UserId = strings.Split(req.UserId[0], ",")
	}
	if len(req.UserName) == 1 && strings.Contains(req.UserName[0], ",") {
		req.UserName = strings.Split(req.UserName[0], ",")
	}
	if len(req.Email) == 1 && strings.Contains(req.Email[0], ",") {
		req.Email = strings.Split(req.Email[0], ",")
	}
	if len(req.PhoneNumber) == 1 && strings.Contains(req.PhoneNumber[0], ",") {
		req.PhoneNumber = strings.Split(req.PhoneNumber[0], ",")
	}
	if len(req.Status) == 1 && strings.Contains(req.Status[0], ",") {
		req.Status = strings.Split(req.Status[0], ",")
	}

	req.GroupId = strutil.SimplifyStringList(req.GroupId)
	req.UserId = strutil.SimplifyStringList(req.UserId)
	req.UserName = strutil.SimplifyStringList(req.UserName)
	req.Email = strutil.SimplifyStringList(req.Email)
	req.PhoneNumber = strutil.SimplifyStringList(req.PhoneNumber)
	req.Status = strutil.SimplifyStringList(req.Status)

	// limit & offset
	if req.Limit == 0 && req.Offset == 0 {
		req.Limit = 20
		req.Offset = 0
	}
	if req.Limit < 0 {
		req.Limit = 0
	}
	if req.Limit > 200 {
		req.Limit = 200
	}
	if req.Offset < 0 {
		req.Offset = 0
	}

	if !validator.IsValidSearchWord(req.SearchWord) {
		err := status.Errorf(codes.InvalidArgument, "invalid search_word: %v", req.SearchWord)
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}
	if !validator.IsValidSortKey(req.SortKey) {
		err := status.Errorf(codes.InvalidArgument, "invalid sort_key: %v", req.SortKey)
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	if !validator.IsValidId(req.GroupId...) {
		err := status.Errorf(codes.InvalidArgument, "invalid gid: %v", req.GroupId)
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}
	if !validator.IsValidId(req.UserId...) {
		err := status.Errorf(codes.InvalidArgument, "invalid uid: %v", req.UserId)
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}
	if !validator.IsValidEmail(req.Email...) {
		err := status.Errorf(codes.InvalidArgument, "invalid email: %v", req.Email)
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}
	if !validator.IsValidPhoneNumber(req.PhoneNumber...) {
		err := status.Errorf(codes.InvalidArgument, "invalid phone_number: %v", req.PhoneNumber)
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	// 1. get group and sub group id list
	if len(req.GroupId) > 0 {
		allGroupId, err := p.getAllSubGroupIds(ctx, &pbim.GroupIdList{
			GroupId: req.GroupId,
		})
		if err != nil {
			logger.Warnf(ctx, "%+v", err)
			return nil, err
		}

		req.GroupId = allGroupId
	}

	const sqlTmpl = `
		{{if not .GroupId}}
			select distinct {{if IsCountMode}}COUNT(*){{else}}*{{end}} from user where 1=1
				{{if .UserId}}
					and user_id in (
						{{range $i, $v := .UserId}}
							{{if eq $i 0}} '{{$v}}' {{else}} ,'{{$v}}' {{end}}
						{{end}}
					)
				{{end}}
				{{if .UserName}}
					and user_name in (
						{{range $i, $v := .UserName}}
							{{if eq $i 0}} '{{$v}}' {{else}} ,'{{$v}}' {{end}}
						{{end}}
					)
				{{end}}
				{{if .Email}}
					and email in (
						{{range $i, $v := .Email}}
							{{if eq $i 0}} '{{$v}}' {{else}} ,'{{$v}}' {{end}}
						{{end}}
					)
				{{end}}
				{{if .PhoneNumber}}
					and phone_number in (
						{{range $i, $v := .PhoneNumber}}
							{{if eq $i 0}} '{{$v}}' {{else}} ,'{{$v}}' {{end}}
						{{end}}
					)
				{{end}}
				{{if .Status}}
					and status in (
						{{range $i, $v := .Status}}
							{{if eq $i 0}} '{{$v}}' {{else}} ,'{{$v}}' {{end}}
						{{end}}
					)
				{{end}}
				{{if .SearchWord}}
					and (1=0
						OR description LIKE '%{{.SearchWord}}%'
						{{if not .UserId}}
							OR user_id LIKE '%{{.SearchWord}}%'
						{{end}}
						{{if not .UserName}}
							OR user_name LIKE '%{{.SearchWord}}%'
						{{end}}
						{{if not .Email}}
							OR email LIKE '%{{.SearchWord}}%'
						{{end}}
						{{if not .PhoneNumber}}
							OR phone_number LIKE '%{{.SearchWord}}%'
						{{end}}
						{{if not .Status}}
							OR status LIKE '%{{.SearchWord}}%'
						{{end}}
					)
				{{end}}
				{{if not IsCountMode}}
					{{if .SortKey}}
						order by {{.SortKey}} {{if .Reverse}} desc {{end}}
					{{end}}
						limit {{.Limit}} offset {{.Offset}}
				{{end}}
		{{else}}
			select distinct {{if IsCountMode}}COUNT(*){{else}}user.*{{end}} from
				user, user_group, user_group_binding
			where 1=1
				and user_group_binding.user_id=user.user_id
				and user_group_binding.group_id=user_group.group_id

				{{if .GroupId}}
					and user_group.group_id in (
						{{range $i, $v := .GroupId}}
							{{if eq $i 0}} '{{$v}}' {{else}} ,'{{$v}}' {{end}}
						{{end}}
					)
				{{end}}
				{{if .UserId}}
					and user.user_id in (
						{{range $i, $v := .UserId}}
							{{if eq $i 0}} '{{$v}}' {{else}} ,'{{$v}}' {{end}}
						{{end}}
					)
				{{end}}
				{{if .UserName}}
					and user.user_name in (
						{{range $i, $v := .UserName}}
							{{if eq $i 0}} '{{$v}}' {{else}} ,'{{$v}}' {{end}}
						{{end}}
					)
				{{end}}
				{{if .Email}}
					and user.email in (
						{{range $i, $v := .Email}}
							{{if eq $i 0}} '{{$v}}' {{else}} ,'{{$v}}' {{end}}
						{{end}}
					)
				{{end}}
				{{if .PhoneNumber}}
					and user.phone_number in (
						{{range $i, $v := .PhoneNumber}}
							{{if eq $i 0}} '{{$v}}' {{else}} ,'{{$v}}' {{end}}
						{{end}}
					)
				{{end}}
				{{if .Status}}
					and user.status in (
						{{range $i, $v := .Status}}
							{{if eq $i 0}} '{{$v}}' {{else}} ,'{{$v}}' {{end}}
						{{end}}
					)
				{{end}}

				{{if .SearchWord}}
					and (1=0
						OR user.description LIKE '%{{.SearchWord}}%'
						{{if not .UserId}}
							OR user.user_id LIKE '%{{.SearchWord}}%'
						{{end}}
						{{if not .UserName}}
							OR user.user_name LIKE '%{{.SearchWord}}%'
						{{end}}
						{{if not .Email}}
							OR user.email LIKE '%{{.SearchWord}}%'
						{{end}}
						{{if not .PhoneNumber}}
							OR user.phone_number LIKE '%{{.SearchWord}}%'
						{{end}}
						{{if not .Status}}
							OR user.status LIKE '%{{.SearchWord}}%'
						{{end}}
					)
				{{end}}
				{{if not IsCountMode}}
					{{if .SortKey}}
						order by user.{{.SortKey}} {{if .Reverse}} desc {{end}}
					{{end}}
						limit {{.Limit}} offset {{.Offset}}
				{{end}}
		{{end}}
	`

	// count mode
	var query, err = template.Render(sqlTmpl, req,
		template.FuncMap{
			"IsCountMode": func() bool { return true },
		},
	)
	if err != nil {
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	query = strutil.SimplifyString(query)
	logger.Infof(ctx, "count: %s", query)

	var total int
	p.DB.Raw(query).Count(&total)
	if err := p.DB.Error; err != nil {
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}
	if total == 0 {
		return &pbim.ListUsersResponse{}, nil
	}

	// query mode
	query, err = template.Render(sqlTmpl, req,
		template.FuncMap{
			"IsCountMode": func() bool { return false },
		},
	)
	if err != nil {
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	query = strutil.SimplifyString(query)
	logger.Infof(ctx, "query: %s; req: %v", query, req)

	var rows []db_spec.User
	p.DB.Raw(query).Find(&rows)
	if err := p.DB.Error; err != nil {
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}
	if len(rows) == 0 {
		return &pbim.ListUsersResponse{}, nil
	}

	// query group_id
	query, err = template.Render(`
		SELECT * FROM user_group_binding WHERE 1=0
			{{range $i, $v := .}}
				OR user_id='{{$v.UserId}}'
			{{end}}
		`, rows,
	)
	if err != nil {
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	query = strutil.SimplifyString(query)
	logger.Infof(ctx, "query: %s", query)

	var bindRows []db_spec.UserGroupBinding
	p.DB.Raw(query).Find(&bindRows)
	if err := p.DB.Error; err != nil {
		logger.Warnf(ctx, "%+v", err)
		// ignore err
	}

	// convert to pb type
	var sets []*pbim.User
	for _, v := range rows {
		v.Password = "" // ignore Password
		sets = append(sets, v.ToPB())
	}

	// save group_id
	for _, v := range bindRows {
		for j, vj := range sets {
			if v.UserId == vj.UserId {
				sets[j].GroupId = append(sets[j].GroupId, v.GroupId)
			}
		}
	}

	reply := &pbim.ListUsersResponse{
		User:  sets,
		Total: int32(total),
	}

	return reply, nil
}
