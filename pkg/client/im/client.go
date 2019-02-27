// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package im

import (
	pbim "kubesphere.io/im/pkg/pb"

	"openpitrix.io/iam/pkg/global"
	"openpitrix.io/iam/pkg/manager"
)

type Client struct {
	pbim.IdentityManagerClient
}

func NewClient() (*Client, error) {
	conn, err := manager.NewClient(global.Global().Config.IMHost, global.Global().Config.IMPort)
	if err != nil {
		return nil, err
	}

	return &Client{
		IdentityManagerClient: pbim.NewIdentityManagerClient(conn),
	}, nil
}
