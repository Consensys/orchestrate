# Package

Orchestrate is a blockchain **Transaction Orchestration** system that can operate multiple chains simultaneously.
It provides production grade and agnostic mechanisms for transaction crafting, nonce management, transaction signing, transaction receipt listening, transaction receipt decoding, faucet and more.

Orchestrate is a **microservices architecture** composed of APIs & Workers. 
Workers communicate following **publish-subscribe** pattern using *Apache Kafka* as message broker. 
All messages are standardized using **protobuf** and **gRPC**.

Orchestrate is Plug & Play, a user only needs to send a business protobuf message to Orchestrate input topic,
Orchestrate then manages the full lifecycle of the transaction from crafting the transaction to decoding event logs data.

## Goal

Package is a low level library in Orchestrate dependency tree. In particular it implements

- *protobuf* containing all protobuf schemes
- *core* that defines core structural elements of Orchestrate (such as ``types.Context``) 
- *common* which are resources that are shared between multiple Orchestrate services
