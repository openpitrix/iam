// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db

import (
	"context"

	"openpitrix.io/iam/pkg/pb/im"
)

func (p *Database) JoinGroup(ctx context.Context, req *pbim.JoinGroupRequest) (*pbim.Empty, error) {
	panic("todo")
}

func (p *Database) LeaveGroup(ctx context.Context, req *pbim.LeaveGroupRequest) (*pbim.Empty, error) {
	panic("todo")
}
