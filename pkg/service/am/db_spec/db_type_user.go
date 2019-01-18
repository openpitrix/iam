// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db_spec

import (
	"encoding/json"
	"time"

	"github.com/golang/protobuf/ptypes"

	pbam "openpitrix.io/iam/pkg/pb/am"
	"openpitrix.io/logger"
)

type DBUserWithRole struct {
	UserId      string    `db:"user_id" gorm:"primary_key"`
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

func PBUserWithRoleToDB(p *pbam.UserWithRole) *DBUserWithRole {
	if p == nil {
		return new(DBUserWithRole)
	}
	var q = &DBUserWithRole{
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

func (p *DBUserWithRole) ToPB() *pbam.UserWithRole {
	if p == nil {
		return new(pbam.UserWithRole)
	}
	var q = &pbam.UserWithRole{
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

func (p *DBUserWithRole) ValidateForInsert() error {
	return nil
}
func (p *DBUserWithRole) ValidateForUpdate() error {
	return nil
}
