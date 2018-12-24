# Copyright 2018 The OpenPitrix Authors. All rights reserved.
# Use of this source code is governed by a Apache license
# that can be found in the LICENSE file.

PWD:=$(shell pwd)
GOPATH:=$(shell go env GOPATH)

PBFILES=$(sort $(wildcard *.proto))

PROTOC_FLAGS:=-I. -Ithird_party/googleapis -Ithird_party

generate: Makefile
	@mkdir -p ../pkg/pb
	@rm -rf ../pkg/pb/*

	# grpc service
	protoc $(PROTOC_FLAGS) --go_out=plugins=grpc:$(GOPATH)/src ${PBFILES}

	# govalidators
	protoc $(PROTOC_FLAGS) --govalidators_out=$(GOPATH)/src ${PBFILES}

	# gateway
	protoc $(PROTOC_FLAGS) --grpc-gateway_out=../pkg/pb ${PBFILES}

clean:
	rm -rf $(GOPATH)/src/openpitrix.io/iam/openpitrix/pkg/pb
