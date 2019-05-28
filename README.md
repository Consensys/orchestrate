# Tx-Signer

Core-Stack is a blockchain **Transaction Orchestration** system that can operate multiple chains simultaneously.
It provides production grade and agnostic mechanisms for transaction crafting, nonce management, transaction signing, transaction receipt listening, transaction receipt decoding, faucet and more.

Core-Stack is a **microservices architecture** composed of APIs & Workers. 
Workers communicate following **publish-subscribe** pattern using *Apache Kafka* as message broker. 
All messages are standardized using **protobuf** and **GRPC**.

Core-Stack is Plug & Play, a user only needs to send a business protobuf message to Core-Stack input topic,
Core-Stack then manages the full lifecycle of the transaction from crafting the transaction to decoding event logs data.

## Goal

Tx-Signer is a Core-Stack worker responsible to 

- **Sign transaction**

It consumes message from *tx signer* Kafka topic and publish to *tx sender* topic.

## Quick-Start

### Prerequisite

- Having ```docker``` and ```docker-compose``` installed
- Having Go 1.11 installed or upper

### Start the application

To quickly start the application

1. Start e2e env

```sh
$ docker-compose -f e2e/docker-compose.yml up
```

2. Start worker

```sh
$ go run . run
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
      --engine-slots uint               Maximum number of messages the engine can treat concurrently.
                                        Environment variable: "ENGINE_SLOTS" (default 20)
  -h, --help                            help for run
      --http-hostname string            Hostname to expose HTTP server
                                        Environment variable: "HTTP_HOSTNAME" (default ":8080")
      --kafka-address strings           Address of Kafka server to connect to.
                                        Environment variable: "KAFKA_ADDRESS" (default [localhost:9092])
      --kafka-group string              Address of Kafka server to connect to.
                                        Environment variable: "KAFKA_GROUP" (default "group-e2e")
      --secret-pkey strings             Private keys to pre-register in key store
                                        Environment variable: "SECRET_PKEY" (default [56202652FDFFD802B7252A456DBD8F3ECC0352BBDE76C23B40AFE8AEBD714E2E,5FBB50BFF6DFAD35C4A374C9237BA2F7EAED9C6868E0108CB259B62D68029B1A,86B021CCB810F26A30445B85F71E4C1596A11A97DDF9B9E348AC93D1DA6735BC,DD614C3B343E1B6DBD1B2811D4F146CC90337DEEF96AB97C353578E871B19D5E,425D92F63A836F890F1690B34B6A25C2971EF8D035CD8EA8592FD1069BD151C6,C4B172E72033581BC41C36FA0448FCF031E9A31C4A3E300E541802DFB7248307,706CC0876DA4D52B6DCE6F5A0FF210AEFCD51DE9F9CFE7D1BF7B385C82A06B8C,1476C66DE79A57E8AB4CADCECCBE858C99E5EDF3BFFEA5404B15322B5421E18C,A2426FE76ECA2AA7852B95A2CE9CC5CC2BC6C05BB98FDA267F2849A7130CF50D,41B9C5E497CFE6A1C641EFCA314FF84D22036D1480AF5EC54558A5EDD2FEAC03])
      --secret-store string             Type of secret store for private keys (one of "test" "hashicorp")
                                        Environment variable: "SECRET_STORE" (default "test")
      --topic-sender string             Kafka topic for messages waiting to have transaction sent
                                        Environment variable: "KAFKA_TOPIC_TX_SENDER" (default "topic-tx-sender")
      --topic-signer string             Kafka topic for messages waiting to have transaction signed
                                        Environment variable: "KAFKA_TOPIC_TX_SIGNER" (default "topic-tx-signer")
      --vault-addr string               Hashicorp secret path
                                        Environment variable: "VAULT_ADDR" (default "https://127.0.0.1:8200")
      --vault-burst-limit int           Hashicorp secret path
                                        Environment variable: "VAULT_RATE_LIMIT"
      --vault-cacert string             Hashicorp secret path
                                        Environment variable: "VAULT_CACERT"
      --vault-capath string             Hashicorp secret path
                                        Environment variable: "VAULT_CAPATH"
      --vault-client-cert string        Hashicorp secret path
                                        Environment variable: "VAULT_CLIENT_CERT"
      --vault-client-key string         Hashicorp secret path
                                        Environment variable: "VAULT_CLIENT_KEY"
      --vault-client-timeout duration   Hashicorp secret path
                                        Environment variable: "VAULT_CLIENT_TIMEOUT" (default 1m0s)
      --vault-kv-version string         Determine which version of the kv secret engine we will be using
                                        Can be "v1" or "v2".
                                        Environment variable: "VAULT_KV_VERSION"  (default "v2")
      --vault-max-retries int           Hashicorp secret path
                                        Environment variable: "VAULT_MAX_RETRIES"
      --vault-mount-point string        Specifies the mount point used.
                                        Environment variable: "VAULT_MOUNT_POINT"  (default "secret")
      --vault-rate-limit float          Hashicorp secret path
                                        Environment variable: "VAULT_RATE_LIMIT"
      --vault-secret-path string        Hashicorp secret path
                                        Environment variable: "VAULT_SECRET_PATH" (default "default")
      --vault-skip-verify               Hashicorp secret path
                                        Environment variable: "VAULT_SKIP_VERIFY"
      --vault-tls-server-name string    Hashicorp secret path
                                        Environment variable: "VAULT_TLS_SERVER_NAME"

Global Flags:
      --log-format string   Log formatter (one of ["text" "json"]).
                            Environment variable: "LOG_FORMAT" (default "text")
      --log-level string    Log level (one of ["panic" "fatal" "error" "warn" "info" "debug" "trace"]).
                            Environment variable: "LOG_LEVEL" (default "debug")
```

## High Level Architecture

Tx-Signer expect all consumed messages to respect [Core-Stack standard protobuf format](https://gitlab.com/ConsenSys/client/fr/core-stack/core/blob/master/protobuf)

Consumed messages should have 

- ```Chain``` attribute set with ```ID``` of the chain to send the transaction to
- ```Sender``` attribute set with an ```Address```
- ```Tx``` attribute set with the following fields set up:
- ```Nonce```
- ```To```
- ```Value```
- ```GasLimit```
- ```GasPrice```
- ```Data```


1. **Signing**

To sign the transaction Tx-Signer inspects the ```Tx``` entry of input protobuf message 

- it uses a ```Signer``` (as for now a ```StaticSigner```) which should implement the [```TxSigner``` interface](https://gitlab.com/ConsenSys/client/fr/core-stack/core/blob/master/services/signer.go).
- it updates the ```Tx``` attribute with the following fields:
  - ```Raw```
  - ```Hash```
- it sends the signed transaction into ```tx-sender``` Kafka topic


## Technical Architecture

![alt core-stack-architecture](https://gitlab.com/ConsenSys/client/fr/core-stack/doc/blob/master/diagrams/Core_Stack_Architecture.png)
