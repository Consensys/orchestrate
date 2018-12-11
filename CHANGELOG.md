# CHANGELOG

### Version 0.1.0

Features

- [NEW] type ``types.TxData`` to store & manipulate transaction content data
- [NEW] type ``types.Transaction`` to store & manipulate transaction data (including transaction *raw*, *hash* and *from*)
- [NEW] protobuffer ``ethereum.TxData`` & ``ethereum.Transaction``
- [NEW] functions ``protobuf.LoadTransaction(pb *ethpb.Transaction, tx *types.Transaction)`` & ``protobuf.DumpTransaction(tx *types.Transaction, pb *ethpb.Transaction)`` to load/dump protobuffer to/from Core types
- [NEW] function ``CraftPayload(method abi.Method, args ...interface{})`` that allows to craft a transaction payload from a method ABI and arguments