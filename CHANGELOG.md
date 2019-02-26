# CHANGELOG

### Version 0.1.19

*Released on February 26th 2019*

- [FIX] add logs on worker

### Version 0.1.18

*Released on February 24th 2019*

- [FIX] fix `Errors.Error()`
 
### Version 0.1.17

*Released on February 24th 2019*

- [FEAT] Implement `Errors` type

### Version 0.1.16

*Released on February 22th 2019*

- [FIX] Add `debug` log messages in worker

### Version 0.1.15

*Released on February 21th 2019*

- [FEAT] Add `String` method every types on Trace

### Version 0.1.14

*Released on February 18th 2019*

- [FEAT] Add `Close` method on Worker
- [FEAT] Add `logrus.Logger` on `Worker`

### Version 0.1.13

*Released on February 18th 2019*

- [FEAT] update `Context` sctruct with `Logger` as a logrus Entry to store recurrent log fields in ctx

### Version 0.1.12

*Released on February 6th 2019*

- [FIX] generalize `TraceProducer` into `Producer` without requiring to be a Trace

### Version 0.1.11

*Released on February 6th 2019*

- [FIX] add `BlockNumber`, `BlockHash`, & `TxIndex` in receipt

### Version 0.1.10

*Released on January 30th 2019*

- [FEAT] add support for `context.Context` in most services

### Version 0.1.9

*Released on January 27th 2019*

- [FEAT] implement method `SetDecodedData` on `types.Log`


### Version 0.1.8

*Released on January 23th 2019*

- [FIX] add `DecodedData` field to `types.Log`


### Version 0.1.7

*Released on January 21th 2019*

- [FIX] update `services.Producer` interface
- [FIX] update `services.Marshaller`/`services.Unmarshaller` interface

### Version 0.1.6

*Released on January 21th 2019*

- [CHORE] tag verion following an error in tagging


### Version 0.1.5

*Released on January 21th 2019*

- [FIX] Make `services.Marshaller`/`services.Unmarshaller` update protobuf interface
- [FIX] Update NonceManager service (``Unlock`` method signature)

### Version 0.1.4

*Released on January 20th 2019*

- [FIX] Clean closing mechnism of ``Worker``
- [FIX] Make `services.Marshaller`/`services.Unmarshaller` agnostic from protobuf Trace


### Version 0.1.3

*Released on January 19th 2019*

- [FIX] Make ``services.FaucetRequest`` attributes exportable

### Version 0.1.2

*Released on January 19th 2019*

- [DOC] Add paragraph about concurrency when using worker
- [DOC] Add examples

### Version 0.1.1

*Released on January 18th 2019*

- [DOC] Document ``core.Worker`` in README.md

### Version 0.1.0

*Released on January 18th 2019*

- [FEAT] type ``types.TxData`` to store & manipulate transaction content data
- [FEAT] type ``types.Transaction`` to store & manipulate transaction data (including transaction *raw*, *hash* and *from*)
- [FEAT] protobuffer ``ethereum.TxData`` & ``ethereum.Transaction``
- [FEAT] functions ``protobuf.LoadTransaction(pb *ethpb.Transaction, tx *types.Transaction)`` & ``protobuf.DumpTransaction(tx *types.Transaction, pb *ethpb.Transaction)`` to load/dump protobuffer to/from Core types
- [FEAT] type ``types.Trace`` to store & manipulate trace information
- [FEAT] functions ``protobuf.LoadTrace(pb *tracepb.Trace, t *types.Trace)`` & ``protobuf.DumpTrace(t *types.Trace, pb *tracepb.Trace)`` to load/dump protobuffer to/from Core types
- [FEAT] services `ABIRegistry`, `Crafter`, `OffsetMarker`, `TraceProducer`, `Faucet`, `GasEstimator`, `GasPricer`, `Marshaller`, `NonceManager`, `TxSigner`,`TraceStore`,`TxSender`, `Unmarshaller`
- [FEAT] type ``types.Context``
- [FEAT] worker ``core.Worker`` implement worker type
- [CHORE] ``scripts/generate-proto.sh``
- [CHORE] ``Makefile``
- [CHORE] ``scripts/coverage.sh``
- [CHORE] ``gitlab-ci.yml``
