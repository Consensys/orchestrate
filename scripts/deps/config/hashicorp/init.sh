# Store root token in a file so it can be shared with other services through volume
mkdir -p /vault/token
rm -rf /vault/token/*

# Init Vault
curl --request POST --data '{"secret_shares": 1, "secret_threshold": 1}' ${VAULT_ADDR}/v1/sys/init >init.json

# Retrieve root token and unseal key
VAULT_TOKEN=$(cat init.json | jq .root_token | tr -d '"')
UNSEAL_KEY=$(cat init.json | jq .keys | jq .[0])
SHA256SUM=$(cat SHA256SUM)
rm init.json

echo "Vault token ${VAULT_TOKEN}"
echo $VAULT_TOKEN >/vault/token/.root
echo "ROOT_TOKEN: $VAULT_TOKEN"
echo "SHA256SUM: ${SHA256SUM}"

# Unseal Vault
curl --request POST --data '{"key": '${UNSEAL_KEY}'}' ${VAULT_ADDR}/v1/sys/unseal

# Enable secret engine
curl --header "X-Vault-Token: ${VAULT_TOKEN}" --request POST \
  --data '{"type": "kv-v2", "config": {"force_no_cache": true} }' \
  ${VAULT_ADDR}/v1/sys/mounts/secret

# Register orchestrate plugin
curl --header "X-Vault-Token: ${VAULT_TOKEN}" --request POST \
  --data "{\"sha256\": \"${SHA256SUM}\", \"command\": \"orchestrate\" }" \
  ${VAULT_ADDR}/v1/sys/plugins/catalog/secret/orchestrate

# Enable orchestrate secret engine
curl --header "X-Vault-Token: ${VAULT_TOKEN}" --request POST \
  --data '{"type": "plugin", "plugin_name": "orchestrate", "config": {"force_no_cache": true, "passthrough_request_headers": ["X-Vault-Namespace"]} }' \
  ${VAULT_ADDR}/v1/sys/mounts/orchestrate

# Enable role policies
# Instructions taken from https://learn.hashicorp.com/tutorials/vault/getting-started-apis
curl --header "X-Vault-Token: ${VAULT_TOKEN}" --request POST \
  --data '{"type": "approle"}' \
  ${VAULT_ADDR}/v1/sys/auth/approle

curl --header "X-Vault-Token: $VAULT_TOKEN" \
  --request PUT \
  --data '{"policy":" # Dev servers have version 2 of KV secrets engine mounted by default, so will\n# need these paths to grant permissions:\npath \"secret/*\" {\n   capabilities = [\"create\", \"read\", \"update\", \"delete\", \"list\"]\n}\n"}' \
  ${VAULT_ADDR}/v1/sys/policies/acl/allow_secrets

curl --header "X-Vault-Token: $VAULT_TOKEN" \
  --request POST \
  --data '{"policies": ["allow_secrets"]}' \
  ${VAULT_ADDR}/v1/auth/approle/role/tx_signer

curl --header "X-Vault-Token: $VAULT_TOKEN" \
  ${VAULT_ADDR}/v1/auth/approle/role/tx_signer/role-id >role.json
ROLE_ID=$(cat role.json | jq .data.role_id | tr -d '"')
echo $ROLE_ID >/vault/token/role

curl --header "X-Vault-Token: $VAULT_TOKEN" \
  --request POST \
  ${VAULT_ADDR}/v1/auth/approle/role/tx_signer/secret-id >secret.json
SECRET_ID=$(cat secret.json | jq .data.secret_id | tr -d '"')
echo $SECRET_ID >/vault/token/secret
