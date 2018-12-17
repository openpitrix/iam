// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

// OpenPitrix Access Management service database package.
package db

import (
	"regexp"

	"github.com/bmatcuk/doublestar"
	"github.com/jmoiron/sqlx"

	"openpitrix.io/iam/pkg/pb/am"
	"openpitrix.io/logger"
)

type Database struct {
	*sqlx.DB
}

func Open(dbtype, dbpath string) (*Database, error) {
	db, err := sqlx.Open(dbtype, dbpath)
	if err != nil {
		return nil, err
	}

	p := &Database{DB: db}
	if err := p.createTablesIfNotExists(); err != nil {
		logger.Criticalf(nil, "%v", err)
	}

	return p, nil
}

func (p *Database) Close() error {
	return p.DB.Close()
}

func (p *Database) createTablesIfNotExists() error {
	for _, sql := range []string{
		SqlTableSchema_Role,
		SqlTableSchema_RoleBinding,
	} {
		if _, err := p.Exec(sql); err != nil {
			return err
		}
	}
	return nil
}

func (p *Database) CreateRole(role *pbam.Role) error {
	v := NewRoleFrom(role)
	_, err := p.Exec(`INSERT INTO role (name, rule) VALUES ($1, $2);`, v.Name, v.Rule)
	if err != nil {
		return err
	}
	return nil
}

func (p *Database) ModifyRole(role *pbam.Role) error {
	v := NewRoleFrom(role)
	_, err := p.Exec(`UPDATE role SET rule = $2 WHERE name = $1;`, v.Name, v.Rule)
	if err != nil {
		return err
	}
	return nil
}

func (p *Database) DeleteRoleByRoleName(name string) error {
	_, err := p.Exec(`DELETE FROM role WHERE name = $1;`, name)
	if err != nil {
		return err
	}
	return nil
}

func (p *Database) GetRoleByName(name string) (*pbam.Role, error) {
	var v Role
	err := p.Get(&v, `SELECT * FROM role WHERE name=$1;`, name)
	if err != nil {
		return nil, err
	}
	return v.ToPbRole(), nil
}

func (p *Database) ListRoles(filter *pbam.NameFilter) (*pbam.RoleList, error) {
	var (
		roles  = []Role{}
		result = &pbam.RoleList{}
	)
	err := p.Select(&roles, "SELECT * FROM role;")
	if err != nil {
		return nil, err
	}

	// if glob pattern
	if sPathPattern := filter.GetGlobPattern(); sPathPattern != "" {
		for _, role := range roles {
			if ok, _ := doublestar.Match(sPathPattern, role.Name); ok {
				result.Value = append(result.Value, role.ToPbRole())
			}
		}

		// fallthrough
	}

	// if regexp pattern
	if nameRegexp := filter.GetRegexpPattern(); nameRegexp != "" {
		re, err := regexp.Compile(nameRegexp)
		if err != nil {
			return nil, err
		}

		for _, role := range roles {
			if re.MatchString(role.Name) {
				result.Value = append(result.Value, role.ToPbRole())
			}
		}

		// fallthrough
	}

	// else: all
	if filter.GetGlobPattern() == "" && filter.GetRegexpPattern() == "" {
		for _, role := range roles {
			result.Value = append(result.Value, role.ToPbRole())
		}

		// fallthrough
	}

	return result, nil
}

func (p *Database) GetRoleByXidList(xid ...string) (*pbam.RoleList, error) {
	var roles []Role
	err := p.Select(&roles, `
		SELECT * FROM role,role_binding WHERE
			role.name=role_binding.role_name AND
			xid IN(?)
		;`, xid,
	)
	if err != nil {
		return nil, err
	}

	result := &pbam.RoleList{
		Value: make([]*pbam.Role, len(roles)),
	}
	for i, v := range roles {
		result.Value[i] = v.ToPbRole()
	}

	return result, nil
}

func (p *Database) CreateRoleBinding(bindings *pbam.RoleXidBindingList) error {
	/*
		TODO:
		tx, err := db.Begin()
		// check role_name exists
		err = tx.Exec(...)
		err = tx.Commit()
	*/
	for _, v := range bindings.GetValue() {
		_, err := p.Exec(
			`INSERT INTO role_binding (role_name, xid) VALUES ($1, $2);`,
			v.RoleName, v.Xid,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Database) DeleteRoleBinding(xid ...string) error {
	_, err := p.Exec(`DELETE FROM role_binding WHERE xid IN(?);`, xid)
	if err != nil {
		return err
	}
	return nil
}

func (p *Database) GetRoleBindingByRoleName(name string) (*pbam.RoleXidBindingList, error) {
	var bindings []RoleBinding
	err := p.Select(&bindings, `SELECT * FROM role_binding WHERE role_name=$1;`, name)
	if err != nil {
		return nil, err
	}

	result := &pbam.RoleXidBindingList{
		Value: make([]*pbam.RoleXidBinding, len(bindings)),
	}
	for i, v := range bindings {
		result.Value[i] = v.ToPbRoleBinding()
	}

	return result, nil
}

func (p *Database) GetRoleBindingByXidList(xid ...string) (*pbam.RoleXidBindingList, error) {
	var bindings []RoleBinding
	err := p.Select(&bindings, `SELECT * FROM role_binding WHERE xid IN(?);`, xid)
	if err != nil {
		return nil, err
	}

	result := &pbam.RoleXidBindingList{
		Value: make([]*pbam.RoleXidBinding, len(bindings)),
	}
	for i, v := range bindings {
		result.Value[i] = v.ToPbRoleBinding()
	}

	return result, nil
}

func (p *Database) ListRoleBindings(filter *pbam.NameFilter) (*pbam.RoleXidBindingList, error) {
	var (
		bindings = []RoleBinding{}
		result   = &pbam.RoleXidBindingList{}
	)

	err := p.Select(&bindings, `SELECT * FROM role_binding;`)
	if err != nil {
		return nil, err
	}

	// if glob pattern
	if sPathPattern := filter.GetGlobPattern(); sPathPattern != "" {
		for _, v := range bindings {
			if ok, _ := doublestar.Match(sPathPattern, v.RoleName); ok {
				result.Value = append(result.Value, v.ToPbRoleBinding())
			}
		}

		// fallthrough
	}

	// if regexp pattern
	if nameRegexp := filter.GetRegexpPattern(); nameRegexp != "" {
		re, err := regexp.Compile(nameRegexp)
		if err != nil {
			return nil, err
		}

		for _, v := range bindings {
			if re.MatchString(v.RoleName) {
				result.Value = append(result.Value, v.ToPbRoleBinding())
			}
		}

		// fallthrough
	}

	// else: all
	if filter.GetGlobPattern() == "" && filter.GetRegexpPattern() == "" {
		for _, v := range bindings {
			result.Value = append(result.Value, v.ToPbRoleBinding())
		}

		// fallthrough
	}

	return result, nil
}
