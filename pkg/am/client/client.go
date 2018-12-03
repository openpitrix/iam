// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package client

import (
	"fmt"

	"google.golang.org/grpc"
	"openpitrix.io/iam/pkg/pb/am"
)

func DialService(host string, port int) (
	client pbam.AccessManagerClient,
	conn *grpc.ClientConn,
	err error,
) {
	conn, err = grpc.Dial(fmt.Sprintf("%s:%d", host, port), grpc.WithInsecure())
	if err != nil {
		return
	}

	client = pbam.NewAccessManagerClient(conn)
	return
}
