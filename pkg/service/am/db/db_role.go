// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db

import (
	"context"
	"fmt"
	"strings"
	"time"

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
			req.RoleId = genId("role-", 12)
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
		if err := tx.Error; err != nil {
			tx.Rollback()
			return nil, err
		}

		tx.Raw("delete from user_role_binding where role_id in (?)", req.RoleId)
		if err := tx.Error; err != nil {
			tx.Rollback()
			return nil, err
		}

		tx.Raw("delete from role_module_binding where role_id in (?)", req.RoleId)
		if err := tx.Error; err != nil {
			tx.Rollback()
			return nil, err
		}
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

	var m = map[string]interface{}{}
	if req.RoleName != "" {
		m["role_name"] = req.RoleName
	}
	if req.Description != "" {
		m["description"] = req.Description
	}
	if req.Portal != "" {
		m["portal"] = req.Portal
	}

	if len(m) == 0 {
		err := status.Errorf(codes.InvalidArgument, "empty field")
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	m["update_time"] = time.Now()
	if err := p.DB.Table("role").Update(m).Error; err != nil {
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	return p.GetRole(ctx, &pbam.RoleId{RoleId: req.RoleId})
}

func (p *Database) GetRole(ctx context.Context, req *pbam.RoleId) (*pbam.Role, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	var role = Role{
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

	var rows []Role
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

	// fix repeated fileds
	// GET donot support repeated type in grpc-gateway
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

	req.RoleId = trimEmptyString(req.RoleId)
	req.RoleName = trimEmptyString(req.RoleName)
	req.Portal = trimEmptyString(req.Portal)
	req.UserId = trimEmptyString(req.UserId)

	var (
		args                 []interface{}
		sqlRoleIdCondition   string
		sqlRoleNameCondition string
		sqlPortalCondition   string
		sqlUserIdCondition   string
	)

	if len(req.RoleId) > 0 {
		sqlRoleIdCondition = `and t1.role_id in(?)`
		args = append(args, req.RoleId)
	}
	if len(req.RoleName) > 0 {
		sqlRoleNameCondition = `and t1.role_name in (?)`
		args = append(args, req.RoleName)
	}
	if len(req.Portal) > 0 {
		sqlPortalCondition = `and t1.portal in (?)`
		args = append(args, req.Portal)
	}
	if len(req.UserId) > 0 {
		sqlUserIdCondition = `and t1.user_id in (?)`
		args = append(args, req.UserId)
	}

	query := fmt.Sprintf(
		`select distinct t1.* from role t1 where 1=1
			%s -- and t1.role_id in(?)
			%s -- and t1.role_name in (?)
			%s -- and t1.portal in (?) {{/Portal}}
			and t1.role_id in
				(select t2.role_id
					from user_role_binding t1, role t2
					where t1.role_id=t2.role_id
						%s -- and t1.user_id in (?)
				)
		`,
		sqlRoleIdCondition,
		sqlRoleNameCondition,
		sqlPortalCondition,
		sqlUserIdCondition,
	)

	var rows []Role
	p.DB.Raw(query, args...).Find(&rows)

	var sets []*pbam.Role
	for _, v := range rows {
		sets = append(sets, v.ToPB())
	}

	reply := &pbam.RoleList{
		Value: sets,
	}

	return reply, nil
}
