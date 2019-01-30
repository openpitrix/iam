// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db_spec

import (
	"fmt"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	idpkg "openpitrix.io/iam/pkg/id"
	"openpitrix.io/iam/pkg/internal/strutil"
	"openpitrix.io/iam/pkg/validator"
)

func (p *User) IsValidSortKey(key string) bool {
	var validKeys = []string{
		"user_id",
		"user_name",
		"email",
		"phone_number",
		"description",
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

	if p.Password != "" {
		hashedPass, _ := bcrypt.GenerateFromPassword(
			[]byte(p.Password), bcrypt.DefaultCost,
		)
		p.Password = string(hashedPass)
	}

	now := time.Now()
	p.CreateTime = now
	p.UpdateTime = now
	p.StatusTime = now

	return p
}

func (p *User) IsValidForCreate() error {
	if !validator.IsValidId(p.UserId) {
		return fmt.Errorf("User.IsValidForCreate: invalid UserId %q", p.UserId)
	}
	if !validator.IsValidName(p.UserName) {
		return fmt.Errorf("User.IsValidForCreate: invalid UserName %q", p.UserName)
	}
	if !validator.IsValidEmail(p.Email) {
		return fmt.Errorf("User.IsValidForCreate: invalid Email %q", p.Email)
	}

	if p.PhoneNumber != "" {
		if !validator.IsValidPhoneNumber(p.PhoneNumber) {
			return fmt.Errorf("User.IsValidForCreate: invalid PhoneNumber %q", p.PhoneNumber)
		}
	}
	if p.Status != "" {
		if !validator.IsValidStatus(p.Status) {
			return fmt.Errorf("User.IsValidForCreate: invalid Status %q", p.Status)
		}
	}

	if p.Password == "" {
		return fmt.Errorf("User.IsValidForCreate: invalid Password")
	}

	return nil
}

func (p *User) AdjustForUpdate() *User {
	p.UserId = strutil.SimplifyString(p.UserId)

	// skip readonly fields
	p.CreateTime = time.Time{}
	p.UpdateTime = time.Now()

	// adjust data
	p.UserName = strutil.SimplifyString(p.UserName)
	p.Email = strutil.SimplifyString(p.Email)
	p.PhoneNumber = strutil.SimplifyString(p.PhoneNumber)
	p.Description = strutil.SimplifyString(p.Description)
	p.Status = strutil.SimplifyString(p.Status)

	if p.Password != "" {
		hashedPass, _ := bcrypt.GenerateFromPassword(
			[]byte(p.Password), bcrypt.DefaultCost,
		)
		p.Password = string(hashedPass)
	}

	if p.Status != "" {
		p.StatusTime = time.Now()
	}

	return p
}

func (p *User) IsValidForUpdate() error {
	if !validator.IsValidId(p.UserId) {
		return fmt.Errorf("User.IsValidForUpdate: invalid UserId %q", p.UserId)
	}

	// check readonly fields
	if p.CreateTime != (time.Time{}) {
		return fmt.Errorf("User.IsValidForUpdate: CreateTime is readonly")
	}

	// check updated fields
	if p.UserName != "" {
		if !validator.IsValidName(p.UserName) {
			return fmt.Errorf("User.IsValidForUpdate: invalid UserName %q", p.UserName)
		}
	}
	if p.Email != "" {
		if !validator.IsValidEmail(p.Email) {
			return fmt.Errorf("User.IsValidForUpdate: invalid Email %q", p.Email)
		}
	}
	if p.PhoneNumber != "" {
		if !validator.IsValidPhoneNumber(p.PhoneNumber) {
			return fmt.Errorf("User.IsValidForUpdate: invalid PhoneNumber %q", p.PhoneNumber)
		}
	}
	if p.Status != "" {
		if !validator.IsValidStatus(p.Status) {
			return fmt.Errorf("User.IsValidForUpdate: invalid Status %q", p.Status)
		}
	}

	if p.Status != "" {
		if !validator.IsValidStatus(p.Status) {
			return fmt.Errorf("User.IsValidForUpdate: invalid Status %q", p.Status)
		}
		if p.StatusTime == (time.Time{}) {
			return fmt.Errorf("User.IsValidForUpdate: invalid StatusTime %q", p.StatusTime)
		}
	}

	return nil
}
