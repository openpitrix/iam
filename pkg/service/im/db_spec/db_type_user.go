// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db_spec

import (
	"encoding/json"
	"time"

	"github.com/golang/protobuf/ptypes"

	"openpitrix.io/iam/pkg/pb/im"
	"openpitrix.io/logger"
)

type DBUser struct {
	Uid         string    `db:"user_id"`
	Name        string    `db:"user_name"`
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
		Uid:         p.Uid,
		Name:        p.Name,
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
		Uid:         p.Uid,
		Name:        p.Name,
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
	return nil
}
func (p *DBUser) ValidateForUpdate() error {
	return nil
}
