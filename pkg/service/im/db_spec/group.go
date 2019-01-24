// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db_spec

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/golang/protobuf/ptypes"

	"openpitrix.io/iam/pkg/internal/strutil"
	"openpitrix.io/iam/pkg/pb/im"
)

type UserGroup struct {
	ParentGroupId string `gorm:"type:varchar(50);not null"`
	GroupId       string `gorm:"primary_key"`
	GroupPath     string `gorm:"type:varchar(255);not null;unique"`
	GroupName     string `gorm:"type:varchar(50);not null"`
	Description   string `gorm:"type:varchar(1000);not null"`
	Status        string `gorm:"type:varchar(10);not null"`
	CreateTime    time.Time
	UpdateTime    time.Time
	StatusTime    time.Time
	Extra         *string `gorm:"type:JSON"`

	// internal
	GroupPathLevel int
}

func NewUserGroupFromPB(p *pbim.Group) *UserGroup {
	if p == nil {
		return new(UserGroup)
	}
	var q = &UserGroup{
		ParentGroupId: p.ParentGroupId,
		GroupId:       p.GroupId,
		GroupPath:     p.GroupPath,
		GroupName:     p.GroupName,
		Description:   p.Description,
		Status:        p.Status,

		GroupPathLevel: 0,
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

	if p.GroupPath != "" {
		q.GroupPathLevel = strings.Count(p.GroupPath, ".") + 1
	}

	return q
}

func (p *UserGroup) ToPB() *pbim.Group {
	q, _ := p.ToProtoMessage()
	return q
}

func (p *UserGroup) ToProtoMessage() (*pbim.Group, error) {
	if p == nil {
		return new(pbim.Group), nil
	}
	var q = &pbim.Group{
		ParentGroupId: p.ParentGroupId,
		GroupId:       p.GroupId,
		GroupPath:     p.GroupPath,
		GroupName:     p.GroupName,
		Description:   p.Description,
		Status:        p.Status,
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
