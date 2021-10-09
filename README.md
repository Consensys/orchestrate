[![Website](https://img.shields.io/website?label=documentation&url=https%3A%2F%2Fdocs.orchestrate.consensys.net%2F)](https://docs.orchestrate.consensys.net/)
[![Website](https://img.shields.io/website?url=https%3A%2F%2Fconsensys.net%2Forchestrate%2F)](https://consensys.net/quorum/)

[![CircleCI](https://img.shields.io/circleci/build/gh/ConsenSys/orchestrate?token=7062612dcd5a98913aa1b330ae48b6a527be52eb)](https://circleci.com/gh/ConsenSys/orchestrate)
[![Go Report Card](https://goreportcard.com/badge/github.com/ConsenSys/orchestrate)](https://goreportcard.com/report/github.com/ConsenSys/orchestrate)

# Orchestrate

Orchestrate is a blockchain **Transaction Orchestration** system that can operate multiple chains simultaneously.
It provides production grade and agnostic mechanisms for transaction crafting, nonce management, transaction signing, transaction receipt listening, transaction receipt decoding, faucet and more.

Orchestrate is a **microservices architecture** composed of APIs & Workers. 
Workers communicate following **publish-subscribe** pattern using *Apache Kafka* as message broker. 

Orchestrate is Plug & Play, a user only needs to send a business protobuf message to Orchestrate input topic,
Orchestrate then manages the full lifecycle of the transaction from crafting the transaction to decoding event logs data.

## Useful links

* [User Documentation](http://docs.orchestrate.consensys.net/)
* [Orchestrate OpenAPI](http://localhost:8031/swagger)
* [GitHub Project](https://github.com/ConsenSys/orchestrate)
* [issues](https://github.com/ConsenSys/orchestrate/issues)
* [Changelog](https://github.com/ConsenSys/orchestrate/blob/main/CHANGELOG.md)
* [Quorum Key Manager](https://github.com/orchestrate/quorum-key-manager)
* [Helm Charts](https://github.com/ConsenSys/orchestrate-helm)
* [Kubernetes deployment example](https://github.com/ConsenSys/orchestrate-kubernetes)

## Run Orchestrate

Now launch Orchestrate service using docker-compose with the following command:

```bash
docker-compose up -d api tx-sender tx-listener
```

Orchestrate is connected to [Quorum Key Manager(QKM)](https://github.com/ConsenSys/quorum-key-manager) to perform every
wallet actions (creation, importing, signing, etc...). In this case for QKM setup we are 
 connecting to a local Hashicorp Vault Server which is enhanced with [Quorum Hashicorp Vault Pluging](https://github.com/ConsenSys/quorum-hashicorp-vault-plugin)
 to support private keys under hashicorp storage. 

## Build from source

### Prerequisites

To build binary locally requires Go (version 1.16 or later) and C compiler. 

### Build

After downloading dependencies (ie `go mod download`) you can run following command to compile the binary

```bash
go build -o ./build/bin/orchestrate
```

Binary will be located in `./build/bin/orchestrate

## License

Orchestrate is licensed under the BSL 1.1.

Please refer to the [LICENSE file](LICENSE) for a detailed description of the license.

Please contact [orchestrate@consensys.net](mailto:orchestrate@consensys.net) if you need to purchase a license for a production use-case.  


