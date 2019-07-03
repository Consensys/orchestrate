# Transaction Signer configuration

## Goal

Tx-Signer is a CoreStack worker responsible for: 

- **Signing transaction**

It consumes messages from *tx signer* Kafka topic and publishs to *tx sender* topic.

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

