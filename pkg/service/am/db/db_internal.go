// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db

import (
	"context"

	"openpitrix.io/iam/pkg/service/am/db_spec"

	pbam "openpitrix.io/iam/pkg/pb/am"
)

func (p *Database) GetUser(ctx context.Context, req *pbam.String) (*pbam.UserWithRole, error) {
	var user db_spec.DBUserWithRole
	if err := p.ormDB.Table("user").Where("user_id = ?", req.Value).First(&user).Error; err != nil {
		return nil, err
	}
	return user.ToPB(), nil
}
