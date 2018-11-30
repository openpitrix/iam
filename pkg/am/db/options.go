// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db

type Options struct {
	Engine    string // MySQL: InnoDB, MyISAM
	Encoding  string // MySQL: utf8
	ParseTime bool   // MySQL: ?parseTime=true, see https://github.com/go-sql-driver/mysql/issues/9
}

func DefaultOptions(dbtype string) *Options {
	if dbtype == "mysql" {
		return &Options{
			Engine:    "InnoDB",
			Encoding:  "utf8",
			ParseTime: true,
		}
	}
	if dbtype == "sqlite3" {
		return &Options{
			Encoding:  "utf8",
			ParseTime: true,
		}
	}

	return &Options{
		Encoding:  "utf8",
		ParseTime: true,
	}
}
