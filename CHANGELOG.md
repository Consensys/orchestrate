# CHANGELOG

### Version 0.3.7

*Released on March 30th 2019*

- [FIX] Ignore `GasPricer` and `GasEstimator` by removing ctx.Next()

### Version 0.3.6

*Released on March 29th 2019*

- [FEAT] Ignore `GasPricer` and `GasEstimator` if gas and gas price are already filled in the Trace

### Version 0.3.5

*Released on March 29th 2019*

- [FIX] No panic when receiving not well formated address

### Version 0.3.4

*Released on March 29th 2019*

- [FIX] No panic in `Crafter` if TraceT.Tx.TxData not defined

### Version 0.3.3

*Released on March 29th 2019*

- [FEAT] No credit faucet transaction if faucet address not known

### Version 0.3.2

*Released on March 29th 2019*

- [FEAT] Load bytecode from ABIRegistry for contract deployment

### Version 0.3.1

*Released on March 24th 2019*

- [FIX] Add support for go context

### Version 0.3.0

*Released on March 24th 2019*

- [CHORE] Update to `pkg` `v0.3.0-alpha.8`

### Version 0.3.0-alpha.1

*Released on March 22nd 2019*

- [REFACTOR] Mobe to `pkg`
- [FEAT] add `Partitioner` on worker
- [FEAT] add support to craft contract deployment transactions

### Version 0.2.3

*Unreleased*

- [DOC] Clean config flags description
- [DOC] Complete ```README.md``` and ```CONTRIBUTING.md```

### Version 0.2.2

*Released on February 27h 2019*

- [FIX] Dockerfile distinguish CMD from ENTRYPOINT

### Version 0.2.1

*Released on February 25h 2019*

- [FIX] Dockerfile now compatible with app launcher
  

### Version 0.2.0

*Released on February 24th 2019*

- [FEAT] re-org of app in multiple package `infra`, `worker`, `app`
- [FEAT] port app on `cobra` + `viper`
- [FEAT] add handler for `signals`
- [FIX] clean log messages
  
### Version 0.1.3

*Released on February 15th 2019*

- [FIX] Change config env variable (KAFKA_GROUP_TX_CRAFTER)

### Version 0.1.2

*Released on February 13th 2019*

- [FEAT] Implement ability to pass ABIs by environment variable


### Version 0.1.1

*Released on February 11th 2019*

- [CORE] Implement build on CI (boilerplate merge)


### Version 0.1.0

*Released on January 30th 2019*

- [FEAT] Implement `infra.NewERC1400ABIRegistry`
- [FEAT] Implement `infra.CreateFaucet`, `infra.SaramaCrediter`, `infra.FaucetConfig`
- [FEAT] Implement `main.go` 