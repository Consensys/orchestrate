# Hard restart the vault if necessary
#docker-compose -f ../docker-compose.yml down | true
#docker-compose -f ../docker-compose.yml up -d --force-recreate
apk add jq curl

# Unseal vault and the propers environment variables
source init-vault.sh
source token.sh

