GOFILES := $(shell find . -name '*.go' | grep -v _test.go)
PACKAGES ?= $(shell go list ./...)

ROOTDIR=$(dir $(abspath $(lastword $(MAKEFILE_LIST))))
PB_FILES=$(shell find . -path '*.pb.go' | grep -v "vendor")
PROTO_FILES=$(shell find . -path '*.proto' | grep -v "vendor")
GRPC_GATEWAY=github.com/grpc-ecosystem/grpc-gateway
BINARIES=$(addprefix bin/,$(COMMANDS))
COMMANDS=protoc-gen-gogotodo
DESTDIR=/usr/local
GENERATE_PROTO=scripts/generate-proto.sh

.PHONY: all protobuf run-coverage coverage fmt fmt-check vet lint misspell-check misspell race tools help

run-coverage: ## Generate global code coverage report
	@sh scripts/coverage.sh $(PACKAGES)

tidy: fmt vet lint misspell mod-tidy ineffassign

coverage: run-coverage
	@xdg-open coverage.html

fmt: ## Formart source
	@gofmt -w $(GOFILES)

fmt-check: ## Check source format
	@gofmt -l ${GOFILES}

vet:
	@go vet $(PACKAGES)

lint: ## Lint code scripts
	@golint -set_exit_status ${PACKAGES}

misspell-check: ## Test misspells
	@misspell -error $(GOFILES)

misspell: ## Correct misspells
	@misspell -w $(GOFILES)

race: ## Run data race detector
	@go test -race -short ${PACKAGES}

mod-tidy:
	@go mod tidy

ineffassign:
	@ineffassign .

tools: ## Install test tools
	@GO111MODULE=off go get golang.org/x/lint/golint
	@GO111MODULE=off go get github.com/client9/misspell/cmd/misspell
	@GO111MODULE=off go get github.com/gordonklaus/ineffassign

setup-proto: ## Install protobuf utilities
	@GO111MODULE=off go get -u github.com/stevvooe/protobuild
	@GO111MODULE=off go get -d $(GRPC_GATEWAY)/...
	@cd $(GOPATH)/src/$(GRPC_GATEWAY)/protoc-gen-grpc-gateway && go install
	@cd $(GOPATH)/src/$(GRPC_GATEWAY)/protoc-gen-swagger && go install

protobuf: ## Generate protobuf stubs
	@sh scripts/generate-proto.sh

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
