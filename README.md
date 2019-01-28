# Tx-Nonce

Tx-Nonce is a Core-Stack worker responsible to set transaction nonce.
- consumes messages from an **Apache Kafka** Topic
- uses **Redis** as a distributed cache for nonce values.
