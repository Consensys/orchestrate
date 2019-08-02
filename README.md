# Tx-Sender

## Goal

Tx-Sender is a CoreStack worker responsible for:

- **Store Transaction Envelope** by sending it to *API-Context-Store*;
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
go run . run --jaeger-service TX-SENDER --http-hostname :8081 --grpc-target-envelope-store :8080 --eth-client http://localhost:8545
```

3. Run producer that will write messages 

// TODO: fix e2e producer, see e2e/producer/main.go

```bash
 go run e2e/producer/main.go
```



```sh
$ docker-compose up
```

## High Level Architecture

Tx-Sender expects all consumed messages to respect a specific CoreStack protobuf format.

Consumed messages should have:

- ```Chain``` entry set;
- ```Tx``` entry set with the ```Raw``` fields.

### Tx-sender can handle two types of transactions:

**1. Standard case: Send signed transaction**

***1.1 Store Transaction Envelope***

Once the Tx-sender worker unmarshall a message from Kafka, the transaction envelope is stored in the *Envelope-Store* with the status `pending`.

***1.2 Send Transaction***

It sends the data located in ```Tx.Raw``` into the ```ETH_CLIENT_URL``` corresponding to the `chainId` located in the envelope.

***1.3 Update Transaction status in the envelope-store***

In the *Envelope-Store*, the transaction status is updated to `pending`.

***2. Quorum case (using `sendTransaction`): Send unsigned transactions to Quorum***

***2.1 Send Transaction***

Once the tx-sender worker unmarshall a message from Kafka, it sends an unsigned transaction to the Quorum node (```ETH_CLIENT_URL```) and retrieves the `txHash`.

***2.2 Store Transaction Envelope***

The envelope is stored in the *Envelope-Store* with the status `pending`.

