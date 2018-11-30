// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

// OpenPitrix Access Management service database package.
package db

import (
	"context"
	"database/sql"
	"sync/atomic"

	"github.com/pkg/errors"
	"gopkg.in/gorp.v2"

	"openpitrix.io/iam/pkg/am/db_spec"
)

type Database struct {
	db               *sql.DB
	dbMap            *gorp.DbMap
	createTablesDone uint32
}

//
// Open Access Manager Database.
//
// MySQL
//	dbtype: mysql
//	dbpath: user:password@tcp(localhost:3306)/dbname
//
// Sqlite3
//	dbtype: sqlite3
//	dbpath: /tmp/db.bin
//
func OpenDatabase(dbtype, dbpath string, opt *Options) (p *Database, err error) {
	if opt == nil {
		opt = DefaultOptions(dbtype)
	}

	if dbtype == "mysql" && opt.ParseTime {
		dbpath += "?parseTime=true"
	}

	db, err := sql.Open(dbtype, dbpath)
	if err != nil {
		err = errors.WithStack(err)
		return nil, err
	}

	dialect := gorp.MySQLDialect{
		Encoding: opt.Encoding,
		Engine:   opt.Engine,
	}

	dbMap := &gorp.DbMap{Db: db, Dialect: dialect}
	for _, v := range db_spec.AllTableSpecList {
		dbMap.AddTableWithNameAndSchema(v,
			v.GetTableSchema(dbtype),
			v.GetTableName(),
		)
	}

	p = &Database{
		db:    db,
		dbMap: dbMap,
	}

	p.initTables()
	return
}

func (p *Database) Close() error {
	err := p.db.Close()
	err = errors.WithStack(err)
	return err
}

func (p *Database) GetRole(ctx context.Context, id string) (*Role, error) {
	p.initTables()
	if v, err := p.dbMap.Get(Role{}, id); err == nil && v != nil {
		return v.(*Role), nil
	} else {
		err = errors.WithStack(err)
		return nil, err
	}
}

func (p *Database) initTables() {
	if atomic.LoadUint32(&p.createTablesDone) == 1 {
		return
	}
	if err := p.dbMap.CreateTablesIfNotExists(); err != nil {
		//logger.Warningf("CreateTablesIfNotExists: %+v", err)
		return
	}
	atomic.StoreUint32(&p.createTablesDone, 1)
}
