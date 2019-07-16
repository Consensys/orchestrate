# Tx-Listener

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

