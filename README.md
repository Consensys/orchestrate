# Transaction Signer configuration

## Goal

Tx-Signer is a CoreStack worker responsible for: 

- **Signing transaction**

It consumes messages from *tx signer* Kafka topic and publishs to *tx sender* topic.

## High Level Architecture

Tx-Signer expects all consumed messages to respect a specific CoreStack protobuf format.

Consumed messages should have:

- ```Chain``` attribute set with ```id``` of the chain to send the transaction to;
- ```From``` attribute set with an ```raw```;
- ```Tx``` attribute set with the following fields set up:
  - ```Nonce```
  - ```To```
  - ```Value```
  - ```GasLimit```
  - ```GasPrice```
  - ```Data```

## Signing Account configuration

In order to sign transactions using a different private key, a user needs to setup an environment variable call SECRET_PKEY.

Multiple values can be set for SECRETE_PKEY, for example:

`SECRET_PKEY="<PRIVATE KEY #1> <PRIVATE KEY #2> <PRIVATE KEY #...> <PRIVATE KEY #n>"`

***Note**: How to storage keys its covered on the storage section.


## Quick-Start

### Prerequisites

- Having ```docker``` and ```docker-compose``` installed;
- Having Go 1.11 installed or upper.

### Start the application

To quickly start the application

**1. Start e2e env**

```sh
$ docker-compose -f e2e/docker-compose.yml up
```

**2. Start worker**

```sh
$ go run . run
```

### Running e2e tests

1. Run testing environment

```bash
docker-compose -f e2e/docker-compose.yml up
```

2. Run test consumer that should read from a topic where a worker is going to write 

```bash
go run e2e/consumer/main.go
```

3. Run worker

```bash
go run . run --jaeger-service TX-SIGNER
```

4. Run producer that will write messages 

```bash
 go run e2e/producer/main.go
```
