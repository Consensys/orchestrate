# Tx-Listener

Core-Stack is a blockchain **Transaction Orchestration** system that can operate multiple chains simultaneously.
It provides production ready and agnostic mechanisms for transaction crafting, nonce management, transaction signing, transaction receipt listening, transaction receipt decoding, faucet and more.

Core-Stack is a **microservices architecture** composed of APIs & Workers. 
Workers communicate following **publish-subscribe** pattern using *Apache Kafka* as message broker. 
All messages are standardized using **protobuf** and **grpc**.

Core-Stack is Plug & Play, a user only needs to send a business protobuf message to the input topic,
Core-Stack then manages the full lifecycle of the transaction from crafting the transaction to decoding event logs data.

## Goal

Tx-Listener is a Core-Stack worker responsible to 

- **Catch transaction receipts** 
- **Load & reconstitute Transaction Envelope** by interogating *API-Context-Store*

It consumes message from *tx signer* Kafka topic.

## Quick-Start

### Prerequisite

Having *docker* and *docker-compose* install

### Start the application

You can start the application with default configuration by running

```sh
$ docker-compose up
```

### Run e2e tests:

1. Set "ETH_CLIENT_URL"

```bash
export ETH_CLIENT_URL="https://ropsten.infura.io/v3/81e039ce6c8a465180822b525e3644d7 https://rinkeby.infura.io/v3/bfc9d6e51fbc4d3db54bea58d1094f9c https://kovan.infura.io/v3/bfc9d6e51fbc4d3db54bea58d1094f9c https://mainnet.infura.io/v3/bfc9d6e51fbc4d3db54bea58d1094f9c"
```

2. Run environment for the worker

```bash
docker-compose -f e2e/docker-compose.yml up
```

3. Run envelope store

First go to the `service/envelope-store` and then run the following command:

```bash
DB_DATABASE=envelope-store DB_HOST=localhost DB_USER=envelope-store go run . run
```

4. Run worker

```bash
go run . run --http-hostname ':8081' --grpc-store-target localhost:8080
```

It will start a worker on port `8081` and it will connect to the envelope store running on port `8080`.

### Configure application

Application can be configured through flags or environment variables, you can run the ```help run``` command line

```text
Run application

Usage:
  app run [flags]

Flags:
      --engine-slots uint                 Maximum number of messages the engine can treat concurrently.
                                          Environment variable: "ENGINE_SLOTS" (default 20)
      --eth-client strings                Ethereum client url
                                          Environment variable: "ETH_CLIENT_URL"
      --grpc-store-target string          GRPC Context Store target (See https://github.com/grpc/grpc/blob/master/doc/naming.md)
                                          Environment variable: "GRPC_STORE_TARGET"
  -h, --help                              help for run
      --http-hostname string              Hostname to expose HTTP server
                                          Environment variable: "HTTP_HOSTNAME" (default ":8080")
      --jaeger-disabled                   Disable Jaeger reporting
                                          Environment variable: "JAEGER_DISABLED"
      --jaeger-endpoint string            Jaeger collector endpoint to send spans to
                                          Environment variable: "JAEGER_ENDPOINT"
      --jaeger-host string                Jaeger host.
                                          Environment variable: "JAEGER_AGENT_HOST" (default "localhost")
      --jaeger-password string            Jaeger collector password
                                          Environment variable: "JAEGER_PASSWORD"
      --jaeger-port int                   Jaeger port
                                          Environment variable: "JAEGER_AGENT_PORT" (default 6831)
      --jaeger-rpc-metrics                Enable Jaeger RPC metrics
                                          Environment variable: "JAEGER_RPC_METRICS"
      --jaeger-sampler-param int          Jaeger sampler
                                          Environment variable: "JAEGER_SAMPLER_PARAM" (default 1)
      --jaeger-sampler-type string        Jaeger sampler
                                          Environment variable: "JAEGER_SAMPLER_TYPE" (default "const")
      --jaeger-service string             Jaeger ServiceName to use on the tracer
                                          Environment variable: "JAEGER_SERVICE_NAME" (default "jaeger")
      --jaeger-user string                Jaeger collector User
                                          Environment variable: "JAEGER_USER"
      --kafka-address strings             Address of Kafka server to connect to.
                                          Environment variable: "KAFKA_ADDRESS" (default [localhost:9092])
      --kafka-group string                Address of Kafka server to connect to.
                                          Environment variable: "KAFKA_GROUP" (default "group-e2e")
      --listener-block-backoff duration   Backoff time to wait before retrying after failing to find a mined block
                                          Environment variable: "LISTENER_BLOCK_BACKOFF" (default 1s)
      --listener-block-limit int          Limit number of block that can be prefetched while listening
                                          Environment variable: "LISTENER_BLOCK_LIMIT" (default 40)
      --listener-start strings            Position listener should start listening from (format <chainID>:<blockNumber>-<txIndex> or <chainID>:<blockNumber>) (e.g. 0x2a:2348721-5 or 0x3:latest)
                                          Environment variable: "LISTENER_START"
      --listener-start-default string     Default block position listener should start listening from (one of 'latest', 'oldest', 'genesis')
                                          Environment variable: "LISTENER_START_DEFAULT" (default "oldest")
      --listener-tracker-depth int        Depth at which we consider a block final (to avoid falling into a re-org)
                                          Environment variable: "LISTENER_TRACKER_DEPTH"
      --topic-decoder string              Kafka topic for messages waiting to have receipt decoded
                                          Environment variable: "KAFKA_TOPIC_TX_DECODER" (default "topic-tx-decoder")

Global Flags:
      --log-format string   Log formatter (one of ["text" "json"]).
                            Environment variable: "LOG_FORMAT" (default "text")
      --log-level string    Log level (one of ["panic" "fatal" "error" "warn" "info" "debug" "trace"]).
                            Environment variable: "LOG_LEVEL" (default "debug")
```

## High Level Architecture

1. **Catch transaction Receipts**

Tx-Listener listen to chains and retrieve transaction receipts as they are mined

2. **Load Transaction Envelope**

Each time it sees a new receipt Tx-Listener interogates *API-Context-Store* in order to possibly retrieve an `Envelope` associated to the transaction and reconstitute it.

## Technical Architecture

![alt core-stack-architecture](https://gitlab.com/ConsenSys/client/fr/core-stack/doc/raw/master/diagrams/Core_Stack_Architecture.png)
