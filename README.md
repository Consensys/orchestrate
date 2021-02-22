# Package

Orchestrate is a blockchain **Transaction Orchestration** system that can operate multiple chains simultaneously.
It provides production grade and agnostic mechanisms for transaction crafting, nonce management, transaction signing, transaction receipt listening, transaction receipt decoding, faucet and more.

Orchestrate is a **microservices architecture** composed of APIs & Workers. 
Workers communicate following **publish-subscribe** pattern using *Apache Kafka* as message broker. 
All messages are standardized using **protobuf** and **gRPC**.

Orchestrate is Plug & Play, a user only needs to send a business protobuf message to Orchestrate input topic,
Orchestrate then manages the full lifecycle of the transaction from crafting the transaction to decoding event logs data.

## Goal

Package is a low level library in Orchestrate dependency tree. In particular it implements

- *protobuf* containing all protobuf schemes
- *core* that defines core structural elements of Orchestrate (such as ``types.Context``) 
- *common* which are resources that are shared between multiple Orchestrate services

## Errors

Internal errors are built on top of protobuf and contain

- `string` **message** describing the error
- `uint64` **error code** that should be interpreted as five nibbles hex code (e.g. `4096` <=> `01000` and `989956` <=> `F1B04`)
- `string` **component** indicating in system the error has been raised
- `map<string, string>` **extra** holding extra information to diagnose the error

### Error codes

Error codes are `uint64` that maps to five hex character code

| Class | Subclass | Error Code | Condition                     | Comment                                    |
|-------|----------|------------|-------------------------------|--------------------------------------------|
| 01XXX |          |    01000   | warning                       | Raised to indicate a warning               |
| 01XXX |   011XX  |    01100   | retry_warning                 | Error occurred system retries              |
| 01XXX |   012XX  |    01200   | faucet_warning                | Faucet credit has been denied              |
| 01XXX |   013XX  |    01300   | invalid_nonce_warning         | Exposure to send tx with invalid nonce     |
| 01XXX |   013XX  |    01301   | nonce_too_high_warning        | Exposure to send tx with nonce too high    |
| 01XXX |   013XX  |    01302   | nonce_too_high_low            | Exposure to send tx with nonce too low     |
| 08XXX |          |    08000   | connection_exception          | Failed connecting to an external service   |
| 08XXX |   081XX  |    08100   | kafka_connection_exception    | Failed connecting to Kafka                 |
| 08XXX |   082XX  |    08200   | http_connection_exception     | Failed connecting to an HTTP service       |
| 08XXX |   083XX  |    08300   | ethereum_connection_exception | Failed connecting to Ethereum jsonRPC API  |
| 08XXX |   084XX  |    08400   | grpc_connection_exception     | Failed connecting to a GRPC API            |
| 08XXX |   085XX  |    08500   | redis_connection_exception    | Failed connecting to Redis                 |
| 09XXX |          |    09000   | authentication_exception      | Unauthorized                               |
| 09XXX |          |    09001   | unauthenticated               | Invalid credentials                        |
| 09XXX |          |    09002   | permission_denied             | Operation not permitted                    |
| 0AXXX |          |    0A000   | feature_not_supported         | Feature is not supported                   |
| 24XXX |          |    24000   | invalid_state                 | System in invalid state                    |
| 24XXX |  241XX   |    24100   | failed_precondition           | System not in required state for operation |
| 24XXX |  242XX   |    24200   | conflicted                    | Op. conflicted with system current state   |
| 42XXX |          |    42000   | invalid_data                  | Failed to process data                     |
| 42XXX |          |    42001   | out_of_range                  | operation attempted past valid range       |
| 42XXX |   421XX  |    42100   | invalid_encoding              | Failed to decode a message                 |
| 42XXX |   422XX  |    42200   | invalid_solidity_data         | Failed to process Solidity related data    |
| 42XXX |   422XX  |    42201   | invalid_method_signature      | Invalid Solidity method signature          |
| 42XXX |   422XX  |    42202   | invalid_args_count            | Invalid args count provided                |
| 42XXX |   422XX  |    42203   | invalid_arg                   | Invalid arg provided                       |
| 42XXX |   422XX  |    42204   | invalid_topics_count          | Invalid topics count in event log          |
| 42XXX |   422XX  |    42205   | invalid_event_data            | Invalid data in event log                  |
| 42XXX |   423XX  |    42300   | invalid_format                | Data does not match expected format        |
| 42XXX |   424XX  |    42400   | invalid_parameter             | Invalid parameter provided                 |
| 53XXX |          |    53000   | insufficient_resources        | System can not handle more operations      |
| 57XXX |          |    57000   | operator_intervention         | Operator interfered with operation         |
| 57XXX |          |    57001   | operation_canceled            | Operation canceled (typically by caller)   |
| C0XXX |          |    C0000   | crypto_operation_exception    | Failed a cryptographical operation         |
| DBXXX |          |    DB000   | storage_exception             | Failed accessing stored data               |
| DBXXX |   DB1XX  |    DB100   | constraint_violated           | Data constraint violated                   |
| DBXXX |   DB1XX  |    DB101   | already_exists                | Resource with same index already existed   |
| DBXXX |   DB2XX  |    DB200   | not_found                     | No data found for given parameters         |
| F0XXX |          |    F0000   | invalid_config                | Invalid configuration                      |
| FFXXX |          |    FF000   | internal_error                | Internal error                             |
| FFXXX |   FF1XX  |    FF100   | data_corrupted                | Data is corrupted                          |
| BEXXX |   BEXXX  |    BE1000  | ethereum_nonce_too_low        | Nonce is too low                           |

## Documentation

### Product

This project documentation is available at https://docs.orchestrate.consensys.net/

The documentation source repos is https://github.com/ConsenSys/doc.orchestrate

### API

Generated API documentation (using a Github Action Workflow in this repos) is available at https://consensys.github.io/orchestrate/

## Local Development Interfaces

### APIs

[Orchestrate OpenAPI](http://localhost:8031/swagger)

### Tools

[Portainer UI](http://localhost:9000)

[PGAdmin UI](http://localhost:9001)

[Jaeger UI](http://localhost:16686)
