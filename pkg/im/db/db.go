// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

// OpenPitrix Identity Management service database package.
package db

import (
	"github.com/jmoiron/sqlx"

	"openpitrix.io/iam/pkg/pb/im"
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
		SqlTableSchema_User,
		SqlTableSchema_Group,
	} {
		if _, err := p.Exec(sql); err != nil {
			return err
		}
	}
	return nil
}

func (p *Database) CreateUser(arg *pbim.User) error {
	panic("TODO")
}
func (p *Database) CreateGroup(arg *pbim.Group) error {
	panic("TODO")
}

func (p *Database) GetUser(arg *pbim.Id) (*pbim.User, error) {
	panic("TODO")
}
func (p *Database) GetUserByGroupId(arg *pbim.Id) (*pbim.UserList, error) {
	panic("TODO")
}

func (p *Database) GetGroup(arg *pbim.Id) (*pbim.Group, error) {
	panic("TODO")
}
func (p *Database) GetRootGroup(arg *pbim.Empty) (*pbim.Group, error) {
	panic("TODO")
}
func (p *Database) GetGroupTree(arg *pbim.Id) (*pbim.GroupTree, error) {
	panic("TODO")
}

func (p *Database) ListUesrs(arg *pbim.Range) (*pbim.ListUesrsResponse, error) {
	panic("TODO")
}
func (p *Database) ListGroups(arg *pbim.Range) (*pbim.ListGroupsResponse, error) {
	panic("TODO")
}

func (p *Database) ModifyUser(arg *pbim.ModifyUsersRequest) error {
	panic("TODO")
}
func (p *Database) ModifyGroup(arg *pbim.ModifyGroupsRequest) error {
	panic("TODO")
}

func (p *Database) ComparePassword(arg *pbim.Password) error {
	panic("TODO")
}
func (p *Database) ModifyPassword(arg *pbim.Password) error {
	panic("TODO")
}

func (p *Database) DeleteUsers(arg *pbim.IdList) error {
	panic("TODO")
}
func (p *Database) DeleteGroups(arg *pbim.IdList) error {
	panic("TODO")
}
