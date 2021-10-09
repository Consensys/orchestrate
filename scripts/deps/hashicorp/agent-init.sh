until [ -f ${VAULT_IS_READY} ]; do
  echo "[AGENT] waiting for vault to be ready..."
  sleep 1
done

VAULT_TOKEN=$(cat "${ROOT_TOKEN_PATH}")
        
curl -s --header "X-Vault-Token: ${VAULT_TOKEN}" --request POST \
  --data '{"type": "approle"}' \
  ${VAULT_ADDR}/v1/sys/auth/approle

curl -s --header "X-Vault-Token: $VAULT_TOKEN" \
  --request PUT \
  --data '{ "policy":"path \"'"${PLUGIN_MOUNT_PATH}/*"'\" { capabilities = [\"create\", \"read\", \"update\", \"delete\", \"list\"] }" }' \
  ${VAULT_ADDR}/v1/sys/policies/acl/allow_secrets

curl -s --header "X-Vault-Token: $VAULT_TOKEN" \
  --request POST \
  --data '{"policies": ["allow_secrets"]}' \
  ${VAULT_ADDR}/v1/auth/approle/role/key-manager

curl -s --header "X-Vault-Token: $VAULT_TOKEN" \
  ${VAULT_ADDR}/v1/auth/approle/role/key-manager/role-id > role.json
ROLE_ID=$(cat role.json | jq .data.role_id | tr -d '"')
echo $ROLE_ID > ${ROLE_FILE_PATH}
rm role.json

curl -s --header "X-Vault-Token: $VAULT_TOKEN" \
  --request POST \
  ${VAULT_ADDR}/v1/auth/approle/role/key-manager/secret-id > secret.json
SECRET_ID=$(cat secret.json | jq .data.secret_id | tr -d '"')
echo $SECRET_ID > ${SECRET_FILE_PATH}
rm secret.json
