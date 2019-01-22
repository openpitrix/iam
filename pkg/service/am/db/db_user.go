// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"openpitrix.io/iam/pkg/internal/funcutil"
	pbam "openpitrix.io/iam/pkg/pb/am"
	pbim "openpitrix.io/iam/pkg/pb/im"
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

	user := &pbam.UserWithRole{
		UserId:      imUser.UserId,
		UserName:    imUser.UserName,
		Email:       imUser.Email,
		PhoneNumber: imUser.PhoneNumber,
		Description: imUser.Description,
		Status:      imUser.Status,
		Extra:       imUser.Extra,
		CreateTime:  imUser.CreateTime,
		UpdateTime:  imUser.UpdateTime,
		StatusTime:  imUser.StatusTime,
	}

	roleList, err := p.GetRoleListByUserId(ctx, &pbam.UserId{UserId: req.UserId})
	if err != nil {
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	user.Role = roleList.Value
	return user, nil
}
func (p *Database) DescribeUsersWithRole(ctx context.Context, req *pbam.DescribeUsersWithRoleRequest) (*pbam.DescribeUsersWithRoleResponse, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	// fix repeated fileds
	// GET donot support repeated type in grpc-gateway
	if len(req.RoleId) == 1 && strings.Contains(req.RoleId[0], ",") {
		req.RoleId = strings.Split(req.RoleId[0], ",")
	}
	if len(req.UserId) == 1 && strings.Contains(req.UserId[0], ",") {
		req.UserId = strings.Split(req.UserId[0], ",")
	}

	req.RoleId = simplifyStringList(req.RoleId)
	req.UserId = simplifyStringList(req.UserId)

	if len(req.RoleId) == 0 && len(req.UserId) == 0 {
		err := status.Errorf(codes.InvalidArgument, "empty user_id or role_id")
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	limit := 20
	offset := 0
	if req.Limit > 0 || req.Offset > 0 {
		limit = int(req.Limit)
		offset = int(req.Offset)
	}

	type Result struct {
		UserId string
		RoleId string
	}

	var rows []Result
	err := p.DB.Raw(
		"select * from user_role_binding where user_id in (?) AND role_id in (?)",
		req.UserId, req.RoleId,
	).Limit(limit).Offset(offset).Find(&rows).Error
	if err != nil {
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	var userIdList, roleIdList []string
	for _, v := range rows {
		userIdList = append(userIdList, v.UserId)
		roleIdList = append(roleIdList, v.RoleId)
	}
	rawUsers, err := p.getUserList(ctx, userIdList...)
	if err != nil {
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	type RoleEx struct {
		UserId string
		Role
	}

	var roles []RoleEx
	err = p.DB.Raw(
		`select t2.user_id, t1.*
			from role t1, user_role_binding t2
			where t1.role_id=t2.role_id and t2.user_id in (?)
		`,
		userIdList,
	).Find(&roles).Error
	if err != nil {
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	var users = make([]*pbam.UserWithRole, len(rawUsers))
	for i := 0; i < len(rawUsers); i++ {
		var pRole *Role
		for j := 0; j < len(roles); j++ {
			if rawUsers[i].UserId == roles[j].UserId {
				pRole = &roles[j].Role
			}
		}

		users[i] = &pbam.UserWithRole{
			UserId:      rawUsers[i].UserId,
			UserName:    rawUsers[i].UserName,
			Email:       rawUsers[i].Email,
			PhoneNumber: rawUsers[i].PhoneNumber,
			Description: rawUsers[i].Description,
			Status:      rawUsers[i].Status,
			Extra:       rawUsers[i].Extra,
		}
		if pRole != nil {
			users[i].Role = []*pbam.Role{pRole.ToPB()}
		}
	}

	reply := &pbam.DescribeUsersWithRoleResponse{
		User:  users,
		Total: int32(len(rows)),
	}

	return reply, nil
}

func (p *Database) getUserList(ctx context.Context, uid ...string) ([]*pbim.User, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", p.cfg.ImHost, p.cfg.ImPort), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := pbim.NewAccountManagerClient(conn)
	reply, err := client.ListUsers(ctx, &pbim.ListUsersRequest{UserId: uid})
	if err != nil {
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	return reply.User, nil
}
