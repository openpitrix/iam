// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db_spec

import (
	"time"

	"github.com/golang/protobuf/ptypes"

	"openpitrix.io/iam/pkg/pb/im"
)

type DBGroup struct {
	GroupId       string    `db:"group_id"`
	GroupName     string    `db:"group_name"`
	ParentGroupId string    `db:"parent_group_id"`
	GroupPath     string    `db:"group_path"`
	Level         int32     `db:"level"`
	Status        string    `db:"status"`
	CreateTime    time.Time `db:"create_time"`
	StatusTime    time.Time `db:"status_time"`
	UpdateTime    time.Time `db:"update_time"`
	Extra         string    `db:"extra"` // JSON
}

func PBGroupToDB(p *pbim.Group) *DBGroup {
	if p == nil {
		return new(DBGroup)
	}
	var q = &DBGroup{
		//GroupId:       p.GroupId,
		//GroupName:     p.GroupName,
		//ParentGroupId: p.ParentGroupId,
		//GroupPath:     p.GroupPath,
		//Level:         p.Level,
		//SeqOrder:      p.SeqOrder,
		//Owner:         p.Owner,
		//OwnerPath:     p.OwnerPath,
		//CreateTime:    p.CreateTime,
		//UpdateTime:    p.UpdateTime,
	}

	q.CreateTime, _ = ptypes.Timestamp(p.CreateTime)
	q.UpdateTime, _ = ptypes.Timestamp(p.UpdateTime)

	return q
}

func (p *DBGroup) ToPB() *pbim.Group {
	if p == nil {
		return new(pbim.Group)
	}
	var q = &pbim.Group{
		//GroupId:       p.GroupId,
		//GroupName:     p.GroupName,
		//ParentGroupId: p.ParentGroupId,
		//GroupPath:     p.GroupPath,
		//Level:         p.Level,
		//SeqOrder:      p.SeqOrder,
		//Owner:         p.Owner,
		//OwnerPath:     p.OwnerPath,
		//CreateTime: p.CreateTime,
		//UpdateTime: p.UpdateTime,
	}

	q.CreateTime, _ = ptypes.TimestampProto(p.CreateTime)
	q.UpdateTime, _ = ptypes.TimestampProto(p.UpdateTime)

	return q
}

func (p *DBGroup) ValidateForInsert() error {
	return nil
}
func (p *DBGroup) ValidateForUpdate() error {
	return nil
}
