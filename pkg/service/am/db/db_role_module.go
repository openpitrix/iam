// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db

import (
	"context"

	pbam "openpitrix.io/iam/pkg/pb/am"
)

func (p *Database) GetRoleModule(context.Context, *pbam.RoleId) (*pbam.RoleModule, error) {
	panic("todo")
}
func (p *Database) ModifyRoleModule(context.Context, *pbam.RoleModule) (*pbam.RoleModule, error) {
	panic("todo")
}
