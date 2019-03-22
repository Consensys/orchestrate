# CHANGELOG

### Version 0.3.0

*Released on March 22th 2019*

- [REFACTOR] Update types and methods on pkg
- [FEAT] add `Partitioner` on worker

### Version 0.2.5

*Released on March 5th 2019*

- [FIX] RedisNonceExpirationTime back in repository (instead of common)

### Version 0.2.4

*Released on March 1th 2019*

- [FEAT] handle case when nonce is too old in cache

### Version 0.2.3

*Released on February 26th 2019*

- [FIX] add error handling on `App.Run()`

### Version 0.2.2

*Released on February 26th 2019*

- [FIX] Dockerfile compatible with v0.2


### Version 0.2.1

- [CHORE] update dependencies

### Version 0.2.0

*Released on February 24th 2019*

- [FEAT] re-org of app in multiple package `infra`, `worker`, `app`
- [FEAT] port app on `cobra` + `viper`
- [FEAT] add handler for `signals`
- [FIX] clean log messages
  
### Version 0.1.3

*Released on February 15th 2019*

- [FIX] change env variable in config (KAFKA_GROUP_TX_NONCE)

### Version 0.1.2

*Released on February 11th 2019*

- [FEAT] add docker image build and push in gitlab CI

### Version 0.1.1

*Released on January 30th 2019*

- [FEAT] update config default output queue to tx signer

### Version 0.1.0

*Released on January 29th 2019*

- [FEAT] implement `handlers.NonceHandler`
- [FEAT] implement `handlers.Logger`
- [FEAT] implement `Config`
- [FEAT] implement `main.go` to start worker