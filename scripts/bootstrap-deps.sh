#!/usr/bin/env bash

# Exit on error
set -Eeu

CONTAINERS=(deps_kafka_1 deps_vault_1 deps_vault-agent_1)

RETRY=20

for CONTAINER in "${CONTAINERS[@]}"
do
    for i in $(seq 1 1 $RETRY)
    do
        HEALTH=$(docker inspect --format='{{json .State.Health}}' ${CONTAINER} | jq '.Status')

        if [[ $HEALTH == '"healthy"' ]] ; then
            echo "...${CONTAINER} is ready."
            break
        fi
        echo "Attempt $i/$RETRY (retry in 5 seconds) - $CONTAINER is not ready"

        if [ $i = $RETRY ]; then
            echo "Stopping ..."
            exit 1
        fi
        # Sleep 5 seconds if not succeded
        sleep 5
    done
done

echo "
=======================================================================
Yay! You are ready to Orchestrate :D
=======================================================================
"
