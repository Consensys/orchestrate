#!/usr/bin/env bash

# Exit on error
set -Ee

TOKEN_HEADER="Circle-Token: ${CIRCLECI_TOKEN}"

#Pass parameters to the Circle CI pipeline
PARAMETERS=""
[ "$ORCHESTRATE_NAMESPACE" ] && export PARAMETERS=$PARAMETERS,\"orchestrate_namespace\":\"$ORCHESTRATE_NAMESPACE\"
[ "$ORCHESTRATE_TAG" ] && export PARAMETERS=$PARAMETERS,\"orchestrate_tag\":\"$ORCHESTRATE_TAG\"
[ "$TEST_TAG" ] && export PARAMETERS=$PARAMETERS,\"test_tag\":\"$TEST_TAG\"
[ "$ORCHESTRATE_REPOSITORY" ] && export PARAMETERS=$PARAMETERS,\"orchestrate_repository\":\"$ORCHESTRATE_REPOSITORY\"
[ "$TEST_REPOSITORY" ] && export PARAMETERS=$PARAMETERS,\"test_repository\":\"$TEST_REPOSITORY\"
[ "$REGISTRY_URL" ] && export PARAMETERS=$PARAMETERS,\"registry_url\":\"$REGISTRY_URL\"
[ "$REGISTRY_USERNAME" ] && export PARAMETERS=$PARAMETERS,\"registry_username\":\"$REGISTRY_USERNAME\"
[ "$REGISTRY_PASSWORD" ] && export PARAMETERS=$PARAMETERS,\"registry_password\":\"$REGISTRY_PASSWORD\"
[ "$PARAMETERS" ] && PARAMETERS=${PARAMETERS:1}

echo "Pipeline parameters: $PARAMETERS"

#Create CircleCI pipeline
RESPONSE=$(curl -s --request POST --header "${TOKEN_HEADER}" --header "Content-Type: application/json" --data '{"branch":"'${BRANCH-master}'","parameters":{'${PARAMETERS}'}}' https://circleci.com/api/v2/project/github/ConsenSys/orchestrate-kubernetes/pipeline)
echo $RESPONSE
ID=$(echo $RESPONSE | jq '.id' -r)
NUMBER=$(echo $RESPONSE | jq '.number' -r)

echo "Circle CI pipeline created: $ID"

#Timeout after 4*450 seconds = 30min
SLEEP=4
RETRY=450

for i in $(seq 1 1 $RETRY)
do
    sleep $SLEEP

    # Get pipeline status
    STATUS=$(curl -s --request GET --header "${TOKEN_HEADER}" --header "Content-Type: application/json" https://circleci.com/api/v2/pipeline/${ID}/workflow | jq '.items[0].status' -r)
    echo "$i/$RETRY - $STATUS"

    if [[ $STATUS != 'running' ]] ; then
        break
    fi

    if [ $i = $RETRY ]; then
        echo "Timeout"
    fi
done

echo "Final status: ${STATUS}"
echo "See the pipeline https://app.circleci.com/pipelines/github/ConsenSys/orchestrate-kubernetes"

if [ "$STATUS" != "success" ]; then
  exit 1;
fi
