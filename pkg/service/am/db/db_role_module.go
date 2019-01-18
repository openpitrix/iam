// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db

import (
	"context"

	"openpitrix.io/iam/pkg/internal/funcutil"
	pbam "openpitrix.io/iam/pkg/pb/am"
	"openpitrix.io/logger"
)

func (p *Database) ModifyRoleModuleBinding(ctx context.Context, req *pbam.ModifyRoleModuleBindingRequest) (*pbam.ActionList, error) {
	logger.Infof(ctx, funcutil.CallerName(1))

	logger.Infof(ctx, "TODO")
	return nil, nil
}
