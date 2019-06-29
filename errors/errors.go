package errors

import (
	err "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/error"
)

// Warning  (hex code 01XXX)
var warningCode uint64 = 1 << 12

// Warning are raised to indicate
func Warning(msg string) *err.Error {
	return err.New(msg).SetCode(warningCode)
}

// Connection errors (hex code 08XXX)
//
// Connection errors are raised when failing to connect to an external service
var connectionErrCode uint64 = 8 << 12

// ConnectionError is raised when failing to connect to an external service
func ConnectionError(msg string) *err.Error {
	return err.New(msg).SetCode(connectionErrCode)
}

// Data Errors (hex code 42XXX)

// Data Errors are raised when a provided data does not match expected format
var dataErrCode uint64 = 4<<16 + 2<<12

// DataError is raised when a provided data does not match expected format  (code 03000)
func DataError(msg string) *err.Error {
	return err.New(msg).SetCode(dataErrCode)
}

// IsDataError indicate whether an error is a data error
func IsDataError(e *err.Error) bool {
	return is(e.GetCode(), dataErrCode)
}

// Encoding errors (hex code 421XX)
var encodingErrCode = dataErrCode + 1<<8

// EncodingError are raised when failing to decode a message
func EncodingError(msg string) *err.Error {
	return err.New(msg).SetCode(encodingErrCode)
}

// Solidity Errors (hex code 422XX)
var solidityErrCode = dataErrCode + 2<<8

// SolidityError is raised when a data related in transaction crafing is incorrect
func SolidityError(msg string) *err.Error {
	return err.New(msg).SetCode(solidityErrCode)
}

// IsSolidityError indicate whether an error is a Solidity error
func IsSolidityError(e *err.Error) bool {
	return is(e.GetCode(), solidityErrCode)
}

// Invalid Signature Error (hex code 42201)
var invalidSigErrCode = solidityErrCode + 1

// InvalidSigError is raised when a solidity method signature is invalid
func InvalidSigError(sig string) *err.Error {
	return err.Errorf("%q is an invalid Solidity method signature (example of valid signature: transfer(address,uint256))", sig).
		SetCode(invalidSigErrCode)
}

// Format Error (hex code 423XX)
var invalidFormatErrCode = dataErrCode + 3<<8

// InvalidFormatError is raised when a data does not match an expected format
func InvalidFormatError(msg string) *err.Error {
	return err.New(msg).
		SetCode(invalidFormatErrCode)
}

// InvalidFormatErrorf is raised when a data does not match an expected format
func InvalidFormatErrorf(format string, a ...interface{}) *err.Error {
	return err.Errorf(format, a...).
		SetCode(invalidFormatErrCode)
}

// Configuration errors (hex code F0XXX)
//
// Configuration errors are raised when an error is encountered while loading configuration format
var configErrCode uint64 = 15 << 16

// ConfigError is raised when an error is encountered while loading configuration (code 01000)
func ConfigError(msg string) *err.Error {
	return err.New(msg).SetCode(configErrCode)
}

// is returns wether code belongs to base family
func is(code, base uint64) bool {
	return (base^code)&(255<<12+15<<8&base) == 0
}
