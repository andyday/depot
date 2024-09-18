SHELL=/bin/bash

GOBIN ?= $$(go env GOPATH)/bin
OS ?= $(shell go env GOOS)
ARCH ?= $(shell go env GOARCH)
TEST_PKGS := $(shell go list ./... | grep -v e2e)

HELP_FUN = \
    %help; while(<>){push@{$$help{$$2//'options'}},[$$1,$$3] \
    if/^([\w-_]+)\s*:.*\#\#(?:@(\w+))?\s(.*)$$/}; \
    print"$$_:\n", map"  \033\[36m$$_->[0]".(" "x(30-length($$_->[0])))."\033[0m$$_->[1]\n",\
    @{$$help{$$_}},"\n" for keys %help; \

help: ##@Miscellaneous Show this help.
	@echo -e "Usage: make [target] ...\n"
	@perl -e '$(HELP_FUN)' $(MAKEFILE_LIST)

.PHONY: all
all: help

.PHONY: pre-commit
pre-commit: tidy generate test lint  ##@Desktop Run pre-commit checks.

.PHONY: tools
tools: install-lint install-mockery install-coverage  ##@Desktop Install all prerequisite tools.

.PHONY: install-lint
install-lint:  ##@CI Install the linting tool.
	brew install golangci-lint
	golangci-lint --version

.PHONY: install-mockery
install-mockery:  ##@CI Install mockery which is used to generate mocks for unit testing.
	go install github.com/vektra/mockery/v2@v2.42.3

.PHONY: install-coverage
install-coverage:
	go install github.com/vladopajic/go-test-coverage/v2@latest

.PHONY: hooks
hooks:  ##@Desktop Install the git hooks.
	ln -sf ../../scripts/pre-commit.sh .git/hooks/pre-commit

.PHONY: setup
setup: tools hooks  ##@Desktop Run the setup for this project.

.PHONY: tidy
tidy:  ##@Build Run go tidy and vendor.
	@echo "*** Tidy Dependencies..."
	@go mod tidy -v

.PHONY: generate
generate:  ##@Build Generate all generated files using go generate.
	@echo "*** Generate API and Mocks..."
	@mockery
	@go generate ./...

.PHONY: test
test:  ##@Build Run the unit tests.
	@echo "*** Unit Testing..."
	go test ./...

.PHONY: lint
lint:  ##@Build Run the lint on all go code.
	@echo "*** Linting..."
	golangci-lint run

.PHONY: build
build:
	@echo "*** Building..."
	@go build ./...

.PHONY: cover
cover:      ##@Build Run unit tests with coverage.
	go test $(TEST_PKGS) -coverprofile=./cover.out -covermode=atomic -coverpkg=$(TEST_PKGS)

check-coverage: cover install-coverage
	go-test-coverage --config=./.coverage.yml

.PHONY: init-deploy
init-deploy:
	@echo "*** Init Terraform..."
	pushd e2e/terraform && terraform init && popd

.PHONY: deploy
deploy: init-deploy
	pushd e2e/terraform && \
	terraform apply --auto-approve -var="project=$(PROJECT)" && \
	popd

.PHONY: e2e
e2e:
	export FIRESTORE_PROJECT_ID=$(PROJECT) && 	go test -count=1 -v ./e2e
