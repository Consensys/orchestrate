# CHANGELOG

### Version 0.5.2

*Unreleased*

- [FEAT] Use one main producer instead of composite

### Version 0.5.1

*Released on July 5th 2019*

- [FIX] Wallet - Remove key partition producer generating wallet
- [FIX] Wallet - Avoid panic when generating wallet with an empty envelope
  
### Version 0.5.0

*Released on June 27th 2019*

- [FEAT] Update `pkg` and protos types

### Version 0.4.5

*Released on May 28th 2019*

- [FEAT] Update `tx-signer` to support vault kv2

### Version 0.4.4

*Released on May 20th 2019*

- [FEAT] Update `pkg` and protos types

### Version 0.4.3

*Released on May 13th 2019*

- [FEAT] Shrink `txctx.envelope.tx.raw` field in logger for debuging

### Version 0.4.2

*Released on May 12th 2019*

- [FIX] Bug on linter to make files compatible to 'goimports'-ed

### Version 0.4.1

*Released on May 12th 2019*

- [REFACTOR] Follow repos re-org

### Version 0.4.0

*Released on May 10th 2019*

- [REFACTOR] Major re-org

### Version 0.3.1

*Released on March 29th 2019*

- [FEAT] Sign new contract creation transaction 

### Version 0.3.0-alpha.3

*Released on March 25th 2019*

- [CHORE] Update to `keystore` `v0.3.0-alpha.5`

### Version 0.3.0-alpha.2

*Released on March 23rd 2019*

- [CHORE] Update to `pkg` `v0.3.0-alpha.8`
- [FEAT] Update signer behavior to process even if sender account is unknwon (so we support signing in Ethereum node)

### Version 0.3.0-alpha.1

*Released on March 23rd 2019*

- [REFACTOR] Mobe to `pkg`
- [FEAT] Use vaults as a secret store
- [FEAT] Stores vault credentials in AWS
- [DEV] Adds a vault container in E2E
- [FEAT] add `Partitioner` on worker
- [FEAT] add support to craft contract deployment transactions
  
### Version 0.2.3

*Unreleased*

- [DOC] Complete ```README.md``` and ```CONTRIBUTING.md```
- [FEAT] enable e2e tests adding ```docker-compose.yml``` to start Kafka locally and ```consumer``` to check that messages are published on ```tx-sender``` kafka topic

### Version 0.2.2

*Released on February 27th 2019*

- [FIX] Dockerfile compatible with v0.2.X

### Version 0.2.1

*Released on February 24th 2019*

- [CHORE] update dependencies

### Version 0.2.0

*Released on February 24th 2019*

- [FEAT] re-org of app in multiple package `infra`, `worker`, `app`
- [FEAT] port app on `cobra` + `viper`
- [FEAT] add handler for `signals`
- [FIX] clean log messages
  
### Version 0.1.1

*Released on February 11th 2019*

- [CORE] Implement build on CI (boilerplate merge)

### Version 0.1.0

*Released on January 30th 2019*

- [FEAT] implement `handlers.Signer`
- [FEAT] implement `handlers.Logger`
- [FEAT] implement `main.go`