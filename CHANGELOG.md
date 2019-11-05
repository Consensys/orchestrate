# CHANGELOG

### Version 1.0.0

*Unreleased*

### BREAKING CHANGES
* **config**Rename `GRPC_TARGET_CONTRACT_REGISTRY` to `CONTRACT_REGISTRY_URL`
* **config**Rename `GRPC_TARGET_ENVELOPE_STORE` to `ENVELOPE_STORE_URL`
* **config**Rename `KAFKA_ADDRESS` to `KAFKA_URL`
* **config**Rename `REDIS_ADDRESS` to `REDIS_URL`
* **config**Rename `TESSERA_ENDPOINTS` to `TESSERA_URL`
* **config** Rename `mock` options to `in-memory`
* **config** Remove `--engine-slot`
* **config** Rename `FAUCET` to `FAUCET_TYPE` & remove `--faucet`
* **config** Rename `VAULT_TOKEN_FILEPATH` to `VAULT_TOKEN_FILE`
* **config** Rename `KAFKA_TOPIC_TX_XXX` to `TOPIC_TX_XXXX`
* **config** Rename `KAFKA_TLS_CLIENTCERTFILEPATH` to `KAFKA_TLS_CLIENT_CERT_FILE`
* **config** Rename `KAFKA_TLS_CLIENTKEYFILEPATH` to `KAFKA_TLS_CLIENT_KEY_FILE`
* **config** Rename `KAFKA_TLS_CACERTFILEPATH` to `KAFKA_TLS_CA_CERT_FILE`
* **config** Rename `KAFKA_TLS_INSECURESKIPVERIFY` to `KAFKA_TLS_INSECURE_SKIP_VERIFY`
* **config** Rename `KAFKA_TLS_ENABLE` to `KAFKA_TLS_ENABLED`
* **config** Rename `NONCE_MANAGER` to `NONCE_MANAGER_TYPE`
* **config** Rename `CONTRACT_REGISTRY` to `CONTRACT_REGISTRY_TYPE`
* **config** Rename `ENVELOPE_STORE` to `ENVELOPE_STORE_TYPE`
* **config** Rename Envelope-store option `pg` to `postgres`
