// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db

type UserGroupBinding struct {
	Id      string `gorm:"primary_key"`
	GroupId string
	UserId  string
}
