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
go run . run
```

4. Run producer that will write messages 

```bash
 go run e2e/producer/main.go
```
