// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db

import (
	"context"
	"fmt"
	"strings"

	"github.com/chai2010/template"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"openpitrix.io/iam/pkg/internal/funcutil"
	"openpitrix.io/iam/pkg/internal/strutil"
	pbam "openpitrix.io/iam/pkg/pb/am"
	pbim "openpitrix.io/iam/pkg/pb/im"
	"openpitrix.io/iam/pkg/service/am/db_spec"
	"openpitrix.io/iam/pkg/validator"
	"openpitrix.io/logger"
)

func (p *Database) GetUserWithRole(ctx context.Context, req *pbam.UserId) (*pbam.UserWithRole, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", p.cfg.ImHost, p.cfg.ImPort), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := pbim.NewAccountManagerClient(conn)
	imUser, err := client.GetUser(ctx, &pbim.UserId{UserId: req.UserId})
	if err != nil {
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	query, err := template.Render(`
		SELECT DISTINCT role.* FROM role, user_role_binding WHERE 1=1
			AND user_role_binding.role_id=role.role_id
			AND user_role_binding.user_id='{{.UserId}}'
		`, imUser,
	)
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
		// ignore err
	}

	user := &pbam.UserWithRole{
		UserId:      imUser.UserId,
		Username:    imUser.Username,
		Email:       imUser.Email,
		PhoneNumber: imUser.PhoneNumber,
		Description: imUser.Description,
		Status:      imUser.Status,
		Extra:       imUser.Extra,
		CreateTime:  imUser.CreateTime,
		UpdateTime:  imUser.UpdateTime,
		StatusTime:  imUser.StatusTime,
	}

	for _, v := range rows {
		user.Role = append(user.Role, v.ToPB())
	}

	return user, nil
}
func (p *Database) DescribeUsersWithRole(ctx context.Context, req *pbam.DescribeUsersWithRoleRequest) (*pbam.DescribeUsersWithRoleResponse, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	// fix repeated fileds
	// GET donot support repeated type in grpc-gateway
	if len(req.RoleId) == 1 && strings.Contains(req.RoleId[0], ",") {
		req.RoleId = strings.Split(req.RoleId[0], ",")
	}
	if len(req.GroupId) == 1 && strings.Contains(req.GroupId[0], ",") {
		req.GroupId = strings.Split(req.GroupId[0], ",")
	}
	if len(req.UserId) == 1 && strings.Contains(req.UserId[0], ",") {
		req.UserId = strings.Split(req.UserId[0], ",")
	}
	if len(req.Username) == 1 && strings.Contains(req.Username[0], ",") {
		req.Username = strings.Split(req.Username[0], ",")
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

	req.RoleId = strutil.SimplifyStringList(req.RoleId)
	req.GroupId = strutil.SimplifyStringList(req.GroupId)
	req.UserId = strutil.SimplifyStringList(req.UserId)
	req.Username = strutil.SimplifyStringList(req.Username)
	req.Email = strutil.SimplifyStringList(req.Email)
	req.PhoneNumber = strutil.SimplifyStringList(req.PhoneNumber)
	req.Status = strutil.SimplifyStringList(req.Status)

	if len(req.RoleId) > 0 {
		if !validator.IsValidId(req.RoleId...) {
			err := status.Errorf(codes.InvalidArgument, "invalid RoleId: %v", req.RoleId)
			logger.Warnf(ctx, "%+v", err)
			return nil, err
		}
	}
	if len(req.GroupId) > 0 {
		if !validator.IsValidId(req.GroupId...) {
			err := status.Errorf(codes.InvalidArgument, "invalid GroupId: %v", req.GroupId)
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
	if len(req.Username) > 0 {
		if !validator.IsValidName(req.Username...) {
			err := status.Errorf(codes.InvalidArgument, "invalid UserName: %v", req.Username)
			logger.Warnf(ctx, "%+v", err)
			return nil, err
		}
	}
	if len(req.Email) > 0 {
		if !validator.IsValidEmail(req.Email...) {
			err := status.Errorf(codes.InvalidArgument, "invalid Email: %v", req.Email)
			logger.Warnf(ctx, "%+v", err)
			return nil, err
		}
	}
	if len(req.PhoneNumber) > 0 {
		if !validator.IsValidPhoneNumber(req.PhoneNumber...) {
			err := status.Errorf(codes.InvalidArgument, "invalid PhoneNumber: %v", req.PhoneNumber)
			logger.Warnf(ctx, "%+v", err)
			return nil, err
		}
	}
	if len(req.Status) > 0 {
		if !validator.IsValidStatus(req.Status...) {
			err := status.Errorf(codes.InvalidArgument, "invalid Status: %v", req.Status)
			logger.Warnf(ctx, "%+v", err)
			return nil, err
		}
	}

	// 1. if len(RoleId) > 0, get newUserId by RoleId
	// 2. merge req.UserId and newUserId
	// 3. query from IM server
	// 4. range users from IM, get role info from role table

	// 1. if len(RoleId) > 0, get newUserId by RoleId
	// 2. merge req.UserId and newUserId
	if len(req.RoleId) > 0 {
		query, err := template.Render(`
		SELECT DISTINCT * FROM user_role_binding WHERE 1=1
			AND user_role_binding.role_id in (
				{{range $i, $v := .RoleId}}
					{{if eq $i 0}} '{{$v}}' {{else}} ,'{{$v}}' {{end}}
				{{end}}
			)
		`, req,
		)
		if err != nil {
			logger.Warnf(ctx, "%+v", err)
			return nil, err
		}

		query = strutil.SimplifyString(query)
		logger.Infof(ctx, "query: %s", query)

		var rows []db_spec.UserRoleBinding
		p.DB.Raw(query).Find(&rows)
		if err := p.DB.Error; err != nil {
			logger.Warnf(ctx, "%+v", err)
			return nil, err
		}

		// update req.UserId
		if len(req.UserId) > 0 {
			var newUserIdList []string
			for _, uid := range req.UserId {
				for _, v := range rows {
					if uid == v.UserId {
						newUserIdList = append(newUserIdList, uid)
						break
					}
				}
			}
			req.UserId = newUserIdList
		} else {
			for _, v := range rows {
				req.UserId = append(req.UserId, v.UserId)
			}
		}
	}

	// query users from IM server
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", p.cfg.ImHost, p.cfg.ImPort), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := pbim.NewAccountManagerClient(conn)
	imReply, err := client.ListUsers(ctx, &pbim.ListUsersRequest{
		GroupId:     req.GroupId,
		UserId:      req.UserId,
		Username:    req.Username,
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
		Status:      req.Status,

		Limit:  req.Limit,
		Offset: req.Offset,

		SearchWord: req.SearchWord,
		SortKey:    req.SortKey,
		Reverse:    req.Reverse,
	})
	if err != nil {
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	imUserList := imReply.User
	if len(imUserList) == 0 {
		return &pbam.DescribeUsersWithRoleResponse{}, nil
	}

	query, err := template.Render(`
			SELECT DISTINCT user_role_binding.user_id, role.* FROM role, user_role_binding WHERE 1=1
				AND user_role_binding.role_id=role.role_id
				AND user_role_binding.user_id in (
					{{range $i, $v := .}}
						{{if eq $i 0}} '{{$v.UserId}}' {{else}} ,'{{$v.UserId}}' {{end}}
					{{end}}
				)
			`, imUserList,
	)
	if err != nil {
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	query = strutil.SimplifyString(query)
	logger.Infof(ctx, "query: %s", query)

	var roleExList []struct {
		UserId string
		db_spec.Role
	}
	p.DB.Raw(query).Find(&roleExList)
	if err := p.DB.Error; err != nil {
		logger.Warnf(ctx, "%+v", err)
		// ignore err
	}

	var userList []*pbam.UserWithRole
	for _, imUser := range imUserList {
		user := &pbam.UserWithRole{
			UserId:      imUser.UserId,
			Username:    imUser.Username,
			Email:       imUser.Email,
			PhoneNumber: imUser.PhoneNumber,
			Description: imUser.Description,
			Status:      imUser.Status,
			Extra:       imUser.Extra,
			CreateTime:  imUser.CreateTime,
			UpdateTime:  imUser.UpdateTime,
			StatusTime:  imUser.StatusTime,
		}

		for _, roleEx := range roleExList {
			if roleEx.UserId == imUser.UserId {
				user.Role = append(user.Role, roleEx.Role.ToPB())
			}
		}

		userList = append(userList, user)
	}

	reply := &pbam.DescribeUsersWithRoleResponse{
		User:  userList,
		Total: imReply.Total,
	}

	return reply, nil
}
