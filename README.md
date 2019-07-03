# Transaction Decoder configuration

## Goal

Tx-Decoder is a CoreStack worker responsible for **Decoding raw events logs from transactions into a human readable mapping of strings** 

It consumes messages from *tx decoder* Kafka topic and publishs to *tx decoded* topic.

## Quick-Start

### Prerequisites

Having *docker* and *docker-compose* installed.

### Start the application

You can start the application with default configuration by running.

```sh
$ docker-compose up
```

### Run e2e tests

**1. Set "ETH_CLIENT_URL"**

```bash
export ETH_CLIENT_URL="https://ropsten.infura.io/v3/81e039ce6c8a465180822b525e3644d7 https://rinkeby.infura.io/v3/bfc9d6e51fbc4d3db54bea58d1094f9c https://kovan.infura.io/v3/bfc9d6e51fbc4d3db54bea58d1094f9c https://mainnet.infura.io/v3/bfc9d6e51fbc4d3db54bea58d1094f9c"
```

**2. Run Kafka and Zookeeper**

```bash
docker-compose -f e2e/docker-compose.yml up
```

**3. Create Kafka topics**

```bash
bash e2e/initTestTopics.sh 
```

This script will fetch topic ids using Ethereum JSON RPC and create necessary topics

**4. Run worker**

```bash
go run . run
```

**5. Run producer**

```bash
 go run e2e/producer/main.go
```
