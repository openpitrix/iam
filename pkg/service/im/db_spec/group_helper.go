// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db_spec

import (
	"fmt"
	"time"

	"openpitrix.io/iam/pkg/validator"
)

func (p *UserGroup) AdjustForCreate() error {
	return nil
}
func (p *UserGroup) AdjustForUpdate() error {
	return nil
}

func (p *UserGroup) IsValidForCreate() error {
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
