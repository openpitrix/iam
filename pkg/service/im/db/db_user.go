// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db

import (
	"context"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"openpitrix.io/iam/pkg/internal/funcutil"
	"openpitrix.io/iam/pkg/pb/im"
	"openpitrix.io/logger"
)

func (p *Database) CreateUser(ctx context.Context, req *pbim.User) (*pbim.User, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	var dbUser = NewUserFromPB(req)
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

	req.GroupId = trimEmptyString(req.GroupId)
	req.UserId = trimEmptyString(req.UserId)
	req.UserName = trimEmptyString(req.UserName)
	req.Email = trimEmptyString(req.Email)
	req.PhoneNumber = trimEmptyString(req.PhoneNumber)
	req.Status = trimEmptyString(req.Status)

	if err := p.validateListUsersReq(req); err != nil {
		return nil, err
	}

	if len(req.GroupId) > 0 {
		return p.listUsers_with_gid(ctx, req)
	} else {
		return p.listUsers_no_gid(ctx, req)
	}
}
