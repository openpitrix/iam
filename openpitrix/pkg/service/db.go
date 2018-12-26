// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package service

import (
	"database/sql"
	"fmt"

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
	for _, v := range InitSqlList {
		if _, err := p.Exec(v.Sql); err != nil {
			logger.Warnf(nil, "%v", err)
		}
	}

	if err := p.DB.Ping(); err != nil {
		logger.Warnf(nil, "%#v", err)
	}

	stats := p.DB.Stats()
	logger.Infof(nil, "DB stats: %s", fmt.Sprintf("%#v", stats))

	return p, nil
}

func (p *Database) Close() error {
	return p.DB.Close()
}
