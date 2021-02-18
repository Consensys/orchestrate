#!/bin/bash

# Exit on error
set -Eeu

TX_SENDER=tx-sender
#TX_DECODED=tx-decoded
#TX_RECOVER=tx-recover

echo "Creating topics if not exist..."
RETRY=10
PARTITIONS=5
#TOPICS=($TX_SENDER $TX_DECODED $TX_RECOVER)
TOPICS=($TX_SENDER)

for NAME in ${TOPICS[@]}
do
    for i in $(seq 1 1 $RETRY)
    do
	    docker-compose -f scripts/deps/docker-compose.yml exec kafka kafka-topics --create --bootstrap-server kafka:9092 --partitions $PARTITIONS --replication-factor 1 --topic topic-$NAME && break
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
