SHELL=/bin/bash

OS ?= $(shell go env GOOS)
ARCH ?= $(shell go env GOARCH)

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
tools: install-lint install-mockery  ##@Desktop Install all prerequisite tools.

.PHONY: install-lint
install-lint:  ##@CI Install the linting tool.
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -d v1.57.2
	bin/golangci-lint --version

.PHONY: install-mockery
install-mockery:  ##@CI Install mockery which is used to generate mocks for unit testing.
	go install github.com/vektra/mockery/v2@v2.42.3

.PHONY: hooks
hooks:  ##@Desktop Install the git hooks.
	ln -sf ../../scripts/pre-commit.sh .git/hooks/pre-commit

.PHONY: setup
setup: tools hooks  ##@Desktop Run the setup for this project.

.PHONY: tidy
tidy:  ##@Build Run go tidy and vendor.
	@echo "*** Tidy Dependencies..."
	@go mod tidy -v
	@#go mod vendor -v

.PHONY: generate
generate:  ##@Build Generate all generated files using go generate.
	@echo "*** Generate API and Mocks..."
	@go generate ./...

.PHONY: test
test:  ##@Build Run the unit tests.
	@echo "*** Unit Testing..."
	#go test -mod vendor $(TEST_PKGS)
	go test ./...

.PHONY: lint
lint:  ##@Build Run the lint on all go code.
	@echo "*** Linting..."
	@#bin/golangci-lint run -v --modules-download-mode vendor
	golangci-lint run

.PHONY: build
build:
	@echo "*** Building..."
	@go build ./...