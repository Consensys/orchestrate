GOFILES := $(shell find . -name '*.go' | egrep -v "^\./\.go" | grep -v _test.go)
PACKAGES ?= $(shell go list ./... | go list ./... | grep -Fv -e e2e -e examples )
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
	@GOOS=linux GOARCH=amd64 go build -i -o ./build/bin/corestack

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
	@./build/bin/corestack help tx-crafter | grep -A 9999 "Global Flags:" | head -n -2 > build/cmd/global.txt
	@for cmd in $(CMD_RUN); do \
		./build/bin/corestack help $$cmd run | grep -B 9999 "Global Flags:" | tail -n +3 | head -n -2 > build/cmd/$$cmd-run.txt; \
	done
	@for cmd in $(CMD_MIGRATE); do \
		./build/bin/corestack help $$cmd migrate | grep -B 9999 "Global Flags:" | tail -n +3 | head -n -2 > build/cmd/$$cmd-migrate.txt; \
	done

gen-help-docker: docker-build
	@mkdir -p build/cmd
	@docker-compose run worker help tx-crafter | grep -A 9999 "Global Flags:" | head -n -3 > build/cmd/global.txt
	@for cmd in $(CMD_RUN); do \
		docker-compose run worker help $$cmd run | grep -B 9999 "Global Flags:" | tail -n +3 | head -n -3 > build/cmd/$$cmd-run.txt; \
	done
	@for cmd in $(CMD_MIGRATE); do \
		docker-compose run worker help $$cmd migrate | grep -B 9999 "Global Flags:" | tail -n +3 | head -n -3 > build/cmd/$$cmd-migrate.txt; \
	done

gobuild-e2e:
	@GOOS=linux GOARCH=amd64 go build -i -o ./build/bin/e2e ./tests/cmd 

corestack: gobuild
	@docker-compose -f docker-compose.dev.yml up -d $(CMD_RUN)

stop-corestack:
	@docker-compose -f docker-compose.dev.yml stop $(CMD_RUN)

deps:
	@docker-compose -f scripts/deps/docker-compose.yml up -d

quorum:
	@docker-compose -f scripts/deps/docker-compose.quorum.yml up -d

stop-quorum:
	@docker-compose -f scripts/deps/docker-compose.quorum.yml stop

e2e: gobuild-e2e
	@docker-compose -f docker-compose.dev.yml up e2e
	@docker-compose -f scripts/report/docker-compose.yml up
