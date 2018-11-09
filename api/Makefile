# Copyright 2018 The OpenPitrix Authors. All rights reserved.
# Use of this source code is governed by a Apache license
# that can be found in the LICENSE file.

PWD:=$(shell pwd)
GOPATH:=$(shell go env GOPATH)

IAM_AM_PBFILES=$(sort $(wildcard iam/am/*.proto))
IAM_IM_PBFILES=$(sort $(wildcard iam/im/*.proto))

generate: Makefile
	@mkdir -p ../pkg/pb
	@rm -rf ../pkg/pb/*

	# grpc service
	protoc $(PROTOC_FLAGS) --go_out=plugins=grpc:$(GOPATH)/src ${IAM_AM_PBFILES}
	protoc $(PROTOC_FLAGS) --go_out=plugins=grpc:$(GOPATH)/src ${IAM_IM_PBFILES}

clean:
	rm -rf $(GOPATH)/src/openpitrix.io/iam/pkg/pb
