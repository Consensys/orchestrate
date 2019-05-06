# CHANGELOG

### Version 0.4.8

*Released on May 6th 2019*

- [FEAT] Add function returning pg.NewOptions from config 

### Version 0.4.6

*Released on May 4th 2019*

- [FEAT] Add exported for `opentracing` handler

### Version 0.4.5

*Released on April 30th 2019*

- [FEAT] Update `opentracing` configuration

### Version 0.4.4

*Released on April 30th 2019*

- [FEAT] Dependency Injection for `http.ServerMux` in `http`

### Version 0.4.3

*Released on April 30th 2019*

- [FIX] Remove necessity to delcare `encoding` marshallers
- [FIX] Remove deprecated `engine` flags
- [REFACTOR] Update `http` organisation
- [FEAT] Implement `http.ListenAndServe`
- [FEAT] Implement `sarama.Consume`
- [TEST] Implement `e2e` test for Broker `sarama`

### Version 0.4.2

*Release on April 22th 2019*

- [FIX] Update `Chain.SetID` method

### Version 0.4.1

*Released on April 21th 2019*

- [FIX] Update `http.App` interface
- [FIX] Update *Dependency Injection* pattern

### Version 0.4.0

*Released on April 19th 2019*

- [REFACTOR] Major reorganisation 
- [FEAT] Implement  `opentracing`
- [REFACTOR] Port `infra/sarama` to `pkg`
- [REFACTOR] Update ChainID manipulation in base 10 decimal format
- [FIX] Remove `Keys` on `engine.TxContext`

### Version 0.3.4

*Released on April 18th 2019*

- [FEAT] Implement handler `Logger`

### Version 0.3.3

*Released on April 10th 2019*

- [FIX] Make `engine.handleMessage` use a `context.Context`

### Version 0.3.2

*Released on April 9th 2019*

- [FIX] Rename `exported.Use` to `exported.Register`

### Version 0.3.1

*Released on April 9th 2019*

- [FEAT] Implement a `testutils` suite for `HandlerFunc`
  
### Version 0.3.0

*Released on April 8th 2019*

- [REFACTOR] Update naming `engine.Use()` to `engine.Register()`
- [REFACTOR] Update naming `worker.Context` to `engine.TxContext`
- [REFACTOR] Update naming `trace.Trace` to `envelope.Envelope`
- [REFACTOR] Update naming `Worker` to `Engine`
- [REFACTOR] Go context usage on `engine.Context`
- [FEAT] Implement `engine.TimeoutHandler`
- [FIX] Deactivate usage of `Timeout` in worker configuration
- [FEAT] Implement *Depency Injection* pattern for `Worker`

### Version 0.2.0-alpha.13

*Released on March 30th 2019*

- [FEAT] Add `CraftConstructor` in Crafter interface

### Version 0.2.0-alpha.12

*Released on March 29th 2019*

- [FEAT] Add `common.utils.ShortString` to shrink a string size and keep the tails

### Version 0.2.0-alpha.11

*Released on March 29th 2019*

- [FIX] Not panic when pass invalid address in `common.Account.Address()`

### Version 0.2.0-alpha.10

*Released on March 29th 2019*

- [FEAT] Update `ABIRegistry` interface with `ABIRegistry.GetBytecodeByID`

### Version 0.2.0-alpha.9

*Released on March 27th 2019*

- [FEAT] Add method `SetAddress` on `Account` proto
- [FEAT] Implement proper flags for Kafka tpics

### Version 0.2.0-alpha.8

*Released on March 23th 2019*

- [FEAT] Add go `Context` support in worker
- [FEAT] Extend `protos` to support `Quorum` `privateFor` transactions
- [FIX] Register config `viper` default & env variables on config `init()`

### Version 0.2.0-alpha.7

*Released on March 18th 2019*

- [FEAT] implement `config` flags for `jaeger`
- [FIX] update flag declaration usng `viper.SetDefault`

### Version 0.2.0-alpha.6

*Released on March 17th 2019*

- [FEAT] Implement `grpc` protobuf interface for `context-store`

### Version 0.2.0-alpha.5

*Released on March 14th 2019*

- [FIX] update `core.services.ABIRegistry.RegisterContract` interface to take proto `Contract` as input

### Version 0.2.0-alpha.4

*Released on March 14th 2019*

- [FEAT] implement partitionning mechanism on worker
- [FEAT] update worker configuration setup

### Version 0.2.0-alpha.3

*Released on March 13th 2019*

- [FIX] fix `trace` package naming in `protos`

### Version 0.2.0-alpha.2

*Released on March 13th 2019*

- [FEAT] extend ```protos.abi``` to support bytecode and ABI

### Version 0.2.0-alpha.1

*Released on March 12th 2019*

- [FEAT] Proto for ```abi```

### Version 0.2.0-alpha

*Released on March 12th 2019*

- [CHORE] Merge of ```core``` and ```common```
- [FEAT] Extend protos and make it basic types
- [FEAT] Implement utility ```Keys```
- [FEAT] Implement ```config``` flags for databases

### Version 0.1.21

*Release on March 10th 2019*

- [FEAT] Add ```ID``` in trace matadata
- [FEAT] remove flag `RedisNonceExpirationTime` in config (should be moved in worker directly)
- [FIX] clean `Bridge` flags description 
- [FEAT] add `Bridge` config and Kafka Bridge group

### Version 0.1.20

*Released on February 28th 2019*

- [FEAT] NonceManager services updated (GetNonce now returns the last time the nonce was touched)
- [FEAT] implement flag `RedisNonceExpirationTime` in config

### Version 0.1.19

*Released on February 26th 2019*

- [FIX] add logs on worker

### Version 0.1.18

*Released on February 24th 2019*

- [FIX] fix `Errors.Error()`
- [FEAT] implement flag `RedisLockTimeout`

### Version 0.1.17

*Released on February 24th 2019*

- [FEAT] Implement `Errors` type

### Version 0.1.16

*Released on February 22th 2019*

- [FIX] Add `debug` log messages in worker

### Version 0.1.15

*Released on February 21th 2019*

- [FEAT] Add `String` method every types on Trace
- [FEAT] implement flag `RedisLockTimeout`
- [FEAT] Log readable trace after unmarshalling from `Loader` 

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
- [FIX] generalize `TraceProducer` services into `Producer` without requiring to be a Trace

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
- [FIX] fix flags `worker-out`
- [FEAT] implement ``handlers.Loader``, ``handlers.Marker``, ``handlers.Producer``
- [FEAT] implement ``infra.TracePbMarshaller``, ``infra.TracePbnmarshaller``

### Version 0.1.4

*Released on January 20th 2019*

- [FIX] Clean closing mechnism of ``Worker``
- [FIX] Make `services.Marshaller`/`services.Unmarshaller` agnostic from protobuf Trace
- [FEAT] log unmarshall message and errors on  `Loader`

### Version 0.1.3

*Released on January 19th 2019*

- [FIX] Make ``services.FaucetRequest`` attributes exportable
- [FEAT] implement `SignalListener`
- [FEAT] implement flag for `ethereum`, `http`, `kafka`, `logger`, `redis`, `worker`
- [FEAT] clean config organisation based on `pflag`, `viper` & `cobra`

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