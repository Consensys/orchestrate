# Tx-Sender

Core-Stack is a blockchain **Transaction Orchestration** system that can operate multiple chains simultaneously.
It provides production ready and agnostic mechanisms for transaction crafting, nonce management, transaction signing, transaction receipt listening, transaction receipt decoding, faucet and more.

Core-Stack is a **microservices architecture** composed of APIs & Workers. 
Workers communicate following **publish-subscribe** pattern using *Apache Kafka* as message broker. 
All messages are standardized using **protobuf** and **grpc**.

Core-Stack is Plug & Play, a user only needs to send a business protobuf message to the input topic,
Core-Stack then manages the full lifecycle of the transaction from crafting the transaction to decoding event logs data.

## Goal

Tx-Sender is a Core-Stack worker responsible to 

- **Store Transaction Trace**
- **Send Transaction to Ethereum node**

It consumes message from *tx signer* Kafka topic.

## Quick-Start

### Prerequisite

Having *docker* and *docker-compose* install

### Start the application

You can start the application with default configuration by running

```sh
$ docker-compose up
```

### Configure application

Application can be configured through flags or environment variables, you can run the ```help run``` command line


```text
Run application

Usage:
  app run [flags]

Flags:
      --eth-client strings      Ethereum client URLs.
                                Environment variable: "ETH_CLIENT_URL" (default [https://ropsten.infura.io/v3/81e039ce6c8a465180822b525e3644d7,https://rinkeby.infura.io/v3/bfc9d6e51fbc4d3db54bea58d1094f9c,https://kovan.infura.io/v3/bfc9d6e51fbc4d3db54bea58d1094f9c,https://mainnet.infura.io/v3/bfc9d6e51fbc4d3db54bea58d1094f9c])
  -h, --help                    help for run
      --http-hostname string    Hostname to expose healthchecks and metrics.
                                Environment variable: "HTTP_HOSTNAME" (default ":8080")
      --kafka-address strings   Address of Kafka server to connect to.
                                Environment variable: "KAFKA_ADDRESS" (default [localhost:9092])
      --worker-group string     Kafka consumer group. 
                                Environment variable: "KAFKA_SENDER_GROUP" (default "tx-sender-group")
      --worker-in string        Kafka topic to consume message from.
                                Environment variable: "KAFKA_TOPIC_TX_SENDER" (default "topic-tx-sender")
      --worker-slots uint       Maximum number of messages the worker can treat in parallel.
                                Environment variable: "WORKER_SLOTS" (default 100)

Global Flags:
      --log-format string   Log formatter (one of ["text" "json"]).
                            Environment variable: "LOG_FORMAT" (default "text")
      --log-level string    Log level (one of ["panic" "fatal" "error" "warn" "info" "debug" "trace"]).
                            Environment variable: "LOG_LEVEL" (default "debug")
```

## High Level Architecture

Tx-Sender expects all consumed messages to respect the [Core-Stack standard protobuf format](https://gitlab.com/ConsenSys/client/fr/core-stack/core/blob/master/protobuf)

Consumed messages should have 

- ```Chain``` entry set
- ```Tx``` entry set with the ```Raw``` fields

1. **Send Transaction**

Once the tx sender worker unmarshall a message from Kafka, it sends the data located in ```Tx.Raw``` into the ```ETH_CLIENT_URL``` corresponding to the chainId located in the trace.

2. **Store Transaction Trace**

TODO

It request a credit to *Faucet*

## Technical Architecture

![alt core-stack-architecture](https://gitlab.com/ConsenSys/client/fr/core-stack/doc/raw/master/diagrams/Core_Stack_Architecture.png)
