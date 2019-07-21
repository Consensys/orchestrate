# Tx-Sender

## Goal

Tx-Sender is a CoreStack worker responsible for:

- **Store Transaction Trace** by sending it to *API-Context-Store*;
- **Send Transaction to Ethereum node**.

It consumes messages from *tx signer* Kafka topic.

## Quick-Start

### Prerequisites

Having *docker* and *docker-compose* install

### Start the application

You can start the application with default configuration by running

```sh
$ go run . run
```

### Running e2e tests

1. Run testing environment

```bash
docker-compose -f e2e/docker-compose.yml up
```

2. Create topic in your local kafka

```bash
./e2e/initTestTopic.sh
```
3. Run worker

```bash
go run . run --jaeger-service TX-SENDER --http-hostname :8081 --grpc-store-target :8080 --eth-client http://localhost:8545
```

3. Run producer that will write messages 

// TODO: fix e2e producer, see e2e/producer/main.go

```bash
 go run e2e/producer/main.go
```