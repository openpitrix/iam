// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package web

import (
	"context"
	"runtime/debug"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/grpc-ecosystem/go-grpc-middleware/validator"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/status"

	"openpitrix.io/iam/pkg/pb/am"
	"openpitrix.io/iam/pkg/pb/im"
	"openpitrix.io/logger"
)

var (
	_ GrpcServer = (*amGrpcServer)(nil)
	_ GrpcServer = (*imGrpcServer)(nil)
)

func WithAccountManager(s pbim.AccountManagerServer) GrpcServer {
	return &imGrpcServer{s}
}
func WithAccessManager(s pbam.AccessManagerServer) GrpcServer {
	return &amGrpcServer{s}
}

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

type amGrpcServer struct {
	s pbam.AccessManagerServer
}

func (p *amGrpcServer) RegisterGrpcServer(s *grpc.Server) {
	pbam.RegisterAccessManagerServer(s, p.s)
}
func (p *amGrpcServer) RegisterGrpcHandlerFromEndpoint(
	ctx context.Context, mux *runtime.ServeMux,
	endpoint string, opts []grpc.DialOption,
) error {
	return pbam.RegisterAccessManagerHandlerFromEndpoint(
		ctx, mux, endpoint, opts,
	)
}

type imGrpcServer struct {
	s pbim.AccountManagerServer
}

func (p *imGrpcServer) RegisterGrpcServer(s *grpc.Server) {
	pbim.RegisterAccountManagerServer(s, p.s)
}
func (p *imGrpcServer) RegisterGrpcHandlerFromEndpoint(
	ctx context.Context, mux *runtime.ServeMux,
	endpoint string, opts []grpc.DialOption,
) error {
	return pbim.RegisterAccountManagerHandlerFromEndpoint(
		ctx, mux, endpoint, opts,
	)
}
