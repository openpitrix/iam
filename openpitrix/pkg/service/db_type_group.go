// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package service

import (
	"time"

	"github.com/golang/protobuf/ptypes"

	"openpitrix.io/iam/openpitrix/pkg/pb"
)

type DBGroup struct {
	GroupId       string    `db:"group_id"`
	GroupName     string    `db:"group_name"`
	ParentGroupId string    `db:"parent_group_id"`
	GroupPath     string    `db:"group_path"`
	Level         int32     `db:"level"`
	SeqOrder      int32     `db:"seq_order"`
	Owner         string    `db:"owner"`
	OwnerPath     string    `db:"owner_path"`
	CreateTime    time.Time `db:"create_time"`
	UpdateTime    time.Time `db:"update_time"`
}

func pbGroupToDB(p *pb.Group) *DBGroup {
	var q = &DBGroup{
		GroupId:       p.GroupId,
		GroupName:     p.GroupName,
		ParentGroupId: p.ParentGroupId,
		GroupPath:     p.GroupPath,
		Level:         p.Level,
		SeqOrder:      p.SeqOrder,
		Owner:         p.Owner,
		OwnerPath:     p.OwnerPath,
		//CreateTime:    p.CreateTime,
		//UpdateTime:    p.UpdateTime,
	}

	q.CreateTime, _ = ptypes.Timestamp(p.CreateTime)
	q.UpdateTime, _ = ptypes.Timestamp(p.UpdateTime)

	return q
}

func (p *DBGroup) ToPb() *pb.Group {
	var q = &pb.Group{
		GroupId:       p.GroupId,
		GroupName:     p.GroupName,
		ParentGroupId: p.ParentGroupId,
		GroupPath:     p.GroupPath,
		Level:         p.Level,
		SeqOrder:      p.SeqOrder,
		Owner:         p.Owner,
		OwnerPath:     p.OwnerPath,
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
