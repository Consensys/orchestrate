GOFILES := $(shell find . -name '*.go' | egrep -v "^\./\.go" | grep -v _test.go)
PACKAGES ?= $(shell go list ./...)
BOILERPLATE_REPOSITORY=git@gitlab.com:ConsenSys/client/fr/core-stack/boilerplate-worker.git

UNAME_S := $(shell uname -s)
ifeq ($(UNAME_S),Linux)
	OPEN = xdg-open
endif
ifeq ($(UNAME_S),Darwin)
	OPEN = open
endif

.PHONY: all run-coverage coverage fmt fmt-check vet lint misspell-check misspell race tools help

# Testing
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

tidy: mod-tidy lint-fix

gocache: ## Create gocache folder for building image locally
	mkdir .gocache

# Tools
tools: ## Install test tools
	@GO111MODULE=off go get -u github.com/client9/misspell/cmd/misspell
	@GO111MODULE=off go get -u github.com/golangci/golangci-lint/cmd/golangci-lint

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

merge-boilerplate:
	@git remote add boilerplate $(BOILERPLATE_REPOSITORY) || true
	@git fetch boilerplate master
	@git merge boilerplate/master
