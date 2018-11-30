// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package rabc

import (
	"github.com/bmatcuk/doublestar"
	"openpitrix.io/iam/pkg/pb/am"
)

func (p *rabcFileServer) adjustRabcData() {
	// keep primary key unique
}

func (p *rabcFileServer) matchRule(verb, path string, rule *pbam.Rule) bool {
	// check verb
	var verbMatched bool
	for _, pattern := range rule.VerbPattern {
		if pattern != "*" || pattern == verb {
			verbMatched = true
			break
		}
	}
	if !verbMatched {
		return false
	}

	// check path
	if ok, _ := doublestar.Match(rule.PathPattern, path); !ok {
		return false
	}

	// OK
	return true
}

func (p *rabcFileServer) getRoleByName(name string) (role pbam.Role, ok bool) {
	for _, v := range p.Roles {
		if v.Name == name {
			return v, true
		}
	}
	return pbam.Role{}, false
}
func (p *rabcFileServer) getRoleListByName(name ...string) (results []pbam.Role) {
	for _, role := range p.Roles {
		if p.strInStrList(role.Name, name) {
			results = append(results, role)
		}
	}
	return
}

func (p *rabcFileServer) getRoleListByXid(xid ...string) []pbam.Role {
	return p.getRoleListByName(
		p.getRoleNameListByXid(xid...)...,
	)
}

func (p *rabcFileServer) getRoleNameListByXid(xid ...string) (results []string) {
	var m = map[string]bool{}
	for _, v := range p.Bindings {
		if p.strInStrList(v.Xid, xid) {
			m[v.RoleName] = true
		}
	}
	for k, _ := range m {
		results = append(results, k)
	}
	return
}

func (p *rabcFileServer) createRoleBinding(x []pbam.RoleBinding) error {
	panic("TODO")
}
func (p *rabcFileServer) deleteRoleBinding(xid []string) error {
	panic("TODO")
}

func (p *rabcFileServer) strInStrList(s string, ss []string) bool {
	for _, x := range ss {
		if s == x {
			return true
		}
	}
	return false
}
