# Copyright 2018 The OpenPitrix Authors. All rights reserved.
# Use of this source code is governed by a Apache license
# that can be found in the LICENSE file.

# OPENPITRIX_IAM_HOST=localhost:9115 make

DEFAULT_HOST:=localhost:9115

SERVER_HOST := $(if ${OPENPITRIX_IAM_HOST},${OPENPITRIX_IAM_HOST},${DEFAULT_HOST})

default:
	@echo ${SERVER_HOST}

server:
	go run main.go

info:
	curl ${SERVER_HOST}/hello
	@echo

	curl ${SERVER_HOST}/static/version
	@echo
	@echo

	grpcurl -plaintext ${SERVER_HOST} openpitrix.iam.IAMManager/GetVersion
	@echo

swagger:
	-curl ${SERVER_HOST}/static/spec/iam.swagger.json

list-method:
	grpcurl -plaintext ${SERVER_HOST} list
	grpcurl -plaintext ${SERVER_HOST} list openpitrix.iam.IAMManager

list-group:
	# TODO

clean:
