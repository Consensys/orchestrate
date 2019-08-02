# Tx-Crafter

## Goal

Tx-Crafter is a CoreStack worker responsible for:

- **Crafting the transaction payload**;  
- **Setting Gas Price** of the transaction;
- **Setting Gas Limit** of the transaction;
- **Requesting Faucet crediting**.

It consumes messages from *tx crafting* Kafka topic and publishs to *tx nonce* topic.

## Quick-Start

### Prerequisites

- Having ```docker``` and ```docker-compose``` installed;
- Having Go 1.11 installed or upper.

### Start the application

To quickly start the application

1. Start Kafka broker, ganache and jaeger

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
go run . run  --eth-client http://localhost:8545 --jaeger-service TX-CRAFTER
```

4. Run producer that will write messages 

```bash
 go run e2e/producer/main.go
```


## High Level Architecture

Tx-Crafter expects all consumed messages to respect a specific CoreStack protobuf format.

Consumed messages should have:

- ```Chain``` attribute set with ```ID``` of the chain to send the transaction to;
- ```Call``` attribute set, or Tx-Worker will consider the transaction as a basic Ethereum transaction with no payload and will skip crafting.


**1. Crafting**
To craft transaction payload, Tx-Worker inspects the ```Call``` entry of input protobuf message:

- it expects ```Call.ID``` to be a string formatted as ```<method>@<contract_name>``` (e.g. in case of an ERC20 transfer: "transfer@ERC20") (Note: this will evolve to handle versioning of contracts).
- it expects ```Call.Args``` entry to be an ordered list of arguments to provide to the transaction call in ```string``` format (e.g. in case of an ERC20 transfer: ["0x6009608a02a7a15fd6689d6dad560c44e9ab61ff", "0x34fde"] for *to* and *value* args)

By reading the ```Call.ID```, Tx-Worker requests the required ABI from the *Contract registry*, then it casts ```Call.Args``` arguments in the expected Solidity type and crafts the payload.

**2. Gas Price**
Tx-Crafter interrogates *Ethereum client* by calling jsonRPC ```eth_gasPrice``` on chain ```Chain.ID```.

**3. Gas Cost**
Tx-Crafter interrogates *Ethereum client* by calling jsonRPC ```eth_estimateGas``` on chain ```Chain.ID```.

**4. Faucet**
It requests a credit to *Faucet*.
