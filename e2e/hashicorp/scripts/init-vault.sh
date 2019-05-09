# INIT the vault
curl --request POST --data '{"secret_shares": 1, "secret_threshold": 1}' ${VAULT_ADDR}/v1/sys/init > init.json

# UNSEAL the vault
export VAULT_UNSEAL_KEY=$(cat init.json | jq .keys | jq .[0])
curl --request POST --data '{"key": '${VAULT_UNSEAL_KEY}'}' ${VAULT_ADDR}/v1/sys/unseal

# Set the VAULT_TOKEN as environment variable
export VAULT_TOKEN=$(cat init.json | jq .root_token | tr -d '"')

# Enable secret engine
curl --header "X-Vault-Token: ${VAULT_TOKEN}" --request POST \
     --data '{"type": "kv-v2", "config": {"force_no_cache": true} }' \
    ${VAULT_ADDR}/v1/sys/mounts/secret

rm init.json