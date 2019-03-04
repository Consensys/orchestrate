# CHANGELOG

### Version 0.1.9

*Released on March 4st 2019*

- [FIX] clean `Bridge` flags description 

### Version 0.1.8

*Released on March 1st 2019*

- [FEAT] add `Bridge` config and Kafka Bridge group

### Version 0.1.7

*Released on February 28th 2019*

- [FEAT] implement flag `RedisNonceExpirationTime` in config

### Version 0.1.6

*Released on February 24th 2019*

- [FEAT] implement flag `RedisLockTimeout`
  
### Version 0.1.5

*Released on February 21th 2019*

- [FEAT] implement flag `RedisLockTimeout`
  
- [FEAT] Log readable trace after unmarshalling from `Loader` 

### Version 0.1.4

*Released on February 21th 2019*

- [FIX] fix flags `worker-out`

### Version 0.1.3

*Released on February 20th 2019*

- [FEAT] log unmarshall message and errors on  `Loader` 

### Version 0.1.2

*Released on February 19th 2019*

- [FEAT] implement `SignalListener`
- [FEAT] implement flag for `ethereum`, `http`, `kafka`, `logger`, `redis`, `worker`
- [FEAT] clean config organisation based on `pflag`, `viper` & `cobra`

### Version 0.1.1

*Released on February 6th 2019*

- [FIX] generalize `TraceProducer` services into `Producer` without requiring to be a Trace

### Version 0.1.0

*Released on January 21th 2019*

- [FEAT] implement ``handlers.Loader``, ``handlers.Marker``, ``handlers.Producer``
- [FEAT] implement ``infra.TracePbMarshaller``, ``infra.TracePbnmarshaller``