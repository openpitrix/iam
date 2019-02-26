// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package am

import (
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"openpitrix.io/iam/pkg/config"
	"openpitrix.io/iam/pkg/global"
	"openpitrix.io/iam/pkg/manager"
	"openpitrix.io/iam/pkg/pb"
	"openpitrix.io/logger"
)

type Server struct {
}

func Serve(cfg *config.Config) {
	global.SetGlobal(cfg)
	s := new(Server)
	if cfg.TlsEnabled {
		creds, err := credentials.NewServerTLSFromFile(cfg.TlsCertFile, cfg.TlsKeyFile)
		if err != nil {
			logger.Criticalf(nil, "Constructs TLS credentials failed: %+v", err)
			os.Exit(1)
		}
		manager.NewGrpcServer(cfg.AMHost, cfg.AMPort).
			Serve(func(server *grpc.Server) {
				pb.RegisterAccessManagerServer(server, s)
				grpc.Creds(creds)
			})
	} else {
		manager.NewGrpcServer(cfg.AMHost, cfg.AMPort).
			Serve(func(server *grpc.Server) {
				pb.RegisterAccessManagerServer(server, s)
			})
	}
}
