// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/golang/protobuf/ptypes"

	"openpitrix.io/iam/pkg/pb/im"
	"openpitrix.io/logger"
)

type UserGroup struct {
	GroupId     string    `db:"group_id"`
	GroupPath   string    `db:"group_path"`
	GroupName   string    `db:"group_name"`
	Description string    `db:"description"`
	Status      string    `db:"status"`
	CreateTime  time.Time `db:"create_time"`
	UpdateTime  time.Time `db:"update_time"`
	StatusTime  time.Time `db:"status_time"`
	Extra       string    `db:"extra"` // JSON

	// DB internal fields
	ParentGroupId  string `db:"parent_group_id"`
	GroupPathLevel int    `db:"level"`
}

func NewUserGroupFromPB(p *pbim.Group) *UserGroup {
	if p == nil {
		return new(UserGroup)
	}
	var q = &UserGroup{
		GroupId:     p.GroupId,
		GroupPath:   p.GroupPath,
		GroupName:   p.GroupName,
		Description: p.Description,
		Status:      p.Status,

		// DB internal fields
		ParentGroupId:  "<nil>",
		GroupPathLevel: -1,
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

func (p *UserGroup) ToPB() *pbim.Group {
	if p == nil {
		return new(pbim.Group)
	}
	var q = &pbim.Group{
		GroupId:     p.GroupId,
		GroupPath:   p.GroupPath,
		GroupName:   p.GroupName,
		Description: p.Description,
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

func (p *UserGroup) ValidateForInsert() error {
	if !reGroupId.MatchString(p.GroupId) {
		return fmt.Errorf("invalid GroupId: %q", p.GroupId)
	}

	if p.Extra != "" {
		var m = make(map[string]string)
		if err := json.Unmarshal([]byte(p.Extra), &m); err != nil {
			return fmt.Errorf("invalid extra")
		}
	}

	return nil
}
func (p *UserGroup) ValidateForUpdate() error {
	if !reGroupId.MatchString(p.GroupId) {
		return fmt.Errorf("invalid GroupId: %q", p.GroupId)
	}

	if p.Extra != "" {
		var m = make(map[string]string)
		if err := json.Unmarshal([]byte(p.Extra), &m); err != nil {
			return fmt.Errorf("invalid extra")
		}
	}

	return nil
}
