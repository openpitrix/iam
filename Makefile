# Copyright 2019 The OpenPitrix Authors. All rights reserved.
# Use of this source code is governed by a Apache license
# that can be found in the LICENSE file.

PWD:=$(shell pwd)

TARG.Name:=iam
TRAG.Gopkg:=openpitrix.io/iam
TRAG.Version:=$(TRAG.Gopkg)/pkg/version

GO_FMT:=goimports -l -w -e -local=openpitrix -srcdir=/go/src/$(TRAG.Gopkg)
GO_MOD_TIDY:=go mod tidy
GO_RACE:=go build -race
GO_VET:=go vet
GO_FILES:=./cmd ./pkg
GO_PATH_FILES:=./cmd/... ./pkg/...
COMPOSE_DB_CTRL=im-db-ctrl am-db-ctrl

BUILDER_IMAGE=openpitrix/openpitrix-builder:release-v0.2.3
RUN_IN_DOCKER:=docker run -it -v `pwd`:/go/src/$(TRAG.Gopkg) -v `pwd`/tmp/cache:/root/.cache/go-build  -w /go/src/$(TRAG.Gopkg) -e GOBIN=/go/src/$(TRAG.Gopkg)/tmp/bin -e USER_ID=`id -u` -e GROUP_ID=`id -g` $(BUILDER_IMAGE)

define get_diff_files
    $(eval DIFF_FILES=$(shell git diff --name-only --diff-filter=ad | grep -e "^(cmd|pkg)/.+\.go" -e "go.mod"))
endef

CMD?=...
comma:= ,
empty:=
space:= $(empty) $(empty)
CMDS=$(subst $(comma),$(space),$(CMD))

.PHONY: build-flyway
build-flyway: ## Build custom flyway image
	docker build -t $(TARG.Name):flyway -f ./pkg/db/Dockerfile ./pkg/db/

.PHONY: build
build: build-flyway ## Build all im images
	docker build -t $(TARG.Name) -f ./Dockerfile .
	docker image prune -f 1>/dev/null 2>&1
	@echo "build done"

.PHONY: test
test: ## Run all tests
	make unit-test
	make e2e-test
	@echo "test done"

.PHONY: unit-test
unit-test: ## Run unit tests
	env GO111MODULE=on go test -a -tags="unit" ./...
	@echo "unit-test done"

.PHONY: e2e-test
e2e-test: ## Run integration tests
	env GO111MODULE=on go test -a -tags="integration" ./test/e2e/...
	@echo "e2e-test done"

.PHONY: compose-migrate-db
compose-migrate-db: ## Migrate db in docker compose
	until docker-compose exec iam-db bash -c "echo 'SELECT VERSION();' | mysql -uroot -ppassword"; do echo "waiting for mysql"; sleep 2; done;
	docker-compose up $(COMPOSE_DB_CTRL)

.PHONY: compose-up
compose-up: ## Launch im in docker compose
	docker-compose up -d iam-db
	make compose-migrate-db
	docker-compose up -d
	@echo "compose-up done"

.PHONY: compose-update
compose-update: build compose-up ## Update service in docker compose
	@echo "compose-update done"

.PHONY: compose-down
compose-down: ## Shutdown docker compose
	docker-compose down
	@echo "compose-down done"

.PHONY: generate-in-local
generate-in-local: ## Generate code from protobuf file in local
	cd api && make

.PHONY: generate
generate: ## Generate code from protobuf file in docker
	$(RUN_IN_DOCKER) make generate-in-local
	@echo "generate done"

.PHONY: fmt-all
fmt-all: ## Format all code
	$(RUN_IN_DOCKER) $(GO_FMT) $(GO_FILES)
	@echo "fmt done"

.PHONY: tidy
tidy: ## Tidy go.mod
	env GO111MODULE=on $(GO_MOD_TIDY)
	@echo "go mod tidy done"

.PHONY: fmt-check
fmt-check: fmt-all tidy ## Check whether all files be formatted
	$(call get_diff_files)
	$(if $(DIFF_FILES), \
		exit 2 \
	)

.PHONY: check
check: ## go vet and race
	env GO111MODULE=on $(GO_RACE) $(GO_PATH_FILES)
	env GO111MODULE=on $(GO_VET) $(GO_PATH_FILES)

build-image-%: ## build docker image
	@if [ "$*" = "latest" ];then \
	docker build -t openpitrix/iam:latest .; \
	docker build -t openpitrix/iam:flyway -f ./pkg/db/Dockerfile ./pkg/db/; \
	elif [ "`echo "$*" | grep -E "^v[0-9]+\.[0-9]+\.[0-9]+"`" != "" ];then \
	docker build -t openpitrix/iam:$* .; \
	docker build -t openpitrix/iam:flyway-$* -f ./pkg/db/Dockerfile ./pkg/db/; \
	fi

push-image-%: ## push docker image
	@if [ "$*" = "latest" ];then \
	docker push openpitrix/iam:latest; \
	docker push openpitrix/iam:flyway; \
	elif [ "`echo "$*" | grep -E "^v[0-9]+\.[0-9]+\.[0-9]+"`" != "" ];then \
	docker push openpitrix/iam:$*; \
	docker push openpitrix/iam:flyway-$*; \
	fi
