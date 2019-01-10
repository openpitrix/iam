# Copyright 2018 The OpenPitrix Authors. All rights reserved.
# Use of this source code is governed by a Apache license
# that can be found in the LICENSE file.

DEFAULT_HOST:=localhost:9115

SERVER_HOST := $(if ${OPENPITRIX_IAM_HOST},${OPENPITRIX_IAM_HOST},${DEFAULT_HOST})

default:
	go fmt ./...
	go vet ./...
	@echo ${SERVER_HOST}

server:
	go fmt ./...
	go vet ./...
	go run main.go

docker-run:
	docker run --rm -it -p 9115:9115 -v `pwd`:/root \
		openpitrix/iam:v0.0.3-dev iam \
		-config=/root/config.json

info:
	curl ${SERVER_HOST}/hello
	@echo

	curl ${SERVER_HOST}/v1.1/version:iam
	@echo
	@echo

	grpcurl -plaintext ${SERVER_HOST} openpitrix.iam.IAMManager/GetVersion
	@echo

swagger:
	curl ${SERVER_HOST}/static/swagger/iam.swagger.json | jq .

list-method:
	grpcurl -plaintext ${SERVER_HOST} list
	grpcurl -plaintext ${SERVER_HOST} list openpitrix.iam.IAMManager

list-group:
	grpcurl -plaintext ${SERVER_HOST} openpitrix.iam.IAMManager/DescribeGroups
	@echo

	curl ${SERVER_HOST}/v1.1/groups | jq .
	@echo
	@echo


test:
	make generate
	cd ./api && make

	go fmt ./...
	go vet ./...
	go test ./...

dev:
	git describe --tags --always > ./_version
	git describe --exact-match || echo latest > ./_version

docker:
	docker build -t openpitrix/iam-dev -f ./Dockerfile .
	docker images openpitrix/iam-dev

generate:
	cd api && make
	go generate ./...

tools:
	# 1. install protoc from https://github.com/protocolbuffers/protobuf/releases
	# 2. install Go1.11+

	go get github.com/golang/protobuf/protoc-gen-go@v1.2

clean:
	cd api && make clean
