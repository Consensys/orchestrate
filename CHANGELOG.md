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
