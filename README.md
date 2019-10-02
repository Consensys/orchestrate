# Package

Core-Stack is a blockchain **Transaction Orchestration** system that can operate multiple chains simultaneously.
It provides production grade and agnostic mechanisms for transaction crafting, nonce management, transaction signing, transaction receipt listening, transaction receipt decoding, faucet and more.

Core-Stack is a **microservices architecture** composed of APIs & Workers. 
Workers communicate following **publish-subscribe** pattern using *Apache Kafka* as message broker. 
All messages are standardized using **protobuf** and **gRPC**.

Core-Stack is Plug & Play, a user only needs to send a business protobuf message to Core-Stack input topic,
Core-Stack then manages the full lifecycle of the transaction from crafting the transaction to decoding event logs data.

## Goal

Package is a low level library in Core-Stack dependency tree. In particular it implements

- *protobuf* containing all protobuf schemes
- *core* that defines core structural elements of Core-Stack (such as ``types.Context``) 
- *common* which are resources that are shared between multiple Core-Stack services
