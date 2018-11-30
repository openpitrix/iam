// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

// +build ignore

// OpenPitrix Access Management service package.
package am

import (
	"context"

	"openpitrix.io/iam/pkg/pb/am"
)

var (
	_ pbam.AccessManagerServer = (*AccessManager)(nil)
)

type AccessManager struct {
	p int
}

func NewAccessManager() *AccessManager {
	return nil
}

func (p *AccessManager) CreateRole(ctx context.Context, req *pbam.Role) (*pbam.Role, error) {
	panic("TODO")
}

func (p *AccessManager) ModifyRole(ctx context.Context, req *pbam.Role) (*pbam.Role, error) {
	panic("TODO")
}

func (p *AccessManager) DeleteRoleByName(ctx context.Context, req *pbam.String) (*pbam.Bool, error) {
	panic("TODO")
}

func (p *AccessManager) GetRoleByRoleName(ctx context.Context, req *pbam.String) (*pbam.Role, error) {
	panic("TODO")
}

func (p *AccessManager) GetRoleByXidList(ctx context.Context, req *pbam.XidList) (*pbam.RoleList, error) {
	panic("TODO")
}

func (p *AccessManager) ListRoles(ctx context.Context, req *pbam.RoleNameRegexp) (*pbam.RoleList, error) {
	panic("TODO")
}

func (p *AccessManager) CreateRoleBinding(ctx context.Context, req *pbam.RoleBindingList) (*pbam.Bool, error) {
	panic("TODO")
}

func (p *AccessManager) DeleteRoleBinding(ctx context.Context, req *pbam.RoleBindingList) (*pbam.Bool, error) {
	panic("TODO")
}

func (p *AccessManager) DeleteAllRoleBindings(ctx context.Context, req *pbam.XidList) (*pbam.Bool, error) {
	panic("TODO")
}

func (p *AccessManager) CanDo(ctx context.Context, req *pbam.Action) (*pbam.Bool, error) {
	panic("TODO")
}
