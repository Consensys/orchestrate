#!/bin/bash

# Exit on error
set -Eeu

echo "Creating topics if not exist..."

TOPICS=(tx-sender)
TX_DECODER=tx-decoder

for ENDPOINT in ${ETH_CLIENT_URL-http://localhost:8545}
do
    TOPICS+=($TX_DECODER-$(curl -H "Content-Type: application/json" -X POST --data '{"jsonrpc":"2.0","method":"net_version","params":[],"id":1}' $ENDPOINT | jq .result | jq tonumber))
done

for NAME in ${TOPICS[@]}
do
	docker exec -it e2e_kafka_1 kafka-topics --create --partitions 1 --replication-factor 1 --if-not-exists --zookeeper zookeeper:32181 --topic topic-$NAME
done
echo "...topics created."