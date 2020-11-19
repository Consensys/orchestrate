# Codefi Orchestrate Release Notes

## v2.6.0 (Unreleased)

### 🆕 Features
* Support for enable/disable metric modules

### ⚠ BREAKING CHANGES

* Remove account-generator and account-generated topic

## v2.5.3 (Unreleased)

### 🛠 Bug fixes

* Fix chain registration issue with Kaleido/Infura when multitenancy is enabled
* Retry on worker messages when connection errors occurred
* Fix missing error communication on edge cases

## v2.5.2 (2020-11-09)

### 🛠 Bug fixes

* Allow usage of OS certificate bundle on Redis TLS connections
* Fix a bug on transaction recovering for invalid nonce 

## v2.5.1 (2020-10-23)

### 🆕 Features

* Add the `CHAIN_REGISTRY_MAXIDLECONNSPERHOST` to control the maximum of open HTTP connections to a chain proxied 

### 🛠 Bug fixes

* Ability to cache gzip HTTP responses in the chain-registry

## v2.5.0 (2020-10-19)

### 🆕 Features

* Enhance service health check endpoint (/ready) to validate external and internal dependencies
* Add support for TLS connection to Redis. Add flags and environment variables:
    * `REDIS_TLS_CERT`: PEM certificate to connect to the database
    * `REDIS_TLS_KEY`: PEM key to connect to the database
    * `REDIS_TLS_CA`: PEM trusted CA that issued certificate
    * `REDIS_TLS_SKIP_VERIFY`: TLS host verification should be enabled/disabled
* Add Redis AUTH credential and database selection. Add flags and environment variables:
    * `REDIS_USERNAME`: Redis AUTH username
    * `REDIS_PASSWORD`: Redis AUTH password
    * `REDIS_DATABASE`: Redis CLI database number

## v2.4.2 (Unreleased)

### 🛠 Bug fixes

* Fix chain registration issue with Kaleido/Infura when multitenancy is enabled

## v2.4.1 (2020-11-09)

### 🛠 Bug fixes
* Remove duplicated tx-recover messages on transaction retries failing on sending

## v2.4.0 (2020-10-19)

### 🆕 Features

* Add the new Transaction Sentry service.
    * Users can now add a `retryPolicy` inside the `gasPricePolicy` settings when publishing transactions by specifying an `interval`, `increment`, and a `limit`.
    * The Transaction Sentry will watch this transaction until it's mined and after each `interval`, will resend the transaction with a gasPrice increased by `increment`(%), capped by `limit` (%).
* Add a Caching mechanism that can be enabled to cache every identical request going from Orchestrate to the Ethereum node. TTL can be configured using the environment variable `CHAIN_REGISTRY_CACHE_TTL`
    * This feature is especially useful for cases where:
        * where multiple chains (belonging to the same tenant or not) calls the same Ethereum node.
        * and using a node with low capabilities or behind a rate limiter (Infura/Kaleido)  

### ⚠ BREAKING CHANGES

* Schedule level was removed from transaction api responses
* Use transaction UUID instead of job UUID as ID for tx-decoded messages

### 🛠 Bug fixes

* Fix a bug when registering an overloaded & rate-limited chain with a configuration starting block to "latest", the synchronization could start from block 0.
* Fix a bug where the process was not failing when DB migration failed
* Improve retry policy on Eth Clients, it's now failing quicker in some cases and more accurately
* Fix a bug with quorum's private transaction signature where the payload is sign with v=37,38 instead of 27,28
* Fix a bug with besu's gas estimation for private transaction on Besu>=1.5.4
* Fix a bug with missing envelope message ID on `account-generated` topic

## v2.3.2 (2020-09-22)

### 🛠 Bug fixes

* Fix a bug when submitting a contract transaction with method arguments containing arrays

## v2.3.1 (2020-09-15)

### 🛠 Bug fixes

* Fix a bug when registering an overloaded & rate-limited chain with a configuration starting block to "latest", the synchronization could start from block 0.
* Fix a bug where the process was not failing when DB migration failed

## v2.3.0 (2020-09-02)

### 🆕 Features

* Add the new `tx-scheduler` API microservice. This new API:
    * replaces the `envelope-store` and serves the same internal purpose
    * is the new API that is used to POST every transactions (they are no longer sent on the tx-crafter Kafka topic).
        * POST `/transactions/deploy-contract`: Creates and sends a new contract deployment. Supports one-time key & private transactions.
        * POST `​/transactions​/send`: Creates and sends a new contract transaction. Supports one-time key & private transactions.
        * POST `/transactions​/send-raw`: Creates and sends a raw transaction
        * POST ​`/transactions​/transfer`: Creates and sends a transfer transaction
    * exposes GET endpoints to fetch transaction details
        * GET ​`/transactions`: Search transaction requests by provided filters
        * GET `​/transactions​/{uuid}`: Fetch a transaction request by uuid
* Added support for Transaction priority. Users can now specify the gas priority of the transaction (`very-high`, `high`, `medium`, `low`, `very-low`), the transaction Gas Price will be adjusted automatically based on network activity.
* Update internal logger to follow [Elastic Common Schema](https://www.elastic.co/guide/en/ecs/current/index.html) when using logs in JSON format
* Private transactions (Tessera/Orion) are now performed in two separate jobs

### 🛠 Bug fixes

* Properly renew HashiCorp client token
* Fix a bug limiting the amount of Ether that can be send to 9.2 ETH
* Tenant wildcard support to access private keys stored in the Secret Storage

### ⚠ BREAKING CHANGES

* `envelope-store` has been removed.

### Migration from v2.2.0

* Remove the envelope-store API, DB and volume and add the transaction-scheduler API, DB and volume. Data from the envelope-store DB doesn't need to be migrated to the new DB. Follow [this diff](https://github.com/PegaSysEng/orchestrate-kubernetes/compare/559bd13ea1dd68faf4e57a826028e1deeea9dfb1...e99443e20049400acf9ba8f33f76e5e661909f9d) to upgrade to the new configuration.
* Update your application to use the [SDK](https://github.com/PegaSysEng/orchestrate-node) v3.1.0. This SDK will now use the REST API of the transaction scheduler to publish transactions instead of using the Kafka queues.

## v2.2.2 (2020-09-15)

### 🛠 Bug fixes

* Fix a bug when registering an overloaded & rate-limited chain with a configuration starting block to "latest", the synchronization could start from block 0.
* Fix a bug where the process was not failing when DB migration failed

## v2.2.1 (2020-08-31)

### 🛠 Bug fixes

* Properly renew HashiCorp client token
* Tenant wildcard support to access private keys stored in the Secret Storage

## v2.2.0 (2020-07-15)

### 🆕 Features

* Add support for 4 configuration modes for TLS connection to Postgres databases. Add flag and environment variable `DB_TLS_SSLMODE` that can be: `disable`, `require`, `verify-ca`, `verify-full`.
* Add support for wildcard authentication, allowing operators to perform any API (especially useful for chains) operations by providing both:
    * a JWT with a tenant_id="*"
    * an HTTP header containing the targeted tenant_id

### 🛠 Bug fixes

* Fix a bug of nonce management when registering multiple chains of the same network but using an identical account for transactions

## v2.1.1 (2020-05-27)

### 🛠 Bug fixes

* Fix a bug where the tx-listener fails to reach the transaction scheduler MS if it is not deployed.

## v2.1.0 (2020-06-05)

### 🆕 Features

* Add support for Quorum+Tessera private transactions by registering the Tessera node to the Chain Registry. Includes sending and listening of transactions
* Add support for Besu+Orion private transactions. Includes sending of transactions and listening of public & private receipts.
* Add support for Revert Reason when fetching receipt from Besu nodes.
* Optimize receipt fetching from the chain when external transactions are disabled.
* Add support for One-time key signature.
* Add support for TLS connection to Postgres. Add flags and environment variables:
    * `DB_TLS_CERT`: PEM certificate to connect to the database
    * `DB_TLS_KEY`: PEM key to connect to the database
    * `DB_TLS_CA`: PEM trusted CA that issued certificate

### 🛠 Bug fixes

* Add chain information into SDK tx-response
* Fix a casting issue on indexed strings in the events decoded by the tx-listener
* Properly exit workers when a critical failure happens
* Fix a bug that may lead to skipping a block if an error occurred while processing the block
* Fix a bug when the node thrown a `missing trie node` leading to and error in the tx-listener
* Reactivate the metrics/liveness/readiness endpoint on the tx-listener
* Fix a bug when listening sessions stopped in the tx-listener when the HTTP call to the node failed.

## v2.0.2 (2020-04-07)

### 🛠 Bug fixes

* Fix a transaction crafting issue to able to craft uint and int arguments in decimal or hexadecimal
* Fix a contract registry issue to be able to register a contract without event

## v2.0.1 (2020-04-01)

### 🛠 Bug fixes

* Fix a security issue where, if a user is authenticated, a transaction could be sent with any tenant.
* Enable one way TLS communication to Kafka to allow connection to Azure Event Hub

## v2.0.0 (2020-03-11)

### 🆕 Multi-tenancy & JWT Authentication

* Add handler into `tx-crafter`, `tx-decoder`, `tx-nonce`, `tx-sender`, `tx-signer`
    * Authenticate (Verify and Validate) the Envelope using the ID/Access Token (JWT) present in the Metadata
    * Extract the tenantID from the ID/Access Token (JWT)
* Add Interceptor into gRPC API into `contract-registry` and `envelope-store`
    * Authenticate (Verify and Validate) the request with the ID/Access Token (JWT) present in the HTTP Header
    * Extract the tenantID from the ID/Access Token (JWT)
* Store private keys based on the tenantID and the address of the keys
* Add flag and environment variable:
    * `MULTI_TENANCY_ENABLED` to enable multi-tenancy.
    * `AUTH_JWT_CERTIFICATE` to provision trusted certificate to control signature of ID / Access Token (JWT)
    * `AUTH_JWT_CLAIMS_NAMESPACE` to provision the namespace to retrieve Orchestrate AUth element in OpenId or Access Token (JWT) (in particular multitenancy information)
    * `AUTH_API_KEY` secret allowing to bypass JWT authentication (useful for some microservice to microservice communications)

### 🆕 Chain-registry and tx-listener

* Add the chain-registry microservice that:
    * Serves an API to store a list of ethereum chains with their configurations (URLs, tenantID, name, block depth, block position, backoff duration). The API allows to dynamically update chains configuration instead of passing them in environment variable.
        * GET `/chains`: get the list of chains registered
        * GET `/chains/{uuid}`: get the chain configuration given by the {uuid}
        * GET `/chains/{tenantID}`: get the list of chains registered for a given {tenantID}
        * GET `/chains/{tenantID}/{name}`: get the chain configuration given by the {tenantID} and {name}
        * POST `/chains/{tenantID}`: create a new chain for a given {tenantID}
        * PATCH `/chains/{uuid}`: modify the chain configuration for a given {uuid}
        * PATCH `/chains/{tenantID}/{name}`: modify the chain configuration for a given {tenantID} and {name}
        * DELETE `/chains/{uuid}`: delete the chain configuration given by the {uuid}
        * DELETE `/chains/{tenantID}/{name}`: delete the chain configuration given by the {tenantID} and {name}
    * Serves an API to store a list of faucets rules their configurations (a creditor account, max balance, amount, cooldown). The API allows to dynamically update faucets configuration instead of passing them in environment variable.
        * GET `/faucets`: get the list of faucets registered
        * GET `/faucets/{uuid}`: get the faucet configuration given by the {uuid}
        * POST `/faucets`: create a new faucet for a given {tenantID}
        * PATCH `/faucets/{uuid}`: modify the faucets configuration for a given {uuid}
        * DELETE `/faucets/{uuid}`: delete the faucets configuration given by the {uuid}
    * Proxy the ethereum chains
* Add flag and environment variable:
    * `CHAIN_REGISTRY_INIT` to initialize the chain registry with some specific chains
    * `CHAIN_REGISTRY_PROVIDER_CHAINS_REFRESH_INTERVAL` to set the time interval for refreshing the list of chains from storage
    * `CHAIN_REGISTRY_URL` to set the URL to reach the chain-registry
    * `TX_LISTENER_PROVIDER_REFRESH_INTERVAL` to set the time interval for refreshing the list of chains from the chain registry
* Add rate limiter on chain registry to avoid burst traffic on underlying chains (in particular when using Infura or Kaleido)

### ⚠ BREAKING CHANGES

#### Infrastructure

* Merge the `tx-decoder` microservice into `tx-listener` microservice. The `tx-listener` publishes transactions directly in the `topic-tx-decoded`
* The `tx-listener` produces kafka messages exclusively in the topic `topic-tx-decoded` instead of the `topic-tx-decoder-{chainID}`
* Merge the `tx-nonce` microservice into `tx-crafter` microservice. The `tx-crafter` publishes transactions directly in the `topic-tx-signer`
* All microservices, now, have to go through the `chain-registry` microservice to communicate with any Blockchain

#### Configuration

* Rename the default topic names from `topic-wallet-generator` and `topic-wallet-generated` to `topic-account-generator` and `topic-account-generated` respectively
* Move environment variables `NONCE_MANAGER_TYPE` `REDIS_URL` `REDIS_LOCKTIMEOUT` from the `tx-nonce` to the `tx-crafter`
* Remove environment variable `ETH_CLIENT_URL`, the chains urls have to be set at start-up in `CHAIN_REGISTRY_INIT` of `chain-registry` microservice or dynamically using the chain-registry API.
* Remove environment variables `FAUCET_CREDIT_AMOUNT`, `FAUCET_BLACKLIST`, `FAUCET_COOLDOWN_TIME`, `FAUCET_CREDITOR_ADDRESS`, the faucets configurations are stored in the chain registry using its API.
* Add the environment variable `CHAIN_REGISTRY_URL` to the `tx-listener`, `tx-crafter`, `tx-sender`
* Remove environment variable `DISABLE_EXTERNAL_TX` in the `tx-listener` and `tx-decoder`. The same feature can be found in the Chain-Registry API

#### API

* Remove `/v1` prefix in the HTTP REST path for the `envelope-store` and the `chain-registry`
* Instead of producing and consuming envelopes to Orchestrate, a user will produce `TxRequest` and only consume `TxResponse`

## v1.2.2 (2020-01-09)

### 🛠 Bug fixes

* Fix incorrect filtering on "name" argument on the GetTags method of the Contract Registry

## v1.2.1 (2019-12-23)

### 🛠 Bug fixes

* Upgrade retry policy when getting `NotFoundError` on JSON-RPC request. In particular it allows the transaction listener to effectively handle Infura endpoint that have sync discrepancies.

## v1.2.0 (2019-12-13)

### 🆕 Features

* Add new flag and environment variable `REDIS_EXPIRATION` to configure Redis entry expiration. It is useful for `tx-nonce` and `tx-sender` workers to expire keys and force a nonce recalibration from chain after inactivity of a sender
* Add new flag and environment variable `NONCE_CHECKER_MAX_RECOVERY` to configure max number of nonce recoveries to perform on a given envelope on `tx-sender`
* Nonce checker on `tx-sender` ignores envelopes with metadata entry `tx.mode` set to `raw`
* Add new flag and environment variable `DISABLE_EXTERNAL_TX` in the tx-listener to filter transactions not sent through Orchestrate

### 🛠 Bug fixes

* Fix connection issue when trying to connect to some Infura endpoints
* On `envelopestore`, when storing a 2nd envelope with same `tx_hash` and `chain_id` but a distinct `metadata.id` overwrites the first one
* Fix exposition of Swagger-UI in Docker images
* Fix crafting transactions with other types from uint256 and int256

### ⚠ BREAKING CHANGES

* **config** Rename `mock` options to `in-memory` for the `NONCE_MANAGER_TYPE` of the `tx-nonce` and `tx-sender`

## 1.0.2 (2019-12-18)

### 🛠 Bug fixes

* Fix issue when registering a contract with no methods and/or no events

## v1.1.0 (2019-12-10)

### 🆕 Features

* Add a new server for APIs services exposing a REST endpoint that will redirect queries to the gRPC endpoint and a swagger UI + documentation
* Add the new flag `KAFKA_CONSUMER_MAX_WAIT_TIME` to configure the maximum waiting time to consume message if messages do not exceed the size`Consumer.Fetch.Min.Byte` (default=20ms)

### 🛠 Bug fixes

* Clean logs: downgrade `OK` and `NotFound` logs to debug level in grpc server and have logger handler in debug level for tx-listener and tx-decoder
* e2e: use `CUCUMBER_STEPS_TIMEOUT` to put a timeout in cucumber steps before failing
* Makefile: add bootstrap stage to wait quorum, geth and kafka to start

### ⚠ BREAKING CHANGES

* **config** grpc & metrics server have been split. Default port for
    * grpc server remains 8080
    * newly rest server is 8081
    * metrics server has changed and is now 8082
* **config** Rename `KAFKA_SASL_ENABLE` to `KAFKA_SASL_ENABLED`

## 1.0.1 (2019-12-10)

### 🛠 Bug fixes

* Fix issue when registering a contract with methods/events having the same name and different signatures

## 1.0.0 (2019-11-07)

This is the first stable release of Orchestrate.

### 🆕 Features

* **Transaction Management**, automatically manage transactions lifecycle
    * Transaction crafting: craft transactions and deploy contract based on the bytecode in the Contract Registry
    * Faucet: get accounts credited with Ether
    * Gas management: gas price and limit can be provided or estimated by the Ethereum client
    * Transaction nonce management: set transactions nonce automatically, avoid Nonce too high scenario
    * Transaction and Event listening: process mined transactions’ receipts *at least once* and attempt to decode them with ABIs in the Contract Registry
    * External transaction signing & Key storage: store keys in HashiCorp Vault
    * External Private Transaction Signature on Besu-Orion and Quorum-Tessera.
* **Contract-Registry**: Dynamically register Smart Contract Artifacts (ABIs, bytecode & deployedBytecode) referenced by a (name, tag). Currently available in Postgres
* Connect to multiple blockchain networks simultaneously (multiple jsonRPC endpoints)
* Connect to EVM based blockchain (public/private, [Hyperledger Besu](https://github.com/hyperledger/besu), Go-Ethereum, Parity, Quorum+Tessera/Constellation)
* Add Logging using Logrus, Prometheus metrics exporter and OpenTracing exporter capabilities

### ⚠ BREAKING CHANGES

This is the list of breaking change with the Beta release.

* **config** Rename `JAEGER_DISABLED` to `JAEGER_ENABLED`
* **config** Rename `GRPC_TARGET_CONTRACT_REGISTRY` to `CONTRACT_REGISTRY_URL`
* **config** Rename `GRPC_TARGET_ENVELOPE_STORE` to `ENVELOPE_STORE_URL`
* **config** Rename `KAFKA_ADDRESS` to `KAFKA_URL`
* **config** Rename `REDIS_ADDRESS` to `REDIS_URL`
* **config** Rename `TESSERA_ENDPOINTS` to `TESSERA_URL`
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
