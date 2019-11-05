# Create a token
curl --header "X-Vault-Token: ${ROOT_TOKEN}" \
     --request POST \
     --data '{"ttl": "5m", "renewable": true}' \
    ${VAULT_ADDR}/v1/auth/token/create > token.json

token=$(cat token.json | jq .auth | jq .client_token | tr -d '"')

mkdir -p /vault/token 
touch /vault/token/.vault-token
echo $token > /vault/token/.vault-token

rm token.json
