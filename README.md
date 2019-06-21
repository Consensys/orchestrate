# Tx-Crafter

Core-Stack is a blockchain **Transaction Orchestration** system that can operate multiple chains simultaneously.
It provides production grade and agnostic mechanisms for transaction crafting, nonce management, transaction signing, transaction receipt listening, transaction receipt decoding, faucet and more.

Core-Stack is a **microservices architecture** composed of APIs & Workers. 
Workers communicate following **publish-subscribe** pattern using *Apache Kafka* as message broker. 
All messages are standardized using **protobuf** and **GRPC**.

Core-Stack is Plug & Play, a user only needs to send a business protobuf message to Core-Stack input topic,
Core-Stack then manages the full lifecycle of the transaction from crafting the transaction to decoding event logs data.

## Goal

Tx-Crafter is a Core-Stack worker responsible to 

- **Craft transaction payload**  
- **Set Gas Price** of the transaction
- **Set Gas Limit** of the transaction
- **Request Faucet crediting**

It consumes message from *tx crafting* Kafka topic and publish to *tx nonce* topic.

## Quick-Start

### Prerequisite

- Having ```docker``` and ```docker-compose``` installed
- Having Go 1.11 installed or upper

### Start the application

To quickly start the application

1. Start Kafka broker

```sh
$ docker-compose -f e2e/docker-compose.yml up
```


2. Start worker

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
go run . run  --eth-client HTTP://127.0.0.1:8545
```

4. Run producer that will write messages 

```bash
 go run e2e/producer/main.go
```

### Configure application

Application can be configured through flags or environment variables, you can run the ```help run``` command line

```sh
$ go run . help run
```

```text
Run application

Usage:
  app run [flags]

Flags:
      --abi strings                  Smart Contract ABIs to register for crafting (expected format <contract>:<abi>:<bytecode>)
                                     Environment variable: "ABI"
      --engine-slots uint            Maximum number of messages the engine can treat concurrently.
                                     Environment variable: "WORKER_SLOTS" (default 20)
      --eth-client strings           Ethereum client url
                                     Environment variable: "ETH_CLIENT_URL"
      --faucet string                Type of Faucet (one of ["mock" "sarama"])
                                     Environment variable: "FAUCET" (default "mock")
      --faucet-amount string         Amount to credit when calling Faucet (Wei in decimal format)
                                     Environment variable: "FAUCET_CREDIT_AMOUNT" (default "100000000000000000")
      --faucet-blacklist strings     Blacklisted address (format <chainID>-<Address>)
                                     Environment variable: "FAUCET_BLACKLIST"
      --faucet-cooldown duration     Faucet minimum to wait before crediting an address again
                                     Environment variable: "FAUCET_COOLDOWN_TIME" (default 1m0s)
      --faucet-creditor strings      Address of Faucet on each chain (format <chainID>:<Address>)
                                     Environment variable: "FAUCET_CREDITOR_ADDRESS"
      --faucet-max string            Max balance (Wei in decimal format)
                                     Environment variable: "FAUCET_MAX_BALANCE" (default "200000000000000000")
  -h, --help                         help for run
      --http-hostname string         Hostname to expose healthchecks and metrics.
                                     Environment variable: "HTTP_HOSTNAME" (default ":8080")
      --jaeger-host string           Jaeger host.
                                     Environment variable: "JAEGER_HOST" (default "jaeger")
      --jaeger-port int              Jaeger port
                                     Environment variable: "JAEGER_PORT" (default 6831)
      --jaeger-sampler-param int     Jaeger sampler
                                     Environment variable: "JAEGER_SAMPLER_PARAM" (default 1)
      --jaeger-sampler-type string   Jaeger sampler
                                     Environment variable: "JAEGER_SAMPLER_TYPE" (default "const")
      --kafka-address strings        Address of Kafka server to connect to.
                                     Environment variable: "KAFKA_ADDRESS" (default [localhost:9092])
      --kafka-group string           Address of Kafka server to connect to.
                                     Environment variable: "KAFKA_GROUP" (default "group-e2e")
      --topic-crafter string         Kafka topic for messages waiting to have transaction payload crafted
                                     Environment variable: "KAFKA_TOPIC_TX_CRAFTER" (default "topic-tx-crafter")
      --topic-nonce string           Kafka topic for messages waiting to have transaction nonce set
                                     Environment variable: "kafka.topic.nonce" (default "topic-tx-nonce")

Global Flags:
      --log-format string   Log formatter (one of ["text" "json"]).
                            Environment variable: "LOG_FORMAT" (default "text")
      --log-level string    Log level (one of ["panic" "fatal" "error" "warn" "info" "debug" "trace"]).
                            Environment variable: "LOG_LEVEL" (default "debug")
```

## High Level Architecture

Tx-Crafter expect all consumed messages to respect [Core-Stack standard protobuf format](https://gitlab.com/ConsenSys/client/fr/core-stack/core/blob/master/protobuf)

Consumed messages should have 

- ```Chain``` attribute set with ```ID``` of the chain to send the transaction to
- ```Call``` attribute set or Tx-Worker will consider the transaction as a basic Ethereum transaction with no payload and will skip crafting

1. **Crafting**

To craft transaction payload Tx-Worker inspects the ```Call``` entry of input protobuf message 
 
- it expects ```Call.ID``` to be a string formated as ```<method>@<contract_name>``` (e.g. in case of an ERC20 transfer: "transfer@ERC20") (Note: this will evolve to handle versioning of contracts) 
- it expects ```Call.Args``` entry to be the ordered list of arguments to provide to the transaction call in ```string``` format (e.g. in case of an ERC20 transfer: ["0x6009608a02a7a15fd6689d6dad560c44e9ab61ff", "0x34fde"] for *to* and *value* args)

By basing on the ```Call.ID```, Tx-Worker requests the required  ABI from the *ABI registry*, then it casts ```Call.Args``` arguments in the expected Solidity type and craft payload.

2. **Gas Price**

Tx-Crafter interogates *Ethereum client*  by calling jsonRPC ```eth_gasPrice``` on chain ```Chain.ID```

3. **Gas Cost**

Tx-Crafter interogates *Ethereum client*  by calling jsonRPC ```eth_estimateGas``` on chain ```Chain.ID```

4. **Faucet**

It request a credit to *Faucet*

## Technical Architecture

![alt core-stack-architecture](https://gitlab.com/ConsenSys/client/fr/core-stack/doc/blob/master/diagrams/Core_Stack_Architecture.png)
