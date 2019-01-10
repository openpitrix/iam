// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package web

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"golang.org/x/tools/godoc/vfs"
	"golang.org/x/tools/godoc/vfs/httpfs"
	"golang.org/x/tools/godoc/vfs/mapfs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	staticSwaggerUI "openpitrix.io/iam/pkg/service/swagger-ui"
	"openpitrix.io/logger"
)

type GrpcServer interface {
	RegisterGrpcServer(s *grpc.Server)
	RegisterGrpcHandlerFromEndpoint(
		ctx context.Context, mux *runtime.ServeMux,
		endpoint string, opts []grpc.DialOption,
	) error
}

func ListenAndServe(addr string,
	grpcServices []GrpcServer,
	defaultHandler http.Handler,
	opts ...grpc.ServerOption,
) error {
	if len(opts) == 0 {
		opts = MakeDefaultGrpcServerOptions()
	}

	grpcServer := grpc.NewServer(opts...)
	reflection.Register(grpcServer)
	for _, svc := range grpcServices {
		svc.RegisterGrpcServer(grpcServer)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "hello", time.Now()) // just for test
	})
	mux.HandleFunc("/index.html", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, homepage)
	})

	if defaultHandler == nil {
		defaultHandler = MainHandler(grpcServer, addr, grpcServices, mux)
	}

	// https://github.com/grpc/grpc-go/issues/555#issuecomment-443293451
	server := &http.Server{
		Addr:    addr,
		Handler: h2c.NewHandler(defaultHandler, &http2.Server{}),
	}

	return server.ListenAndServe()
}

func ListenAndServeTLS(addr, certFile, keyFile string,
	grpcServices []GrpcServer,
	defaultHandler http.Handler,
	opts ...grpc.ServerOption,
) error {
	if len(opts) == 0 {
		opts = MakeDefaultGrpcServerOptions()
	}

	grpcServer := grpc.NewServer(opts...)
	reflection.Register(grpcServer)
	for _, svc := range grpcServices {
		svc.RegisterGrpcServer(grpcServer)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "hello", time.Now()) // just for test
	})
	mux.HandleFunc("/index.html", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, homepage)
	})

	if defaultHandler == nil {
		defaultHandler = MainHandler(grpcServer, addr, grpcServices, mux)
	}

	server := &http.Server{
		Addr:    addr,
		Handler: defaultHandler,
	}

	return server.ListenAndServeTLS(certFile, keyFile)
}

func MainHandler(
	grpcServer *grpc.Server,
	grpcServerAddress string,
	grpcServices []GrpcServer,
	mux *http.ServeMux,
) http.Handler {
	var gwmux = runtime.NewServeMux()
	var opts = []grpc.DialOption{grpc.WithInsecure()}
	var err error

	for _, svc := range grpcServices {
		err = svc.RegisterGrpcHandlerFromEndpoint(
			context.Background(), gwmux, grpcServerAddress, opts,
		)
		if err != nil {
			logger.Warnf(nil, "%v", err)
		}
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
			if r.URL.Path == "/" {
				http.Redirect(w, r, "/index.html", http.StatusFound)
			} else {
				mux.ServeHTTP(w, r)
			}
		}
	})
}
