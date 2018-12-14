# CHANGELOG

### Version 0.1.0

Features

- [NEW] type ``types.TxData`` to store & manipulate transaction content data
- [NEW] type ``types.Transaction`` to store & manipulate transaction data (including transaction *raw*, *hash* and *from*)
- [NEW] protobuffer ``ethereum.TxData`` & ``ethereum.Transaction``
- [NEW] functions ``protobuf.LoadTransaction(pb *ethpb.Transaction, tx *types.Transaction)`` & ``protobuf.DumpTransaction(tx *types.Transaction, pb *ethpb.Transaction)`` to load/dump protobuffer to/from Core types
- [NEW] function ``CraftPayload(method abi.Method, args ...interface{})`` that allows to craft a transaction payload from a method ABI and arguments
- [NEW] type ``types.Trace`` to store & manipulate trace information
- [NEW] functions ``protobuf.LoadTrace(pb *tracepb.Trace, t *types.Trace)`` & ``protobuf.DumpTrace(t *types.Trace, pb *tracepb.Trace)`` to load/dump protobuffer to/from Core types


Chore

- [NEW] ``scripts/generate-proto.sh``
- [NEW] ``Makefile``
- [NEW] ``scripts/coverage.sh``
- [NEW] ``gitlab-ci.yml``