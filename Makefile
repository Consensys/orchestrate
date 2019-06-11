GOFILES := $(shell find . -name '*.go' | egrep -v "^\./\.go" | grep -v _test.go)
PACKAGES ?= $(shell go list ./...)

UNAME_S := $(shell uname -s)
ifeq ($(UNAME_S),Linux)
	OPEN = xdg-open
endif
ifeq ($(UNAME_S),Darwin)
	OPEN = open
endif


GRPC_GATEWAY=github.com/grpc-ecosystem/grpc-gateway

.PHONY: all run-coverage coverage fmt fmt-check vet lint misspell-check misspell race tools help

# Linters
run-coverage: ## Generate global code coverage report
	@sh scripts/coverage.sh $(PACKAGES)

coverage: run-coverage ## Generate and open coverage report
	$(OPEN) coverage.html

race: ## Run data race detector
	@go test -race -short ${PACKAGES}

mod-tidy: 
	@go mod tidy

lint-fix:
	@misspell -w $(GOFILES)
	@golangci-lint run --fix

lint:
	@misspell -error $(GOFILES)
	@golangci-lint run

lint-ci:
	@misspell -error $(GOFILES)
	@GOGC=20 golangci-lint -v run

tidy: mod-tidy lint-fix protobuf


# Tools
tools: ## Install test tools
	@GO111MODULE=off go get -u github.com/client9/misspell/cmd/misspell
	@GO111MODULE=off go get -u github.com/golangci/golangci-lint/cmd/golangci-lint

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

setup-proto: ## Install protobuf utilities
    # Install protoc-gen-go
	@GO111MODULE=off go get -d -u github.com/golang/protobuf/protoc-gen-go
	@GO111MODULE=off go install github.com/golang/protobuf/protoc-gen-go

	@GO111MODULE=off go get -u github.com/stevvooe/protobuild
	@GO111MODULE=off go get -d $(GRPC_GATEWAY)/...
	@cd $(GOPATH)/src/$(GRPC_GATEWAY)/protoc-gen-grpc-gateway && go install
	@cd $(GOPATH)/src/$(GRPC_GATEWAY)/protoc-gen-swagger && go install

protobuf: ## Generate protobuf stubs
	@sh scripts/generate-proto.sh
