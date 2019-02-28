# Tx-Crafter

Core-Stack is a blockchain Transaction Orchestration system that can operate multiple chains simultaneously.
It provides robust and agnostic mechanisms for transaction crafting, nonce management, transaction signing, transaction receipt listening, transaction receipt decoding, faucet and more.

Core-Stack is a *microservices architecture* composed of APIs & Workers. 
Workers communicate following **publish-subscribe** pattern using *Apache Kafka* as message broker. 
All messages are standardized using *protobuf* and *grpc*.

Core-Stack is Plug & Play, a user only needs to send a business protobuf message to an input topic,
Core-Stack then manages the full lifecycle of the transaction from crafting to caraching to transaction and decoding logs data.

# Goal

Tx-Crafter is a Core-Stack Worker responsible to 

#. Craft transaction payload
#. Set Gas Price for the transaction
#. Set Gas Limit for the transaction
#. Request Faucet crediting

It
- consumes messages from an **Apache Kafka** Topic
