// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db

import (
	"bytes"
	"encoding/gob"

	pbam "openpitrix.io/iam/pkg/pb/am"
)

// keep same pbam.Action
type DBRecord struct {
	RoleId   string
	RoleName string
	Portal   string

	ModuleId   string
	ModuleName string

	FeatureId   string
	FeatureName string
	DataLevel   string

	ActionId      string
	ActionName    string
	ActionEnabled string

	ApiId          string
	ApiMethod      string
	ApiDescription string

	Url       string
	UrlMethod string
}

func NewRecordFromPB(p *pbam.Action) *DBRecord {
	if p == nil {
		return new(DBRecord)
	}

	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(p); err != nil {
		// return nil, err
	}

	var q = new(DBRecord)
	if err := gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(q); err != nil {
		// return nil, err
	}

	return q
}

func (p *DBRecord) ToPB() *pbam.Action {
	q, err := p.ToProtoMessage()
	if err != nil {
		panic(err) // unreachable
	}
	return q
}

func (p *DBRecord) ToProtoMessage() (*pbam.Action, error) {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(p); err != nil {
		return nil, err
	}

	var q = new(pbam.Action)
	if err := gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(q); err != nil {
		return nil, err
	}

	return q, nil
}
