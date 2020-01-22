#!/bin/bash

# Exit on error
set -Eeu

TX_CRAFTER=tx-crafter
TX_NONCE=tx-nonce
TX_SIGNER=tx-signer
TX_SENDER=tx-sender
TX_LISTENER=tx-listener
TX_DECODER=tx-decoder
TX_DECODED=tx-decoded
TX_RECOVER=tx-recover
WALLET_GENERATED=wallet-generated
WALLET_GENERATOR=wallet-generator

echo "Creating topics if not exist..."
RETRY=10
TOPICS=($TX_CRAFTER $TX_NONCE $TX_SIGNER $TX_SENDER $TX_DECODER $TX_DECODED $TX_RECOVER $WALLET_GENERATED $WALLET_GENERATOR)

for NAME in ${TOPICS[@]}
do
    # Retry 10 times if could not create topic
    n=10
    for i in $(seq 1 1 $RETRY)
    do
	    docker-compose -f scripts/deps/docker-compose.yml exec kafka kafka-topics --create --partitions 1 --replication-factor 1 --if-not-exists --zookeeper zookeeper:32181 --topic topic-$NAME && break
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
