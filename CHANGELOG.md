# CHANGELOG

### Version 0.2.1

*Released on February 26th 2019*

- [FIX] Dockerfile compatible with v0.2

### Version 0.2.0

*Released on February 26th 2019*

- [FEAT] re-org of app in multiple package `infra`, `worker`, `app`
- [FEAT] port app on `cobra` + `viper`
- [FEAT] add handler for `signals`
- [FIX] clean log messages

### Version 0.1.1

*Released on February 15th 2019*

- [FEAT] Update `infra/ethereum` version to integrate TxListener with Retry policy on Ethereum client fetch
- [FEAT] `config.go` Add configuration to make listener starting position configurable

### Version 0.1.0

*Released on February 12th 2019*

- [FEAT] Worker using a multi tx-listener
- [CHORE] Implement build on CI (boilerplate merge)
