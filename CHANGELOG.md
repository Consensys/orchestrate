# CHANGELOG

### Version 0.4.3

*Released on July 19th 2019*

- [FIX] Bump `pkg` version to `0.7.1`

### Version 0.4.2

*Released on July 19th 2019*

- [FEAT] Error Management

### Version 0.4.1

*Released on June 24th 2019*

- [FIX] Updated package to use new messages types from `pkg@v.0.6.1`

### Version 0.4.0

*Released on June 20th 2019*

- [FIX] Updated package to use new messages types from `pkg@v.0.6.0`

### Version 0.3.3

*Released on May 20th 2019*

- [FEAT] Update `pkg` and protos types

### Version 0.3.2

*Released on May 12th 2019*

- [FIX] Bug on linter to make files compatible to 'goimports'-ed

### Version 0.3.1

*Released on May 12th 2019*

- [REFACTOR] Follow repo re-org

### Version 0.3.0

*Released on May 6th 2019*

- [REFACTOR] Adopt to new pattern, migrate to pkg v0.4.x

### Version 0.2.1

*Released on May 2nd 2019*

- [TEST] Add postgres service to run tests

### Version 0.2.0

*Released on May 1st 2019*

- [REFACTOR] Refactor worker with new pkg version, linting & ci

### Version 0.1.5

*Released on March 28th 2019*

- [FIX] Add Where conditions when querying via `LoadByTxHash` and `LoadByID`

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
- [FEAT] implement `testutils.PGTestHelper`, `testutils.EnvelopeStoreTestSuite`
- [FEAT] implement `app`
- [FEAT] implement `cmd`
  