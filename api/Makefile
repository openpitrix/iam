# Copyright 2018 The OpenPitrix Authors. All rights reserved.
# Use of this source code is governed by a Apache license
# that can be found in the LICENSE file.

PWD:=$(shell pwd)
GOPATH:=$(shell go env GOPATH)

IM_PBFILES+=./openpitrix/iam/im/im.proto
IM_PBFILES+=./openpitrix/iam/im/filter.proto

AM_PBFILES+=./openpitrix/iam/am/am.proto

PROTOC_FLAGS:=-I. -Ithird_party/googleapis -Ithird_party

generate: Makefile
	@rm -rf ./_output/*
	@mkdir -p ./_output

	@rm -rf ../pkg/pb/*
	@mkdir -p ../pkg/pb
	@mkdir -p ../pkg/pb/am
	@mkdir -p ../pkg/pb/im

	make gen_im
	make gen_am

gen_im:
	# grpc service
	protoc $(PROTOC_FLAGS) --go_out=plugins=grpc:./_output ${IM_PBFILES}
	cp ./_output/openpitrix.io/iam/pkg/pb/im/*.go ../pkg/pb/im/

	# gateway
	protoc $(PROTOC_FLAGS) --grpc-gateway_out=logtostderr=true,allow_delete_body=true:./_output ${IM_PBFILES}
	cp ./_output/openpitrix/iam/im/*.go ../pkg/pb/im/

	# swagger
	protoc $(PROTOC_FLAGS) --swagger_out=logtostderr=true,allow_delete_body=true:./_output ${IM_PBFILES}
	cp ./_output/openpitrix/iam/im/*.json ../pkg/pb/im/

gen_am:
	# grpc service
	protoc $(PROTOC_FLAGS) --go_out=plugins=grpc:./_output ${AM_PBFILES}
	cp ./_output/openpitrix.io/iam/pkg/pb/am/*.go ../pkg/pb/am/

	# gateway
	protoc $(PROTOC_FLAGS) --grpc-gateway_out=logtostderr=true,allow_delete_body=true:./_output ${AM_PBFILES}
	cp ./_output/openpitrix/iam/am/*.go ../pkg/pb/am/

	# swagger
	protoc $(PROTOC_FLAGS) --swagger_out=logtostderr=true,allow_delete_body=true:./_output ${AM_PBFILES}
	cp ./_output/openpitrix/iam/am/*.json ../pkg/pb/am/

clean:
	rm -rf $(GOPATH)/src/openpitrix.io/iam/pkg/pb
