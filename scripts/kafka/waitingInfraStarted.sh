#!/bin/bash

# Exit on error
set -Eeu

source .env

echo "Waiting kafka starting..."
RETRY=10

for ENDPOINT in ${ETH_CLIENT_URL-http://localhost:8545}
do
    # Retry 10 times if could not call blockchain endpoint
    for i in $(seq 1 1 $RETRY)
    do
        CHAINID=$(curl -H "Content-Type: application/json" -X POST --data '{"jsonrpc":"2.0","method":"net_version","params":[],"id":1}' $ENDPOINT | jq .result | jq tonumber)
        # Check if CHAINID is an number
        re='^[0-9]+$'
        if [[ $CHAINID =~ $re ]] ; then
            break
        fi
        echo "
=======================================================================
Attempt $i/$RETRY (retry in 2 seconds) - could not retrieve chain id from $ENDPOINT
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
echo "...Eth Client started."

# Retry 10 times if could not reach topic list
for i in $(seq 1 1 $RETRY)
do
    docker-compose exec kafka kafka-topics --list --bootstrap-server kafka:9092 && break
    echo "
=======================================================================
Attempt $i/$RETRY (retry in 2 seconds) - could not list topic
=======================================================================
        "
    if [ $i = $RETRY ]; then
        echo "Stopping ..."
        exit
    fi
    # Sleep 2 seconds if not succeded
    sleep 2
done
echo "...kafka started."


topic-tx-crafter topic-tx-nonce topic-tx-signer topic-tx-sender topic-tx-decoded topic-tx-recover
