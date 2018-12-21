// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package service

import (
	"context"

	"openpitrix.io/iam/openpitrix/pkg/pb"
)

func (p *Server) DescribeActions(context.Context, *pb.Range) (*pb.ActionList, error) {
	panic("TODO")
}

func (p *Server) GetOwnerPath(context.Context, *pb.GetOwnerPathRequest) (*pb.String, error) {
	panic("TODO")
}

func (p *Server) GetAccessPath(context.Context, *pb.GetAccessPathRequest) (*pb.String, error) {
	panic("TODO")
}

func (p *Server) CanDoAction(context.Context, *pb.CanDoActionRequest) (*pb.CanDoActionResponse, error) {
	panic("TODO")
}
