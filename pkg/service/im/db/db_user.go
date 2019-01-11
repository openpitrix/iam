// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db

import (
	"context"
	"fmt"

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

	if err := req.Validate(); err != nil {
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

	if req == nil || len(req.Uid) == 0 || !isValidIds(req.Uid...) {
		err := status.Errorf(codes.InvalidArgument, "empty field")
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	sql := pkgBuildSql_Delete(
		db_spec.UserTableName,
		db_spec.UserPrimaryKeyName,
		req.Uid...,
	)

	_, err := p.DB.ExecContext(ctx, sql)
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

		_, err := p.DB.ExecContext(ctx, sql, uid)
		if err != nil {
			logger.Warnf(ctx, "%v", sql)
			logger.Warnf(ctx, "%+v", err)
			return nil, err
		}
		if err != nil {
			logger.Warnf(ctx, "uid = %v, err = %+v", uid, err)
		}
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

	// ignore Password
	dbUser.Password = ""

	if err := dbUser.ValidateForUpdate(); err != nil {
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	sql, values := pkgBuildSql_Update(
		db_spec.UserTableName, dbUser,
		db_spec.UserPrimaryKeyName,
	)

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

func (p *Database) ListUsers(ctx context.Context, req *pbim.Range) (*pbim.ListUsersResponse, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	if req.GetVerbose() {
		err := status.Errorf(codes.Unimplemented, "unsupport range.Verbose")
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	if req.GetSearchWord() == "" {
		return p._ListUsers_all(ctx, req)
	} else {
		return p._ListUsers_bySearchWord(ctx, req)
	}
}
