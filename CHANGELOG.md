# CHANGELOG

### Version 0.4.1

*Released on May 12th 2019*

- [REFACTOR] Follow repos re-org

### Version 0.4.0

*Released on May 10th 2019*

- [REFACTOR] Major refactor

### Version 0.3.0-alpha.4

*Released on April 04th 2019*

- [CHORE] Update pkg to v0.3.3 and improve block listening and cursor update

### Version 0.3.0-alpha.3

*Released on March 24th 2019*

- [CHORE] Upgrade to `pkg` `v0.3.0-alpha.8`
- [FIX] Update trace status to `mined` after reconstituting context

### Version 0.3.0-alpha.2

*Released on March 21th 2019*

- [FIX] clean keys & topics for produced messages in kafka

### Version 0.3.0-alpha.1

*Released on March 20th 2019*

- [FEAT] implement handler `TraceLoader`
- [FEAT] implement infra `TraceStore`

### Version 0.2.2

*Released on February 27th 2019*

- [FIX] Fix logger issue (remove sarama msg based logger)

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
