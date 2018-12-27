// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package service

import (
	"context"

	"openpitrix.io/iam/openpitrix/pkg/pb"
	"openpitrix.io/iam/openpitrix/pkg/version"
)

func (p *Server) GetVersion(ctx context.Context, req *pb.Empty) (*pb.String, error) {
	reply := &pb.String{Value: version.GetVersionString()}
	return reply, nil
}
