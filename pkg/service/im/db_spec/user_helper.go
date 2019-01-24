// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db_spec

import (
	"time"

	idpkg "openpitrix.io/iam/pkg/id"
	"openpitrix.io/iam/pkg/internal/strutil"
)

func (p *User) AdjustForCreate() *User {
	p.UserId = strutil.SimplifyString(p.UserId)
	p.UserName = strutil.SimplifyString(p.UserName)
	p.Email = strutil.SimplifyString(p.Email)
	p.PhoneNumber = strutil.SimplifyString(p.PhoneNumber)
	p.Description = strutil.SimplifyString(p.Description)
	p.Status = strutil.SimplifyString(p.Status)

	if p.UserId == "" {
		p.UserId = idpkg.GenId("uid-")
	}

	now := time.Now()
	p.CreateTime = now
	p.UpdateTime = now
	p.StatusTime = now

	return p
}

func (p *User) IsValidForCreate() error {
	return nil
}

func (p *User) AdjustForUpdate() *User {
	return p
}

func (p *User) IsValidForUpdate() error {
	return nil
}
