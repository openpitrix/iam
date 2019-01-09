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

type DBGroup struct {
	Gid         string    `db:"group_id"`
	GroupPath   string    `db:"group_path"`
	Name        string    `db:"group_name"`
	Description string    `db:"description"`
	Status      string    `db:"status"`
	CreateTime  time.Time `db:"create_time"`
	UpdateTime  time.Time `db:"update_time"`
	StatusTime  time.Time `db:"status_time"`
	Extra       string    `db:"extra"` // JSON

	// DB internal fields
	parent_gid     string `db:"parent_group_id"`
	GroupPathLevel string `db:"level"`
}

func PBGroupToDB(p *pbim.Group) *DBGroup {
	if p == nil {
		return new(DBGroup)
	}
	var q = &DBGroup{
		Gid:         p.Gid,
		GroupPath:   p.GroupPath,
		Name:        p.Name,
		Description: p.Description,
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

func (p *DBGroup) ToPB() *pbim.Group {
	if p == nil {
		return new(pbim.Group)
	}
	var q = &pbim.Group{
		Gid:         p.Gid,
		GroupPath:   p.GroupPath,
		Name:        p.Name,
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

func (p *DBGroup) ValidateForInsert() error {
	return nil
}
func (p *DBGroup) ValidateForUpdate() error {
	return nil
}