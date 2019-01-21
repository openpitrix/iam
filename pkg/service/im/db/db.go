// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"openpitrix.io/iam/pkg/service/im/config"
	"openpitrix.io/iam/pkg/service/im/db_spec"
	"openpitrix.io/logger"
)

type Database struct {
	cfg *config.Config
	*gorm.DB
}

type Options struct {
	SqlInitDB    []string
	SqlInitTable []string
	SqlInitData  []string
}

func OpenDatabase(cfg *config.Config, opt *Options) (*Database, error) {
	cfg = cfg.Clone()

	logger.Infof(nil, "DB config: begin")
	logger.Infof(nil, "\tType: %s", cfg.DB.Type)
	logger.Infof(nil, "\tHost: %s", cfg.DB.Host)
	logger.Infof(nil, "\tPort: %d", cfg.DB.Port)
	logger.Infof(nil, "\tUser: %s", cfg.DB.User)
	logger.Infof(nil, "\tDatabase: %s", cfg.DB.Database)
	logger.Infof(nil, "DB config: end")

	orm, err := gorm.Open(cfg.DB.Type, cfg.DB.GetUrl())
	if err != nil {
		return nil, err
	}

	p := &Database{
		cfg: cfg,
		DB:  orm,
	}
	for _, v := range db_spec.TableMap {
		if !isValidDatabaseTableName(v.Name) {
			logger.Warnf(nil, "invalid table name %s", v.Name)
		}
		if err := p.DB.Exec(v.Sql).Error; err != nil {
			logger.Warnf(nil, "%s: %v", v.Name, err)
		}
	}

	return p, nil
}

func (p *Database) Close() error {
	return p.DB.Close()
}
