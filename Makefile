GOFILES := $(shell find . -name '*.go' | egrep -v "^\./\.go" | grep -v _test.go)
PACKAGES ?= $(shell go list ./... | go list ./... | grep -Fv -e e2e -e examples )
CMD_RUN = tx-crafter tx-nonce tx-signer tx-sender tx-listener tx-decoder contract-registry envelope-store
CMD_MIGRATE = contract-registry envelope-store

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

run-e2e: gobuild-e2e
	@docker-compose up e2e
	@docker-compose -f scripts/report/docker-compose.yml up

e2e: run-e2e
	$(OPEN) build/report/report.html

clean: mod-tidy lint-fix protobuf

gocache:
	mkdir .gocache

generate-mocks:
	mockgen -destination=mocks/mock_client.go -package=mocks \
	gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/ethereum/rpc Client

	mockgen -destination=mocks/mock_enclave_endpoint.go -package=mocks \
	gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/ethereum/tessera EnclaveEndpoint

# Tools
tools: ## Install test tools
	@GO111MODULE=off go get -u github.com/client9/misspell/cmd/misspell
	@GO111MODULE=off go get -u github.com/golangci/golangci-lint/cmd/golangci-lint
	@GO111MODULE=off go get -u github.com/DATA-DOG/godog/cmd/godog
	@GO111MODULE=off go get -u github.com/golang/mock/gomock

# Help
help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

gen-help: gobuild
	@mkdir -p build/cmd
	@./build/bin/orchestrate help tx-crafter | grep -A 9999 "Global Flags:" | head -n -2 > build/cmd/global.txt
	@for cmd in $(CMD_RUN); do \
		./build/bin/orchestrate help $$cmd run | grep -B 9999 "Global Flags:" | tail -n +3 | head -n -2 > build/cmd/$$cmd-run.txt; \
	done
	@for cmd in $(CMD_MIGRATE); do \
		./build/bin/orchestrate help $$cmd migrate | grep -B 9999 "Global Flags:" | tail -n +3 | head -n -2 > build/cmd/$$cmd-migrate.txt; \
	done

gen-help-docker: docker-build
	@mkdir -p build/cmd
	@docker run orchestrate help tx-crafter | grep -A 9999 "Global Flags:" | head -n -3 > build/cmd/global.txt
	@for cmd in $(CMD_RUN); do \
		docker run orchestrate help $$cmd run | grep -B 9999 "Global Flags:" | tail -n +3 | head -n -3 > build/cmd/$$cmd-run.txt; \
	done
	@for cmd in $(CMD_MIGRATE); do \
		docker run orchestrate help $$cmd migrate | grep -B 9999 "Global Flags:" | tail -n +3 | head -n -3 > build/cmd/$$cmd-migrate.txt; \
	done

# Protobuf
protobuf: ## Generate protobuf stubs
	@docker-compose -f scripts/docker-compose.yml up

# Create kafka topics
topics:
	@bash scripts/kafka/initTopics.sh

gobuild:
	@GOOS=linux GOARCH=amd64 go build -i -o ./build/bin/orchestrate

docker-build:
	@DOCKER_BUILDKIT=1 docker build -t orchestrate .

gobuild-e2e:
	@GOOS=linux GOARCH=amd64 go build -i -o ./build/bin/e2e ./tests/cmd

orchestrate: gobuild
	@docker-compose up -d $(CMD_RUN)

stop-orchestrate:
	@docker-compose stop $(CMD_RUN)

down-orchestrate:
	@docker-compose down --volumes --timeout 0

deps:
	@docker-compose -f scripts/deps/docker-compose.yml up -d

down-deps:
	@docker-compose -f scripts/deps/docker-compose.yml down --volumes --timeout 0

quorum:
	@docker-compose -f scripts/deps/docker-compose.quorum.yml up -d

stop-quorum:
	@docker-compose -f scripts/deps/docker-compose.quorum.yml stop

down-quorum:
	@docker-compose -f scripts/deps/docker-compose.quorum.yml down --volumes --timeout 0

up-all: deps quorum orchestrate

down-all: down-orchestrate down-quorum down-deps
