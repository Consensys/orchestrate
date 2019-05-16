# E2E Tests

Core-Stack is a blockchain **Transaction Orchestration** system that can operate multiple chains simultaneously.
It provides production grade and agnostic mechanisms for transaction crafting, nonce management, transaction signing, transaction receipt listening, transaction receipt decoding, faucet and more.

Core-Stack is a **microservices architecture** composed of APIs & Workers. 
Workers communicate following **publish-subscribe** pattern using *Apache Kafka* as message broker. 
All messages are standardized using **protobuf** and **GRPC**.

Core-Stack is Plug & Play, a user only needs to send a business protobuf message to Core-Stack input topic,
Core-Stack then manages the full lifecycle of the transaction from crafting the transaction to decoding event logs data.
