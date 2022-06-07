# Orchestrate Release Notes

## v21.12.7 (2022-05-24)
### 🛠 Bug fixes
* Stop printing chain-proxy access logs errors when `ACCESSLOG_ENABLED=false`.  
* Fixed issue in `tx-sender` forwarding jwt when `KEY_MANAGER_API_KEY` is set.

### 🛠 Enhancements
* Reduced `tx-listener` services request to `orchestrate-api` can be reduced by usage of an 
optional in-memory cache. To be enabled set a duration using `API_CACHE_TTL` environment variable.
* Reduced database I/O usage by ~60%.

## v21.12.6 (2022-05-04)
### 🛠 Bug fixes
* Added missing delete account endpoint into HTTP API.
* Fixed data migration of transaction request data.
* Remove unnecessary and verbose logging.
* Update block number only after fetching 3 blocks

## v21.12.5 (2022-03-30)
### 🛠 Bug fixes
* Fixed `externalTxEnabled` chain update.

## v21.12.4 (2022-03-23)
### 🛠 Bug fixes
* Forward user's JWT token to the transaction-sender microservice

## v21.12.3 (2022-03-17)
### 🛠 Bug fixes
* Fixed panic on tx-sender updating status of jobs with `owner_id` different than nil.
* Migrated faucet DB table to use TEXT for every VARCHAR column.
* Validate existence of `creditorAccount` and `chainRule` before creating a new Faucet

## v21.12.2 (2022-02-17)
### 🛠 Bug fixes
* Fixed Chain Proxy issues caused by additional header being added to forwarded request.  
* Fixed wrong error code returned importing duplicate accounts
* Fixed missing `nonce` attribute in transaction request payloads.
* Fixed go-web3 panic error passing HEX as bytes.

## v21.12.1 (2022-01-13)
### 🆕 Features
* Compatibility with all versions of Solidity <= 0.8.11.

### ⚠ BREAKING CHANGES
* The ABI of smart contracts must now be registered in the contract registry before they can be used in transactions.
    * `contractName` is now a mandatory argument of contract transactions.
    * `contractTag` is now an optional argument of contract transactions.
* Removed authentication fallback behaviour to token `subject` when custom claims are enabled.

### 🛠 Bug fixes
* Fixed issue where Job's, Transaction's and Account's properties are removed when values are not set in update request payload. 
* Fixed ACCESS_LOG enabling/disabling toggle feature.
* Fixed lowercase ethereum addresses in response payloads.
* Fixed issue where smart contracts using Solidity structs could not be registered in the contract registry.
* Fixed unintended persisted claims over ongoing requests when using custom claims.

## v21.12.0 LTS (2021-12-16)
### 🆕 Features
* Support for `username` as additional constraint to control access over resources. Impersonation would be allowed only via API-KEY.
* Support for nested tenants in custom claims, for example tenant `tenantOne:groupOne:subGroupOne` will have
  access to resources owned by `tenantOne` and `tenantOne:groupOne` and `tenantOne:groupOne:subGroupOne` would be
  able to impersonate same tenants.
* Support Token Issuer Servers to validate JWTs. Environment variable `AUTH_JWT_ISSUER_URL`
* Support for new transaction pricing mechanism (eip-1559)
* Support for go-quorum privacy privacy enhancements: `privacyFlags`, `mandatoryFor`
* Support for go-quorum private transaction with optional `privateFrom`.
* Integration of Quorum Key Manager as replacement of Orchestrate Key Manager service
* Attach contract name and tag into transaction receipts when bytecode matches to one of the registered contracts.
* Attach contract information into transaction receipts on every new contract deployment and contract events.
* Quorum Key Manager StoreID can be defined on every account creation.

### ⚠ BREAKING CHANGES
* `Orion` was removed in favor of `EEA` as *PrivateTxManager* in chain APIs
* Following ETH transaction properties types has been BigInt updated:
    - `value` expects an HEX value prefix by "0x" instead of BigInt string.
    - `gasPrice` expects an HEX value prefix by "0x" instead of BigInt string.
    - `nonce` expects an uint64 instead of Integer string.
    - `gas` expects an uint64.
* Following Faucet request params has been modified:
    - `amount` expects an HEX value prefix by "0x" instead of BigInt string.
    - `maxBalance` expects an HEX value prefix by "0x" instead of BigInt string.
* In case of empty Orchestrate custom claims, token subject `sub` is used as `tenant_id:username`.
* Command `migration init` is merged into `migration up`.
* Removed usage of `AUTH_JWT_CERTIFICATE` in favor of `AUTH_JWT_ISSUER_URL` and `AUTH_JWT_AUDIENCE`
* Renamed `AUTH_JWT_CLAIMS_NAMESPACE` by `AUTH_JWT_ORCHESTRATE_CLAIMS`.
* In case of empty Orchestrate custom claims token subject, `sub` is used as `tenant_id`.
* Removed endpoints `/accounts/{address}/sign` and `/accounts/{address}/verify-signature` in favor of `/accounts/{address}/sign-message` and `/accounts/verify-message` accordingly to EIP-191 standards
* Removed support of zk-snarks account in favor of Quorum Key Manager implementation

### 🛠 Bug fixes
* Removed `warning` log removed when the events of the receipt are not found in the contract registry
* Fix contract deployment bug where arguments of the constructor are not parsed correctly

## v21.1.15 (2022-01-19)
### 🛠 Bug fixes
* Fixed Chain Proxy issues caused by additional header being added to forwarded request.

## v21.1.14 (2021-12-20)
### 🛠 Bug fixes
* Fix sequence of primary key when a DB copy is performed using the `copy-db` command

## v21.1.13 (2021-12-14)
### 🛠 Bug fixes
* Commit the offset to Kafka broker every time a message is processed

## v21.1.12 (2021-11-23)
### 🛠 Bug fixes
* Migrations fail when key-manager is disabled
* Key Manager fails with incorrect error code when key-manager is disabled

## v21.1.11 (2021-11-23)
### 🛠 Bug fixes
* Tx-sender exits updating jobs already in final status
* Tx-sender does not send message on `tx-recover` topic when there are persistent connectivity issues with RPC nodes

## v21.1.10 (2021-10-28)
### 🛠 Bug fixes
* Transaction `priority` is applied as expected
* Sender is not funded in raw transactions

## v21.1.9 (2021-10-21)
### 🛠 Bug fixes
* Incorrect server name verification using Postgres in `verify-ca` mode
* Tx-sender exits sending Tessera private transaction with invalid 'from'
* Added logging in key-manager microservice

## v21.1.8 (2021-08-25)
### 🛠 Bug fixes
* Wrong tenant assigment when API_KEY was not defined
* Missing decoded logs in kafka receipts for private contract events

## v21.1.7 (2021-07-06)
### 🛠 Bug fixes
* Database overload querying for registered chains

## v21.1.6 (2021-06-25)
### 🛠 Bug fixes
* Incorrect transition to FAILED status on rpc node connectivity issues
* Tx-listener do not exit when it fails to fetch private receipt from Besu node

## v21.1.5 (2021-06-02)
### 🆕 Features
* Support for metadata on chains

## v21.1.4 (2021-04-07)
### 🛠 Bug fixes
* Signing and verifying payload for zk-snarks accounts
* Hexadecimal string validation for signing endpoints

## v21.1.3 (2021-04-07)
### 🛠 Bug fixes
* Renew token with the Vault Agent where the Key Manager is watching "VAULT_TOKEN_FILE". The Key Manager supports plaintext token and wrapped-token
* Metric value for job status update CREATED to STARTED
* Improve Tx Listener performance to update transaction status to MINED

## v21.1.2 (2021-02-25)
### 🆕 Features
* New environment variable, `KAFKA_NUM_CONSUMERS`, to launch multiple kafka consumer in `tx-sender`
* Support for new Postgres setting `DB_POOL_TIMEOUT`
* Major API and DB performance improvements

### 🛠 Bug fixes
* Prevent unnecessary HTTP retries on internal API calls

## v21.1.1 (2021-02-19)

### 🛠 Bug fixes
* Hanging issue during synchronization from block 0 
* Tx-listener crashes on heavy load over API
* Import identities from connected KeyManager Vault
* Fail to send raw transaction with not empty data field

### ⚠ BREAKING CHANGES
* Rename deprecated naming from application metrics `orchestrate_transaction_scheduler_*` to `orchestrate_api_*`

## v21.1.0 LTS (2021-01-25)

### 🆕 Features

#### Orchestrate simplification
* Merge all previous APIs into a single service: `orchestate-api`, encapsulating every individual previous API services
* Merge `tx-crafter` and `tx-signer` into the `tx-sender` worker to reduce maintenance complexity
* Support usage of `in-memory` as storage for Nonce Manager

#### Identity Management API
* Release the Identity API on top of the `orchestate-api`, allowing dynamic CRUD operation over accounts whose keys are stored in Vault
* Integrate [Orchestrate HashiCorp Vault plugin](https://github.com/consensys/orchestrate-hashicorp-vault-plugin) to enhance security

#### Metrics & logging
* Add application metrics:
    * `orchestrate_transaction_scheduler_job_latency_seconds`: Histogram of job latency between status (second). Except PENDING and MINED (Histogram)
    * `orchestrate_transaction_scheduler_mined_latency_seconds` Histogram of latency between PENDING and MINED (Histogram)
    * `orchestrate_transaction_listener_current_block`: Last block processed by each listening session (Counter)
* Support for enable/disable metric modules
* Harmonize and improve logging across all services

#### Miscellaneous
* Ability set a custom keep alive interval for Postgres clients
* New environment variable `KAFKA_CONSUMER_GROUP_NAME` to set the Kafka consumer group name

### 🛠 Bug fixes
* Incorrect metrics counting for 429 http responses

### ⚠ BREAKING CHANGES
* Remove `account-generator` and `account-generated` topics
* Worker services `tx-crafter` and `tx-signer` were removed along with topics `tx-crafter` and `tx-sender`
* Jaeger reporting disabled by default
* Remove support for environment variable `ABI` to register solidity contract at start
* Remove support for environment variable `SECRET_PKEY` to import ethereum keys to key vault at start
* Remove support for environment variable `CHAIN_REGISTRY_INIT` to import chains at start
* Remove support for GRPC contract API 
* Remove API services `contract-registry`, `transaction-scheduler` and `chain-registry`
* Replace support of `kv-v2` HashiCorp engine by `orchestrate` engine.
* Environment variable `CHAIN_REGISTRY_CACHE_TTL` renamed to `PROXY_CACHE_TTL`
* Environment variable `TRANSACTION_SCHEDULER_URL` replaced by `API_URL`
* Environment variable `CONTRACT_REGISTRY_URL` replaced by `API_URL`
* Environment variable `CHAIN_REGISTRY_URL` replaced by `API_URL`

### Migrate steps from v2.5.x to v21.1.x

> IMPORTANT ! In order to perform this migration, Orchestrate has to be running on the latest minor version of v2.5.x
and been migrated to latest v21.1.x  

#### HashiCorp keys
In order to migrate your keys from `kv-v2` engine to `orchestrate` engine you need to follow the next steps:

1. Instantiate HashiCorp with both engines enabled: `kv-v2` and `orchestrate`
1. Initialize the following environment variables: 
    - `VAULT_ADDR`: HashiCorp host URL 
    - `VAULT_TOKEN_FILE`:  Disk path to token file valid for orchestrate engine
    - `VAULT_MOUNT_POINT`: Mounting point of orchestrate engine
    - `VAULT_V2_SECRET_PATH`: Path where keys are stored in kv-v2 engine 
    - `VAULT_V2_MOUNT_POINT`: Mounting point of kv-v2 engine
    - `VAULT_V2_TOKEN_FILE`:  Disk path to token file valid for kv-v2 engine
1. Execute command: 
```
$> orchestrate key-manager migrate import-secrets
```

#### Orchestrate Service Data
In previous versions of orchestrate each of the API service data was stored in a independent postgres DB. 
Therefore to update to `v21.1.x` you need to import each of service's data by following the next steps for
each of the service DBs you intend to migrate:

1. Initialize the following:
    - `DB_MIGRATION_SERVICE`: Source DB service name. Values are: "chain-registry", "transaction-scheduler" and "contract-registry"
    - `DB_MIGRATION_ADDRESS`: Source DB URL
    - `DB_MIGRATION_DATABASE`: Source DB name
    - `DB_MIGRATION_USERNAME`: Source DB username
    - `DB_MIGRATION_PASSWORD`: Source DB password
1. Execute command:
```
$> orchestrate api migrate copy-db
```
