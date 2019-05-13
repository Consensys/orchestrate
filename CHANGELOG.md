# CHANGELOG

### Version 0.4.3

*Released on May 13th 2019*

- [FIX] Shrink `txctx.raw` data length in logger for debuging

### Version 0.4.2

*Released on May 12th 2019*

- [REFACTOR] Follow repos re-org

### Version 0.4.1

*Released on May 11th 2019*

- [FIX] Correct initialization of the gRPC client in the `Sender` handler

### Version 0.4.0

*Released on May 3rd 2019*

- [REFACTOR] Refactor using new pkg version, linting & CI

### Version 0.3.4

*Released on March 20th 2019*

- [FEAT] Make `tx-sender` compatible with `PrivateFor` transactions

### Version 0.3.3

*Released on March 20th 2019*

- [FIX] add Partitioning key on worker

### Version 0.3.2

*Released on March 20th 2019*

- [FIX] clean logs
 
### Version 0.3.1

*Released on March 20th 2019*

- [CHORE] clean dependencies

### Version 0.3.0

*Released on March 20th 2019*

- [CHORE] move to `pkg`
- [FEAT] add connexion to `API-Context-Store` and store trace

### Version 0.2.2

*Released on February 27th 2019*

- [FIX] Separate CMD from ENTRYPOINT in Dockerfile

### Version 0.2.1

*Released on February 26th 2019*

- [FIX] Dockerfile compatible with v0.2

### Version 0.2.0

*Released on February 24th 2019*

- [FEAT] re-org of app in multiple package `infra`, `worker`, `app`
- [FEAT] port app on `cobra` + `viper`
- [FEAT] add handler for `signals`
- [FIX] clean log messages

### Version 0.1.1

*Released on February 11th 2019*

- [CORE] Implement build on CI (boilerplate merge)

### Version 0.1.0

*Released on January 30th 2019*

- [FEAT] Implement `handlers.Logger`, `handlers.Sender`, `handlers.Store`
- [FEAT] Implement `main.go`