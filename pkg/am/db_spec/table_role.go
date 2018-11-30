// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db_spec

import (
	"encoding/json"
	"log"
	"strings"

	"openpitrix.io/iam/pkg/pb/am"
)

const (
	RoleTableName = "iam_role"

	RoleTableSchema = `
		CREATE TABLE IF NOT EXISTS ` + RoleTableName + ` (
			name VARCHAR(50)  NOT NULL,
			rule TEXT         NOT NULL,

			PRIMARY KEY (name)
		);
	`
	RoleTableSchema_MySQL = `
		CREATE TABLE IF NOT EXISTS ` + RoleTableName + ` (
			name VARCHAR(50)  NOT NULL,
			rule JSON         NOT NULL, -- use JSON type

			PRIMARY KEY (name)
		);
	`
)

type Role struct {
	Name string
	Rule string // JSON string
}

func (Role) GetTableName() string {
	return RoleTableName
}

func (Role) GetTableSchema(dbtype string) string {
	if strings.EqualFold(dbtype, "mysql") {
		return RoleTableSchema_MySQL
	} else {
		return RoleTableSchema
	}
}

func NewRoleFrom(x *pbam.Role) *Role {
	return &Role{
		Name: x.Name,
		Rule: encodeRuleList(x.Rule),
	}
}

func (p *Role) ToPbRole() *pbam.Role {
	return &pbam.Role{
		Name: p.Name,
		Rule: decodeRuleList(p.Rule),
	}
}

func encodeRuleList(x []*pbam.Rule) string {
	if x == nil {
		x = []*pbam.Rule{}
	}

	data, err := json.Marshal(x)
	if err != nil {
		log.Panic(err)
	}

	return string(data)
}

func decodeRuleList(jsonValue string) []*pbam.Rule {
	if jsonValue == "" {
		return []*pbam.Rule{}
	}

	var x []*pbam.Rule
	err := json.Unmarshal([]byte(jsonValue), &x)
	if err != nil {
		log.Panic(err)
	}

	return x
}
