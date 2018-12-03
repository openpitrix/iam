// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package am

import (
	"context"

	"openpitrix.io/iam/pkg/pb/am"
)

func (p *Server) CreateRole(ctx context.Context, req *pbam.Role) (*pbam.Role, error) {
	panic("TODO")
}

func (p *Server) ModifyRole(ctx context.Context, req *pbam.Role) (*pbam.Role, error) {
	panic("TODO")
}

func (p *Server) DeleteRoleByName(ctx context.Context, req *pbam.String) (*pbam.Bool, error) {
	panic("TODO")
}

func (p *Server) GetRoleByName(ctx context.Context, req *pbam.String) (*pbam.Role, error) {
	panic("TODO")
}

func (p *Server) GetRoleByXidList(ctx context.Context, req *pbam.XidList) (*pbam.RoleList, error) {
	panic("TODO")
}

func (p *Server) ListRoles(ctx context.Context, req *pbam.RoleNameFilter) (*pbam.RoleList, error) {
	panic("TODO")
}

func (p *Server) CreateRoleBinding(ctx context.Context, req *pbam.RoleBindingList) (*pbam.Bool, error) {
	panic("TODO")
}

func (p *Server) DeleteRoleBinding(ctx context.Context, req *pbam.XidList) (*pbam.Bool, error) {
	panic("TODO")
}

func (p *Server) GetRoleBindingByRoleName(ctx context.Context, req *pbam.String) (*pbam.RoleBindingList, error) {
	panic("TODO")
}
func (p *Server) GetRoleBindingByXidList(ctx context.Context, req *pbam.XidList) (*pbam.RoleBindingList, error) {
	panic("TODO")
}
func (p *Server) ListRoleBindings(ctx context.Context, req *pbam.RoleNameFilter) (*pbam.RoleBindingList, error) {
	panic("TODO")
}

func (p *Server) CanDo(ctx context.Context, req *pbam.Action) (*pbam.Bool, error) {
	panic("TODO")
}
