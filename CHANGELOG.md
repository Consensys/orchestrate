# CHANGELOG

### Version 0.5.5

*Released on June 12th 2019*

- [FIX] Ethclient: add missing methods in the mock
  
### Version 0.5.4

*Released on June 12th 2019*

- [CHORE] Registry: rename event selector to event sigHash

### Version 0.5.3

*Released on June 5th 2019*

- [FIX] Replace `ERC20` by `SimpleToken`

### Version 0.5.2

*Released on June 4th 2019*

- [FEAT] Replace abi registry with contract registry

### Version 0.5.1

*Released on May 20th 2019*

- [FIX] Update `abi.registry.GetBytecodeByID` without useless signature

### Version 0.5.0

*Released on May 16th 2019*

- [FEAT] Update interfaces to query methods and events in the `abi.registry`
- [FEAT] Add ERC20 ABI as default in the `abi.registry`

### Version 0.4.17

*Released on May 12th 2019*

- [FIX] Bug on linter to make files compatible to 'goimports'-ed

### Version 0.4.16

*Released on May 12th 2019*

- [REFACTOR] Following repo re-org

### Version 0.4.15

*Released on May 9th 2019*

- [FIX] Fix `listener.TxListener` to avoid panic on graceful close

### Version 0.4.14

*Released on May 9th 2019*

- [FIX] Fix `sarama.Handler.GetInitialPosition()` when no message was produce

### Version 0.4.13

*Released on May 9th 2019*

- [FIX] Update `sarama.Handler.GetInitialPosition()`

### Version 0.4.12

*Released on May 9th 2019*

- [FIX] Make `mock.Client` support multiple chains
- [FEAT] Add methods to `mock.Client` 
- [REFACTOR] Major re-org of `tx-listener`

### Version 0.4.11

*Released on May 7th 2019*

- [FIX] Fix method parsing when method is Method()

### Version 0.4.10

*Released on May 5th 2019*

- [FEAT] Decoder: able to decode slices and structs
- [REFACTOR] Extract decoder from crafter package

### Version 0.4.9

*Released on May 2nd 2019*

- [FEAT] Update client so all routes are covered with retry bakoff
- [REFACTOR] Major refactor of tx-listener

### Version 0.4.8

*Unreleased*

- [CHORE] Update to `pkg` `v0.4.4`

### Version 0.4.7

*Released on April 22th 2019*

- [FIX] Update `MultiClient.EstimateGas` interface

### Version 0.4.6

*Unreleased*

- [CHORE] Update to `pkg` `v0.4.3`

### Version 0.4.5

*Released on April 22th 2019*

- [REFACTOR] Rename `abi` packages

### Version 0.4.4

*Released on April 21th 2019*

- [CHORE] Update to `pkg` `v0.4.1`

### Version 0.4.3

*Released on April 20th 2019*

- [FIX] Update dependency injection pattern

### Version 0.4.2

*Released on April 20th 2019*

- [CHORE] Update to Geth `1.18.26`

### Version 0.4.2

*Released on April 20th 2019*

- [REFACTOR] Update to `pkg` `v0.4.0`
- [REFACTOR] Rename `MultiEthClient` to `MultiClient`
- [REFACTOR] Re-org `abi`
- [FEAT] Update `crafter.CraftConstructor`

### Version 0.4.1

*Released on April 18th 2019*

- [FEAT] Implement dependency injection

### Version 0.4.0

*Released on April 17th 2019*

- [FEAT] Update `MultiEthClient` for dynamic Dial of Ethereum clients
- [FEAT] Implement elements for *Dependy Injection* pattern for `MultiEthClient`
- [FIX] Fix infinite recursivity error in `txlistener.BlockMissingError` and `txlistener.ReceiptMissingError`
- [REFACTOR] Update `txlistener.BlockCursor` main loop to use `time.Ticker` and `trigger`
- [TEST] Enhance tests for `txlistener.BlockCursor`
- [FIX] Rollback fix from `v0.3.3`

### Version 0.3.3

*Released on April 04th 2019*

- [FIX] Update ```HighestBlock``` to correctly move forward to the cursor

### Version 0.3.2

*Released on March 30th 2019*

- [FEAT] Add `CraftConstructor` for crafting constructor arguments in contract creation

### Version 0.3.1

*Released on March 29th 2019*

- [FEAT] Get bytecode for contract creation

### Version 0.3.0

*Released on March 24th 2019*

- [FEAT] Add support for `PrivateFor` transactions
- [FEAT] Update `TxSender` interface
- [FEAT] Make `MultiDial` dialing clients concurrently for faster start
- [FIX] Clean `config`

### Version 0.2.0-alpha.2

*Released on March 15th 2019*

- [FIX] Keep common.hash as arg input of ```abi.FormatIndexedArg```

### Version 0.2.0-alpha.1

*Released on March 14th 2019*

- [FIX] Decode proto Logs instead of geth Logs in ```abi.Decode```

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
