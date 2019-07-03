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
$ docker-compose up
```