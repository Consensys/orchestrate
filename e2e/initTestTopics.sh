#!/bin/bash

# Exit on error
set -Eeu

echo "Creating topics if not exist..."

TOPICS=(tx-sender)

for NAME in ${TOPICS[@]}
do
	docker exec -it e2e_kafka_1 kafka-topics --create --partitions 1 --replication-factor 1 --if-not-exists --zookeeper zookeeper:32181 --topic topic-$NAME
done
echo "...topics created."