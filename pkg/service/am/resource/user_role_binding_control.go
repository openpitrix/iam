// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package resource

import (
	"context"

	"openpitrix.io/iam/pkg/constants"
	"openpitrix.io/iam/pkg/gerr"
	"openpitrix.io/iam/pkg/global"
	"openpitrix.io/iam/pkg/models"
	"openpitrix.io/iam/pkg/pb"
)

func GetUserRoleBindings(ctx context.Context, userIds, roleIds []string) ([]*models.UserRoleBinding, error) {
	var userRoleBindings []*models.UserRoleBinding
	if err := global.Global().Database.Table(constants.TableUserRoleBinding).
		Where(constants.ColumnRoleId+" in (?)", roleIds).
		Where(constants.ColumnUserId+" in (?)", userIds).
		Find(&userRoleBindings).
		Error; err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}

	return userRoleBindings, nil
}

func GetRoleIdsByUserIds(ctx context.Context, userIds []string) ([]string, error) {
	rows, err := global.Global().Database.Table(constants.TableUserRoleBinding).
		Select(constants.ColumnRoleId).
		Where(constants.ColumnUserId+" in (?)", userIds).
		Rows()
	defer rows.Close()
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}
	var roleIds []string
	for rows.Next() {
		var roleId string
		rows.Scan(&roleId)
		roleIds = append(roleIds, roleId)
	}
	return roleIds, nil
}

func GetUserIdsByRoleIds(ctx context.Context, roleIds []string) ([]string, error) {
	rows, err := global.Global().Database.Table(constants.TableUserRoleBinding).
		Select(constants.ColumnUserId).
		Where(constants.ColumnRoleId+" in (?)", roleIds).
		Rows()
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}
	var userIds []string
	for rows.Next() {
		var userId string
		rows.Scan(&userId)
		userIds = append(userIds, userId)
	}
	return userIds, nil
}

func BindUserRole(ctx context.Context, req *pb.BindUserRoleRequest) (*pb.BindUserRoleResponse, error) {
	if len(req.UserId) == 0 || len(req.RoleId) == 0 {
		return nil, gerr.New(ctx, gerr.InvalidArgument, gerr.ErrorMissingParameter, "userId, roleId")
	}

	if len(req.RoleId) > 1 {
		return nil, gerr.New(ctx, gerr.PermissionDenied, gerr.ErrorCannotBindRole)
	}

	// one user can only bind to one role
	roleIds, err := GetRoleIdsByUserIds(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	if len(roleIds) > 0 {
		return nil, gerr.New(ctx, gerr.PermissionDenied, gerr.ErrorCannotBindRole)
	}

	tx := global.Global().Database.Begin()
	{
		for _, roleId := range req.RoleId {
			for _, userId := range req.UserId {
				if err := tx.Create(models.NewUserRoleBinding(userId, roleId)).Error; err != nil {
					tx.Rollback()
					return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
				}
			}
		}
	}
	if err := tx.Commit().Error; err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}

	return &pb.BindUserRoleResponse{
		UserId: req.UserId,
		RoleId: req.RoleId,
	}, nil
}

func UnbindUserRole(ctx context.Context, req *pb.UnbindUserRoleRequest) (*pb.UnbindUserRoleResponse, error) {
	if len(req.UserId) == 0 || len(req.RoleId) == 0 {
		return nil, gerr.New(ctx, gerr.InvalidArgument, gerr.ErrorMissingParameter, "userId, roleId")
	}

	// check user bound to role
	userRoleBindings, err := GetUserRoleBindings(ctx, req.UserId, req.RoleId)
	if err != nil {
		return nil, err
	}
	if len(userRoleBindings) != len(req.UserId)*len(req.RoleId) {
		return nil, gerr.New(ctx, gerr.PermissionDenied, gerr.ErrorCannotUnbindGroup)
	}

	if err := global.Global().Database.
		Where(constants.ColumnRoleId+" in (?)", req.RoleId).
		Where(constants.ColumnUserId+" in (?)", req.UserId).
		Delete(models.UserRoleBinding{}).Error; err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}

	return &pb.UnbindUserRoleResponse{
		UserId: req.UserId,
		RoleId: req.RoleId,
	}, nil
}
