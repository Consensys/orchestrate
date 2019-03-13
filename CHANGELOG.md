# CHANGELOG

### Version 0.2.0-alpha

*Released on March 14th 2019*

- [FEAT] Update dependency to use ```pkg```
- [FIX] Clean ```abi``` logics
- [FIX] Stabilize ```tx-listener``` main loop to decrease CPU consumption
- [FEAT] Add support for ```viper``` configuration

### Version 0.1.15

*Released on February 20th 2019*

- [FIX] Add checks in `Decode` to make sure that the number of Topics corresponds to arguments indexed in the abi

### Version 0.1.14

*Released on February 15th 2019*

- [FEAT] add ability to craft transactions arrays as input

### Version 0.1.13

*Released on February 14th 2019*

- [DOC] add example
  
### Version 0.1.12

*Released on February 14th 2019*

- [FIX] Add retry backoff strategy on tx-listener when fetching client
- [FIX] Clean tx-listener logs
- [FIX] Clean tx-listener closing mechanism

### Version 0.1.11

*Released on February 12th 2019*

- [CHORE] Dependencies update to Core v0.1.12
  
### Version 0.1.11

*Released on February 12th 2019*

- [FIX] fix `ChainTracker` to use `HeaderByNumber` instead of `SyncProgress`


### Version 0.1.10

*Released on February 11th 2019*

- [FEAT] Impleemnt method `Networks` on `MultiEthClient`

### Version 0.1.9

*Released on February 11th 2019*

- [FIX] Fix typo in `MultiDial`
  
### Version 0.1.8

*Released on February 11th 2019*

- [FEAT] Implement `BlockByNumber` & `SyncProgress` on multichain client
  
### Version 0.1.7

*Released on February 11th 2019*

- [FEAT] Implement `listener.TxListener` a multichain transaction listener

### Version 0.1.6

*Released on February 11th 2019*

 - [FEAT] ``Decode`` is able to decode uint, int, bool, address and bytes arrays into string

### Version 0.1.5

*Released on February 5th 2019*

 - [FIX] add `SetPos` on tx-listener (but will be soon updated)


### Version 0.1.4

*Released on February 4th 2019*

 - [FEAT] implement first version of ``Decode`` a log decoder to string service

### Version 0.1.3

*Released on January 30th 2019*

 - [CHORE] update to core v0.1.10 and add golang context management

### Version 0.1.2

*Released on January 27th 2019*

 - [FEAT] implement ``PendingNonceAt`` and ``PendingBalanceAt`` on Multi Client
### Version 0.1.1

*Released on January 27th 2019*

- [FIX] update crafter to support `bytes32`, `bytes16`, `bytes8`, `bytes1` ABI types

### Version 0.1.0

*Released on January 25th 2019*

- [FEAT] implement ``ContractABIRegistry``
- [FEAT] implement ``PayloadCrafter``
- [FEAT] implement ``SingleChainSender`` and ``MultiChainSender``
- [FEAT] implement ``EthGasManager``
- [FEAT] implement ``StaticSigner``
- [FEAT] impement ``listener.TxListener``
- [FEAT] implement ``ethclient.MultiEthClient``
