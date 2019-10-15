GOFILES := $(shell find . -name '*.go' | egrep -v "^\./\.go" | grep -v _test.go)
PACKAGES ?= $(shell go list ./...)
CMD_RUN = tx-crafter tx-nonce tx-signer tx-sender tx-listener tx-decoder contract-registry envelope-store
CMD_MIGRATE = envelope-store

UNAME_S := $(shell uname -s)
ifeq ($(UNAME_S),Linux)
	OPEN = xdg-open
endif
ifeq ($(UNAME_S),Darwin)
	OPEN = open
endif

.PHONY: all run-coverage coverage fmt fmt-check vet lint misspell-check misspell race tools help

# Linters
run-coverage: ## Generate global code coverage report
	@sh scripts/coverage.sh $(PACKAGES)

gobuild:
	@go build -o ./build/corestack .

docker-build:
	@docker-compose build

coverage:
	@docker-compose -f e2e/docker-compose.yml up -d postgres
	@sh scripts/coverage.sh $(PACKAGES)
	@docker-compose -f e2e/docker-compose.yml stop postgres
	$(OPEN) build/coverage/coverage.html

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

clean: mod-tidy lint-fix protobuf

gocache:
	mkdir .gocache

generate-mocks:
	mockgen -destination=mocks/mock_client.go -package=mocks \
	gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/ethereum/rpc Client

	mockgen -destination=mocks/mock_enclave_endpoint.go -package=mocks \
	gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/ethereum/tessera EnclaveEndpoint

# Tools
tools: ## Install test tools
	@GO111MODULE=off go get -u github.com/client9/misspell/cmd/misspell
	@GO111MODULE=off go get -u github.com/golangci/golangci-lint/cmd/golangci-lint
	@GO111MODULE=off go get -u github.com/DATA-DOG/godog/cmd/godog
	@GO111MODULE=off go get -u github.com/golang/mock/gomock

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

protobuf: ## Generate protobuf stubs
	@docker-compose -f scripts/docker-compose.yml up

report:
	@docker-compose -f report/docker-compose.yml up
	$(OPEN) report/output/report.html

gen-help: gobuild
	@mkdir -p build/cmd
	@for cmd in $(CMD_RUN); do \
		./build/corestack help $$cmd run | tail -n +3 > build/cmd/$$cmd-run.md; \
	done
	@for cmd in $(CMD_MIGRATE); do \
		./build/corestack help $$cmd migrate | tail -n +3 > build/cmd/$$cmd-migrate.md; \
	done

gen-help-docker: docker-build
	@mkdir -p build/cmd
	@for cmd in $(CMD); do \
		docker-compose run worker help $$cmd run | tail -n +3 | head -n -1 > build/cmd/$$cmd-run.md; \
	done
	@for cmd in $(CMD); do \
		docker-compose run worker help $$cmd migrate | tail -n +3 | head -n -1 > build/cmd/$$cmd-migrate.md; \
	done
