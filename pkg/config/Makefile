# Copyright 2019 The OpenPitrix Authors. All rights reserved.
# Use of this source code is governed by a Apache license
# that can be found in the LICENSE file.

default:
	@make tag
	@make file
	@make env
	@make env-file

tag:
	@echo from tag
	@go run hello.go

file:
	@echo from file
	@go run hello.go -conf=_config.json

env:
	@echo from env
	@\
	OPENPITRIX_IAM_DB_TYPE=db-type-9527 \
	OPENPITRIX_IAM_DB_HOST=db-host-9527 \
	OPENPITRIX_IAM_DB_PORT=9528 \
	OPENPITRIX_IAM_DB_USER=chai \
	OPENPITRIX_IAM_DB_PASSWORD=chai-9527 \
	OPENPITRIX_IAM_DB_DATABASE=chai-iam-db \
	\
	OPENPITRIX_IAM_HOST=iam-9527 \
	OPENPITRIX_IAM_PORT=9527 \
	\
	OPENPITRIX_IAM_TLS_ENABLED=true \
	OPENPITRIX_IAM_TLS_CERT_FILE=9527.cert \
	OPENPITRIX_IAM_TLS_KEY_FILE=9527.key \
	\
	OPENPITRIX_IAM_LOG_LEVEL=DEBUG \
	\
	go run hello.go

env-file:
	@echo from env and file
	@\
	OPENPITRIX_IAM_DB_TYPE=db-type-9527 \
	OPENPITRIX_IAM_DB_HOST=db-host-9527 \
	OPENPITRIX_IAM_DB_PORT=9528 \
	OPENPITRIX_IAM_DB_USER=chai \
	OPENPITRIX_IAM_DB_PASSWORD=chai-9527 \
	OPENPITRIX_IAM_DB_DATABASE=chai-iam-db \
	\
	go run hello.go -conf=_config.json

clean:
