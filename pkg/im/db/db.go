// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

// OpenPitrix Identity Management service database package.
package db

import (
	"database/sql"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/gorp.v2"

	"openpitrix.io/iam/pkg/pb/im"
)

// https://github.com/openpitrix/openpitrix/blob/b3542cadf7c893e3098eb1d5d876d7c16a4a8531/pkg/db/app/app.go

type Database struct {
	db    *sql.DB
	dbMap *gorp.DbMap
}

func Open(dbtype, dbpath string) (*Database, error) {
	// https://github.com/go-sql-driver/mysql/issues/9
	if strings.EqualFold(dbtype, "mysql") {
		dbpath += "?parseTime=true"
	}

	db, err := sql.Open(dbtype, dbpath)
	if err != nil {
		return nil, err
	}

	var dialect = func() gorp.Dialect {
		switch {
		case strings.EqualFold(dbtype, "mysql"):
			return gorp.MySQLDialect{
				Encoding: "mysql",
				Engine:   "utf8",
			}
		case strings.EqualFold(dbtype, "postgres"):
			return gorp.PostgresDialect{}
		case strings.EqualFold(dbtype, "sqlite3"):
			return gorp.SqliteDialect{}
		default:
			return gorp.SqliteDialect{}
		}
	}()

	dbMap := &gorp.DbMap{Db: db, Dialect: dialect}
	for _, v := range _TableSchemaMap {
		if _, err := db.Exec(v.Schema); err != nil {
			return nil, err
		}
		dbMap.AddTableWithNameAndSchema(
			v.Value, v.Schema, v.Name,
		)
	}

	p := &Database{db: db, dbMap: dbMap}
	return p, nil
}

func (p *Database) Close() error {
	return p.db.Close()
}

func (p *Database) createTablesIfNotExists() error {
	for _, sql := range []string{
		SqlTableSchema_User,
		SqlTableSchema_Group,
	} {
		if _, err := p.db.Exec(sql); err != nil {
			return err
		}
	}
	return nil
}

func (p *Database) CreateUser(arg *pbim.User) error {
	//v := NewUserFrom(arg)

	/*
			p1 := newPost("Go 1.1 released!", "Lorem ipsum lorem ipsum")
		    p2 := newPost("Go 1.2 released!", "Lorem ipsum lorem ipsum")

		    // insert rows - auto increment PKs will be set properly after the insert
		    err = dbmap.Insert(&p1, &p2)
		    checkErr(err, "Insert failed")
	*/

	/*
		_, err := p.Exec(`INSERT INTO user (name, rule) VALUES ($1, $2);`, v.Name, v.Rule)
		if err != nil {
			return err
		}
		return nil

				v := NewRoleFrom(role)
			_, err := p.Exec(`INSERT INTO role (name, rule) VALUES ($1, $2);`, v.Name, v.Rule)
			if err != nil {
				return err
			}
			return nil
	*/
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
func (p *Database) GetGroupPath(arg *pbim.Id) (*pbim.GroupPath, error) {
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
