# CHANGELOG

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