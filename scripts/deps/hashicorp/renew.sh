ROLE_ID=$(cat $TOKEN_PATH/role)
SECRET_ID=$(cat $TOKEN_PATH/secret)
echo "ROLE_ID: $ROLE_ID"
echo "SECRET_ID: $SECRET_ID"

while true; do 
  curl --request POST \
    --data "{\"role_id\": \"$ROLE_ID\", \"secret_id\": \"${SECRET_ID}\"}" \
    ${VAULT_ADDR}/v1/auth/approle/login > login.json
    
  ORCHESTRATE_CLI_TOKEN=$(cat login.json | jq .auth.client_token | tr -d '"')
  echo "New client token: $ORCHESTRATE_CLI_TOKEN"
  echo $ORCHESTRATE_CLI_TOKEN > $TOKEN_FILE_PATH
  sleep 10
done
