# CHANGELOG

### Version 0.2.0

*Released on February 25th 2019*

- [FEAT] re-org of app in multiple package `infra`, `worker`, `app`
- [FEAT] port app on `cobra` + `viper`
- [FEAT] add handler for `signals`
- [FIX] clean log messages

### Version 0.1.7

*Released on February 20th 2019*

- [FIX] Not decoding logs without Topics


### Version 0.1.6

*Released on February 18th 2019*

- [FIX] Clean logs and use logger in `ctx.Logger`


### Version 0.1.5

*Released on February 15th 2019*

- [FIX] Correct chainID inserted in kafka topic to listen to


### Version 0.1.4

*Released on February 15th 2019*

- [FIX] Change config env variable(KAFKA_GROUP_TX_DECODER)


### Version 0.1.3

*Released on February 13th 2019*

- [FEAT] Listen multi topics depending on the chainIDs
- [FEAT] Implement ability to pass ABIs by environment variable


### Version 0.1.2

*Released on February 11th 2019*

- [CHORE] Implement build on CI (boilerplate merge)


### Version 0.1.1

*Released on February 11th 2019*

- [CHORE] Update Infra/Ethereum package to decode events with array

### Version 0.1.0

*Released on February 05th 2019*

- [FEAT] Implement `handlers.TransactionDecoder`and `handlers.LogDecoder`
