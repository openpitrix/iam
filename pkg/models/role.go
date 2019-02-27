// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

import (
	"time"

	"github.com/golang/protobuf/ptypes"

	"openpitrix.io/iam/pkg/constants"
	"openpitrix.io/iam/pkg/pb"
	"openpitrix.io/iam/pkg/util/idutil"
)

type Role struct {
	RoleId      string `gorm:"type:varchar(50);primary_key"`
	RoleName    string `gorm:"type:varchar(200);not null"`
	Description string `gorm:"type:varchar(200)"`
	Portal      string `gorm:"type:varchar(50);not null"`
	Owner       string `gorm:"type:varchar(50)"`
	OwnerPath   string `gorm:"type:varchar(50)"`
	Status      string `gorm:"type:varchar(50)"`
	Controller  string `gorm:"type:varchar(50)"`

	CreateTime time.Time
	UpdateTime time.Time
	StatusTime time.Time
}

type RoleWithUser struct {
	Role    *Role
	UserIds []string
}

func NewRole(roleName, description, portal, owner, ownerPath string) *Role {
	now := time.Now()
	role := &Role{
		RoleId:      idutil.GetUuid(constants.PrefixRoleId),
		RoleName:    roleName,
		Description: description,
		Portal:      portal,
		Owner:       owner,
		OwnerPath:   ownerPath,
		Status:      constants.StatusActive,
		Controller:  constants.ControllerSelf,
		CreateTime:  now,
		UpdateTime:  now,
		StatusTime:  now,
	}
	return role
}

func (p *Role) ToPB() *pb.Role {
	if p == nil {
		return new(pb.Role)
	}
	var q = &pb.Role{
		RoleId:      p.RoleId,
		RoleName:    p.RoleName,
		Description: p.Description,
		Portal:      p.Portal,
		Owner:       p.Owner,
		OwnerPath:   p.OwnerPath,
		Status:      p.Status,
		Controller:  p.Controller,
	}

	q.CreateTime, _ = ptypes.TimestampProto(p.CreateTime)
	q.UpdateTime, _ = ptypes.TimestampProto(p.UpdateTime)
	q.StatusTime, _ = ptypes.TimestampProto(p.StatusTime)

	return q
}

func (p *RoleWithUser) ToPB() *pb.RoleWithUser {
	if p == nil {
		return new(pb.RoleWithUser)
	}

	return &pb.RoleWithUser{
		Role:      p.Role.ToPB(),
		UserIdSet: p.UserIds,
	}
}
