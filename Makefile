# Copyright 2018 The OpenPitrix Authors. All rights reserved.
# Use of this source code is governed by a Apache license
# that can be found in the LICENSE file.

DEFAULT_HOST:=localhost:9115

SERVER_HOST := $(if ${OPENPITRIX_IAM_HOST},${OPENPITRIX_IAM_HOST},${DEFAULT_HOST})

default:
	go fmt ./...
	go vet ./...
	@echo ${SERVER_HOST}

im-server:
	go fmt ./...
	go vet ./...
	go run ./cmd/im/main.go

am-server:
	go fmt ./...
	go vet ./...
	go run ./cmd/am/main.go

mysql-up:
	docker run --rm --name mysql-dev -p 3306:3306 -e MYSQL_ROOT_PASSWORD=password -d mysql:5.7

mysql-down:
	docker stop --name mysql-dev

docker-build:
	docker build -t openpitrix/iam:latest -f ./Dockerfile .
	docker images openpitrix/iam:latest

docker-run-macos:
	OPENPITRIX_IAM_DB_HOST=docker.for.mac.localhost \
		docker run --rm -it -p 9115:9115 openpitrix/iam

docker-run-linux:
	OPENPITRIX_IAM_DB_HOST=172.17.0.1 \
		docker run --rm -it -p 9115:9115 openpitrix/iam

docker-run-windows:
	OPENPITRIX_IAM_DB_HOST=docker.for.win.localhost \
		docker run --rm -it -p 9115:9115 openpitrix/iam

docker-run:
	docker run --rm -it -p 9115:9115 -v `pwd`:/root \
		openpitrix/iam iam -config=/root/config.json

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
	grpcurl -plaintext ${SERVER_HOST} list openpitrix.iam.im.AccountManager

list-group:
	grpcurl -plaintext ${SERVER_HOST} openpitrix.iam.im.AccountManager/ListGroups
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
