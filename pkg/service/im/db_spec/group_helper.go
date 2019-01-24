// Copyright 2018 The OpenPitrix Authors. All rights reserved.
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

func (p *UserGroup) AdjustForCreate() {
	if p.GroupId == "" {
		p.GroupId = idpkg.GenId("gid-", 12)
	}

	// fix parent id
	if p.ParentGroupId == "" && p.GroupPath != "" {
		switch {
		case p.GroupPath == p.GroupId:
			p.ParentGroupId = ""
		case strings.HasSuffix(p.GroupPath, "."+p.GroupId):
			// parentGroupPath: a.b.ParentGroupId
			parentGroupPath := p.GroupPath[:len(p.GroupPath)-len(p.GroupId)-1]
			if idx := strings.LastIndex(parentGroupPath, "."); idx >= 0 {
				p.ParentGroupId = parentGroupPath[idx:]
			} else {
				p.ParentGroupId = parentGroupPath
			}
		default:
			// invalid group_path
		}
	}

	// fix root group_path
	if p.ParentGroupId == "" && p.GroupPath == "" {
		p.GroupPath = p.GroupId
	}

	p.GroupName = strutil.SimplifyString(p.GroupName)
	p.Description = strutil.SimplifyString(p.Description)
	p.Status = strutil.SimplifyString(p.Status)

	now := time.Now()
	p.CreateTime = now
	p.UpdateTime = now
	p.StatusTime = now

	// fix group_path_level
	p.GroupPathLevel = strings.Count(p.GroupPath, ".") + 1
}

func (p *UserGroup) IsValidForCreate() error {
	if p.ParentGroupId != "" {
		if !validator.IsValidId(p.ParentGroupId) {
			return fmt.Errorf("UserGroup.IsValidForCreate: invalid ParentGroupId %q", p.ParentGroupId)
		}
	}

	if !validator.IsValidId(p.GroupId) {
		return fmt.Errorf("UserGroup.IsValidForCreate: invalid GroupId %q", p.GroupId)
	}
	if !validator.IsValidGroupPath(p.GroupPath) {
		return fmt.Errorf("UserGroup.IsValidForCreate: invalid GroupId %q", p.GroupId)
	}
	if !strings.Contains(p.GroupPath, p.GroupId) {
		return fmt.Errorf(
			"UserGroup.IsValidForCreate: GroupPath(%q) must contain GroupId(%q)",
			p.GroupPath, p.GroupId,
		)
	}

	if p.ParentGroupId == "" { // root group
		if p.GroupPath != p.GroupId {
			return fmt.Errorf(
				"UserGroup.IsValidForCreate: RootGroup: GroupPath(%q) != GroupId(%q)",
				p.GroupPath, p.GroupId,
			)
		}
	} else { // sub group
		if !strings.HasSuffix(p.GroupPath, "."+p.GroupId) {
			return fmt.Errorf(
				"UserGroup.IsValidForCreate: invalid GroupPath(%q), GroupId(%q)",
				p.GroupPath, p.GroupId,
			)
		}
		if !strings.Contains(p.GroupPath, p.ParentGroupId) {
			return fmt.Errorf(
				"UserGroup.IsValidForCreate: GroupPath(%q) must contain ParentGroupId(%q)",
				p.GroupPath, p.ParentGroupId,
			)
		}
	}

	if !validator.IsValidName(p.GroupName) {
		return fmt.Errorf("UserGroup.IsValidForCreate: invalid GroupName %q", p.GroupName)
	}

	if p.Status != "" {
		if !validator.IsValidStatus(p.Status) {
			return fmt.Errorf("UserGroup.IsValidForCreate: invalid Status %q", p.Status)
		}
	}

	if p.GroupPathLevel <= 0 {
		return fmt.Errorf("UserGroup.IsValidForCreate: invalid GroupPathLevel %d", p.GroupPathLevel)
	}

	return nil
}

func (p *UserGroup) AdjustForUpdate() error {
	return nil
}

func (p *UserGroup) IsValidForUpdate() error {
	if !validator.IsValidId(p.GroupId) {
		return fmt.Errorf("UserGroup.IsValidForUpdate: invalid GroupId %q", p.GroupId)
	}

	// check readable
	if p.ParentGroupId != "" {
		return fmt.Errorf("UserGroup.IsValidForUpdate: ParentGroupId is readonly!")
	}
	if p.GroupPath != "" {
		return fmt.Errorf("UserGroup.IsValidForUpdate: GroupPath is readonly!")
	}
	if p.CreateTime != (time.Time{}) {
		return fmt.Errorf("UserGroup.IsValidForUpdate: CreateTime is readonly!")
	}

	// check updated fields
	if p.Status != "" {
		if !validator.IsValidStatus(p.Status) {
			return fmt.Errorf("UserGroup.IsValidForUpdate: invalid Status %q", p.Status)
		}
	}

	// OK
	return nil
}
