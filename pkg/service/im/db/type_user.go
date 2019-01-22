// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db

import (
	"encoding/json"
	"fmt"
	"regexp"
	"time"

	"github.com/golang/protobuf/ptypes"
	"golang.org/x/crypto/bcrypt"

	"openpitrix.io/iam/pkg/pb/im"
	"openpitrix.io/logger"
)

type User struct {
	UserId      string `gorm:"primary_key"`
	UserName    string `gorm:"type:varchar(50)"`
	Email       string `gorm:"type:varchar(50)"`
	PhoneNumber string `gorm:"type:varchar(50)"`
	Description string `gorm:"type:varchar(1000)"`
	Password    string `gorm:"type:varchar(128)"`
	Status      string `gorm:"type:varchar(10)"`
	CreateTime  time.Time
	UpdateTime  time.Time
	StatusTime  time.Time
	Extra       *string `gorm:"type:JSON"`
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

		q.Extra = newString(string(data))
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

	if p.Extra != nil && *p.Extra != "" {
		if q.Extra == nil {
			q.Extra = make(map[string]string)
		}
		err := json.Unmarshal([]byte(*p.Extra), &q.Extra)
		if err != nil {
			logger.Warnf(nil, "%+v", err)
			return q
		}
	}
	return q
}

func (p *User) BeforeCreate() (err error) {
	if p.UserId == "" {
		p.UserId = genId("uid-", 12)
	} else {
		var re = regexp.MustCompile(`^[a-zA-Z0-9-_]+$`)
		if !re.MatchString(p.UserId) {
			return fmt.Errorf("invalid UserId: %s", p.UserId)
		}
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

	now := time.Now()
	p.CreateTime = now
	p.UpdateTime = now
	p.StatusTime = now

	return
}
func (p *User) BeforeUpdate() (err error) {
	if p.UpdateTime == (time.Time{}) {
		p.UpdateTime = time.Now()
	}

	// ignore readonly fields
	p.CreateTime = time.Time{}
	p.Password = ""

	return
}

func (p *User) ValidateForInsert() error {
	if !isValidIds(p.UserId) {
		return fmt.Errorf("invalid uid: %q", p.UserId)
	}
	if p.Password == "" {
		return fmt.Errorf("invalid password")
	}

	if p.Extra != nil && *p.Extra != "" {
		var m = make(map[string]string)
		if err := json.Unmarshal([]byte(*p.Extra), &m); err != nil {
			return fmt.Errorf("invalid extra")
		}
	}

	return nil
}
func (p *User) ValidateForUpdate() error {
	if !isValidIds(p.UserId) {
		return fmt.Errorf("invalid uid: %q", p.UserId)
	}

	if p.Extra != nil && *p.Extra != "" {
		var m = make(map[string]string)
		if err := json.Unmarshal([]byte(*p.Extra), &m); err != nil {
			return fmt.Errorf("invalid extra")
		}
	}

	return nil
}
