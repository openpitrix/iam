// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package am

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
)

func (p *Server) ListenAndServe(addr string, opt ...grpc.ServerOption) error {
	if p.webServer != nil {
		return fmt.Errorf("web server is running")
	}

	// https://github.com/grpc/grpc-go/issues/555#issuecomment-443293451
	h2Handler := h2c.NewHandler(p.mainHandler(), &http2.Server{})

	p.grpcServer = grpc.NewServer(opt...)
	p.webServer = &http.Server{Addr: addr, Handler: h2Handler}

	return p.webServer.ListenAndServe()
}

func (p *Server) ListenAndServeTLS(addr, certFile, keyFile string, opt ...grpc.ServerOption) error {
	if p.webServer != nil {
		return fmt.Errorf("web server is running")
	}

	p.grpcServer = grpc.NewServer(opt...)
	p.webServer = &http.Server{Addr: addr, Handler: p.mainHandler()}

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

func (p *Server) mainHandler() http.Handler {
	mux := http.NewServeMux()

	// mux.HandleFunc("/", ...)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO(tamird): point to merged gRPC code rather than a PR.
		// This is a partial recreation of gRPC's internal checks
		// https://github.com/grpc/grpc-go/issues/555#issuecomment-443293451
		// https://github.com/grpc/grpc-go/pull/514/files#diff-95e9a25b738459a2d3030e1e6fa2a718R61
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			p.grpcServer.ServeHTTP(w, r)
		} else {
			mux.ServeHTTP(w, r)
		}
	})
}
