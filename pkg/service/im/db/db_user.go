// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/golang/protobuf/ptypes"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"openpitrix.io/iam/pkg/internal/funcutil"
	"openpitrix.io/iam/pkg/pb/im"
	"openpitrix.io/iam/pkg/service/im/db_spec"
	"openpitrix.io/logger"
)

func (p *Database) CreateUser(ctx context.Context, req *pbim.User) (*pbim.User, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	if req == nil {
		err := status.Errorf(codes.InvalidArgument, "empty field")
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}
	if req != nil {
		if req.Uid == "" {
			req.Uid = genUid()
		}

		if isZeroTimestamp(req.CreateTime) {
			req.CreateTime = ptypes.TimestampNow()
		}
		if isZeroTimestamp(req.UpdateTime) {
			req.UpdateTime = ptypes.TimestampNow()
		}
		if isZeroTimestamp(req.StatusTime) {
			req.StatusTime = ptypes.TimestampNow()
		}
	}

	if req.Uid == "" || req.Password == "" {
		err := status.Errorf(codes.InvalidArgument, "empty uid or password")
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	// TODO: check group_path valid

	var dbUser = db_spec.PBUserToDB(req)
	if err := dbUser.ValidateForInsert(); err != nil {
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	hashedPass, err := bcrypt.GenerateFromPassword(
		[]byte(dbUser.Password), bcrypt.DefaultCost,
	)
	if err != nil {
		err := status.Errorf(codes.Internal, "bcrypt failed")
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}
	dbUser.Password = string(hashedPass)

	sql, values := pkgBuildSql_InsertInto(
		db_spec.UserTableName,
		dbUser,
	)
	if len(values) == 0 {
		err := status.Errorf(codes.InvalidArgument, "empty field")
		logger.Warnf(ctx, "%v, %v", sql, values)
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	_, err = p.DB.ExecContext(ctx, sql, values...)
	if err != nil {
		logger.Warnf(ctx, "%v, %v", sql, values)
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	return req, nil
}

func (p *Database) DeleteUsers(ctx context.Context, req *pbim.UserIdList) (*pbim.Empty, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	if len(req.Uid) == 1 && strings.Contains(req.Uid[0], ",") {
		req.Uid = strings.Split(req.Uid[0], ",")
	}

	if req == nil || len(req.Uid) == 0 || !isValidUids(req.Uid...) {
		err := status.Errorf(codes.InvalidArgument, "empty field")
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	sql := pkgBuildSql_Delete(
		db_spec.UserTableName,
		db_spec.UserPrimaryKeyName,
		req.Uid...,
	)

	tx, err := p.DB.Beginx()
	if err != nil {
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	_, err = tx.ExecContext(ctx, sql)
	if err != nil {
		logger.Warnf(ctx, "%v", sql)
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	// delete binding
	for _, uid := range req.Uid {
		sql := fmt.Sprintf(
			`delete from %s where user_id=?`,
			db_spec.UserGroupBindingTableName,
		)

		_, err := tx.ExecContext(ctx, sql, uid)
		if err != nil {
			logger.Warnf(ctx, "%v", sql)
			logger.Warnf(ctx, "%+v", err)
			return nil, err
		}
		if err != nil {
			logger.Warnf(ctx, "uid = %v, err = %+v", uid, err)
		}
	}

	err = tx.Commit()
	if err != nil {
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	reply := &pbim.Empty{}
	return reply, nil
}

func (p *Database) ModifyUser(ctx context.Context, req *pbim.User) (*pbim.User, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	if req == nil || req.Uid == "" {
		err := status.Errorf(codes.InvalidArgument, "empty field")
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	var dbUser = db_spec.PBUserToDB(req)

	// ignore read only field
	{
		dbUser.Password = ""

		dbUser.CreateTime = time.Time{}
		dbUser.UpdateTime = time.Now()

		switch {
		case dbUser.Status == "":
			dbUser.StatusTime = time.Time{}
		default:
			dbUser.StatusTime = time.Now()
		}
	}

	if err := dbUser.ValidateForUpdate(); err != nil {
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	sql, values := pkgBuildSql_Update(
		db_spec.UserTableName, dbUser,
		db_spec.UserPrimaryKeyName,
	)
	if len(values) == 0 {
		return p.GetUser(ctx, &pbim.UserId{Uid: req.Uid})
	}

	_, err := p.DB.ExecContext(ctx, sql, values...)
	if err != nil {
		logger.Warnf(ctx, "%v, %v", sql, values)
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	return p.GetUser(ctx, &pbim.UserId{Uid: req.Uid})
}

func (p *Database) GetUser(ctx context.Context, req *pbim.UserId) (*pbim.User, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	var query = fmt.Sprintf(
		"SELECT * FROM %s WHERE %s=? LIMIT 1 OFFSET 0;",
		db_spec.UserTableName,
		db_spec.UserPrimaryKeyName,
	)

	var v = db_spec.DBUser{}
	err := p.DB.GetContext(ctx, &v, query, req.GetUid())
	if err != nil {
		logger.Warnf(ctx, "%v", query)
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
	if len(req.Gid) == 1 && strings.Contains(req.Gid[0], ",") {
		req.Gid = strings.Split(req.Gid[0], ",")
	}
	if len(req.Uid) == 1 && strings.Contains(req.Uid[0], ",") {
		req.Uid = strings.Split(req.Uid[0], ",")
	}
	if len(req.Name) == 1 && strings.Contains(req.Name[0], ",") {
		req.Name = strings.Split(req.Name[0], ",")
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

	if err := p.validateListUsersReq(req); err != nil {
		return nil, err
	}

	if len(req.Gid) > 0 {
		return p.listUsers_with_gid(ctx, req)
	} else {
		return p.listUsers_no_gid(ctx, req)
	}
}
