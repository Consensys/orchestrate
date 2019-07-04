package errors

import (
	ierror "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/error"
)

// Warning  (hex code 01XXX)
const warningCode uint64 = 1 << 12

// Warning are raised to indicate
func Warning(format string, a ...interface{}) *ierror.Error {
	return Errorf(format, a...).SetCode(warningCode)
}

// IsWarning indicate whether an error is a warning
func IsWarning(err error) bool {
	return isErrorClass(FromError(err).GetCode(), warningCode)
}

// Warning retry (hex code 01100)
const retryWarningCode = warningCode + 1<<8

// RetryWarning are raised when failing to connect to a service and retrying
func RetryWarning(format string, a ...interface{}) *ierror.Error {
	return Errorf(format, a...).SetCode(retryWarningCode)
}

// Faucet warning (hex code 01200)
const faucetWarningCode = warningCode + 2<<8

// FaucetWarning are raised when a faucet credit has been denied
func FaucetWarning(format string, a ...interface{}) *ierror.Error {
	return Errorf(format, a...).SetCode(faucetWarningCode)
}

// IsFaucetWarning indicate whether an error is a faucet warning
func IsFaucetWarning(err error) bool {
	return isErrorClass(FromError(err).GetCode(), faucetWarningCode)
}

// Connection errors (hex code 08XXX)
//
// Connection errors are raised when failing to connect to an external service
const connectionErrCode uint64 = 8 << 12

// ConnectionError is raised when failing to connect to an external service
func ConnectionError(format string, a ...interface{}) *ierror.Error {
	return Errorf(format, a...).SetCode(connectionErrCode)
}

// IsConnectionError indicate whether an error is a connection error
func IsConnectionError(err error) bool {
	return isErrorClass(FromError(err).GetCode(), connectionErrCode)
}

// Kafka connection errors are raised when failing to connect to Kafka
const kafkaConnectionErrCode = connectionErrCode + 1<<8

// KafkaConnectionError is raised when failing to connect to Kafka
func KafkaConnectionError(format string, a ...interface{}) *ierror.Error {
	return Errorf(format, a...).SetCode(kafkaConnectionErrCode)
}

// HTTP connection errors are raised when failing to connect over HTTP
const httpConnectionErrCode = connectionErrCode + 2<<8

// HTTPConnectionError is raised when failing to connect over http
func HTTPConnectionError(format string, a ...interface{}) *ierror.Error {
	return Errorf(format, a...).SetCode(httpConnectionErrCode)
}

// Ethereum connection errors are raised when failing to connect to Ethereum client jsonRPC API
const ethConnectionErrCode = connectionErrCode + 3<<8

// EthConnectionError is raised when failing to connect to Ethereum client jsonRPC API
func EthConnectionError(format string, a ...interface{}) *ierror.Error {
	return Errorf(format, a...).SetCode(ethConnectionErrCode)
}

// Feature Not Supported Errors (hex code 0AXXX)
const featureNotSupportedErrCode uint64 = 10 << 12

// FeatureNotSupportedError is raised when using a feature which is not implemented
func FeatureNotSupportedError(format string, a ...interface{}) *ierror.Error {
	return Errorf(format, a...).SetCode(featureNotSupportedErrCode)
}

// Data Errors (hex code 42XXX)

// Data Errors are raised when a provided data does not match expected format
const dataErrCode uint64 = 4<<16 + 2<<12

// DataError is raised when a provided data does not match expected format  (code 03000)
func DataError(format string, a ...interface{}) *ierror.Error {
	return Errorf(format, a...).SetCode(dataErrCode)
}

// IsDataError indicate whether an error is a data error
func IsDataError(err error) bool {
	return isErrorClass(FromError(err).GetCode(), dataErrCode)
}

// Encoding errors (hex code 421XX)
const encodingErrCode = dataErrCode + 1<<8

// EncodingError are raised when failing to decode a message
func EncodingError(format string, a ...interface{}) *ierror.Error {
	return Errorf(format, a...).SetCode(encodingErrCode)
}

// Solidity Errors (hex code 422XX)
const solidityErrCode = dataErrCode + 2<<8

// SolidityError is raised when a data related in transaction crafing is incorrect
func SolidityError(format string, a ...interface{}) *ierror.Error {
	return Errorf(format, a...).SetCode(solidityErrCode)
}

// IsSolidityError indicate whether an error is a Solidity error
func IsSolidityError(err error) bool {
	return isErrorClass(FromError(err).GetCode(), solidityErrCode)
}

// Invalid Signature Error (hex code 42201)
const invalidSigErrCode = solidityErrCode + 1

// InvalidSigError is raised when a solidity method signature is invalid
func InvalidSigError(sig string) *ierror.Error {
	return Errorf("%q is an invalid Solidity method signature (example of valid signature: transfer(address,uint256))", sig).
		SetCode(invalidSigErrCode)
}

// Invalid Arg Error (hex code 42202)
const invalidArgsCountErrCode = solidityErrCode + 2

// InvalidArgsCountError is raised when invalid arguments count is provided to craft a transaction
func InvalidArgsCountError(format string, a ...interface{}) *ierror.Error {
	return Errorf(format, a...).SetCode(invalidArgsCountErrCode)
}

// Invalid Arg Error (hex code 42203)
const invalidArgErrCode = solidityErrCode + 3

// InvalidArgError is raised when invalid argument is provided to craft a transaction
func InvalidArgError(format string, a ...interface{}) *ierror.Error {
	return Errorf(format, a...).SetCode(invalidArgErrCode)
}

// Invalid topic Error (hex code 42204)
const invalidTopicsCountErrCode = solidityErrCode + 4

// InvalidTopicsCountError is raised when topics count is in receipt
func InvalidTopicsCountError(format string, a ...interface{}) *ierror.Error {
	return Errorf(format, a...).SetCode(invalidTopicsCountErrCode)
}

// Invalid EventData Error (hex code 42205)
const invalidEventDataErrCode = solidityErrCode + 5

// InvalidEventDataError is raised when event data is invalid
func InvalidEventDataError(format string, a ...interface{}) *ierror.Error {
	return Errorf(format, a...).SetCode(invalidEventDataErrCode)
}

// Format Error (hex code 423XX)
const invalidFormatErrCode = dataErrCode + 3<<8

// InvalidFormatError is raised when a data does not match an expected format
func InvalidFormatError(format string, a ...interface{}) *ierror.Error {
	return Errorf(format, a...).SetCode(invalidFormatErrCode)
}

// Storage Error (hex code DBXXX)

// Storage errors are raised when an error is encountered while accessing stored data
const storageErrCode uint64 = 13<<16 + 11<<12

// StorageError is raised when an error is encountered while accessing stored data
func StorageError(format string, a ...interface{}) *ierror.Error {
	return Errorf(format, a...).SetCode(storageErrCode)
}

// IsStorageError indicate whether an error is a storage error
func IsStorageError(err error) bool {
	return isErrorClass(FromError(err).GetCode(), storageErrCode)
}

// Data constratin violated error (hex code DB1XX)
const constraintViolatedErrCode = storageErrCode + 1<<8

// ConstraintViolatedError is raised when a data constraint has been violated
func ConstraintViolatedError(format string, a ...interface{}) *ierror.Error {
	return Errorf(format, a...).SetCode(constraintViolatedErrCode)
}

// IsConstraintViolatedError indicate whether an error is a constraint violated error
func IsConstraintViolatedError(err error) bool {
	return isErrorClass(FromError(err).GetCode(), constraintViolatedErrCode)
}

// Not found error (hex code DB2XX)
const notFoundErrCode = storageErrCode + 2<<8

// NoDataFoundError is raised when accessing a missing data
func NotFoundError(format string, a ...interface{}) *ierror.Error {
	return Errorf(format, a...).SetCode(notFoundErrCode)
}

// IsNotFoundError indicate whether an error is a no data found error
func IsNotFoundError(err error) bool {
	return isErrorClass(FromError(err).GetCode(), notFoundErrCode)
}

// Data corrupted (hex code DB3XX)
const dataCorruptedErrCode = storageErrCode + 3<<8

// DataCorruptedError is raised loading a corrupted data
func DataCorruptedError(format string, a ...interface{}) *ierror.Error {
	return Errorf(format, a...).SetCode(dataCorruptedErrCode)
}

// Configuration errors (hex code F0XXX)
//
// Configuration errors are raised when an error is encountered while loading configuration format
const configErrCode uint64 = 15 << 16

// ConfigError is raised when an error is encountered while loading configuration (code 01000)
func ConfigError(format string, a ...interface{}) *ierror.Error {
	return Errorf(format, a...).SetCode(configErrCode)
}

// Internal errors (hex code FFXXX)
//
// Configuration errors are raised when an error is encountered while loading configuration format
const internalErrCode uint64 = 15<<16 + 15<<12

// InternalError is raised when an error that is not expected to be outputes is encountered
func InternalError(format string, a ...interface{}) *ierror.Error {
	return Errorf(format, a...).SetCode(internalErrCode)
}
