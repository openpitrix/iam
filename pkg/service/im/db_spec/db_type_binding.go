// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db_spec

type DBUsertGroupBinding struct {
	Id      string `db:"id"`
	UserId  string `db:"user_id"`
	GroupId string `db:"group_id"`
}
