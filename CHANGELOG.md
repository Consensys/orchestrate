# CHANGELOG

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

### Version 0.1.1

*Released on January 18th 2019*

- [DOC] Document ``core.Worker`` in README.md

### Version 0.1.2

*Released on January 19th 2019*

- [DOC] Add paragraph about concurrency when using worker
- [DOC] Add examples

### Version 0.1.3

*Released on January 19th 2019*

- [FIX] Make ``services.FaucetRequest`` attributes exportable

### Version 0.1.4

*Released on January 20th 2019*

- [FIX] Clean closing mechnism of ``Worker``
- [FIX] Make `services.Marshaller`/`services.Unmarshaller` agnostic from protobuf Trace

### Version 0.1.5

*Released on January 21th 2019*

- [FIX] Make `services.Marshaller`/`services.Unmarshaller` update protobuf interface
- [FIX] Update NonceManager service (``Unlock`` method signature)

### Version 0.1.6

*Released on January 21th 2019*

- [CHORE] tag verion following an error in tagging

### Version 0.1.7

*Released on January 21th 2019*

- [FIX] update `services.Producer` interface
- [FIX] update `services.Marshaller`/`services.Unmarshaller` interface

### Version 0.1.8

*Released on January 23th 2019*

- [FIX] add `DecodedData` field to `types.Log`

### Version 0.1.9

*Released on Januarty 27th 2019*

- [FEAT] implement method `SetDecodedData` on `types.Log`