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

	req.GroupId = simplifyStringList(req.GroupId)
	req.UserId = simplifyStringList(req.UserId)
	req.UserName = simplifyStringList(req.UserName)
	req.Email = simplifyStringList(req.Email)
	req.PhoneNumber = simplifyStringList(req.PhoneNumber)
	req.Status = simplifyStringList(req.Status)

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

	var (
		inKeys   []string
		inValues []interface{}
	)

	if len(req.GroupId) > 0 {
		inKeys = append(inKeys, "user_group.group_id in(?)")
		inValues = append(inValues, req.GroupId)
	}
	if len(req.UserId) > 0 {
		inKeys = append(inKeys, "user.user_id in(?)")
		inValues = append(inValues, req.UserId)
	}
	if len(req.UserName) > 0 {
		inKeys = append(inKeys, "user.user_name in(?)")
		inValues = append(inValues, req.UserName)
	}
	if len(req.Email) > 0 {
		inKeys = append(inKeys, "user.email in(?)")
		inValues = append(inValues, req.Email)
	}
	if len(req.PhoneNumber) > 0 {
		inKeys = append(inKeys, "user.phone_number in(?)")
		inValues = append(inValues, req.PhoneNumber)
	}
	if len(req.Status) > 0 {
		inKeys = append(inKeys, "user.status in(?)")
		inValues = append(inValues, req.Status)
	}

	if req.SearchWord != "" {
		var likeKey = "%" + req.SearchWord + "%"

		var likeSql = `(1=0`
		likeSql += " OR user.user_name LIKE ?"
		likeSql += " OR user.email LIKE ?"
		likeSql += " OR user.phone_number LIKE ?"
		likeSql += " OR user.description LIKE ?"
		likeSql += " OR user.status LIKE ?"
		likeSql += ")"

		inKeys = append(inKeys, likeSql)
		inValues = append(inValues,
			likeKey, // user_name
			likeKey, // email
			likeKey, // phone_number
			likeKey, // description
			likeKey, // status
		)
	}

	var query = ""
	if len(inKeys) > 0 {
		if len(req.UserId) > 0 {
			query += "SELECT user.* from user, user_group, user_group_binding"
			query += " WHERE "
			query += " user_group_binding.user_id=user.user_id AND"
			query += " user_group_binding.group_id=user_group.group_id AND"
			query += strings.Join(inKeys, " AND ")
		} else {
			query += "SELECT user.* from user_group"
			query += " WHERE "
			query += strings.Join(inKeys, " AND ")
		}
	} else {
		query = "SELECT * from user WHERE 1=1"
		inValues = nil
	}

	logger.Infof(ctx, "query: %s", query)
	logger.Infof(ctx, "inValues: %v", inValues)

	var rows []User
	p.DB.Limit(req.Limit).Offset(req.Offset).Where(query, inValues...).Find(&rows)
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
