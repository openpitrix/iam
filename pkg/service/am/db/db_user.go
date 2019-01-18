// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db

import (
	"context"
	"fmt"

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
	imUser, err := client.GetUser(ctx, &pbim.UserId{Uid: req.UserId})
	if err != nil {
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	user := &pbam.UserWithRole{
		UserId:      imUser.Uid,
		UserName:    imUser.Name,
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

	if len(req.RoleId) > 1 && len(req.UserId) > 1 {
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
	err := p.DB.Table("user_role_binding").Where("user_id in (?) AND role_id in (?)",
		req.UserId, req.RoleId,
	).Limit(limit).Offset(offset).Find(&rows).Error
	if err != nil {
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	var users []*pbam.UserWithRole
	for i := 0; i < len(rows); i++ {
		v, err := p.GetUserWithRole(ctx, &pbam.UserId{UserId: rows[i].UserId})
		if err != nil {
			logger.Warnf(ctx, "%+v", err)
			return nil, err
		}
		users = append(users, v)
	}

	reply := &pbam.DescribeUsersWithRoleResponse{
		User:  users,
		Total: int32(len(rows)),
	}

	return reply, nil
}
