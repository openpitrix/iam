// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package am

import (
	"context"

	pbam "openpitrix.io/iam/pkg/pb/am"
)

var _ pbam.InternalAccessManagerServer = (*Server)(nil)

func (p *Server) GetUser(ctx context.Context, req *pbam.String) (*pbam.UserWithRole, error) {
	return p.db.GetUser(ctx, req)
}
