# Tx-Nonce

Tx-Nonce is a Core-Stack worker responsible to set transaction nonce.
- consumes messages from an **Apache Kafka** Topic
- uses **Redis** as a distributed cache for nonce values.

## Running e2e tests

1. Run testing environment

```bash
docker-compose -f e2e/docker-compose.yml up
```

2. Run worker

```bash
REDIS_ADDRESS=localhost:6379 REDIS_LOCKTIMEOUT=1500 ETH_CLIENT_URL=http://localhost:8545 go run . run
```

3. Run producer that will write messages 

```bash
 go run e2e/producer/main.go
```