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
	*db.Database
}

func openDatabase(dbtype, dbpath string) (*rabcDbServer, error) {
	db, err := db.Open(dbtype, dbpath)
	if err != nil {
		return nil, err
	}
	p := &rabcDbServer{
		Database: db,
	}
	return p, nil
}

func (p *rabcDbServer) Close() error {
	return p.Database.Close()
}

func (p *rabcDbServer) CreateRole(role pbam.Role) error {
	panic("TODO")
}
func (p *rabcDbServer) ModifyRole(role pbam.Role) error {
	panic("TODO")
}
func (p *rabcDbServer) DeleteRoleByName(name string) error {
	panic("TODO")
}

func (p *rabcDbServer) GetRoleByName(name string) (*pbam.Role, error) {
	panic("TODO")
}
func (p *rabcDbServer) GetRoleByXidList(xid ...string) (*pbam.RoleList, error) {
	panic("TODO")
}
func (p *rabcDbServer) ListRoles(filter *pbam.RoleNameFilter) (*pbam.RoleList, error) {
	panic("TODO")
}

func (p *rabcDbServer) CreateRoleBinding(x []pbam.RoleBinding) error {
	panic("TODO")
}
func (p *rabcDbServer) DeleteRoleBinding(xid []string) error {
	panic("TODO")
}

func (p *rabcDbServer) GetRoleBindingByRoleName(name string) (*pbam.RoleBindingList, error) {
	panic("TODO")
}
func (p *rabcDbServer) GetRoleBindingByXidList(xid ...string) (*pbam.RoleBindingList, error) {
	panic("TODO")
}
func (p *rabcDbServer) ListRoleBindings(filter *pbam.RoleNameFilter) (*pbam.RoleBindingList, error) {
	panic("TODO")
}

func (p *rabcDbServer) CanDo(x pbam.Action) (bool, error) {
	panic("TODO")
}
