// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package web

import (
	"fmt"
	"net/http"
	"time"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"openpitrix.io/iam/pkg/pb/am"
	"openpitrix.io/iam/pkg/pb/im"
)

type GrpcServer func(grpcServer *grpc.Server)

func WithAccountManager(s pbim.AccountManagerServer) GrpcServer {
	return func(grpcServer *grpc.Server) {
		pbim.RegisterAccountManagerServer(grpcServer, s)
	}
}
func WithAccessManager(s pbam.AccessManagerServer) GrpcServer {
	return func(grpcServer *grpc.Server) {
		pbam.RegisterAccessManagerServer(grpcServer, s)
	}
}

func ListenAndServe(addr string, servers []GrpcServer,
	opts ...grpc.ServerOption,
) error {
	if len(opts) == 0 {
		opts = MakeDefaultGrpcServerOptions()
	}

	grpcServer := grpc.NewServer(opts...)
	reflection.Register(grpcServer)
	for _, fn := range servers {
		fn(grpcServer)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "hello", time.Now()) // just for test
	})

	// https://github.com/grpc/grpc-go/issues/555#issuecomment-443293451
	server := &http.Server{
		Addr: addr,
		Handler: h2c.NewHandler(
			MainHandler(grpcServer, addr, mux),
			&http2.Server{},
		),
	}

	return server.ListenAndServe()
}

func ListenAndServeTLS(addr, certFile, keyFile string, servers []GrpcServer,
	opts ...grpc.ServerOption,
) error {
	if len(opts) == 0 {
		opts = MakeDefaultGrpcServerOptions()
	}

	grpcServer := grpc.NewServer(opts...)
	reflection.Register(grpcServer)
	for _, fn := range servers {
		fn(grpcServer)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "hello", time.Now()) // just for test
	})

	server := &http.Server{Addr: addr,
		Handler: MainHandler(grpcServer, addr, mux),
	}

	return server.ListenAndServeTLS(certFile, keyFile)
}
