# Tx-Crafter

Core-Stack is a blockchain *Transaction Orchestration* system that can operate multiple chains simultaneously.
It provides robust and agnostic mechanisms for transaction crafting, nonce management, transaction signing, transaction receipt listening, transaction receipt decoding, faucet and more.

Core-Stack is a *microservices architecture* composed of APIs & Workers. 
Workers communicate following **publish-subscribe** pattern using *Apache Kafka* as message broker. 
All messages are standardized using *protobuf* and *grpc*.

Core-Stack is Plug & Play, a user only needs to send a business protobuf message to the input topic,
Core-Stack then manages the full lifecycle of the transaction from crafting the transaction to decoding event logs data.

## Goal

Tx-Crafter is a Core-Stack worker responsible to 

- *Craft transaction payload*  
- *Set Gas Price* of the transaction
- *Set Gas Limit* of the transaction
- *Request Faucet crediting*

It consumes message from *tx crafting* Kafka topic and publish to *tx nonce* topic.


## High Level Architecture

Tx-Crafter expect all consumed messages to respect [Core-Stack standard protobuf format](https://gitlab.com/ConsenSys/client/fr/core-stack/core/blob/master/protobuf/trace/trace.proto)

Consumed messages should have 

- **Chain** entry set
- **Call** entry set or Tx-Worker will consider the transaction as a basic Eth transaction with no payload.

1. Crafting

To craft transaction payload Tx-Worker inspects the **Call** entry of input protobuf message 
 
- it expect the **ID** entry formated as **<method>@<contract_name>** (e.g. "transfer@ERC20") (this will evolve to handle versioning of contracts) 
- it expects the **Args** entry to be a list of expected ordered arguments in ```string``` format (e.g. ["0x6009608a02a7a15fd6689d6dad560c44e9ab61ff", "0xabced"] for **to** and **value** )

By basing on the **ID**, Tx-Worker requests method ABI from ABI registry, then it casts ```string``` arguments in the expected Solidity type to craft payload.

2. Gas Price

Tx-Crafter interogates the Ethereum client using the identifier in **Chain** by calling jsonRPC ```eth_gasPrice```

3. Gas Cost

Tx-Crafter interogates the Ethereum client using the identifier in **Chain** by calling jsonRPC ```eth_estimateGas```

4. Faucet

It request a credit to *Faucet*

