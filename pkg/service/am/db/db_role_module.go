// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db

import (
	"context"

	pbam "openpitrix.io/iam/pkg/pb/am"
	"openpitrix.io/iam/pkg/util/funcutil"
	"openpitrix.io/logger"
)

func (p *Database) GetRoleModule(ctx context.Context, req *pbam.RoleId) (*pbam.RoleModule, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	panic("todo")
}
func (p *Database) ModifyRoleModule(ctx context.Context, req *pbam.RoleModule) (*pbam.RoleModule, error) {
	panic("todo")
}
