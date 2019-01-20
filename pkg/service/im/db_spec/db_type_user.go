// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db_spec

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/golang/protobuf/ptypes"

	"openpitrix.io/iam/pkg/pb/im"
	"openpitrix.io/logger"
)

type DBUser struct {
	UserId      string    `db:"user_id"`
	UserName    string    `db:"user_name"`
	Email       string    `db:"email"`
	PhoneNumber string    `db:"phone_number"`
	Description string    `db:"description"`
	Password    string    `db:"password"`
	Status      string    `db:"status"`
	CreateTime  time.Time `db:"create_time"`
	UpdateTime  time.Time `db:"update_time"`
	StatusTime  time.Time `db:"status_time"`
	Extra       string    `db:"extra"` // JSON
}

func PBUserToDB(p *pbim.User) *DBUser {
	if p == nil {
		return new(DBUser)
	}
	var q = &DBUser{
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

func (p *DBUser) ToPB() *pbim.User {
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

func (p *DBUser) ValidateForInsert() error {
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
func (p *DBUser) ValidateForUpdate() error {
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
