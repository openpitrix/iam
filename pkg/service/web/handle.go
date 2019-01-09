// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package web

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/grpc-ecosystem/go-grpc-middleware/validator"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/tools/godoc/vfs"
	"golang.org/x/tools/godoc/vfs/httpfs"
	"golang.org/x/tools/godoc/vfs/mapfs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/status"

	"openpitrix.io/iam/pkg/pb/am"
	"openpitrix.io/iam/pkg/pb/im"
	staticSwaggerUI "openpitrix.io/iam/pkg/service/swagger-ui"
	"openpitrix.io/logger"
)

func MakeDefaultGrpcServerOptions() []grpc.ServerOption {
	return []grpc.ServerOption{
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             10 * time.Second,
			PermitWithoutStream: true,
		}),
		grpc_middleware.WithUnaryServerChain(
			grpc_validator.UnaryServerInterceptor(),
			grpc_recovery.UnaryServerInterceptor(
				grpc_recovery.WithRecoveryHandler(func(p interface{}) error {
					logger.Criticalf(nil, "GRPC server recovery with error: %+v", p)
					logger.Criticalf(nil, string(debug.Stack()))
					return status.Errorf(codes.Internal, "panic")
				}),
			),
		),
		grpc_middleware.WithStreamServerChain(
			grpc_recovery.StreamServerInterceptor(
				grpc_recovery.WithRecoveryHandler(func(p interface{}) error {
					logger.Criticalf(nil, "GRPC server recovery with error: %+v", p)
					logger.Criticalf(nil, string(debug.Stack()))
					return status.Errorf(codes.Internal, "panic")
				}),
			),
		),
	}
}

func MainHandler(
	grpcServer *grpc.Server,
	grpcServerAddress string,
	mux *http.ServeMux,
) http.Handler {
	var gwmux = runtime.NewServeMux()
	var opts = []grpc.DialOption{grpc.WithInsecure()}
	var err error

	err = pbim.RegisterAccountManagerHandlerFromEndpoint(context.Background(),
		gwmux, grpcServerAddress, opts,
	)
	if err != nil {
		logger.Warnf(nil, "%v", err)
	}
	err = pbam.RegisterAccessManagerHandlerFromEndpoint(context.Background(),
		gwmux, grpcServerAddress, opts,
	)
	if err != nil {
		logger.Warnf(nil, "%v", err)
	}

	// swagger file
	// GET /static/swagger/iam.swagger.json
	ns := vfs.NameSpace{}
	ns.Bind("/swagger", mapfs.New(staticSwaggerUI.Files), "/", vfs.BindAfter)
	if appPath, err := os.Executable(); err == nil {
		pubDir := filepath.Join(filepath.Dir(appPath), "public")
		if fi, _ := os.Lstat(pubDir); fi != nil && fi.IsDir() {
			ns.Bind("/", vfs.OS(pubDir), "/", vfs.BindAfter)
		}
	}

	mux.Handle("/", gwmux)
	mux.Handle("/static/", http.StripPrefix("/static", http.FileServer(httpfs.New(ns))))

	// grpc
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO(tamird): point to merged gRPC code rather than a PR.
		// This is a partial recreation of gRPC's internal checks
		// https://github.com/grpc/grpc-go/issues/555#issuecomment-443293451
		// https://github.com/grpc/grpc-go/pull/514/files#diff-95e9a25b738459a2d3030e1e6fa2a718R61
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			mux.ServeHTTP(w, r)
		}
	})
}
