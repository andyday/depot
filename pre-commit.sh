#!/bin/bash

set -e

go mod tidy -v
go generate ./...
go build ./...
golangci-lint run
go test ./...
go test ./... -covermode=atomic -coverprofile=$(COVERAGE_OUT) -coverpkg=./...
git add -u