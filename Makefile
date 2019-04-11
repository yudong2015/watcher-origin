# Copyright 2018 The OpenPitrix Authors. All rights reserved.
# Use of this source code is governed by a Apache license
# that can be found in the LICENSE file.

TARG.Name:=watcher
TRAG.Gopkg:=openpitrix.io/$(TRAG.Name)
TRAG.Version:=$(TRAG.Gopkg)/pkg/version

GO_FMT:=goimports -l -w -e -local=openpitrix -srcdir=/go/src/$(TRAG.Gopkg)
GO_RACE:=go build -race
GO_VET:=go vet
GO_FILES:=./pkg
GO_PATH_FILES:=./pkg/...

DOCKER_TAGS=latest
BUILDER_IMAGE=openpitrix/openpitrix-builder:release-v0.2.3
RUN_IN_DOCKER:=docker run -it -v `pwd`:/go/src/$(TRAG.Gopkg) -v `pwd`/tmp/cache:/root/.cache/go-build  -w /go/src/$(TRAG.Gopkg) -e GOBIN=/go/src/$(TRAG.Gopkg)/tmp/bin -e USER_ID=`id -u` -e GROUP_ID=`id -g` $(BUILDER_IMAGE)

define get_diff_files
    $(eval DIFF_FILES=$(shell git diff --name-only --diff-filter=ad | grep -e "^(cmd|pkg)/.+\.go"))
endef
# Get project build flags
define get_build_flags
    $(eval SHORT_VERSION=$(shell git describe --tags --always --dirty="-dev"))
    $(eval SHA1_VERSION=$(shell git show --quiet --pretty=format:%H))
	$(eval DATE=$(shell date +'%Y-%m-%dT%H:%M:%S'))
	$(eval BUILD_FLAG= -X $(TRAG.Version).ShortVersion="$(SHORT_VERSION)" \
		-X $(TRAG.Version).GitSha1Version="$(SHA1_VERSION)" \
		-X $(TRAG.Version).BuildDate="$(DATE)")
endef

CMD?=...
comma:= ,
empty:=
space:= $(empty) $(empty)
CMDS=$(subst $(comma),$(space),$(CMD))

.PHONY: fmt-all
fmt-all: ## Format all code
	$(RUN_IN_DOCKER) $(GO_FMT) $(GO_FILES)
	@echo "fmt done"

.PHONY: fmt-check
fmt-check: fmt-all ## Check whether all files be formatted
	$(call get_diff_files)
	$(if $(DIFF_FILES), \
		exit 2 \
	)

.PHONY: check
check: ## go vet and race
	env GO111MODULE=on $(GO_RACE) $(GO_PATH_FILES)
	env GO111MODULE=on $(GO_VET) $(GO_PATH_FILES)

.PHONY: build
build: fmt-all ## Build notification image
	mkdir -p ./tmp/bin
	$(call get_build_flags)
	$(RUN_IN_DOCKER) env GO111MODULE=on time go install -tags netgo -v -ldflags '$(BUILD_FLAG)' pkg/watch/main.go
	docker build -t $(TARG.Name) -f ./Dockerfile.dev ./tmp/bin
	docker image prune -f 1>/dev/null 2>&1
	@echo "build done"

.PHONY: compose-up
compose-up:
	docker-compose up -d
	@echo "compose-up done"

build-image-%: ## build docker image
	@if [ "$*" = "latest" ];then \
	docker build -t openpitrix/watcher:latest .; \
	elif [ "`echo "$*" | grep -E "^v[0-9]+\.[0-9]+\.[0-9]+"`" != "" ];then \
	docker build -t openpitrix/watcher:$* .; \
	fi

push-image-%: ## push docker image
	@if [ "$*" = "latest" ];then \
	docker push openpitrix/watcher:latest; \
	elif [ "`echo "$*" | grep -E "^v[0-9]+\.[0-9]+\.[0-9]+"`" != "" ];then \
	docker push openpitrix/watcher:$*; \
	fi

.PHONY: test
test: ## Run all tests
	make unit-test
	@echo "test done"

.PHONY: unit-test
unit-test: ## Run unit tests
	$(DB_TEST) $(ETCD_TEST) go test -a -tags="etcd db" ./...
	@echo "unit-test done"



