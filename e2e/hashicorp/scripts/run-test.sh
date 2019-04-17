# Hard restart the vault if necessary
docker-compose -f ../docker-compose.yml down | true
docker-compose -f ../docker-compose.yml up -d --force-recreate
sleep 2

# Unseal vault and the propers environment variables
source init-vault.sh

# Run the tests passing the variables to the go command
VAULT_TOKEN=${VAULT_TOKEN} VAULT_ADDR=${VAULT_ADDR} \
    go test ../../../...

# Cleanly shut down the vault container
docker-compose -f ../docker-compose.yml down
