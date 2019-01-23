// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db

import (
	"context"
	"strings"

	"github.com/chai2010/template"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"openpitrix.io/iam/pkg/internal/funcutil"
	"openpitrix.io/iam/pkg/pb/im"
	"openpitrix.io/logger"
)

func (p *Database) CreateUser(ctx context.Context, req *pbim.User) (*pbim.User, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	var dbUser = NewUserFromPB(req)
	if dbUser.Password == "" {
		err := status.Errorf(codes.InvalidArgument, "empty password")
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	if dbUser.Password != "" {
		hashedPass, err := bcrypt.GenerateFromPassword(
			[]byte(dbUser.Password), bcrypt.DefaultCost,
		)
		if err != nil {
			logger.Warnf(ctx, "%+v", err)
			return nil, err
		}
		dbUser.Password = string(hashedPass)
	}

	if err := p.DB.Create(dbUser).Error; err != nil {
		logger.Warnf(ctx, "%+v, %v", err, dbUser)
		return nil, err
	}

	return req, nil
}

func (p *Database) DeleteUsers(ctx context.Context, req *pbim.UserIdList) (*pbim.Empty, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	if req == nil || len(req.UserId) == 0 || !isValidIds(req.UserId...) {
		err := status.Errorf(codes.InvalidArgument, "empty field")
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	tx := p.DB.Begin()
	{
		tx.Delete(User{}, `user_id in (?)`, req.UserId)
		if err := tx.Error; err != nil {
			tx.Rollback()
			return nil, err
		}

		tx.Delete(User{}, `user_id in (?)`, req.UserId)
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

	if req.UserId == "" {
		err := status.Errorf(codes.InvalidArgument, "empty field")
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	req.Password = ""

	var dbUser = NewUserFromPB(req)
	if err := p.DB.Model(dbUser).Updates(dbUser).Error; err != nil {
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	return p.GetUser(ctx, &pbim.UserId{UserId: req.UserId})
}

func (p *Database) GetUser(ctx context.Context, req *pbim.UserId) (*pbim.User, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	var v = User{UserId: req.UserId}
	if err := p.DB.Model(User{}).Take(&v).Error; err != nil {
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

	req.GroupId = simplifyStringList(req.GroupId)
	req.UserId = simplifyStringList(req.UserId)
	req.UserName = simplifyStringList(req.UserName)
	req.Email = simplifyStringList(req.Email)
	req.PhoneNumber = simplifyStringList(req.PhoneNumber)
	req.Status = simplifyStringList(req.Status)

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

	if !isValidSearchWord(req.SearchWord) {
		err := status.Errorf(codes.InvalidArgument, "invalid search_word: %v", req.SearchWord)
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}
	if !isValidSortKey(req.SortKey) {
		err := status.Errorf(codes.InvalidArgument, "invalid sort_key: %v", req.SortKey)
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	if !isValidIds(req.GroupId...) {
		err := status.Errorf(codes.InvalidArgument, "invalid gid: %v", req.GroupId)
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}
	if !isValidIds(req.UserId...) {
		err := status.Errorf(codes.InvalidArgument, "invalid uid: %v", req.UserId)
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}
	if !isValidEmails(req.Email...) {
		err := status.Errorf(codes.InvalidArgument, "invalid email: %v", req.Email)
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}
	if !isValidPhoneNumbers(req.PhoneNumber...) {
		err := status.Errorf(codes.InvalidArgument, "invalid phone_number: %v", req.PhoneNumber)
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	var query, err = template.Render(
		`{{if not .GroupId}}
			select * from user where 1=1
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
				{{if .SortKey}}
					order by {{.SortKey}} {{if .Reverse}} desc {{end}}
				{{end}}
				limit {{.Limit}} offset {{.Offset}}
		{{else}}
			select user.* from user, user_group, user_group_binding where 1=1
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
				{{if .SortKey}}
					order by user.{{.SortKey}} {{if .Reverse}} desc {{end}}
				{{end}}
				limit {{.Limit}} offset {{.Offset}}
		{{end}}
		`, req,
	)
	if err != nil {
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	query = simplifyString(query)
	logger.Infof(ctx, "query: %s", query)

	var rows []User
	p.DB.Raw(query).Find(&rows)
	if err := p.DB.Error; err != nil {
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	var sets []*pbim.User
	for _, v := range rows {
		sets = append(sets, v.ToPB())
	}

	reply := &pbim.ListUsersResponse{
		User: sets,
	}

	return reply, nil
}
