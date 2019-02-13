// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db

import (
	"context"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	idpkg "openpitrix.io/iam/pkg/id"
	"openpitrix.io/iam/pkg/internal/funcutil"
	"openpitrix.io/iam/pkg/internal/strutil"
	pbam "openpitrix.io/iam/pkg/pb/am"
	"openpitrix.io/iam/pkg/service/am/db_spec"
	"openpitrix.io/logger"
)

func (p *Database) BindUserRole(ctx context.Context, req *pbam.BindUserRoleRequest) (*pbam.Empty, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	if len(req.RoleId) == 1 && strings.Contains(req.RoleId[0], ",") {
		req.RoleId = strings.Split(req.RoleId[0], ",")
	}
	if len(req.UserId) == 1 && strings.Contains(req.UserId[0], ",") {
		req.UserId = strings.Split(req.UserId[0], ",")
	}

	req.RoleId = strutil.SimplifyStringList(req.RoleId)
	req.UserId = strutil.SimplifyStringList(req.UserId)

	if len(req.UserId) == 0 || len(req.RoleId) == 0 {
		err := status.Errorf(codes.InvalidArgument, "empty user_id or role_id")
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}
	if !(len(req.UserId) == 1 || len(req.RoleId) == 1 || len(req.UserId) == len(req.RoleId)) {
		err := status.Errorf(codes.InvalidArgument,
			"user_id and role_id donot math: user_id = %v, role_id = %v",
			req.UserId, req.RoleId,
		)
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	// get id list
	idPairList := p.getUserRoleBindList(ctx, req.UserId, req.RoleId)
	if len(idPairList) == 0 {
		err := status.Errorf(codes.InvalidArgument,
			"empty UserId or RoleId: %v, %v",
			req.UserId, req.RoleId,
		)
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}
	if len(idPairList) == 0 {
		return &pbam.Empty{}, nil
	}
	logger.Warnf(ctx, "%+v", idPairList)

	// id list to id map
	var (
		idPairMap = make(map[string]string)
		// roleIdMap = make(map[string]string)
	)
	for _, v := range idPairList {
		if _, ok := idPairMap[v.UserId]; ok {
			err := status.Errorf(codes.InvalidArgument,
				"a user can only have one role: %v, %v",
				v.UserId, v.RoleId,
			)
			logger.Warnf(ctx, "%+v", err)
			// ignore err
		}

		idPairMap[v.UserId] = v.RoleId
		// roleIdMap[v.RoleId]= v.RoleId
	}
	logger.Warnf(ctx, "%+v", idPairMap)

	// TODO: check user_id valid
	// TODO: check role_id valid

	/*
		// get exists bind_id
		var query, err = template.Render(`
			SELECT * FROM user_role_binding WHERE 1=0
				{{range $i, $v := .}}
					OR (role_id='{{$v.RoleId}}' AND user_id='{{$v.UserId}}')
				{{end}}
			`, idPairList,
		)
		if err != nil {
			logger.Warnf(ctx, "%+v", err)
			return nil, err
		}

		query = strutil.SimplifyString(query)
		logger.Infof(ctx, "query: %s", query)

		var rows []db_spec.UserRoleBinding
		p.DB.Raw(query).Find(&rows)
		if err := p.DB.Error; err != nil {
			logger.Warnf(ctx, "%+v", err)
			return nil, err
		}

		// read bind_id
		for _, vi := range rows {
			for j, vj := range idPairList {
				if vi.UserId == vj.UserId && vi.RoleId == vj.RoleId {
					idPairList[j].Id = vi.Id
				}
			}
		}
	*/

	// 2. insert new bind
	tx := p.DB.Begin()
	{
		// delete old bind
		for _, v := range idPairList {
			tx.Exec(
				`DELETE FROM user_role_binding WHERE user_id=? and role_id=?`,
				v.UserId, v.RoleId,
			)
			if err := tx.Error; err != nil {
				tx.Rollback()
				return nil, err
			}
		}

		// insert new bind
		for _, v := range idPairList {
			tx.Exec(
				`INSERT INTO user_role_binding (id, user_id, role_id) VALUES (?,?,?)`,
				idpkg.GenId("xid-"), v.UserId, v.RoleId,
			)
			if err := tx.Error; err != nil {
				tx.Rollback()
				return nil, err
			}
		}
	}
	if err := tx.Commit().Error; err != nil {
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	return &pbam.Empty{}, nil
}

func (p *Database) UnbindUserRole(ctx context.Context, req *pbam.UnbindUserRoleRequest) (*pbam.Empty, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	if len(req.RoleId) == 1 && strings.Contains(req.RoleId[0], ",") {
		req.RoleId = strings.Split(req.RoleId[0], ",")
	}
	if len(req.UserId) == 1 && strings.Contains(req.UserId[0], ",") {
		req.UserId = strings.Split(req.UserId[0], ",")
	}

	req.RoleId = strutil.SimplifyStringList(req.RoleId)
	req.UserId = strutil.SimplifyStringList(req.UserId)

	if len(req.UserId) == 0 || len(req.RoleId) == 0 {
		err := status.Errorf(codes.InvalidArgument, "empty user_id or role_id")
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}
	if !(len(req.UserId) == 1 || len(req.RoleId) == 1 || len(req.UserId) == len(req.RoleId)) {
		err := status.Errorf(codes.InvalidArgument,
			"user_id and role_id donot math: user_id = %v, role_id = %v",
			req.UserId, req.RoleId,
		)
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	// get id list
	idPairList := p.getUserRoleBindList(ctx, req.UserId, req.RoleId)
	if len(idPairList) == 0 {
		err := status.Errorf(codes.InvalidArgument,
			"empty UserId or RoleId: %v, %v",
			req.UserId, req.RoleId,
		)
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	tx := p.DB.Begin()
	{
		for _, v := range idPairList {
			tx.Delete(db_spec.UserRoleBinding{},
				`user_id=? and role_id=?`,
				v.UserId, v.RoleId,
			)
			if err := tx.Error; err != nil {
				tx.Rollback()
				return nil, err
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}

	return &pbam.Empty{}, nil
}

func (p *Database) getUserRoleBindList(ctx context.Context, userId, roleId []string) (
	s []db_spec.UserRoleBinding,
) {
	if len(userId) == 0 || len(roleId) == 0 {
		return nil
	}
	if !(len(userId) == 1 || len(roleId) == 1 || len(userId) == len(roleId)) {
		return nil
	}

	if len(userId) == len(roleId) {
		for i := 0; i < len(userId); i++ {
			s = append(s, db_spec.UserRoleBinding{
				UserId: userId[i],
				RoleId: roleId[i],
			})
		}
		return
	}
	if len(userId) == 1 {
		for i := 0; i < len(roleId); i++ {
			s = append(s, db_spec.UserRoleBinding{
				UserId: userId[0],
				RoleId: roleId[i],
			})
		}
		return
	}
	if len(roleId) == 1 {
		for i := 0; i < len(userId); i++ {
			s = append(s, db_spec.UserRoleBinding{
				UserId: userId[i],
				RoleId: roleId[0],
			})
		}
		return
	}

	logger.Errorf(ctx, "unreachable, should panic")
	return
}
