# CHANGELOG

### Version 0.5.1

*Released on July 19th 2019*

- [FEAT] Error management

### Version 0.5.0

*Released on June 20th 2019*

- [FIX] Updated package to use new messages types from `pkg@v.0.6.0`

### Version 0.4.4

*Released on May 31th 2019*

- [FIX] Fix the case where the app is provided an infinite TTL token

### Version 0.4.3

*Released on May 28th 2019*

- [FEAT] Support for vault kv2
- [FEAT] Support for token refreshment
- [FEAT] Support for passing token through file 

### Version 0.4.2

*Released on May 12th 2019*

- [FIX] Bug on linter to make files compatible to 'goimports'-ed

### Version 0.4.1

*Released on May 12th 2019*

- [REFACTOR] Follow repo re-org

### Version 0.4.0

*Released on May 7th 2019*

- [REFACTOR] Re-org script
- [FEAT] Add Hashicorp vault secret store 

### Version 0.3.0-alpha.5

*Released on March 25th 2019*

- [REFACTOR] Re-org script
- [FEAT] Implement a `keystore.NewKeyStore`
- [TEST] Add test for `secretstore` and `keystore`

### Version 0.2.1

*Released on March 22nd 2019*

- [FEAT] Implement `InitFlags` for Hashicorpt vault

### Version 0.2.0

*Released on March 22nd 2019*

- [FEAT] Implement `KeyStore` interface
- [FEAT] Implement Hashicorps `KeyStore`
- [FEAT] Impelment configurations & flags
- [FEAT] PreRegister private keys from configuration
- [FEAT] Get Credentials from AWS