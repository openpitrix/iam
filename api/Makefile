# Copyright 2018 The OpenPitrix Authors. All rights reserved.
# Use of this source code is governed by a Apache license
# that can be found in the LICENSE file.

PWD:=$(shell pwd)
GOPATH:=$(shell go env GOPATH)

PBFILES=$(sort $(wildcard *.proto))

PROTOC_FLAGS:=-I. -Ithird_party/googleapis -Ithird_party

generate: Makefile
	@rm -rf ./_output/*
	@mkdir -p ./_output

	@rm -rf ../pkg/pb/*
	@mkdir -p ../pkg/pb

	# grpc service
	protoc $(PROTOC_FLAGS) --go_out=plugins=grpc:./_output ${PBFILES}
	cp ./_output/openpitrix.io/iam/pkg/pb/*.go ../pkg/pb/

	# gateway
	protoc $(PROTOC_FLAGS) --grpc-gateway_out=logtostderr=true,allow_delete_body=true:../pkg/pb ${PBFILES}

	# swagger
	protoc $(PROTOC_FLAGS) --swagger_out=logtostderr=true,allow_delete_body=true:../pkg/service/swagger-ui ${PBFILES}
	make -C ../pkg/service/swagger-ui

clean:
	rm -rf $(GOPATH)/src/openpitrix.io/iam/pkg/pb
