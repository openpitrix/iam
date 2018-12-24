// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package service

import (
	"github.com/jmoiron/sqlx"

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
	if _, err := p.Exec(dbInitSql); err != nil {
		logger.Criticalf(nil, "%v", err)
	}

	return p, nil
}

func (p *Database) Close() error {
	return p.DB.Close()
}
