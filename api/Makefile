# Copyright 2019 The OpenPitrix Authors. All rights reserved.
# Use of this source code is governed by a Apache license
# that can be found in the LICENSE file.

PWD:=$(shell pwd)
GOPATH:=$(shell go env GOPATH)

PROTO_FILES=$(sort $(wildcard ./*.proto))
PROTOC_INC_PATH=/usr/local/include
GOOGLEAPIS_PATH=third_party/googleapis

PROTOC_FLAGS:=-I. -I$(GOOGLEAPIS_PATH) -I$(PROTOC_INC_PATH)

generate: Makefile
	@rm -rf ../pkg/pb/*
	@mkdir -p ../pkg/pb

	# grpc service
	protoc $(PROTOC_FLAGS) --go_out=plugins=grpc:../pkg/pb ${PROTO_FILES}
	# need to exec chown finally
	chown -R $${USER_ID}:$${GROUP_ID} ../pkg
	# format generated files
	goimports -l -w -e -local=openpitrix ../pkg/pb/

clean:
	rm -rf $(GOPATH)/src/openpitrix.io/iam/pkg/pb
