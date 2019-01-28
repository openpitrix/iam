// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db_spec

import (
	"encoding/json"
	"time"

	"github.com/golang/protobuf/ptypes"

	"openpitrix.io/iam/pkg/internal/strutil"
	"openpitrix.io/iam/pkg/pb/im"
)

type User struct {
	UserId      string `gorm:"primary_key"`
	UserName    string `gorm:"type:varchar(50);not null;unique;"`
	Email       string `gorm:"type:varchar(50);not null;unique"`
	PhoneNumber string `gorm:"type:varchar(50);not null"`
	Description string `gorm:"type:varchar(1000);not null"`
	Password    string `gorm:"type:varchar(128);not null"`
	Status      string `gorm:"type:varchar(10);not null"`
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
		UserName:    p.Username,
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
			panic(err) // unreachable
		}

		q.Extra = strutil.NewString(string(data))
	}
	return q
}

func (p *User) ToPB() *pbim.User {
	q, _ := p.ToProtoMessage()
	return q
}

func (p *User) ToProtoMessage() (*pbim.User, error) {
	if p == nil {
		return new(pbim.User), nil
	}
	var q = &pbim.User{
		UserId:      p.UserId,
		Username:    p.UserName,
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
			return q, err
		}
	}
	return q, nil
}
