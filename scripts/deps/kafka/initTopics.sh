#!/bin/bash

# Exit on error
set -Eeu

TX_SENDER=tx-sender
TX_DECODED=tx-decoded
TX_RECOVER=tx-recover

echo "Creating topics if not exist..."
RETRY=10
PARTITIONS=2
TOPICS=($TX_SENDER $TX_DECODED $TX_RECOVER)

for NAME in ${TOPICS[@]}
do
    for i in $(seq 1 1 $RETRY)
    do
	    docker-compose -f scripts/deps/docker-compose.yml exec kafka kafka-topics --create --partitions $PARTITIONS --replication-factor 1 --if-not-exists --zookeeper zookeeper:32181 --topic topic-$NAME && break
        echo "
=======================================================================
Attempt $i/$RETRY (retry in 2 seconds) - could not create topic-$NAME
=======================================================================
        "
        if [ $i = $RETRY ]; then
            echo "Stopping ..."
            exit
        fi
        # Sleep 2 seconds if not succeded
        sleep 2
    done
done

echo "...topics created."
