# Hard restart the vault if necessary
#docker-compose -f ../docker-compose.yml down | true
#docker-compose -f ../docker-compose.yml up -d --force-recreate
sleep 2

# Unseal vault and the propers environment variables
source init-vault.sh

# # Go to Makefile
cd ../../..

echo "make race"
# Run the tests passing the variables to the go command
VAULT_TOKEN=${VAULT_TOKEN} VAULT_ADDR=${VAULT_ADDR} \
    make race

echo "run-coverage"
# Run the tests passing the variables to the go command
VAULT_TOKEN=${VAULT_TOKEN} VAULT_ADDR=${VAULT_ADDR} \
    make run-coverage


# Cleanly shut down the vault container
#docker-compose -f ../docker-compose.yml down
