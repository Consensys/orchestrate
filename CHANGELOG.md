# CHANGELOG

### Version 0.1.8

*Released on August 7th 2019*

- [FIX] Use correct Redis command to set nonce value

### Version 0.1.7

*Released on August 7th 2019*

- [FIX] Use TTL feature to automatically refresh cached nonce value

### Version 0.1.6

*Released on August 7th 2019*

- [FIX] Do not return an error if nonce is missing in cache

### Version 0.1.5

*Released on July 19th 2019*

- [FEAT] Error Management

### Version 0.1.4

*Released on May 13th 2019*

- [FIX] Bug on linter to make files compatible to 'goimports'-ed

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

