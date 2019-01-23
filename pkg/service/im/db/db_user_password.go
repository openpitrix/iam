// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db

import (
	"context"

	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"openpitrix.io/iam/pkg/internal/funcutil"
	"openpitrix.io/iam/pkg/pb/im"
	"openpitrix.io/logger"
)

func (p *Database) ComparePassword(ctx context.Context, req *pbim.Password) (*pbim.Bool, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	var user User
	if err := p.DB.First(&user).Error; err != nil {
		logger.Warnf(ctx, "uid = %s, err = %+v", req.UserId, err)
		return nil, err
	}

	err := bcrypt.CompareHashAndPassword(
		[]byte(user.Password), []byte(req.GetPassword()),
	)
	if err != nil {
		return &pbim.Bool{Value: false}, nil
	}

	// OK
	return &pbim.Bool{Value: true}, nil
}

func (p *Database) ModifyPassword(ctx context.Context, req *pbim.Password) (*pbim.Empty, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	var dbUser = &User{
		UserId:   req.UserId,
		Password: req.Password,
	}

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

	if err := p.DB.Model(dbUser).Updates(dbUser).Error; err != nil {
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	return &pbim.Empty{}, nil
}
