# CHANGELOG

### Version 0.1.0

*Released on January 18th 2019*

Features

- [NEW] type ``types.TxData`` to store & manipulate transaction content data
- [NEW] type ``types.Transaction`` to store & manipulate transaction data (including transaction *raw*, *hash* and *from*)
- [NEW] protobuffer ``ethereum.TxData`` & ``ethereum.Transaction``
- [NEW] functions ``protobuf.LoadTransaction(pb *ethpb.Transaction, tx *types.Transaction)`` & ``protobuf.DumpTransaction(tx *types.Transaction, pb *ethpb.Transaction)`` to load/dump protobuffer to/from Core types
- [NEW] type ``types.Trace`` to store & manipulate trace information
- [NEW] functions ``protobuf.LoadTrace(pb *tracepb.Trace, t *types.Trace)`` & ``protobuf.DumpTrace(t *types.Trace, pb *tracepb.Trace)`` to load/dump protobuffer to/from Core types
- [NEW] services `ABIRegistry`, `Crafter`, `OffsetMarker`, `TraceProducer`, `Faucet`, `GasEstimator`, `GasPricer`, `Marshaller`, `NonceManager`, `TxSigner`,`TraceStore`,`TxSender`, `Unmarshaller`
- [NEW] type ``types.Context``
- [NEW] worker ``core.Worker`` implement worker type


Chore

- [NEW] ``scripts/generate-proto.sh``
- [NEW] ``Makefile``
- [NEW] ``scripts/coverage.sh``
- [NEW] ``gitlab-ci.yml``

### Version 0.1.1

*Released on January 18th 2019*

Doc

- [NEW] Document ``core.Worker`` in README.md
