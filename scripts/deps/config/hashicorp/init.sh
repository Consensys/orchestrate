# Init Vault
    curl --request POST --data '{"secret_shares": 1, "secret_threshold": 1}' ${VAULT_ADDR}/v1/sys/init > init.json

    # UNSEAL the vault
    export VAULT_UNSEAL_KEY=$(cat init.json | jq .keys | jq .[0])
    curl --request POST --data '{"key": '${VAULT_UNSEAL_KEY}'}' ${VAULT_ADDR}/v1/sys/unseal

    # Set the ROOT_TOKEN as environment variable
    export ROOT_TOKEN=$(cat init.json | jq .root_token | tr -d '"')

    # Enable secret engine
    curl --header "X-Vault-Token: ${ROOT_TOKEN}" --request POST \
         --data '{"type": "kv-v2", "config": {"force_no_cache": true} }' \
        ${VAULT_ADDR}/v1/sys/mounts/secret

    rm init.json

# Get Token
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