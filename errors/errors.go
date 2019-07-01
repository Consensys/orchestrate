package errors

import (
	ierror "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/error"
)

// Warning  (hex code 01XXX)
var warningCode uint64 = 1 << 12

// Warning are raised to indicate
func Warning(format string, a ...interface{}) *ierror.Error {
	return Errorf(format, a...).SetCode(warningCode)
}

// Connection errors (hex code 08XXX)
//
// Connection errors are raised when failing to connect to an external service
var connectionErrCode uint64 = 8 << 12

// ConnectionError is raised when failing to connect to an external service
func ConnectionError(format string, a ...interface{}) *ierror.Error {
	return Errorf(format, a...).SetCode(connectionErrCode)
}

// IsConnectionError indicate whether an error is a connection error
func IsConnectionError(err error) bool {
	return is(FromError(err).GetCode(), connectionErrCode)
}

// Kafka connection errors are raised when failing to connect to Kafka
var kafkaConnectionErrCode = connectionErrCode + 1<<8

// KafkaConnectionError is raised when failing to connect to Kafka
func KafkaConnectionError(format string, a ...interface{}) *ierror.Error {
	return Errorf(format, a...).SetCode(kafkaConnectionErrCode)
}

// Kafaka connection errors are raised when failing to connect to Kafka
var httpConnectionErrCode = connectionErrCode + 2<<8

// HTTPConnectionError is raised when failing to connect over http
func HTTPConnectionError(format string, a ...interface{}) *ierror.Error {
	return Errorf(format, a...).SetCode(httpConnectionErrCode)
}

// Feature Not Supported Errors (hex code 0AXXX)
var featureNotSupportedErrCode uint64 = 10 << 12

// FeatureNotSupportedError is raised when using a feature which is not implemented
func FeatureNotSupportedError(format string, a ...interface{}) *ierror.Error {
	return Errorf(format, a...).SetCode(featureNotSupportedErrCode)
}

// Data Errors (hex code 42XXX)

// Data Errors are raised when a provided data does not match expected format
var dataErrCode uint64 = 4<<16 + 2<<12

// DataError is raised when a provided data does not match expected format  (code 03000)
func DataError(format string, a ...interface{}) *ierror.Error {
	return Errorf(format, a...).SetCode(dataErrCode)
}

// IsDataError indicate whether an error is a data error
func IsDataError(err error) bool {
	return is(FromError(err).GetCode(), dataErrCode)
}

// Encoding errors (hex code 421XX)
var encodingErrCode = dataErrCode + 1<<8

// EncodingError are raised when failing to decode a message
func EncodingError(format string, a ...interface{}) *ierror.Error {
	return Errorf(format, a...).SetCode(encodingErrCode)
}

// Solidity Errors (hex code 422XX)
var solidityErrCode = dataErrCode + 2<<8

// SolidityError is raised when a data related in transaction crafing is incorrect
func SolidityError(format string, a ...interface{}) *ierror.Error {
	return Errorf(format, a...).SetCode(solidityErrCode)
}

// IsSolidityError indicate whether an error is a Solidity error
func IsSolidityError(err error) bool {
	return is(FromError(err).GetCode(), solidityErrCode)
}

// Invalid Signature Error (hex code 42201)
var invalidSigErrCode = solidityErrCode + 1

// InvalidSigError is raised when a solidity method signature is invalid
func InvalidSigError(sig string) *ierror.Error {
	return Errorf("%q is an invalid Solidity method signature (example of valid signature: transfer(address,uint256))", sig).
		SetCode(invalidSigErrCode)
}

// Invalid Arg Error (hex code 42202)
var invalidArgsCountErrCode = solidityErrCode + 2

// InvalidArgsCountError is raised when invalid arguments count is provided to craft a transaction
func InvalidArgsCountError(format string, a ...interface{}) *ierror.Error {
	return Errorf(format, a...).SetCode(invalidArgsCountErrCode)
}

// Invalid Arg Error (hex code 42203)
var invalidArgErrCode = solidityErrCode + 3

// InvalidArgError is raised when invalid argument is provided to craft a transaction
func InvalidArgError(format string, a ...interface{}) *ierror.Error {
	return Errorf(format, a...).SetCode(invalidArgErrCode)
}

// Invalid topic Error (hex code 42204)
var invalidTopicsCountErrCode = solidityErrCode + 4

// InvalidTopicsCountError is raised when topics count is in receipt
func InvalidTopicsCountError(format string, a ...interface{}) *ierror.Error {
	return Errorf(format, a...).SetCode(invalidTopicsCountErrCode)
}

// Invalid EventData Error (hex code 42205)
var invalidEventDataErrCode = solidityErrCode + 5

// InvalidEventDataError is raised when event data is invalid
func InvalidEventDataError(format string, a ...interface{}) *ierror.Error {
	return Errorf(format, a...).SetCode(invalidEventDataErrCode)
}

// Format Error (hex code 423XX)
var invalidFormatErrCode = dataErrCode + 3<<8

// InvalidFormatError is raised when a data does not match an expected format
func InvalidFormatError(format string, a ...interface{}) *ierror.Error {
	return Errorf(format, a...).SetCode(invalidFormatErrCode)
}

// Storage Error (hex code DBXXX)

// Storage errors are raised when an error is encountered while accessing stored data
var storageErrCode uint64 = 13<<16 + 11<<12

// StorageError is raised when an error is encountered while accessing stored data
func StorageError(format string, a ...interface{}) *ierror.Error {
	return Errorf(format, a...).SetCode(storageErrCode)
}

// IsStorageError indicate whether an error is a storage error
func IsStorageError(err error) bool {
	return is(FromError(err).GetCode(), storageErrCode)
}

// No data found error (hex code DB2XX)
var noDataFoundErrCode = storageErrCode + 2<<8

// NoDataFoundError is raised when accessing a missing data
func NoDataFoundError(format string, a ...interface{}) *ierror.Error {
	return Errorf(format, a...).SetCode(noDataFoundErrCode)
}

// IsNoDataFoundError indicate whether an error is a no data found error
func IsNoDataFoundError(err error) bool {
	return is(FromError(err).GetCode(), noDataFoundErrCode)
}

// Configuration errors (hex code F0XXX)
//
// Configuration errors are raised when an error is encountered while loading configuration format
var configErrCode uint64 = 15 << 16

// ConfigError is raised when an error is encountered while loading configuration (code 01000)
func ConfigError(format string, a ...interface{}) *ierror.Error {
	return Errorf(format, a...).SetCode(configErrCode)
}

// is returns wether code belongs to base family
func is(code, base uint64) bool {
	return (base^code)&(255<<12+15<<8&base) == 0
}
