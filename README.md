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

## Errors

Error codes are uint64 that maps to five hex character code

| Class | Subclass | Hex Code   |  Code     | Description   |
|-------|----------|------------| --------- | --------------|
| 01XXX |          |    01000   | 4096   | Indicate a warning operation               |
|       |   012XX  |    01200   | 4608   | Faucet credit has been denied              |
| 013XX |          |            |        | Sent tx with invalid nonce                 |
|       |          |    01301   | 4865   | Nonce too high                |
|       |          |    01302   | 4866   | Nonce too low                 |
| 08XXX |          |    08000   | 32768  | Failed connecting to an external service   |
|       |   081XX  |    08100   | 33024  | Connecting to Kafka                 |
|       |   082XX  |    08200   | 33280  | Connecting to an HTTP service       |
|       |   083XX  |    08300   | 33536  | Connecting to Ethereum Node  |
|       |   085XX  |    08500   | 34048  | Connecting to Redis                 |
|       |   086XX  |    08600   | 34304  | Connecting to Postgres                 |
|       |   087XX  |    08700   | 34560  | Connecting to external services                 |
| 09XXX |          |    09000   | 36864  | Unauthorized operation                              |
|       |          |    09001   | 36865  | Invalid credentials                        |
|       |          |    09002   | 36866  | Operation not permitted                    |
| 0AXXX |          |    0A000   | 40960  | Feature is not supported                   |
| 24XXX |          |    24000   | 147456 | Invalid data state                     |
|       |   242XX  |    24200   | 147968 | Conflicted with current data state   |
| 42XXX |          |    42000   | 270336 | Failed to process input data                     |
|       |   421XX  |    42100   | 270592 | Decoding a message                 |
|       |   422XX  |    42200   | 270848 | Processing Solidity related data    |
|       |          |    42201   | 270849 | Invalid Solidity method signature          |
|       |          |    42202   | 270850 | Invalid arguments count provided                |
|       |          |    42203   | 270851 | Invalid provided arguments                        |
|       |          |    42204   | 270852 | Invalid topics count in ABI event log          |
|       |          |    42205   | 270853 | Invalid data in ABI event log                  |
|       |          |    42300   | 271104 | Data does not match expected format        |
|       |          |    42400   | 271360 | Invalid provided parameter                  |
| BEXXX |          |    BE000   | 778240 | Failed a Ethereum operation                   |
|       |          |    BE001   | 778241 | Nonce too low            |
|       |          |    BE002   | 778242 | Invalid Nonce            |
| C0XXX |          |    C0000   | 786432 | Failed a cryptographic operation         |
|       |          |    C0001   | 786433 | Invalid cryptographic signature         |
| DBXXX |          |    DB000   | 897024 | Failed data operation               |
|       |   DB1XX  |    DB100   | 897280 | Data constraint violated                   |
|       |          |    DB101   | 897281 | Resource with same unique index already existed   |
|       |   DB2XX  |    DB200   | 897536 | No data found         |
| F0XXX |          |    F0000   | 983040 | Invalid configuration                      |
| FFXXX |          |    FF000   | 1044480 | Internal error                             |
|       |   FF1XX  |    FF100   | 1044736 | Data is corrupted                          |
|       |   FF2XX  |    FF200   | 1044992 | Dependency failure                          |

## License

Orchestrate is licensed under the BSL 1.1.

Please refer to the [LICENSE file](LICENSE) for a detailed description of the license.

Please contact [orchestrate@consensys.net](mailto:orchestrate@consensys.net) if you need to purchase a license for a production use-case.  


