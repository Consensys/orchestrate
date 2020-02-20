GOFILES := $(shell find . -name '*.go' | grep -v services/chain-registry/genstatic/gen.go | egrep -v "^\./\.go" | grep -v _test.go)
PACKAGES ?= $(shell go list ./... | go list ./... | grep -Fv -e e2e -e examples -e genstatic -e mocks )
CMD_RUN = tx-crafter tx-nonce tx-signer tx-sender tx-listener contract-registry chain-registry envelope-store
CMD_MIGRATE = contract-registry envelope-store chain-registry

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

coverage: postgres run-coverage down-postgres
	@$(OPEN) build/coverage/coverage.html 2>/dev/null

race: ## Run data race detector
	@go test -count=1 -race -short ${PACKAGES}

mod-tidy: ## Run deps cleanup
	@go mod tidy

lint: ## Run linter to fix issues
	@misspell -w $(GOFILES)
	@golangci-lint run --fix

lint-ci: ## Check linting
	@misspell -error $(GOFILES)
	@golangci-lint run

run-e2e: gobuild-e2e
	@docker-compose up e2e
	@docker-compose -f scripts/report/docker-compose.yml up

e2e: run-e2e
	@$(OPEN) build/report/report.html 2>/dev/null

clean: mod-tidy lint-ci protobuf ## Run all clean-up tasks

generate-mocks:
	mockgen -source=services/chain-registry/client/client.go -destination=services/chain-registry/client/mocks/mock_client.go -package=mocks
	mockgen -source=ethereum/ethclient/ethclient.go -destination=ethereum/ethclient/mocks/mock_client.go -package=mocks
	mockgen -source=types/contract-registry/registry.pb.go -destination=types/contract-registry/client/mocks/mock_client.go -package=mocks
	mockgen -source=types/envelope-store/store.pb.go -destination=types/envelope-store/client/mocks/mock_client.go -package=mocks
	mockgen -source=services/chain-registry/store/types/store.go -destination=services/chain-registry/store/mocks/mock_store.go -package=mocks
	mockgen -source=services/multi-vault/keystore/keystore.go -destination=services/multi-vault/keystore/mocks/mock_client.go -package=mocks

# Tools
lint-tools: ## Install linting tools
	@GO111MODULE=off go get -u github.com/client9/misspell/cmd/misspell
	@GO111MODULE=off go get -u github.com/golangci/golangci-lint/cmd/golangci-lint

tools: lint-tools ## Install test tools
	@GO111MODULE=off go get -u github.com/cucumber/godog/cmd/godog
	@GO111MODULE=off go get -u github.com/golang/mock/gomock
	@GO111MODULE=off go get -u github.com/golang/mock/mockgen
	@GO111MODULE=off go get -u github.com/golang/protobuf/protoc-gen-go
	@GO111MODULE=off go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
	@GO111MODULE=off go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
	@GO111MODULE=off go get -u github.com/swaggo/swag/cmd/swag

# Help
help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

gen-help: gobuild ## Generate Command Help file
	@mkdir -p build/cmd
	@./build/bin/orchestrate help tx-crafter | grep -A 9999 "Global Flags:" | head -n -2 > build/cmd/global.txt
	@for cmd in $(CMD_RUN); do \
		./build/bin/orchestrate help $$cmd run | grep -B 9999 "Global Flags:" | tail -n +3 | head -n -2 > build/cmd/$$cmd-run.txt; \
	done
	@for cmd in $(CMD_MIGRATE); do \
		./build/bin/orchestrate help $$cmd migrate | grep -B 9999 "Global Flags:" | tail -n +3 | head -n -2 > build/cmd/$$cmd-migrate.txt; \
	done

gen-help-docker: docker-build ## Generate Command Help file using docker
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
	@docker-compose -f scripts/protobuf/docker-compose.yml up --build

topics: ## Create kafka topics
	@bash scripts/deps/kafka/initTopics.sh

gobuild: ## Build Orchestrate Go binary
	@GOOS=linux GOARCH=amd64 go build -i -o ./build/bin/orchestrate

docker-build: ## Build Orchestrate Docker image
	@DOCKER_BUILDKIT=1 docker build -t orchestrate .

bootstrap: ## Wait for dependencies to be ready
	@bash scripts/bootstrap.sh

gobuild-e2e: ## Build Orchestrate e2e Docker image
	@GOOS=linux GOARCH=amd64 go build -i -o ./build/bin/e2e ./tests/cmd

orchestrate: gobuild ## Start Orchestrate
	@docker-compose up -d $(CMD_RUN)

stop-orchestrate: ## Stop Orchestrate
	@docker-compose stop $(CMD_RUN)

down-orchestrate:## Down Orchestrate
	@docker-compose down --volumes --timeout 0

deps:
	@docker-compose -f scripts/deps/docker-compose.yml up -d

down-deps:
	@docker-compose -f scripts/deps/docker-compose.yml down --volumes --timeout 0

geth:
	@docker-compose -f scripts/geth/docker-compose.yml up -d

stop-geth:
	@docker-compose -f scripts/geth/docker-compose.yml stop

down-geth:
	@docker-compose -f scripts/geth/docker-compose.yml down  --volumes --timeout 0

quorum:
	@docker-compose -f scripts/quorum/docker-compose.yml up -d

stop-quorum:
	@docker-compose -f scripts/quorum/docker-compose.yml stop

down-quorum:
	@docker-compose -f scripts/quorum/docker-compose.yml down --volumes --timeout 0

besu:
	@docker-compose -f scripts/besu/docker-compose.yml up -d

stop-besu:
	@docker-compose -f scripts/besu/docker-compose.yml stop

down-besu:
	@docker-compose -f scripts/besu/docker-compose.yml down --volumes --timeout 0

postgres:
	@docker-compose -f e2e/docker-compose.yml up -d postgres

down-postgres:
	@docker-compose -f e2e/docker-compose.yml rm --force -s -v postgres

up: deps geth besu quorum bootstrap orchestrate ## Start Orchestrate and deps

down: down-orchestrate down-quorum down-geth down-besu down-deps  ## Down Orchestrate and deps

hashicorp-accounts:
	@bash scripts/deps/config/hashicorp/vault.sh kv list secret/default

hashicorp-token-lookup:
	@bash scripts/deps/config/hashicorp/vault.sh token lookup

hashicorp-vault:
	@bash scripts/deps/config/hashicorp/vault.sh $(COMMAND)

pgadmin:
	@docker-compose -f scripts/deps/docker-compose-tools.yml up -d pgadmin

down-pgadmin:
	@docker-compose -f scripts/deps/docker-compose-tools.yml rm --force -s -v pgadmin
