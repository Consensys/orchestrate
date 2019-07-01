# Tx-Decoder

## About CoreStack

**CoreStack** is a blockchain *Transaction Orchestration System* that abstracts blockchain complexity and can operate multiple chains simultaneously. It provides production grade and agnostic mechanisms for transaction crafting, nonce management, transaction signing, transaction receipt listening, transaction receipt decoding, faucet and more.

**CoreStack** is a Plug & Play component, as a user only needs to send a business protobuf message to **CoreStack** input topic. **CoreStack** then manages the full lifecycle of the transaction from crafting the transaction to decoding event logs data.

**CoreStack** target users are Solutions Architects, Developers, Integrators and Operations Engineers developing, deploying and scaling a blockchain solution.


## Goal

Tx-Decoder is a Core-Stack worker responsible to 

- **Decode raw events logs from transactions into a human readable mapping of strings** 

It consumes message from *tx decoder* Kafka topic and publish to *tx decoded* topic.

