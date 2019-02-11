# CHANGELOG

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
