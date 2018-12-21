// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package iam

type Server struct {
	p int
}

func OpenServer(dbtype, dbpath string) (*Server, error) {
	panic("TODO")
}

func (p *Server) Close() error {
	panic("TODO")
}
