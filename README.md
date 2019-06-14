# E2E Tests

Core-Stack is a blockchain **Transaction Orchestration** system that can operate multiple chains simultaneously.
It provides production grade and agnostic mechanisms for transaction crafting, nonce management, transaction signing, transaction receipt listening, transaction receipt decoding, faucet and more.

Core-Stack is a **microservices architecture** composed of APIs & Workers. 
Workers communicate following **publish-subscribe** pattern using *Apache Kafka* as message broker. 
All messages are standardized using **protobuf** and **GRPC**.

Core-Stack is Plug & Play, a user only needs to send a business protobuf message to Core-Stack input topic,
Core-Stack then manages the full lifecycle of the transaction from crafting the transaction to decoding event logs data.

## Goal

The e2e tests worker aims to tests the end-to-end lifecycle transactions in CoreStack. Making sure that each workers and APIs are working well as a microservice system. To that end, this tool is able to: 

- **Send tests transactions** defined in a specific scenario  
- **Subscribe CoreStack Envelopes** that go in and out of each worker
- **Run features and scenarios tests** independently, in parallel, and defined in Cucumber
- **Generate reports** that summuries features and scenarios that succeded or failed and time each steps

## Quick-Start

### Prerequisite

- Having ```docker``` and ```docker-compose``` installed
- Having Go 1.11 installed or upper
- Having CoreStack and its infrastructure running

### Start the test

1. Start worker

```sh
$ go run . run
```

Note: you must have the `ETH_CLIENT_URL` environment variable set

2. Once the test is completed, you can generate reports

```sh
$ make report
```

The HTML report could be visualized in `report/output/report.html`

### Configure application

Application can be configured through flags or environment variables, you can run the ```help run``` command line

```
Run application

Usage:
  app run [flags]

Flags:
      --cucumber-concurrency int           Concurrency rate, not all formatters accepts this
                                           Environment variable: "CUCUMBER_CONCURRENCY" (default 1)
      --cucumber-format string             The formatter name
                                           Environment variable: "CUCUMBER_FORMAT" (default "pretty")
      --cucumber-nocolors                  Forces ansi color stripping
                                           Environment variable: "CUCUMBER_NOCOLORS"
      --cucumber-outputpath string         Where it should print the cucumber output (only works with cucumber format)
                                           Environment variable: "CUCUMBER_OUTPUTPATH"
      --cucumber-paths strings             All feature file paths
                                           Environment variable: "CUCUMBER_PATHS" (default [features])
      --cucumber-randomize int             Seed to randomize feature tests. The default value of -1 means to have a random seed. 0 means do not randomize
                                           Environment variable: "CUCUMBER_RANDOMIZE" (default -1)
      --cucumber-showstepdefinitions       Print step definitions found and exit
                                           Environment variable: "CUCUMBER_SHOWSTEPDEFINITION"
      --cucumber-steps-miningtimeout int   Duration for waiting envelopes to be processed by a blockchain before timeout
                                           Environment variable: "CUCUMBER_STEPS_MININGTIMEOUT" (default 10)
      --cucumber-steps-timeout int         Duration for waiting envelopes to be processed by a step method before timeout
                                           Environment variable: "CUCUMBER_STEPS_TIMEOUT" (default 5)
      --cucumber-stoponfailure             Stops on the first failure
                                           Environment variable: "CUCUMBER_STOPONFAILURE"
      --cucumber-strict                    Fail suite when there are pending or undefined steps
                                           Environment variable: "CUCUMBER_STRICT"
      --cucumber-tags string               Various filters for scenarios parsed from feature files
                                           Environment variable: "CUCUMBER_TAGS"
      --engine-slots uint                  Maximum number of messages the engine can treat concurrently.
                                           Environment variable: "ENGINE_SLOTS" (default 20)
      --eth-client strings                 Ethereum client url
                                           Environment variable: "ETH_CLIENT_URL"
      --grpc-store-target string           GRPC Context Store target (See https://github.com/grpc/grpc/blob/master/doc/naming.md)
                                           Environment variable: "GRPC_STORE_TARGET"
  -h, --help                               help for run
      --http-hostname string               Hostname to expose HTTP server
                                           Environment variable: "HTTP_HOSTNAME" (default ":8080")
      --jaeger-disabled                    Disable Jaeger reporting
                                           Environment variable: "JAEGER_DISABLED"
      --jaeger-endpoint string             Jaeger collector endpoint to send spans to
                                           Environment variable: "JAEGER_ENDPOINT"
      --jaeger-host string                 Jaeger host.
                                           Environment variable: "JAEGER_AGENT_HOST" (default "localhost")
      --jaeger-password string             Jaeger collector password
                                           Environment variable: "JAEGER_PASSWORD"
      --jaeger-port int                    Jaeger port
                                           Environment variable: "JAEGER_AGENT_PORT" (default 6831)
      --jaeger-rpc-metrics                 Enable Jaeger RPC metrics
                                           Environment variable: "JAEGER_RPC_METRICS"
      --jaeger-sampler-param int           Jaeger sampler
                                           Environment variable: "JAEGER_SAMPLER_PARAM" (default 1)
      --jaeger-sampler-type string         Jaeger sampler
                                           Environment variable: "JAEGER_SAMPLER_TYPE" (default "const")
      --jaeger-service string              Jaeger ServiceName to use on the tracer
                                           Environment variable: "JAEGER_SERVICE_NAME" (default "jaeger")
      --jaeger-user string                 Jaeger collector User
                                           Environment variable: "JAEGER_USER"
      --kafka-address strings              Address of Kafka server to connect to.
                                           Environment variable: "KAFKA_ADDRESS" (default [localhost:9092])
      --kafka-group string                 Address of Kafka server to connect to.
                                           Environment variable: "KAFKA_GROUP" (default "group-e2e")
      --topic-crafter string               Kafka topic for messages waiting to have transaction payload crafted
                                           Environment variable: "KAFKA_TOPIC_TX_CRAFTER" (default "topic-tx-crafter")
      --topic-decoded string               Kafka topic for messages which receipt has been decoded
                                           Environment variable: "KAFKA_TOPIC_TX_DECODED" (default "topic-tx-decoded")
      --topic-decoder string               Kafka topic for messages waiting to have receipt decoded
                                           Environment variable: "KAFKA_TOPIC_TX_DECODER" (default "topic-tx-decoder")
      --topic-nonce string                 Kafka topic for messages waiting to have transaction nonce set
                                           Environment variable: "kafka.topic.nonce" (default "topic-tx-nonce")
      --topic-sender string                Kafka topic for messages waiting to have transaction sent
                                           Environment variable: "KAFKA_TOPIC_TX_SENDER" (default "topic-tx-sender")
      --topic-signer string                Kafka topic for messages waiting to have transaction signed
                                           Environment variable: "KAFKA_TOPIC_TX_SIGNER" (default "topic-tx-signer")

Global Flags:
      --log-format string   Log formatter (one of ["text" "json"]).
                            Environment variable: "LOG_FORMAT" (default "text")
      --log-level string    Log level (one of ["panic" "fatal" "error" "warn" "info" "debug" "trace"]).
                            Environment variable: "LOG_LEVEL" (default "debug")
```

## Feature tests

This worker is implementing cucumber Behavior-Driven Development framework with [Godog](https://cucumber.io/docs/installation/golang/). You could find Gerkin files in the `features` folder and their implementations in `service/cucumber/steps/steps.go`

### Important scenario steps

Most of the standard steps are located in `features/deployment.feature` and `features/transaction.feature` which includes the following:

```
    Given I have the following envelope:
      | Alias       | chainId | from                                       | contractName | methodSignature | gas     |
      | SimpleToken | 888     | 0x7E654d251Da770A068413677967F6d3Ea2FeA9E4 | SimpleToken  | constructor()   | 2000000 |
```
This step is parsing a table into envelopes and stored in the `Scenario Context`. The labels are as close as the SDK, except `Alias` where it provides a name to a contract to be deployed. The mapping of labels could be found in `service/cucumber/steps/utils.go`


```    
    When I send these envelopes to CoreStack
```
This step will send to CoreStack every Envelopes in the `Scenario Context`.


```
    Then CoreStack should receive envelopes
```
This step will listen the entry topic of CoreStack, i.e. `topic-tx-crafter` and checks that the number of envelope sent corresponds to the number of envelopes received.

```
    Then the tx-crafter should set the data
```
This step will listen the out topic of the tx-crafter, i.e.`topic-tx-nonce`, and checks that envelopes are enriched with `TxData.Data`

```
    Then the tx-nonce should set the nonce
```
This step will listen the out topic of the tx-nonce, i.e.`topic-tx-signer`, and checks that envelopes are enriched with `TxData.Nonce` and make sure that nonce is well incremented for a specific account.

```
    Then the tx-signer should sign
```
This step will listen the out topic of the tx-signer, i.e.`topic-tx-sender`, and checks that envelopes are enriched with `TxData.Nonce` and make sure that nonce is well incremented for a specific account.

```
    Then the tx-sender should send the tx
```
Not implemented yet.

```
    Then the tx-listener should catch the tx
```
This step will listen the out topic of the tx-listener, i.e.`topic-tx-decoder`, and checks that envelopes are enriched with a `Receipt.TxHash`.


```
    Then the tx-decoder should decode
```
This step will listen the out topic of the tx-decoder, i.e.`topic-tx-decoded`, and checks that envelopes are enriched with a `Receipt.Logs.GetDecodedData` if applicable.
