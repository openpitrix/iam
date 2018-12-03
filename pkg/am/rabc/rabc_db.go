// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package rabc

import (
	"openpitrix.io/iam/pkg/am/db"
	"openpitrix.io/iam/pkg/pb/am"
)

var (
	_ Interface = (*rabcDbServer)(nil)
)

type rabcDbServer struct {
	db *db.Database
}

func openDatabase(dbtype, dbpath string) (*rabcDbServer, error) {
	db, err := db.Open(dbtype, dbpath)
	if err != nil {
		return nil, err
	}
	p := &rabcDbServer{
		db: db,
	}
	return p, nil
}

func (p *rabcDbServer) Close() error {
	return p.db.Close()
}

func (p *rabcDbServer) CreateRole(role *pbam.Role) error {
	return p.db.CreateRole(role)
}
func (p *rabcDbServer) ModifyRole(role *pbam.Role) error {
	return p.db.ModifyRole(role)
}
func (p *rabcDbServer) DeleteRoleByName(name string) error {
	return p.db.DeleteRoleByRoleName(name)
}

func (p *rabcDbServer) GetRoleByName(name string) (*pbam.Role, error) {
	return p.db.GetRoleByName(name)
}
func (p *rabcDbServer) GetRoleByXidList(xid ...string) (*pbam.RoleList, error) {
	return p.db.GetRoleByXidList(xid...)
}
func (p *rabcDbServer) ListRoles(filter *pbam.RoleNameFilter) (*pbam.RoleList, error) {
	return p.db.ListRoles(filter)
}

func (p *rabcDbServer) CreateRoleBinding(x *pbam.RoleBindingList) error {
	return p.db.CreateRoleBinding(x)
}
func (p *rabcDbServer) DeleteRoleBinding(xid ...string) error {
	return p.db.DeleteRoleBinding(xid...)
}

func (p *rabcDbServer) GetRoleBindingByRoleName(name string) (*pbam.RoleBindingList, error) {
	return p.db.GetRoleBindingByRoleName(name)
}
func (p *rabcDbServer) GetRoleBindingByXidList(xid ...string) (*pbam.RoleBindingList, error) {
	return p.db.GetRoleBindingByXidList(xid...)
}
func (p *rabcDbServer) ListRoleBindings(filter *pbam.RoleNameFilter) (*pbam.RoleBindingList, error) {
	return p.db.ListRoleBindings(filter)
}

func (p *rabcDbServer) CanDo(x *pbam.Action) bool {
	for _, name := range x.GetRoleName() {
		if role, err := p.db.GetRoleByName(name); err == nil {
			if canDoAction(x, role.GetRule()) {
				return true
			}
		}
	}

	if roleList, err := p.db.GetRoleByXidList(x.GetXid()...); err == nil {
		for _, role := range roleList.GetValue() {
			if canDoAction(x, role.GetRule()) {
				return true
			}
		}
	}

	return false
}
