// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db

import (
	"bytes"
	"encoding/gob"

	pbam "openpitrix.io/iam/pkg/pb/am"
)

type DBAction struct {
	RoleId   string
	RoleName string
	Portal   string

	ModuleId   string
	ModuleName string
	DataLevel  string

	FeatureId   string
	FeatureName string

	ActionId      string
	ActionName    string
	ActionEnabled string

	ApiId          string
	ApiMethod      string
	ApiDescription string

	Url       string
	UrlMethod string
}

func (p *DBAction) ToPB() *pbam.Action {
	var buf bytes.Buffer
	gob.NewEncoder(&buf).Encode(p)

	var q = new(pbam.Action)
	gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(q)

	return q
}
