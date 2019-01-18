// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db

import (
	"context"

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

	tx := p.DB.Begin()
	{
		tx.Raw("delete from role where role_id in (?)", req.RoleId)
		tx.Raw("delete from user_role_binding where role_id in (?)", req.RoleId)
		tx.Raw("delete from role_module_binding where role_id in (?)", req.RoleId)
	}
	if err := tx.Commit().Error; err != nil {
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	return &pbam.Empty{}, nil
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
func (p *Database) GetRoleListByUserId(ctx context.Context, req *pbam.UserId) (*pbam.RoleList, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	var query = `
		select
			role.*
		from
			role,
			user_role_binding
		where
			role.role_id=user_role_binding.role_id and
			user_role_binding.user_id=?
	`

	var rows []DBRole
	if err := p.DB.Raw(query, req.UserId).Find(&rows).Error; err != nil {
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	var sets []*pbam.Role
	for _, v := range rows {
		sets = append(sets, v.ToPB())
	}

	reply := &pbam.RoleList{
		Value: sets,
	}
	return reply, nil
}

func (p *Database) DescribeRoles(ctx context.Context, req *pbam.DescribeRolesRequest) (*pbam.RoleList, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	var rows []DBRole
	p.DB.Raw(
		`select distinct t1.*
			from  role t1
			where t1.role_id  in(?)
			and t1.role_name in (?)
			and t1.portal in (?)
			and t1.role_id in
				(select t2.role_id
					from user_role_binding t1, role t2
					where  t1.role_id=t2.role_id and t1.user_id in (?)
				)
		`,
		req.RoleId,
		req.RoleName,
		req.Portal,
		req.UserId,
	).Find(&rows)

	var sets []*pbam.Role
	for _, v := range rows {
		sets = append(sets, v.ToPB())
	}

	reply := &pbam.RoleList{
		Value: sets,
	}

	return reply, nil
}
