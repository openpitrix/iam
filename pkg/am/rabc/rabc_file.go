// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package rabc

import (
	"os"
	"strings"
	"sync"

	"openpitrix.io/iam/pkg/internal/copyutil"
	"openpitrix.io/iam/pkg/internal/jsonutil"
	"openpitrix.io/iam/pkg/pb/am"
)

type rabcJsonData struct {
	Bindings []pbam.RoleBinding
	Roles    []pbam.Role
}

type rabcFileServer struct {
	jsonpath string
	rabcJsonData
	sync.Mutex
}

func openFileServer(jsonpath string) (Interface, error) {
	var d rabcJsonData
	if err := jsonutil.Load(jsonpath, &d); err != nil {
		return nil, err
	}

	var p = &rabcFileServer{
		jsonpath:     jsonpath,
		rabcJsonData: d,
	}

	p.adjustRabcData()
	return p, nil
}

func (p *rabcFileServer) Close() error {
	p.Lock()
	defer p.Unlock()

	var bakname = func() string {
		if strings.HasSuffix(p.jsonpath, ".json") {
			return p.jsonpath[:len(p.jsonpath)-len(".json")] + ".bak.json"
		} else {
			return p.jsonpath + ".bak.json"
		}
	}()

	os.Rename(p.jsonpath, bakname)
	if err := jsonutil.Save(p.jsonpath, p.rabcJsonData); err != nil {
		os.Rename(bakname, p.jsonpath)
		return err
	}

	p.jsonpath = ""
	p.rabcJsonData = rabcJsonData{}
	return nil
}

func (p *rabcFileServer) CanDo(x pbam.Action) bool {
	p.Lock()
	defer p.Unlock()

	// 1. check by role name
	for _, role := range p.getRoleListByName(x.RoleName...) {
		for _, rule := range role.Rule {
			if p.matchRule(x.Verb, x.Path, rule) {
				return true
			}
		}
	}

	// 2. check by xid list
	for _, role := range p.getRoleListByXid(x.RoleName...) {
		for _, rule := range role.Rule {
			if p.matchRule(x.Verb, x.Path, rule) {
				return true
			}
		}
	}

	// Failed
	return false
}

func (p *rabcFileServer) AllRoles() []pbam.Role {
	p.Lock()
	defer p.Unlock()

	var d rabcJsonData
	copyutil.MustDeepCopy(&d, p.rabcJsonData)
	return d.Roles
}
func (p *rabcFileServer) AllRoleBindings() []pbam.RoleBinding {
	p.Lock()
	defer p.Unlock()

	var d rabcJsonData
	copyutil.MustDeepCopy(&d, p.rabcJsonData)
	return d.Bindings
}

func (p *rabcFileServer) GetRoleByName(name string) (role pbam.Role, ok bool) {
	p.Lock()
	defer p.Unlock()

	panic("TODO")
}
func (p *rabcFileServer) GetRoleByXid(xid []string) pbam.RoleList {
	p.Lock()
	defer p.Unlock()

	panic("TODO")
}

func (p *rabcFileServer) CreateRole(role pbam.Role) error {
	p.Lock()
	defer p.Unlock()

	panic("TODO")
}
func (p *rabcFileServer) ModifyRole(role pbam.Role) error {
	p.Lock()
	defer p.Unlock()

	panic("TODO")
}
func (p *rabcFileServer) DeleteRole(name string) error {
	p.Lock()
	defer p.Unlock()

	panic("TODO")
}

func (p *rabcFileServer) CreateRoleBinding(x []pbam.RoleBinding) error {
	p.Lock()
	defer p.Unlock()

	panic("TODO")
}
func (p *rabcFileServer) DeleteRoleBinding(xid []string) error {
	p.Lock()
	defer p.Unlock()

	panic("TODO")
}
