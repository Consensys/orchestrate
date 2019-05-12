# CHANGELOG

### Version 0.1.3

*Released on May 12th 2019*

- [REFACTOR] Follow repo re-org

### Version 0.1.2

*Released on May 2nd 2019**

- [FEAT] Add `REDIS_NONCE_EXPIRATION_TIME` variable and setting it to manage nonce expiration

### Version 0.1.0

*Released on May 1st 2019**

- [REFACTOR] Major re-org of code and migrate code from redis repository
- [FEAT] Add `NonceManager`
- [FEAT] `GetNonce` now returns the idle time of the nonce in cache (ie: time elapsed since last modif)

