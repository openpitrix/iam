// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package resource

import (
	"context"
	"time"

	"strings"

	"openpitrix.io/iam/pkg/constants"
	"openpitrix.io/iam/pkg/db"
	"openpitrix.io/iam/pkg/gerr"
	"openpitrix.io/iam/pkg/global"
	"openpitrix.io/iam/pkg/models"
	"openpitrix.io/iam/pkg/pb"
	"openpitrix.io/iam/pkg/util/ctxutil"
	"openpitrix.io/iam/pkg/util/strutil"
	"openpitrix.io/logger"
)

func GetRole(ctx context.Context, roleId string) (*models.Role, error) {
	var role = &models.Role{RoleId: roleId}
	if err := global.Global().Database.Table(constants.TableRole).
		Take(role).Error; err != nil {
		return nil, gerr.New(ctx, gerr.NotFound, gerr.ErrorRoleNotFound, roleId)
	}

	return role, nil
}

func GetRoles(ctx context.Context, roleIds []string) ([]*models.Role, error) {
	var roles []*models.Role
	if err := global.Global().Database.Table(constants.TableRole).
		Where(constants.ColumnRoleId+" in (?)", roleIds).
		Find(&roles).Error; err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}

	return roles, nil
}

func GetSenderPortal(ctx context.Context) (string, error) {
	s := ctxutil.GetSender(ctx)

	var senderRoleId string
	if s.UserId == constants.UserSystem {
		senderRoleId = constants.RoleGlobalAdmin
	} else {
		roleIds, err := GetRoleIdsByUserIds(ctx, []string{s.UserId})
		if err != nil {
			return "", err
		}
		if len(roleIds) == 0 {
			return "", gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
		}
		senderRoleId = roleIds[0]
	}

	senderRole, err := GetRole(ctx, senderRoleId)
	if err != nil {
		return "", err
	}
	return senderRole.Portal, nil
}

func CreateRole(ctx context.Context, req *pb.CreateRoleRequest) (*pb.CreateRoleResponse, error) {
	if !strutil.Contains(constants.PortalSet, req.Portal) {
		return nil, gerr.New(ctx, gerr.InvalidArgument, gerr.ErrorPortalNotFound, req.Portal)
	}

	role := models.NewRole(req.RoleName, req.Description, req.Portal, req.Owner, req.OwnerPath)

	moduleIds, err := GetModuleIds(ctx)
	if err != nil {
		return nil, err
	}

	tx := global.Global().Database.Begin()
	{
		// create new record
		if err := tx.Create(role).Error; err != nil {
			tx.Rollback()
			logger.Errorf(ctx, "Insert role failed: %v", err)
			return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
		}

		// bind modules
		for _, moduleId := range moduleIds {
			isCheckAll := false

			// m0 enabled for all roles
			if moduleId == constants.ModuleIdM0 {
				isCheckAll = true
			}

			roleModuleBinding := models.NewRoleModuleBinding(role.RoleId, moduleId, constants.DataLevelSelf, isCheckAll)
			if err := tx.Create(roleModuleBinding).Error; err != nil {
				tx.Rollback()
				logger.Errorf(ctx, "Insert role module binding failed: %v", err)
				return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
			}
		}
	}
	if err := tx.Commit().Error; err != nil {
		logger.Errorf(ctx, "Create role failed: %+v", err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}

	return &pb.CreateRoleResponse{
		RoleId: role.RoleId,
	}, nil
}

func DeleteRoles(ctx context.Context, req *pb.DeleteRolesRequest) (*pb.DeleteRolesResponse, error) {
	roleIds := req.GetRoleId()

	userIds, err := GetUserIdsByRoleIds(ctx, roleIds)
	if err != nil {
		return nil, err
	}
	if len(userIds) != 0 {
		return nil, gerr.New(ctx, gerr.PermissionDenied, gerr.ErrorStillUserBindingRole)
	}

	now := time.Now()
	attributes := map[string]interface{}{
		constants.ColumnStatusTime: now,
		constants.ColumnUpdateTime: now,
		constants.ColumnStatus:     constants.StatusDeleted,
	}
	if err := global.Global().Database.Table(constants.TableRole).
		Where(constants.ColumnRoleId+" in (?)", roleIds).
		Updates(attributes).Error; err != nil {
		logger.Errorf(ctx, "Update role status failed: %+v", err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}

	return &pb.DeleteRolesResponse{
		RoleId: roleIds,
	}, nil
}

func ModifyRole(ctx context.Context, req *pb.ModifyRoleRequest) (*pb.ModifyRoleResponse, error) {
	roleId := req.GetRoleId()
	_, err := GetRole(ctx, roleId)
	if err != nil {
		return nil, err
	}

	attributes := make(map[string]interface{})
	if req.RoleName != "" {
		attributes[constants.ColumnRoleName] = req.RoleName
	}
	if req.Description != "" {
		attributes[constants.ColumnDescription] = req.Description
	}
	attributes[constants.ColumnUpdateTime] = time.Now()

	if err := global.Global().Database.Table(constants.TableRole).
		Where(constants.ColumnRoleId+" = ?", roleId).
		Updates(attributes).Error; err != nil {
		logger.Errorf(ctx, "Update role [%s] failed: %+v", roleId, err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}

	return &pb.ModifyRoleResponse{
		RoleId: roleId,
	}, nil
}

func DescribeRoles(ctx context.Context, req *pb.DescribeRolesRequest) (*pb.DescribeRolesResponse, error) {
	senderPortal, err := GetSenderPortal(ctx)
	if err != nil {
		return nil, err
	}

	if !strutil.Contains(constants.PortalSet, senderPortal) {
		return nil, gerr.New(ctx, gerr.InvalidArgument, gerr.ErrorPortalNotFound, senderPortal)
	}

	req.RoleId = strutil.SimplifyStringList(req.RoleId)
	req.RoleName = strutil.SimplifyStringList(req.RoleName)
	req.Portal = strutil.SimplifyStringList(req.Portal)
	req.Status = strutil.SimplifyStringList(req.Status)
	req.UserId = strutil.SimplifyStringList(req.UserId)

	var addedRoleIds []string
	if senderPortal == constants.PortalGlobalAdmin {
		addedRoleIds = []string{constants.RoleGlobalAdmin, constants.RoleIsv, constants.RoleUser}
	} else if senderPortal == constants.PortalIsv {
		addedRoleIds = []string{constants.RoleIsv, constants.RoleDeveloper}
	} else if senderPortal == constants.PortalUser {
		addedRoleIds = []string{constants.RoleUser}
	}

	req.Portal = []string{}

	var pbRoles []*pb.Role

	if len(req.UserId) > 0 {
		roleIds, err := GetRoleIdsByUserIds(ctx, req.UserId)
		if err != nil {
			logger.Errorf(ctx, "Get role id by user id failed: %+v", err)
			return nil, err
		}

		if len(req.RoleId) == 0 {
			req.RoleId = roleIds
		} else {
			var inRoleIds []string
			for _, roleId := range req.RoleId {
				if strutil.Contains(roleIds, roleId) {
					inRoleIds = append(inRoleIds, roleId)
				}
			}
			req.RoleId = inRoleIds
		}

		if len(req.RoleId) == 0 {
			return &pb.DescribeRolesResponse{
				RoleSet: pbRoles,
				Total:   uint32(0),
			}, nil
		}
	}

	limit := db.GetLimitFromRequest(req)
	offset := db.GetOffsetFromRequest(req)

	var roles []*models.Role
	var count int

	andConditions := []string{
		constants.ColumnPortal + " = '" + senderPortal + "'",
	}

	orConditions := []string{
		constants.ColumnRoleId + " in ('" + strings.Join(addedRoleIds, "','") + "')",
	}

	if err := db.GetChain(global.Global().Database.Table(constants.TableRole)).
		AddQueryOrderDir(req, constants.ColumnCreateTime).
		BuildOwnerPathFilter(ctx, req, andConditions, orConditions).
		BuildFilterConditions(req, constants.TableRole).
		Offset(offset).
		Limit(limit).
		Find(&roles).Error; err != nil {
		logger.Errorf(ctx, "Describe role failed: %+v", err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}

	if err := db.GetChain(global.Global().Database.Table(constants.TableRole)).
		BuildOwnerPathFilter(ctx, req, andConditions, orConditions).
		BuildFilterConditions(req, constants.TableRole).
		Count(&count).Error; err != nil {
		logger.Errorf(ctx, "Describe role count failed: %+v", err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}

	for _, role := range roles {
		pbRoles = append(pbRoles, role.ToPB())
	}

	return &pb.DescribeRolesResponse{
		RoleSet: pbRoles,
		Total:   uint32(count),
	}, nil
}

func DescribeRolesWithUser(ctx context.Context, req *pb.DescribeRolesRequest) (*pb.DescribeRolesWithUserResponse, error) {
	response, err := DescribeRoles(ctx, req)
	if err != nil {
		return nil, err
	}

	var pbRoleWithUsers []*pb.RoleWithUser
	for _, pbRole := range response.RoleSet {
		userIds, err := GetUserIdsByRoleIds(ctx, []string{pbRole.RoleId})
		if err != nil {
			return nil, err
		}
		pbRoleWithUser := &pb.RoleWithUser{
			Role:      pbRole,
			UserIdSet: userIds,
		}
		pbRoleWithUsers = append(pbRoleWithUsers, pbRoleWithUser)
	}

	return &pb.DescribeRolesWithUserResponse{
		Total:   response.Total,
		RoleSet: pbRoleWithUsers,
	}, nil
}

func GetRoleWithUser(ctx context.Context, req *pb.GetRoleRequest) (*pb.GetRoleWithUserResponse, error) {
	role, err := GetRole(ctx, req.RoleId)
	if err != nil {
		return nil, err
	}
	userIds, err := GetUserIdsByRoleIds(ctx, []string{req.RoleId})
	if err != nil {
		return nil, err
	}

	roleWithUser := &models.RoleWithUser{
		Role:    role,
		UserIds: userIds,
	}

	return &pb.GetRoleWithUserResponse{
		Role: roleWithUser.ToPB(),
	}, nil
}
