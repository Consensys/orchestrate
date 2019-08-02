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
go run . run --jaeger-service TX-DECODER
```

**5. Run producer**

```bash
 go run e2e/producer/main.go
```
## High Level Architecture

Tx-Decoder expects all consumed messages to respect a specific CoreStack protobuf format.

Consumed messages should have:

- ```Chain``` entry set;
- ```Receipt``` entry set with TxHash, Topics and Data fields.

**1. Find Event in ABI**

To decode the raw logs from the blockchain it requires to know the ABI of the event as the arguments are packed in ```Log.Data``` and ```Log.Topics``` without knowing the type expected.

This is why Tx-Decoder is loading the ```ABI``` of interest to decode logs that could be identified in ```Log.Topics[0]```, correspondig to the signature of the event.
 
**2. Decoding**

Once the event is identified, the Tx-Decoder knows exactly the arguments to decode, i.e. their types and which of them are indexded/non-indexed, and could seamlessly decode the raw logs by the following:

- First, it will unpack values from ```Log.Data``` that contains every non-indexed arguments of the event and will return a slice of abstract type `interface{}`.
- Second, as the ```unpackValues``` and ```Log.Topics``` should be in the same order as the event arguments are ordered, the Tx-Decoder will loop through the expected event arguments and pick values from ```unpackValues``` for non-indexed argments and from ```Log.Topics``` for indexed arguments. For non-indexed values, the method in ```core-stack.infra.ethereum.FormatNonIndexedArg``` will transform `interface{}` into string, wheareas for indexed values the method in ```core-stack.infra.ethereum.FormatIndexedArg``` will transform Hash type into string.
- Finally, every arguents strings are mapped and integrated in the Trace context ```ctx.T.Receipt.Logs[i].Decoded```

***Note**: that the Tx-decoder is also decoding any kind of array and will return a string encapsulated into square brakets and delimited by a comma.*
