// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db_spec

import (
	"fmt"
	"strings"
	"time"

	idpkg "openpitrix.io/iam/pkg/id"
	"openpitrix.io/iam/pkg/internal/strutil"
	"openpitrix.io/iam/pkg/validator"
)

func (p *Role) AdjustForCreate() *Role {
	p.RoleId = strutil.SimplifyString(p.RoleId)
	p.RoleName = strutil.SimplifyString(p.RoleName)
	p.Description = strutil.SimplifyString(p.Description)
	p.Portal = strutil.SimplifyString(p.Portal)
	p.Owner = strutil.SimplifyString(p.Owner)
	p.OwnerPath = strutil.SimplifyString(p.OwnerPath)
	p.Status = strutil.SimplifyString(p.Status)

	if p.RoleId == "" {
		p.RoleId = idpkg.GenId("role-")
	}

	now := time.Now()
	p.CreateTime = now
	p.UpdateTime = now
	p.StatusTime = now

	return p
}

func (p *Role) IsValidForCreate() error {
	if !validator.IsValidId(p.RoleId) {
		return fmt.Errorf("Role.IsValidForCreate: invalid RoleId %q", p.RoleId)
	}
	if !validator.IsValidName(p.RoleName) {
		return fmt.Errorf("Role.IsValidForCreate: invalid RoleName %q", p.RoleName)
	}
	if p.Status != "" {
		if !validator.IsValidStatus(p.Status) {
			return fmt.Errorf("Role.IsValidForCreate: invalid Status %q", p.Status)
		}
	}

	return nil
}

func (p *Role) IsValidSortKey(key string) bool {
	var validKeys = []string{
		"role_id",
		"role_name",
		"description",
		"portal",
		"owner",
		"owner_path",
		"status",
		"create_time",
		"update_time",
		"status_time",
	}
	for _, k := range validKeys {
		if strings.EqualFold(k, key) {
			return true
		}
	}
	return false
}

func (p *Role) AdjustForUpdate() *Role {
	p.RoleId = strutil.SimplifyString(p.RoleId)

	// skip readonly fields
	p.CreateTime = time.Time{}
	p.UpdateTime = time.Now()

	// adjust data
	p.RoleName = strutil.SimplifyString(p.RoleName)
	p.Description = strutil.SimplifyString(p.Description)
	p.Portal = strutil.SimplifyString(p.Portal)
	p.Owner = strutil.SimplifyString(p.Owner)
	p.OwnerPath = strutil.SimplifyString(p.OwnerPath)
	p.Status = strutil.SimplifyString(p.Status)

	if p.Status != "" {
		p.StatusTime = time.Now()
	}

	return p
}

func (p *Role) IsValidForUpdate() error {
	if !validator.IsValidId(p.RoleId) {
		return fmt.Errorf("Role.IsValidForUpdate: invalid RoleId %q", p.RoleId)
	}

	// check readonly fields
	if p.CreateTime != (time.Time{}) {
		return fmt.Errorf("Role.IsValidForUpdate: CreateTime is readonly")
	}

	// check updated fields
	if p.RoleName != "" {
		if !validator.IsValidName(p.RoleName) {
			return fmt.Errorf("Role.IsValidForUpdate: invalid RoleName %q", p.RoleName)
		}
	}

	if p.Status != "" {
		if !validator.IsValidStatus(p.Status) {
			return fmt.Errorf("Role.IsValidForUpdate: invalid Status %q", p.Status)
		}
	}

	if p.Status != "" {
		if !validator.IsValidStatus(p.Status) {
			return fmt.Errorf("Role.IsValidForUpdate: invalid Status %q", p.Status)
		}
		if p.StatusTime == (time.Time{}) {
			return fmt.Errorf("Role.IsValidForUpdate: invalid StatusTime %q", p.StatusTime)
		}
	}

	return nil
}
