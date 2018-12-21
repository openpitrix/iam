// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package service

import (
	"context"

	"openpitrix.io/iam/openpitrix/pkg/pb"
)

func (p *Server) ComparePassword(context.Context, *pb.UserPassword) (*pb.Bool, error) {
	panic("TODO")
}
func (p *Server) ModifyPassword(context.Context, *pb.UserPassword) (*pb.Bool, error) {
	panic("TODO")
}
