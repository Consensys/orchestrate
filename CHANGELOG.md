# CHANGELOG

### Version 0.1.5

*Released on March 28th 2019*

- [FIX] Add Where conditions when querying via `LoadByTxHash` and `LoadByTraceID`

### Version 0.1.4

*Released on March 26th 2019*

- [FIX] Disable `SQLAutodiscovery` for migrations to evoid `panic` when running migrations in docker container

### Version 0.1.3

*Released on March 20th 2019*

- [TEST] deactivate race detector for `infra.pg`

### Version 0.1.2

*Released on March 20th 2019*

- [TEST] deactivate race detector for `infra.pg`

### Version 0.1.1

*Released on March 20th 2019*

- [STYLE] clean linting

### Version 0.1.0

*Released on March 19th 2019*

- [FEAT] `infra` implementation of `pg`, `mock`, `grpc`
- [FEAT] implement `testutils.PGTestHelper`, `testutils.TraceStoreTestSuite`
- [FEAT] implement `app`
- [FEAT] implement `cmd`
  