// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package am

import (
	"context"
	"strings"

	"openpitrix.io/iam/pkg/constants"
	"openpitrix.io/iam/pkg/gerr"
	"openpitrix.io/iam/pkg/models"
	"openpitrix.io/iam/pkg/sender"
	"openpitrix.io/iam/pkg/service/am/resource"
	"openpitrix.io/iam/pkg/util/ctxutil"
	"openpitrix.io/logger"
)

func CheckRolesPermission(ctx context.Context, roleIds []string, action string) ([]*models.Role, error) {
	s := ctxutil.GetSender(ctx)

	logger.Debugf(ctx, "Sender user [%s].", s.UserId)

	isBound := false
	userRoleBindings, err := resource.GetUserRoleBindings(ctx, []string{s.UserId}, roleIds)
	if err != nil {
		return nil, err
	}
	if len(userRoleBindings) > 0 {
		isBound = true
	}

	roles, err := resource.GetRoles(ctx, roleIds)
	if err != nil {
		return nil, err
	}

	if len(roleIds) != len(roles) {
		return nil, gerr.New(ctx, gerr.PermissionDenied, gerr.ErrorRoleNotFound, strings.Join(roleIds, ","))
	}

	for _, role := range roles {
		if role.Status != constants.StatusActive {
			return nil, gerr.New(ctx, gerr.PermissionDenied, gerr.ErrorRoleNotInStatus, role.RoleId, constants.StatusActive)
		}

		if role.Controller == constants.ControllerPitrix {
			switch action {
			case constants.ActionModify, constants.ActionDelete:
				return nil, gerr.New(ctx, gerr.PermissionDenied, gerr.ErrorInternalRoleIllegalAction)
			}
		} else if role.Controller == constants.ControllerSelf {
			if isBound && action == constants.ActionDescribe {
				continue
			} else {
				ownerPath := sender.OwnerPath(role.OwnerPath)
				if !ownerPath.CheckPermission(s) {
					return nil, gerr.New(ctx, gerr.PermissionDenied, gerr.ErrorRoleAccessDenied, role.RoleId)
				}
			}

		}
	}
	return roles, nil
}
