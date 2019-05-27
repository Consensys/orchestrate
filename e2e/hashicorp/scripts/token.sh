# Create a token
curl --header "X-Vault-Token: ${ROOT_TOKEN}" \
     --request POST \
     --data '{"ttl": "5m", "renewable": true}' \
    ${VAULT_ADDR}/v1/auth/token/create > token.json

token=$(cat token.json | jq .auth | jq .client_token | tr -d '"')

mkdir -p /auth/vault 
touch /auth/vault/token
echo $token > /auth/vault/token

rm token.json