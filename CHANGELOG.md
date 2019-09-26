# CHANGELOG

### Version 0.4.3

*Released on September 26th 2019*

- [FEAT] Nonce Management v2
- [FEAT] Contract-Registry
- [FEAT] Wallet generator with Faucet
- [FIX] Avoid blocking cucumber scenario executions when not consuming from every topic
 
### Version 0.4.2

*Unreleased*

- [FEAT] Add a new scenario to test deployment of private transactions with Quorum

### Version 0.4.1

*Released on August 25th 2019*

- [FIX] Update `sarama` to v1.22.1 to auto-create topics when start to consuming 

### Version 0.4.0

*Released on August 11th 2019*

- [FEAT] Update `pkg` version to `v0.8.0
- [FEAT] Sarama - Kafka connection with TLS/SASL`

### Version 0.3.0

*Released on July 27th 2019*

- [FIX] Update `pkg` version to `v0.7.2`

### Version 0.2.0

*Released on June 27th 2019*

- [FIX] Update `pkg` version to `v0.6.1`

### Version 0.1.6

*Released on June 18th 2019*

- [FIX] Update `CUCUMBER_CHAINID_PRIMARY` and `CUCUMBER_CHAINID_SECONDARY` to string
- [FEAT] Increase the timeout by steps to 60 seconds by default
- 
### Version 0.1.5

*Released on June 17th 2019*

- [FIX] Update `pkg` version to `v0.5.6`
- 
### Version 0.1.4

*Released on June 17th 2019*

- [FIX] Add `AliaisChainID` to run scenario tests in any chains
- [FIX] Remove `ethclient`
  
### Version 0.1.3

*Released on June 12th 2019*

- [FEAT] Add tests and clean code
- [FEAT] Add stardard steps, utils as the `EnvelopeCrafter`
  
### Version 0.1.2

*Released on June 06th 2019*

- [FIX] Update Dockerfile to parse features files

### Version 0.1.1

*Released on June 05th 2019*

- [FEAT] Generate cucumber `.json` report by default
- [FEAT] Improve sceneario steps requiring a contract deployment

### Version 0.1.0

*Released on June 04th 2019*

- [FEAT] Add initial `deployment` contract test
- [FEAT] Implement `cucumber` BDD and basic steps
- [FEAT] Implement `chanregistry` to make the bridge between the consumer group and cucumber