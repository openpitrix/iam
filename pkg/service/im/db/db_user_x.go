// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db

import (
	"context"

	"openpitrix.io/iam/pkg/pb/im"
)

func (p *Database) GetUsersByGroupId(ctx context.Context, req *pbim.GroupId) (*pbim.UserList, error) {
	panic("todo")
}

func (p *Database) ComparePassword(ctx context.Context, req *pbim.Password) (*pbim.Empty, error) {
	panic("todo")
}
func (p *Database) ModifyPassword(ctx context.Context, req *pbim.Password) (*pbim.Empty, error) {
	panic("todo")
}
