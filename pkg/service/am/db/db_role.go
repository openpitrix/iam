// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db

import (
	"context"
	"strings"

	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"openpitrix.io/iam/pkg/internal/funcutil"
	pbam "openpitrix.io/iam/pkg/pb/am"
	"openpitrix.io/logger"
)

func (p *Database) CreateRole(ctx context.Context, req *pbam.Role) (*pbam.Role, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	if req == nil {
		err := status.Errorf(codes.InvalidArgument, "empty field")
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}
	if req != nil {
		if req.RoleId == "" {
			req.RoleId = genRoleId()
		}

		if isZeroTimestamp(req.CreateTime) {
			req.CreateTime = ptypes.TimestampNow()
		}
		if isZeroTimestamp(req.UpdateTime) {
			req.UpdateTime = ptypes.TimestampNow()
		}
	}

	if !p.DB.NewRecord(PBRoleToDB(req)) {
		// failed
	}
	if err := p.DB.Error; err != nil {
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	return p.GetRole(ctx, &pbam.RoleId{RoleId: req.RoleId})
}
func (p *Database) DeleteRoles(ctx context.Context, req *pbam.RoleIdList) (*pbam.Empty, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	logger.Infof(ctx, "TODO")
	return nil, nil
}

func (p *Database) ModifyRole(ctx context.Context, req *pbam.Role) (*pbam.Role, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	if req == nil || req.RoleId == "" {
		err := status.Errorf(codes.InvalidArgument, "empty field")
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	logger.Infof(ctx, "TODO")
	return nil, nil
}

func (p *Database) GetRole(ctx context.Context, req *pbam.RoleId) (*pbam.Role, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	var role = DBRole{
		RoleId: req.RoleId,
	}

	if err := p.DB.First(&role).Error; err != nil {
		return nil, err
	}
	return role.ToPB(), nil
}

func (p *Database) DescribeRoles(ctx context.Context, req *pbam.DescribeRolesRequest) (*pbam.RoleList, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	// fix repeated fileds
	if len(req.RoleId) == 1 && strings.Contains(req.RoleId[0], ",") {
		req.RoleId = strings.Split(req.RoleId[0], ",")
	}
	if len(req.RoleName) == 1 && strings.Contains(req.RoleName[0], ",") {
		req.RoleName = strings.Split(req.RoleName[0], ",")
	}
	if len(req.Portal) == 1 && strings.Contains(req.Portal[0], ",") {
		req.Portal = strings.Split(req.Portal[0], ",")
	}
	if len(req.UserId) == 1 && strings.Contains(req.UserId[0], ",") {
		req.UserId = strings.Split(req.UserId[0], ",")
	}

	logger.Infof(ctx, "TODO")
	return nil, nil
}
