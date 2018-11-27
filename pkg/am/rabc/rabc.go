// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package rabc

import (
	"openpitrix.io/iam/pkg/pb/am"
)

type Interface interface {
	All() []pbam.Role
	Get(name string) (role pbam.Role, ok bool)
	CanDo(x pbam.Action) bool

	Create(name string, rule []pbam.Rule) error
	Modify(name string, rule []pbam.Rule) error
	Delete(name string) error

	Close() error
}

type DBOptions struct {
	DBType     string // mysql/sqlite3
	DBEngine   string // InnoDB/...
	DBEncoding string // utf8/...
}

func OpenDatabase(dbpath string, opt *DBOptions) (Interface, error) {
	panic("TODO")
}

func OpenFile(jsonFile string) (Interface, error) {
	panic("TODO")
}
