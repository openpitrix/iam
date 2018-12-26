// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package service

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"

	"openpitrix.io/logger"
)

type Database struct {
	*sql.DB
}

func Open(dbtype, dbpath string) (*Database, error) {
	db, err := sql.Open(dbtype, dbpath)
	if err != nil {
		return nil, err
	}

	p := &Database{DB: db}
	if _, err := p.Exec(dbInitSql); err != nil {
		logger.Criticalf(nil, "%v", err)
	}

	return p, nil
}

func (p *Database) Close() error {
	return p.DB.Close()
}
