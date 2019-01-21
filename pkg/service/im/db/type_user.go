// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/golang/protobuf/ptypes"
	"golang.org/x/crypto/bcrypt"

	"openpitrix.io/iam/pkg/pb/im"
	"openpitrix.io/logger"
)

type User struct {
	UserId      string
	UserName    string
	Email       string
	PhoneNumber string
	Description string
	Password    string
	Status      string
	CreateTime  time.Time
	UpdateTime  time.Time
	StatusTime  time.Time
	Extra       string // JSON
}

func NewUserFromPB(p *pbim.User) *User {
	if p == nil {
		return new(User)
	}
	var q = &User{
		UserId:      p.UserId,
		UserName:    p.UserName,
		Email:       p.Email,
		PhoneNumber: p.PhoneNumber,
		Description: p.Description,
		Password:    p.Password,
		Status:      p.Status,
	}

	q.CreateTime, _ = ptypes.Timestamp(p.CreateTime)
	q.UpdateTime, _ = ptypes.Timestamp(p.UpdateTime)
	q.StatusTime, _ = ptypes.Timestamp(p.StatusTime)

	if len(p.Extra) > 0 {
		data, err := json.MarshalIndent(p.Extra, "", "\t")
		if err != nil {
			logger.Warnf(nil, "%+v", err)
			return q
		}
		q.Extra = string(data)
	}
	return q
}

func (p *User) ToPB() *pbim.User {
	if p == nil {
		return new(pbim.User)
	}
	var q = &pbim.User{
		UserId:      p.UserId,
		UserName:    p.UserName,
		Email:       p.Email,
		PhoneNumber: p.PhoneNumber,
		Description: p.Description,
		Password:    p.Password,
		Status:      p.Status,
	}

	q.CreateTime, _ = ptypes.TimestampProto(p.CreateTime)
	q.UpdateTime, _ = ptypes.TimestampProto(p.UpdateTime)
	q.StatusTime, _ = ptypes.TimestampProto(p.StatusTime)

	if p.Extra != "" {
		if q.Extra == nil {
			q.Extra = make(map[string]string)
		}
		err := json.Unmarshal([]byte(p.Extra), &q.Extra)
		if err != nil {
			logger.Warnf(nil, "%+v", err)
			return q
		}
	}
	return q
}

func (p *User) BeforeCreate() (err error) {
	if p.UserId == "" {
		p.UserId = genUid()
	}
	if p.Password != "" {
		hashedPass, err := bcrypt.GenerateFromPassword(
			[]byte(p.Password), bcrypt.DefaultCost,
		)
		if err != nil {
			return err
		}
		p.Password = string(hashedPass)
	}

	if p.CreateTime == (time.Time{}) {
		p.CreateTime = time.Now()
	}

	return
}
func (p *User) BeforeUpdate() (err error) {
	if p.UpdateTime == (time.Time{}) {
		p.UpdateTime = time.Now()
	}
	p.Password = ""
	return
}

func (p *User) ValidateForInsert() error {
	if !reUserId.MatchString(p.UserId) {
		return fmt.Errorf("invalid uid: %q", p.UserId)
	}
	if p.Password == "" {
		return fmt.Errorf("invalid password")
	}

	if p.Extra != "" {
		var m = make(map[string]string)
		if err := json.Unmarshal([]byte(p.Extra), &m); err != nil {
			return fmt.Errorf("invalid extra")
		}
	}

	return nil
}
func (p *User) ValidateForUpdate() error {
	if !reUserId.MatchString(p.UserId) {
		return fmt.Errorf("invalid uid: %q", p.UserId)
	}

	if p.Extra != "" {
		var m = make(map[string]string)
		if err := json.Unmarshal([]byte(p.Extra), &m); err != nil {
			return fmt.Errorf("invalid extra")
		}
	}

	return nil
}
