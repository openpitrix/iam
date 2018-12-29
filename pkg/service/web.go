// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package service

import (
	"context"
	"fmt"
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
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"golang.org/x/tools/godoc/vfs"
	"golang.org/x/tools/godoc/vfs/httpfs"
	"golang.org/x/tools/godoc/vfs/mapfs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	"openpitrix.io/iam/pkg/pb"
	staticSwaggerUI "openpitrix.io/iam/pkg/service/swagger-ui"
	"openpitrix.io/logger"
)

func (p *Server) makeDefaultGrpcServerOptions() []grpc.ServerOption {
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

func (p *Server) ListenAndServe(addr string) error {
	if p.webServer != nil {
		return fmt.Errorf("web server is running")
	}

	var opt = p.makeDefaultGrpcServerOptions()

	// https://github.com/grpc/grpc-go/issues/555#issuecomment-443293451
	h2Handler := h2c.NewHandler(p.mainHandler(addr), &http2.Server{})

	p.grpcServer = grpc.NewServer(opt...)
	reflection.Register(p.grpcServer)
	pb.RegisterIAMManagerServer(p.grpcServer, p)

	p.webServer = &http.Server{Addr: addr, Handler: h2Handler}

	return p.webServer.ListenAndServe()
}

func (p *Server) ListenAndServeTLS(addr, certFile, keyFile string) error {
	if p.webServer != nil {
		return fmt.Errorf("web server is running")
	}

	var opt = p.makeDefaultGrpcServerOptions()

	p.grpcServer = grpc.NewServer(opt...)
	reflection.Register(p.grpcServer)
	pb.RegisterIAMManagerServer(p.grpcServer, p)

	p.webServer = &http.Server{Addr: addr, Handler: p.mainHandler(addr)}

	return p.webServer.ListenAndServeTLS(certFile, keyFile)
}

func (p *Server) Shutdown() error {
	if p.webServer == nil {
		return nil
	}

	p.grpcServer.Stop()

	err := p.webServer.Shutdown(context.Background())
	p.webServer = nil
	if err != nil {
		return err
	}

	return nil
}

func (p *Server) mainHandler(addr string) http.Handler {
	var gwmux = runtime.NewServeMux()
	var opts = []grpc.DialOption{grpc.WithInsecure()}
	var err error

	err = pb.RegisterIAMManagerHandlerFromEndpoint(context.Background(),
		gwmux, addr, opts,
	)
	if err != nil {
		logger.Warnf(nil, "%v", err)
	}

	mux := http.NewServeMux()

	// GET /readme.md
	mux.HandleFunc("/readme.md", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, web_readme_md())
	})

	// just for test
	// GET /hello
	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "hello", time.Now())
	})

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
			p.grpcServer.ServeHTTP(w, r)
		} else {
			if r.URL.Path == "/" {
				fmt.Fprintln(w, web_homepage())
			} else {
				mux.ServeHTTP(w, r)
			}
		}
	})
}
